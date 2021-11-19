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
	UserID  string
	Name    string
	Roles   map[string]bool
	IsAdmin bool
}

func (p *Provider) SaveUserAuth(ctx context.Context, state string, code string) error {
	oauth2Token, err := p.Exchange(ctx, code)
	if err != nil {
		return err
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return errors.New("failed to get id_token")
	}

	idToken, err := p.Verify(ctx, rawIDToken)
	if err != nil {
		return err
	}

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

	log.Logger.Debug("oauth claims", zap.Any("claims", claims))

	return p.saveUser(state, claims)
}

func (p *Provider) ReadUser(userID string) (*UserEntry, error) {
	log.Logger.Info("service provider read user", zap.String("user_id", userID))
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
		return &entry, nil
	}
	return nil, errs.ErrNotFound
}

func (p *Provider) saveUser(userID string, claims map[string]interface{}) error {
	log.Logger.Info("save user with claims", zap.Any("claims", claims))

	ue := &UserEntry{
		UserID:  userID,
		Name:    claims["preferred_username"].(string),
		Roles:   extractClaimMap(claims["roles"]),
		IsAdmin: extractIsAdminRole(claims["roles"]),
	}

	log.Logger.Debug("user entry data", zap.Any("user_entry", ue))
	av, err := dynamodbattribute.MarshalMap(ue)
	if err != nil {
		return err
	}

	_, err = p.userDB.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(p.Cfg.UserTableName),
	})

	return err
}

// TODO: Setup a more "polished" RBAC strategy
// Want to be able to map incoming/dowstream groups with Roles in our system
func (p *Provider) IsAuthorized(cmd commands.EvebotCommand, userEntry *UserEntry) bool {
	reqRole := requestedRole(cmd)
	log.Logger.Info("auth check",
		zap.String("user", cmd.Info().User),
		zap.String("command", cmd.Info().CommandName),
		zap.String("requested_role", reqRole),
		zap.Any("user_entry", userEntry),
	)
	// If User has matching role (and enabled) let them pass
	if enabled, ok := userEntry.Roles[reqRole]; enabled && ok {
		return true
	}
	// Check if admin or cmd is Help, Root or Auth command
	return (cmd.Info().IsHelpRequest || cmd.Info().IsRootCmd || cmd.Info().IsAuthCmd || userEntry.IsAdmin)
}

func requestedRole(cmd commands.EvebotCommand) string {
	tmpRequestedRoleCmd := fmt.Sprintf("eve-%s", cmd.Info().CommandName)
	if strings.Contains(strings.ToLower(extractEnv(cmd.Options())), "prod") {
		tmpRequestedRoleCmd = tmpRequestedRoleCmd + "-prod"
	}
	return tmpRequestedRoleCmd
}
