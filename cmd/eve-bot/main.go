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

// Public/Global Variables Passed in dynamically during Build time
// used to add build metadata into the binary
var (
	// GitCommit is the Full Git Commit SHA
	GitCommit string
	// GitCommitAuthor is the author of the Git Commit
	GitCommitAuthor string
	// GitBranch is the Full Git Branch Name
	GitBranch string
	// BuildDate is the DateTimeStamp during build
	BuildDate string
	// GitDescribe is a way to intentionally describe the version
	GitDescribe string
	// Version is the Full Semantic Version
	Version string
	// VersionPrerelease is the pre-release name (dev,rc-1,alpha,beta,nightly,etc.)
	VersionPrerelease string
	// VersionMetaData is the optional metadata to attach to a version
	VersionMetaData string
	// Builder is the name of the user that builds the artifact (i.e whoami)
	Builder string
	// BuildHost is the name of the host that builds the artifact
	BuildHost string
)

func main() {
	api, err := mux.NewApi(api.Controllers)
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
