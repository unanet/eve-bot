package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/unanet/eve-bot/internal/botcommander/commands"
	errs "github.com/unanet/go/pkg/errors"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"
)

type UserStore interface {
	SaveUserAuth(ctx context.Context, state string, code string) error
	ReadUser(userID string) (*UserEntry, error)
}

// UserEntry struct to hold info about new user item
type UserEntry struct {
	UserID   string
	Name     string
	Roles    []string
	Groups   []string
	MapRoles map[string]bool
}

func (p *Provider) SaveUserAuth(ctx context.Context, state string, code string) error {
	oauth2Token, err := p.oidc.Exchange(ctx, code)
	if err != nil {
		return err
	}

	log.Logger.Info("oauth token details",
		zap.Any("AccessToken", oauth2Token.AccessToken),
		zap.Any("Expiry", oauth2Token.Expiry),
		zap.Any("TokenType", oauth2Token.TokenType),
		zap.Any("RefreshToken", oauth2Token.RefreshToken),
	)

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return errors.New("failed to get id_token")
	}

	idToken, err := p.oidc.Verify(ctx, rawIDToken)
	if err != nil {
		return err
	}

	log.Logger.Info("oauth idToken details",
		zap.Any("idToken", idToken),
		zap.Any("Subject", idToken.Subject),
		zap.Any("Nonce", idToken.Nonce),
		zap.Any("Issuer", idToken.Issuer),
	)

	var idTokenClaims = new(json.RawMessage)
	err = idToken.Claims(&idTokenClaims)
	if err != nil {
		return err
	}
	var claims = make(map[string]interface{})
	b, err := idTokenClaims.MarshalJSON()
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &claims)
	if err != nil {
		return err
	}

	log.Logger.Info("oauth claims",
		zap.Any("claims", claims),
	)

	return p.saveUser(state, claims)
}

func (p *Provider) ReadUser(userID string) (*UserEntry, error) {
	log.Logger.Info("service provider read user", zap.String("user_id", userID))
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if u, ok := p.userCache[userID]; ok {
		log.Logger.Info("user in cache", zap.Any("user", u))
		return &u, nil
	}
	result, err := p.userDB.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(p.Cfg.UserTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(userID),
			},
		},
	})
	if err != nil {
		log.Logger.Error("failed to get user item", zap.Error(err))
		return nil, err
	}
	if result != nil && result.Item != nil {
		log.Logger.Debug("user exists in db", zap.Any("res", result))
		entry := UserEntry{}
		err = dynamodbattribute.UnmarshalMap(result.Item, &entry)
		if err != nil {
			return nil, err
		}
		log.Logger.Debug("setting user cache from db", zap.Any("user_entry", entry))
		p.userCache[entry.UserID] = entry
		return &entry, nil
	}
	return nil, errs.ErrNotFound
}

func (p *Provider) saveUser(userID string, claims map[string]interface{}) error {
	log.Logger.Info("save user with claims", zap.Any("claims", claims))

	ue := &UserEntry{
		UserID:   userID,
		Name:     claims["preferred_username"].(string),
		Roles:    extractClaimSlice(claims["roles"]),
		Groups:   extractClaimSlice(claims["groups"]),
		MapRoles: extractClaimMap(claims["roles"]),
	}

	log.Logger.Debug("user entry data", zap.Any("user_entry", ue))
	av, err := dynamodbattribute.MarshalMap(ue)
	if err != nil {
		return err
	}

	userEntry, err := p.userDB.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(p.Cfg.UserTableName),
	})
	if err != nil {
		return err
	}
	log.Logger.Debug("saved user entry", zap.Any("user_entry", userEntry))
	return nil
}

// TODO: Setup a more "polished" RBAC strategy
// Want to be able to map incoming/dowstream groups with Roles in our system
func (p *Provider) IsAuthorized(cmd commands.EvebotCommand, userEntry *UserEntry) bool {
	// always allow the user to authenticate explicitly (re-login)
	if strings.ToLower(cmd.Info().CommandName) == "auth" {
		return true
	}
	if userEntry.isAdmin() {
		return true
	}
	if enabled, ok := userEntry.MapRoles[requestedRole(cmd)]; enabled && ok {
		return true
	}
	return false
}

func requestedRole(cmd commands.EvebotCommand) string {
	tmpRequestedRoleCmd := fmt.Sprintf("eve-%s", cmd.Info().CommandName)
	if strings.Contains(strings.ToLower(extractEnv(cmd.Options())), "prod") {
		tmpRequestedRoleCmd = tmpRequestedRoleCmd + "-prod"
	}
	return tmpRequestedRoleCmd
}

func (e UserEntry) isAdmin() bool {
	for _, role := range e.Roles {
		if strings.Contains(strings.ToLower(role), "admin") {
			return true
		}
	}
	return false
}
