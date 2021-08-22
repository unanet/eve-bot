package api

import (
	"encoding/json"
	"github.com/go-chi/jwtauth"
	"github.com/unanet/go/pkg/log"
	"github.com/unanet/go/pkg/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/unanet/eve-bot/internal/manager"
	"github.com/unanet/go/pkg/identity"
)

// AuthController is the Controller/Handler for ping routes
type AuthController struct {
	state string
	oidc  *identity.Service
}

// NewAuthController creates a new OIDC controller
func NewAuthController(mgr *manager.Service) *AuthController {
	return &AuthController{
		state: "somestate",
		oidc:  mgr.OpenIDService(),
	}
}

// Setup satisfies the EveController interface for setting up the
func (c AuthController) Setup(r chi.Router) {
	r.Get("/oidc/callback", c.callback)
	r.Get("/auth", c.auth)
}

func (c AuthController) auth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	unknownToken := jwtauth.TokenFromHeader(r)

	if len(unknownToken) == 0 {
		middleware.Log(ctx).Debug("unknown token")
		http.Redirect(w, r, c.oidc.AuthCodeURL(c.state), http.StatusFound)
		return
	}

	verifiedToken, err := c.oidc.Verify(ctx, unknownToken)
	if err != nil {
		middleware.Log(ctx).Debug("invalid token")
		http.Redirect(w, r, c.oidc.AuthCodeURL(c.state), http.StatusFound)
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

	//render.Respond(w, r, fmt.Sprintf("hello %s %s", verifiedToken.Subject, verifiedToken.Audience))
}

func (c AuthController) callback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != c.state {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
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

	//"roles":["default-roles-devops","admin"],
	//"name":"Troy Sampson",
	//"groups":["default-roles-devops","admin"],
	//"preferred_username":"tsampson@plainsight.ai",
	//"given_name":"Troy",
	//"family_name":"Sampson",
	//"email":"tsampson@plainsight.ai",}
	var claims = make(map[string]interface{})
	b,err := idTokenClaims.MarshalJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(b,&claims)
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
		zap.Any("name",name),
		zap.Any("groups",groups),
		zap.Any("roles",roles),
		zap.Any("username",username),
		zap.Any("email",email),
	)

	render.JSON(w, r, TokenResponse{
		AccessToken:  oauth2Token.AccessToken,
		RefreshToken: oauth2Token.RefreshToken,
		TokenType:    oauth2Token.TokenType,
		Expiry:       oauth2Token.Expiry,
		Claims:       idTokenClaims,
	})
}


type TokenResponse struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	TokenType    string           `json:"token_type"`
	Expiry       time.Time        `json:"expiry"`
	Claims       *json.RawMessage `json:"claims"`
}
