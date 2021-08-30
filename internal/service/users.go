package service

import (
	"context"
	"encoding/json"
	"errors"
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
	UserID string
	Email  string
	Name   string
	Roles  []string
	Groups []string
}

func (p *Provider) SaveUserAuth(ctx context.Context, state string, code string) error {
	oauth2Token, err := p.oidc.Exchange(ctx, code)
	if err != nil {
		return err
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return errors.New("failed to get id_token")
	}

	idToken, err := p.oidc.Verify(ctx, rawIDToken)
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
	ue := &UserEntry{
		UserID: userID,
		Email:  claims["email"].(string),
		Name:   claims["preferred_username"].(string),
		Roles:  extractClaimSlice(claims["roles"]),
		Groups: extractClaimSlice(claims["groups"]),
	}

	log.Logger.Debug("user entry data", zap.Any("user_entry", ue))
	av, err := dynamodbattribute.MarshalMap(ue)
	if err != nil {
		return err
	}

	userEntry, err := p.userDB.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("eve-bot-users"),
	})
	if err != nil {
		return err
	}
	log.Logger.Debug("saved user entry", zap.Any("user_entry", userEntry))
	return nil
}

func (p *Provider) isAuthorized(cmd commands.EvebotCommand, userEntry *UserEntry) bool {
	if userEntry.isAdmin() {
		return true
	}
	if userEntry.canWriteAll() {
		return true
	}
	if userEntry.validEnvironment(extractEnv(cmd.Options())) {
		return true
	}
	return false
}

// validEnvironment check
// TODO: refactor this with better RBAC strategy
// Access on actions (deploy, release, set, delete, etc.)
// Access on environment (int,dev,qa,stage,perf,prod)
func (e UserEntry) validEnvironment(env string) bool {
	if env == "" {
		return true
	}
	for _, role := range e.Roles {
		if strings.ToLower(role) == "write-nonprod" {
			// Let's see if they are performing an action to something in the lower environments (int,qa,dev)
			// Most actions can be taken against resources in the lower environments
			// the only action that can't is the `release` command
			switch {
			case strings.Contains(env, "int"), strings.Contains(env, "qa"), strings.Contains(env, "dev"):
				return true
			}
		}
	}
	return false
}

// canWriteAll check
// TODO: refactor this with better RBAC strategy
// Access on actions (deploy, release, set, delete, etc.)
// Access on environment (int,dev,qa,stage,perf,prod)
func (e UserEntry) canWriteAll() bool {
	for _, role := range e.Roles {
		if strings.ToLower(role) == "write-all" {
			return true
		}
	}
	return false
}

func (e UserEntry) isAdmin() bool {
	for _, role := range e.Roles {
		if strings.ToLower(role) == "admin" {
			return true
		}
	}
	return false
}
