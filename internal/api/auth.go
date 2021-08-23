package api

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/unanet/eve-bot/internal/manager"
	"github.com/unanet/go/pkg/identity"
)

// AuthController is the Controller/Handler for ping routes
type AuthController struct {
	state    string
	oidc     *identity.Service
	dynamoDB *dynamodb.DynamoDB
}

// NewAuthController creates a new OIDC controller
func NewAuthController(mgr *manager.Service, svc *dynamodb.DynamoDB) *AuthController {
	return &AuthController{
		state:    "somestate",
		oidc:     mgr.OpenIDService(),
		dynamoDB: svc,
	}
}

// Setup satisfies the EveController interface for setting up the
func (c AuthController) Setup(r chi.Router) {
	r.Get("/oidc/callback", c.callback)
	//r.Get("/auth", c.auth)
}

//func (c AuthController) auth(w http.ResponseWriter, r *http.Request) {
//	ctx := r.Context()
//	unknownToken := jwtauth.TokenFromHeader(r)
//
//	if len(unknownToken) == 0 {
//		middleware.Log(ctx).Debug("unknown token")
//		http.Redirect(w, r, c.oidc.AuthCodeURL(c.state), http.StatusFound)
//		return
//	}
//
//	verifiedToken, err := c.oidc.Verify(ctx, unknownToken)
//	if err != nil {
//		middleware.Log(ctx).Debug("invalid token")
//		http.Redirect(w, r, c.oidc.AuthCodeURL(c.state), http.StatusFound)
//		return
//	}
//
//	var idTokenClaims = new(json.RawMessage)
//	if err := verifiedToken.Claims(&idTokenClaims); err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	render.JSON(w, r, TokenResponse{
//		AccessToken: unknownToken,
//		Expiry:      verifiedToken.Expiry,
//		Claims:      idTokenClaims,
//	})
//}

func (c AuthController) callback(w http.ResponseWriter, r *http.Request) {
	incomingState := r.URL.Query().Get("state")
	if incomingState != c.state {
		log.Logger.Info("mismatching state",
			zap.Any("incoming_state", incomingState),
			zap.Any("expected_state", c.state),
		)
		//http.Error(w, "state did not match", http.StatusBadRequest)
		//return
	}

	ctx := r.Context()

	oauth2Token, err := c.oidc.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		render.Respond(w, r, err)
		return
	}

	idToken, err := c.oidc.Verify(ctx, rawIDToken)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	var idTokenClaims = new(json.RawMessage)
	if err := idToken.Claims(&idTokenClaims); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create struct to hold info about new item
	type UserEntry struct {
		UserID string
		Email  string
		Name   string
		Roles  []string
		Groups []string
	}

	//"roles":["default-roles-devops","admin"],
	//"name":"Troy Sampson",
	//"groups":["default-roles-devops","admin"],
	//"preferred_username":"tsampson@plainsight.ai",
	//"given_name":"Troy",
	//"family_name":"Sampson",
	//"email":"tsampson@plainsight.ai",}
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

	// TODO: Store these in Dynamo Table
	name := claims["name"]
	groups := claims["groups"]
	roles := claims["roles"]
	username := claims["preferred_username"]
	email := claims["email"]

	log.Logger.Info("user data",
		zap.Any("name", name),
		zap.Any("groups", groups),
		zap.Any("roles", roles),
		zap.Any("username", username),
		zap.Any("email", email),
	)

	body := []byte(`
<!doctype html>

<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Successful Auth</title>
</head>
<body>
  	<p> You have successfully Logged In </p>
	<p> You may <a href="#" onclick="close_window();return false;">close the tab</a>  </p>
</body>
</html>
`)

	bodyBytes := ioutil.NopCloser(bytes.NewReader(body))
	r.Body = bodyBytes
	r.ContentLength = int64(len(body))
	r.Header.Set("Content-Length", strconv.Itoa(len(body)))
	r.Header.Set("Content-Type", "text/html")
	render.Respond(w,r,"successful sign-in - close the window")
	return
	//r.Header.Set("Content-Length", len(body))
	//
	//render.SetContentType(render.ContentTypeHTML)
	//render.Respond(w,r, )

	//render.JSON(w, r, TokenResponse{
	//	AccessToken:  oauth2Token.AccessToken,
	//	RefreshToken: oauth2Token.RefreshToken,
	//	TokenType:    oauth2Token.TokenType,
	//	Expiry:       oauth2Token.Expiry,
	//	Claims:       idTokenClaims,
	//})
}

type TokenResponse struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	TokenType    string           `json:"token_type"`
	Expiry       time.Time        `json:"expiry"`
	Claims       *json.RawMessage `json:"claims"`
}
