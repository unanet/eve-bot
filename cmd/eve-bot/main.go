package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"gitlab.unanet.io/devops/eve/pkg/log"
	"gitlab.unanet.io/devops/eve/pkg/mux"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve-bot/internal/api"
)

func main() {
	api, err := mux.NewApi(api.Controllers, mux.Config{
		Port:        8080,
		MetricsPort: 3000,
		ServiceName: "eve-bot",
	})
	if err != nil {
		log.Logger.Panic("Failed to Create Api App", zap.Error(err))
	}
	api.Start()
}

// GenerateNonce is is to generate random bytes for Okta
func generateNonce() (string, error) {
	nonceBytes := make([]byte, 32)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return "", fmt.Errorf("could not generate nonce")
	}

	return base64.URLEncoding.EncodeToString(nonceBytes), nil
}

var state = "ApplicationState"
var nonce = "NonceNotSetYet"

func loginHandler(res http.ResponseWriter, req *http.Request) {
	// nonce, _ = generateNonce()
	// //var redirectPath string

	// q := req.URL.Query()
	// q.Add("client_id", TOKEN)
	// q.Add("response_type", "code")
	// q.Add("response_mode", "query")
	// q.Add("scope", "openid profile email")
	// q.Add("redirect_uri", "http://localhost:3000/authorization-code/callback")
	// q.Add("state", state)
	// q.Add("nonce", nonce)

	//redirectPath = svcFactory.Config.OktaSecrets.IssuerURL + "/v1/authorize?" + q.Encode()

	//svcFactory.Logger.Bg().Fatal("HELLLO", zap.String("redir_url", redirectPath))
	// svcFactory.Logger.For(req.Context()).Fatal("HELLLO", zap.String("redir_url", redirectPath))

	http.Redirect(res, req, "https://google.com", http.StatusMovedPermanently)

}
