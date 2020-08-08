package eveapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"

	"go.uber.org/zap"

	"github.com/dghubble/sling"
	eveerror "gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	evehttp "gitlab.unanet.io/devops/eve/pkg/http"
	evejson "gitlab.unanet.io/devops/eve/pkg/json"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

// Config data structure for the Eve API
// EVEBOT_EVEAPI_BASE_URL
// EVEBOT_EVEAPI_TIMEOUT
// EVEBOT_EVEAPI_CALLBACK_URL
type Config struct {
	EveapiBaseURL     string        `split_words:"true" required:"true"`
	EveapiTimeout     time.Duration `split_words:"true" default:"20s"`
	EveapiCallbackURL string        `split_words:"true" required:"true"`
}

// Client interface for Eve API
// TODO: clean up this interface with more generic calls and interfaces
type Client interface {
	Deploy(ctx context.Context, dp eveapimodels.DeploymentPlanOptions, slackUser, slackChannel, ts string) (*eveapimodels.DeploymentPlanOptions, error)
	GetEnvironmentByID(ctx context.Context, id string) (*eve.Environment, error)
	GetEnvironments(ctx context.Context) (eveapimodels.Environments, error)
	GetNamespacesByEnvironment(ctx context.Context, environmentName string) (eveapimodels.Namespaces, error)
	GetServicesByNamespace(ctx context.Context, namespace string) (eveapimodels.Services, error)
	GetServiceByID(ctx context.Context, id int) (eveapimodels.EveService, error)
	SetServiceMetadata(ctx context.Context, metadata params.MetadataMap, id int) (params.MetadataMap, error)
	DeleteServiceMetadata(ctx context.Context, m string, id int) (params.MetadataMap, error)
	SetServiceVersion(ctx context.Context, version string, id int) (eveapimodels.EveService, error)
	SetNamespaceVersion(ctx context.Context, version string, id int) (eve.Namespace, error)
	GetNamespaceByID(ctx context.Context, id int) (eve.Namespace, error)
	Release(ctx context.Context, payload eve.Release) (eve.Release, error)
}

// client data structure
type client struct {
	cfg   *Config
	sling *sling.Sling
}

// NewClient creates a new eve api client
func NewClient(cfg Config) Client {
	var httpClient = &http.Client{
		Timeout:   cfg.EveapiTimeout,
		Transport: evehttp.LoggingTransport,
	}

	if !strings.HasSuffix(cfg.EveapiBaseURL, "/") {
		cfg.EveapiBaseURL += "/"
	}

	return &client{
		cfg: &cfg,
		sling: sling.New().
			Base(cfg.EveapiBaseURL).
			Client(httpClient).
			Add("User-Agent", "eve-bot").
			ResponseDecoder(evejson.NewJsonDecoder()),
	}
}

// Release method calls the API to move artifacts in feeds
func (c *client) Release(ctx context.Context, payload eve.Release) (eve.Release, error) {
	var success eve.Release
	var failure eveerror.RestError

	r, err := c.sling.New().Post(fmt.Sprintf("release")).BodyJSON(payload).Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api Release request", zap.Error(err))
		return success, err
	}

	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		return success, err
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		return success, nil
	default:
		log.Logger.Debug("an error occurred while trying to call eve-api SetNamespaceVersion", zap.String("error_msg", failure.Message))
		return success, fmt.Errorf(failure.Message)
	}
}

// SetNamespaceVersion sets the version on the namespace
func (c *client) SetNamespaceVersion(ctx context.Context, version string, id int) (eve.Namespace, error) {
	var success eve.Namespace
	var failure eveerror.RestError

	fullNS, err := c.GetNamespaceByID(ctx, id)
	if err != nil {
		return success, err
	}

	// Update the Version
	fullNS.RequestedVersion = version

	r, err := c.sling.New().Post(fmt.Sprintf("namespaces/%v", fullNS.ID)).BodyJSON(fullNS).Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api SetNamespaceVersion request", zap.Error(err))
		return success, err
	}

	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		return success, err
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusPartialContent:
		return success, nil
	default:
		log.Logger.Debug("an error occurred while trying to call eve-api SetNamespaceVersion", zap.String("error_msg", failure.Message))
		return success, fmt.Errorf(failure.Message)
	}
}

// SetServiceVersion sets the version on the service
func (c *client) SetServiceVersion(ctx context.Context, version string, id int) (eveapimodels.EveService, error) {
	var success eveapimodels.EveService
	var failure eveerror.RestError

	fullSvc, err := c.GetServiceByID(ctx, id)
	if err != nil {
		return success, err
	}

	// Update the Version
	fullSvc.OverrideVersion = version

	r, err := c.sling.New().Post(fmt.Sprintf("services/%v", fullSvc.ID)).BodyJSON(fullSvc).Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api SetServiceVersion request", zap.Error(err))
		return success, err
	}

	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		return success, err
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusPartialContent:
		return success, nil
	default:
		log.Logger.Debug("an error occurred while trying to call eve-api DeleteServiceMetadata", zap.String("error_msg", failure.Message))
		return success, fmt.Errorf(failure.Message)
	}
}

// DeleteServiceMetadata deletes a metadata key on a service
func (c *client) DeleteServiceMetadata(ctx context.Context, m string, id int) (params.MetadataMap, error) {
	var success params.MetadataMap
	var failure eveerror.RestError

	r, err := c.sling.New().Delete(fmt.Sprintf("services/%v/metadata/%s", id, m)).Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api DeleteServiceMetadata request", zap.Error(err))
		return nil, err
	}

	log.Logger.Debug("eve-api DeleteServiceMetadata req", zap.Any("metadata_key", m), zap.Int("service", id))
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		log.Logger.Error("error calling eve-api DeleteServiceMetadata", zap.Error(err))
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusPartialContent:
		return success, nil
	default:
		log.Logger.Debug("an error occurred while trying to call eve-api DeleteServiceMetadata", zap.String("error_msg", failure.Message))
		return nil, fmt.Errorf(failure.Message)
	}
}

// GetServiceByID returns a service by an ID
func (c *client) GetServiceByID(ctx context.Context, id int) (eveapimodels.EveService, error) {
	var success eveapimodels.EveService
	var failure eveerror.RestError
	r, err := c.sling.New().Get(fmt.Sprintf("services/%v", id)).Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api GetServiceByID request", zap.Error(err))
		return success, eveerror.Wrap(err)
	}
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		log.Logger.Error("error calling eve-api GetServiceByID", zap.Error(err))
		return success, eveerror.Wrap(err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return success, nil
	default:
		log.Logger.Debug("an error occurred while trying to call eve-api GetServicesByNamespace", zap.String("error_msg", failure.Message))
		return success, fmt.Errorf(failure.Message)
	}
}

// GetServicesByNamespace returns all of the services for a given namespace
func (c *client) GetServicesByNamespace(ctx context.Context, namespace string) (eveapimodels.Services, error) {
	var success eveapimodels.Services
	var failure eveerror.RestError
	r, err := c.sling.New().Get(fmt.Sprintf("namespaces/%s/services", namespace)).Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api GetServicesByNamespace request", zap.Error(err))
		return nil, eveerror.Wrap(err)
	}
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		log.Logger.Error("error calling eve-api GetServicesByNamespace", zap.Error(err))
		return nil, eveerror.Wrap(err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return success, nil
	default:
		log.Logger.Debug("an error occurred while trying to call eve-api GetServicesByNamespace", zap.String("error_msg", failure.Message))
		return nil, fmt.Errorf(failure.Message)
	}
}

// GetEnvironmentByID returns an environment by ID
func (c *client) GetEnvironmentByID(ctx context.Context, id string) (*eve.Environment, error) {
	var success eve.Environment
	var failure eveerror.RestError
	r, err := c.sling.New().Get(fmt.Sprintf("environments/%s", id)).Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api GetEnvironment request", zap.Error(err))
		return nil, eveerror.Wrap(err)
	}
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		log.Logger.Error("error calling eve-api GetEnvironment", zap.Error(err))
		return nil, eveerror.Wrap(err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return &success, nil
	default:
		log.Logger.Debug("an error occurred while trying to call eve-api GetEnvironment", zap.String("error_msg", failure.Message))
		return nil, fmt.Errorf(failure.Message)
	}
}

// GetEnvironments returns all of the environments
func (c *client) GetEnvironments(ctx context.Context) (eveapimodels.Environments, error) {
	var success eveapimodels.Environments
	var failure eveerror.RestError
	r, err := c.sling.New().Get("environments").Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api GetEnvironments request", zap.Error(err))
		return nil, eveerror.Wrap(err)
	}
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		log.Logger.Error("error calling eve-api GetEnvironments", zap.Error(err))
		return nil, eveerror.Wrap(err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return success, nil
	default:
		log.Logger.Debug("an error occurred while trying to call eve-api GetEnvironments", zap.String("error_msg", failure.Message))
		return nil, fmt.Errorf(failure.Message)
	}
}

// GetNamespacesByEnvironment returns all of the namespaces for an environment
func (c *client) GetNamespacesByEnvironment(ctx context.Context, environmentName string) (eveapimodels.Namespaces, error) {
	var success eveapimodels.Namespaces
	var failure eveerror.RestError
	r, err := c.sling.New().Get("namespaces").Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api GetNamespacesByEnvironment request", zap.Error(err))
		return nil, eveerror.Wrap(err)
	}
	r.URL.RawQuery = fmt.Sprintf("environment=%s", environmentName)
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		log.Logger.Error("error calling eve-api GetNamespacesByEnvironment", zap.Error(err))
		return nil, eveerror.Wrap(err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return success, nil
	default:
		log.Logger.Debug("an error occurred while trying to call eve-api GetNamespacesByEnvironment", zap.String("error_msg", failure.Message))
		return nil, fmt.Errorf(failure.Message)
	}
}

// Deploy calls the eve api to deploy resources
func (c *client) Deploy(ctx context.Context, dp eveapimodels.DeploymentPlanOptions, user, channel, ts string) (*eveapimodels.DeploymentPlanOptions, error) {
	var success eveapimodels.DeploymentPlanOptions
	var failure eveerror.RestError

	cbURLVals := url.Values{}
	cbURLVals.Set("user", user)
	cbURLVals.Add("channel", channel)
	cbURLVals.Add("ts", ts)

	dp.CallbackURL = c.cfg.EveapiCallbackURL + "?" + cbURLVals.Encode()

	r, err := c.sling.New().Post("deployment-plans").BodyJSON(dp).Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api Deploy request", zap.Error(err))
		return nil, err
	}

	log.Logger.Debug("eve-api Deploy req", zap.Any("req", dp))
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		log.Logger.Error("error calling eve-api Deploy", zap.Error(err))
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

// SetServiceMetadata sets the metadata on the service
func (c *client) SetServiceMetadata(ctx context.Context, metadata params.MetadataMap, id int) (params.MetadataMap, error) {
	var success params.MetadataMap
	var failure eveerror.RestError

	r, err := c.sling.New().Patch(fmt.Sprintf("services/%v/metadata", id)).BodyJSON(metadata).Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api SetServiceMetadata request", zap.Error(err))
		return nil, err
	}

	log.Logger.Debug("eve-api SetServiceMetadata req", zap.Any("req", metadata))
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		log.Logger.Error("error calling eve-api SetServiceMetadata", zap.Error(err))
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusPartialContent:
		return success, nil
	default:
		log.Logger.Debug("an error occurred while trying to call eve-api SetServiceMetadata", zap.String("error_msg", failure.Message))
		return nil, fmt.Errorf(failure.Message)
	}
}

// GetNamespaceByID returns the namespace by an ID
func (c *client) GetNamespaceByID(ctx context.Context, id int) (eve.Namespace, error) {
	var success eve.Namespace
	var failure eveerror.RestError

	r, err := c.sling.New().Get(fmt.Sprintf("namespaces/%v", id)).Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api GetNamespaceByID request", zap.Error(err))
		return success, err
	}

	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		log.Logger.Error("error calling eve-api GetNamespaceByID", zap.Error(err))
		return success, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return success, nil
	default:
		log.Logger.Debug("an error occurred while trying to call eve-api GetNamespaceByID", zap.String("error_msg", failure.Message))
		return success, fmt.Errorf(failure.Message)
	}
}
