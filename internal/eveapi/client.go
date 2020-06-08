package eveapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/dghubble/sling"
	eveerror "gitlab.unanet.io/devops/eve/pkg/errors"
	evehttp "gitlab.unanet.io/devops/eve/pkg/http"
	evejson "gitlab.unanet.io/devops/eve/pkg/json"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

// EVEBOT_EVEAPI_BASE_URL
// EVEBOT_EVEAPI_TIMEOUT
// EVEBOT_EVEAPI_CALLBACK_URL
type Config struct {
	EveapiBaseUrl     string        `split_words:"true" required:"true"`
	EveapiTimeout     time.Duration `split_words:"true" default:"20s"`
	EveapiCallbackUrl string        `split_words:"true" required:"true"`
}

type Client interface {
	Deploy(ctx context.Context, dp DeploymentPlanOptions, slackUser, slackChannel, ts string) (*DeploymentPlanOptions, error)
}

type client struct {
	cfg   *Config
	sling *sling.Sling
}

func NewClient(cfg Config) Client {
	var httpClient = &http.Client{
		Timeout:   cfg.EveapiTimeout,
		Transport: evehttp.LoggingTransport,
	}

	if !strings.HasSuffix(cfg.EveapiBaseUrl, "/") {
		cfg.EveapiBaseUrl += "/"
	}

	return &client{
		cfg: &cfg,
		sling: sling.New().
			Base(cfg.EveapiBaseUrl).
			Client(httpClient).
			Add("User-Agent", "eve-bot").
			ResponseDecoder(evejson.NewJsonDecoder()),
	}

}

func (c *client) Deploy(ctx context.Context, dp DeploymentPlanOptions, user, channel, ts string) (*DeploymentPlanOptions, error) {
	var success DeploymentPlanOptions
	var failure eveerror.RestError

	cbUrlVals := url.Values{}
	cbUrlVals.Set("user", user)
	cbUrlVals.Add("channel", channel)
	cbUrlVals.Add("ts", ts)
	cbUrlVals.Add("action", dp.Type.String())

	dp.CallbackURL = c.cfg.EveapiCallbackUrl + "?" + cbUrlVals.Encode()

	r, err := c.sling.New().Post("deployment-plans").BodyJSON(dp).Request()
	if err != nil {
		log.Logger.Error("error calling eve-api", zap.Error(err))
		return nil, err
	}

	log.Logger.Debug("eve-api req", zap.Any("req", dp))
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		log.Logger.Error("error calling eve-api", zap.Error(err))
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusPartialContent:
		return &success, nil
	default:
		log.Logger.Debug("an error occurred while trying to call eve-api deploy", zap.String("error_msg", failure.Message))
		return nil, fmt.Errorf(failure.Message)
	}

}
