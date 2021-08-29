package api

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/unanet/eve-bot/internal/service"
	"github.com/unanet/go/pkg/errors"
	"github.com/unanet/go/pkg/log"
	"github.com/unanet/go/pkg/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// AuthController is the Controller/Handler for ping routes
type AuthController struct {
	svc *service.Provider
}

// NewAuthController creates a new OIDC controller
func NewAuthController(svc *service.Provider) *AuthController {
	return &AuthController{
		svc: svc,
	}
}

// Setup satisfies the EveController interface for setting up the
func (c AuthController) Setup(r chi.Router) {
	r.Get("/oidc/callback", c.callback)
	r.Get("/signed-in", c.successfulSignIn)
	r.Get("/auth", c.auth)
}

func (c AuthController) successfulSignIn(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("<!doctype html>\n\n<html lang=\"en\">\n<head>\n <script language=\"javascript\" type=\"text/javascript\">\nfunction windowClose() {\nwindow.open('','_parent','');\nwindow.close();\n}\n</script> <meta charset=\"utf-8\">\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n  <title>Successful Auth</title>\n</head>\n<body>\n  \t<p> You have successfully Signed In. You may close this windows</p>\n</body>\n</html>"))
}

func (c AuthController) auth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	unknownToken := jwtauth.TokenFromHeader(r)

	if len(unknownToken) == 0 {
		middleware.Log(ctx).Debug("unknown token")
		http.Redirect(w, r, c.svc.MgrSvc.AuthCodeURL("empty"), http.StatusFound)
		return
	}

	verifiedToken, err := c.svc.MgrSvc.OpenIDService().Verify(ctx, unknownToken)
	if err != nil {
		middleware.Log(ctx).Debug("invalid token")
		http.Redirect(w, r, c.svc.MgrSvc.AuthCodeURL("empty"), http.StatusFound)
		return
	}

	var idTokenClaims = new(json.RawMessage)
	if err := verifiedToken.Claims(&idTokenClaims); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, TokenResponse{
		AccessToken: unknownToken,
		Expiry:      verifiedToken.Expiry,
		Claims:      idTokenClaims,
	})
}

func (c AuthController) callback(w http.ResponseWriter, r *http.Request) {
	incomingState := r.URL.Query().Get("state")
	log.Logger.Info("incoming state", zap.Any("state", incomingState))

	ctx := r.Context()

	oauth2Token, err := c.svc.MgrSvc.OpenIDService().Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		render.Respond(w, r, err)
		return
	}

	idToken, err := c.svc.MgrSvc.OpenIDService().Verify(ctx, rawIDToken)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	var idTokenClaims = new(json.RawMessage)
	if err := idToken.Claims(&idTokenClaims); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var claims = make(map[string]interface{})
	b, err := idTokenClaims.MarshalJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(b, &claims)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Logger.Info("incoming claims data", zap.Any("claims", claims))

	// Create struct to hold info about new item
	type UserEntry struct {
		UserID string
		Email  string
		Name   string
		Roles  []string
		Groups []string
	}

	ue := &UserEntry{
		UserID: incomingState,
		Email:  claims["email"].(string),
		Name:   claims["preferred_username"].(string),
	}

	switch v := claims["roles"].(type) {
	case []string:
		ue.Roles = v
	case []interface{}:
		log.Logger.Info("incoming claims roles slice of interfaces", zap.Any("v", v))
	case interface{}:
		log.Logger.Info("incoming claims roles interface", zap.Any("v", v))
	default:
		log.Logger.Info("incoming claims roles unknown", zap.Any("v", v))
	}

	switch g := claims["groups"].(type) {
	case []string:
		ue.Groups = g
	case []interface{}:
		log.Logger.Info("incoming claims groups slice of interfaces", zap.Any("g", g))
	case interface{}:
		log.Logger.Info("incoming claims groups interface", zap.Any("g", g))
	default:
		log.Logger.Info("incoming claims groups unknown", zap.Any("g", g))
	}

	log.Logger.Info("user data", zap.Any("user_entry", ue))

	result, err := c.svc.UserDB.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("eve-bot-users"),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String(ue.UserID),
			},
			"Email": {
				S: aws.String(ue.Email),
			},
		},
	})
	if err != nil {
		log.Logger.Error("failed to get user db entry", zap.Any("user_entry", ue), zap.Error(err))
	}

	log.Logger.Info("read bot user result", zap.Any("res", result))

	// User does not exist (lets create an entry)
	if result == nil || result.Item == nil {
		av, err := dynamodbattribute.MarshalMap(ue)
		if err != nil {
			render.Respond(w, r, errors.Wrap(err))
			return
		}
		_, err = c.svc.UserDB.PutItem(&dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String("eve-bot-users"),
		})
		if err != nil {
			render.Respond(w, r, errors.Wrap(err))
			return
		}
	} else {
		entry := UserEntry{}
		err = dynamodbattribute.UnmarshalMap(result.Item, &entry)
		if err != nil {
			render.Respond(w, r, errors.Wrap(err))
			return
		}
	}

	//render.JSON(w, r, TokenResponse{
	//	AccessToken:  oauth2Token.AccessToken,
	//	RefreshToken: oauth2Token.RefreshToken,
	//	TokenType:    oauth2Token.TokenType,
	//	Expiry:       oauth2Token.Expiry,
	//	Claims:       idTokenClaims,
	//})

	// Just redirecting to a different page to prevent id refresh (which throws an error)
	http.Redirect(w, r, "/signed-in", http.StatusFound)
	return

}

type TokenResponse struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	TokenType    string           `json:"token_type"`
	Expiry       time.Time        `json:"expiry"`
	Claims       *json.RawMessage `json:"claims"`
}
