package api

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/unanet/eve-bot/internal/service"
	"github.com/unanet/go/pkg/log"
	"github.com/unanet/go/pkg/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// AuthController is the Controller/Handler for ping routes
type AuthController struct {
	svc   *service.Provider
}

// NewAuthController creates a new OIDC controller
func NewAuthController(svc *service.Provider) *AuthController {
	return &AuthController{
		svc:   svc,
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

	// Create struct to hold info about new item
	type UserEntry struct {
		UserID string
		Email  string
		Name   string
		Roles  []string
		Groups []string
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

	//_ = TokenResponse{
	//	AccessToken:  oauth2Token.AccessToken,
	//	RefreshToken: oauth2Token.RefreshToken,
	//	TokenType:    oauth2Token.TokenType,
	//	Expiry:       oauth2Token.Expiry,
	//	Claims:       idTokenClaims,
	//}

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
