package eveapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gitlab.unanet.io/devops/eve/pkg/json"

	"github.com/dghubble/sling"

	ehttp "gitlab.unanet.io/devops/eve/pkg/http"
)

const (
	userAgent = "eve-bot"
)

// EVEBOT_EVEAPI_BASE_URL
// EVEBOT_EVEAPI_TIMEOUT
type Config struct {
	EveapiBaseUrl string        `split_words:"true" required:"true"`
	EveapiTimeout time.Duration `split_words:"true" default:"20s"`
}

type Client interface {
	Deploy(ctx context.Context, req DeploymentPlanOptions) (*DeployResponse, error)
	Migrate() (*MigrateResponse, error)
}

type client struct {
	sling *sling.Sling
}

type CallbackState struct {
	User    string
	Channel string
}

type EveParams struct {
	State CallbackState `url:"state,omitempty"`
}

func NewClient(cfg Config) Client {
	var httpClient = &http.Client{
		Timeout:   cfg.EveapiTimeout,
		Transport: ehttp.LoggingTransport,
	}

	if !strings.HasSuffix(cfg.EveapiBaseUrl, "/") {
		cfg.EveapiBaseUrl += "/"
	}

	return &client{
		sling: sling.New().
			Base(cfg.EveapiBaseUrl).
			Client(httpClient).
			Add("User-Agent", userAgent).
			ResponseDecoder(json.NewJsonDecoder()),
	}

}

func (c *client) Deploy(ctx context.Context, dp DeploymentPlanOptions) (*DeployResponse, error) {
	var success DeployResponse
	var failure string

	params := &EveParams{State: CallbackState{User: "something", Channel: "something else"}}

	r, err := c.sling.New().Post("deployment-plans").BodyJSON(dp).QueryStruct(params).Request()
	if err != nil {
		return nil, err
	}

	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return &success, nil
	default:
		return nil, fmt.Errorf("an error occurred while trying to call eve-api deploy: %s", failure)
	}

}

func (c *client) Migrate() (*MigrateResponse, error) {
	return nil, nil
}
