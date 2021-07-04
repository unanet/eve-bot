package eveapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/interfaces"

	"github.com/dghubble/sling"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	eveerror "gitlab.unanet.io/devops/go/pkg/errors"
	evehttp "gitlab.unanet.io/devops/go/pkg/http"
	evejson "gitlab.unanet.io/devops/go/pkg/json"
	"gitlab.unanet.io/devops/go/pkg/log"
	"go.uber.org/zap"
)

// Config data structure for the Eve API
// EVEBOT_EVEAPI_BASE_URL
// EVEBOT_EVEAPI_TIMEOUT
// EVEBOT_EVEAPI_CALLBACK_URL
// EVEBOT_EVEAPI_ADMIN_TOKEN
type Config struct {
	EveapiBaseURL     string        `split_words:"true" required:"true"`
	EveapiTimeout     time.Duration `split_words:"true" default:"20s"`
	EveapiCallbackURL string        `split_words:"true" required:"true"`
	EveapiAdminToken  string        `split_words:"true" required:"true"`
}

// Client data structure
type Client struct {
	cfg   *Config
	sling *sling.Sling
}

// New creates a new eve api Client
func New(cfg Config) interfaces.EveAPI {
	var httpClient = &http.Client{
		Timeout:   cfg.EveapiTimeout,
		Transport: evehttp.LoggingTransport,
	}

	if !strings.HasSuffix(cfg.EveapiBaseURL, "/") {
		cfg.EveapiBaseURL += "/"
	}

	return &Client{
		cfg: &cfg,
		sling: sling.New().
			Base(cfg.EveapiBaseURL).
			Client(httpClient).
			Add("User-Agent", "eve-bot").
			Add("Authorization", fmt.Sprintf("Bearer %s", cfg.EveapiAdminToken)).
			ResponseDecoder(evejson.NewJsonDecoder()),
	}
}

func (c *Client) GetNamespaceJobs(ctx context.Context, ns *eve.Namespace) ([]eve.Job, error) {
	var success []eve.Job
	var failure eveerror.RestError

	r, err := c.sling.New().Get(fmt.Sprintf("namespaces/%v/jobs", ns.ID)).Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api GetNamespaceJobs request", zap.Error(err))
		return success, err
	}
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		return success, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return success, nil
	default:
		return success, fmt.Errorf(failure.Message)
	}
}

// DeleteMetadataKey calls the API to delete the metadata KEY (leaves empty {} is no metadata)
func (c *Client) DeleteMetadataKey(ctx context.Context, id int, key string) (eve.Metadata, error) {
	var success eve.Metadata
	var failure eveerror.RestError

	r, err := c.sling.New().Delete(fmt.Sprintf("metadata/%s/%s", strconv.Itoa(id), key)).Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api DeleteMetadataKey request", zap.Error(err))
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
		return success, fmt.Errorf(failure.Message)
	}
}

// UpsertMetadataServiceMap calls the API to upsert (insert/update) the metadata service map record
func (c *Client) UpsertMetadataServiceMap(ctx context.Context, payload eve.MetadataServiceMap) (eve.MetadataServiceMap, error) {
	var success eve.MetadataServiceMap
	var failure eveerror.RestError

	r, err := c.sling.New().Put(fmt.Sprintf("metadata/%s/service-maps", strconv.Itoa(payload.MetadataID))).BodyJSON(payload).Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api UpsertMetadataServiceMap request", zap.Error(err))
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
		return success, fmt.Errorf(failure.Message)
	}
}

// UpsertMergeMetadata calls the API to upsert (insert/update) the metadata record
func (c *Client) UpsertMergeMetadata(ctx context.Context, payload eve.Metadata) (eve.Metadata, error) {
	var success eve.Metadata
	var failure eveerror.RestError

	r, err := c.sling.New().Patch(fmt.Sprintf("metadata")).BodyJSON(payload).Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api UpsertMergeMetadata request", zap.Error(err))
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
		return success, fmt.Errorf(failure.Message)
	}
}

// GetMetadata calls the API to retrieve metadata by key
func (c *Client) GetMetadata(ctx context.Context, key string) (eve.Metadata, error) {
	var success eve.Metadata
	var failure eveerror.RestError

	r, err := c.sling.New().Get(fmt.Sprintf("metadata/%s", key)).Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api GetMetadata request", zap.Error(err))
		return success, err
	}

	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		return success, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return success, nil
	default:
		return success, failure
	}
}

// Release method calls the API to move artifacts in feeds
func (c *Client) Release(ctx context.Context, payload eve.Release) (eve.Release, error) {
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
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusPartialContent:
		return success, nil
	default:
		log.Logger.Warn("failed to release artifact", zap.Error(err), zap.Any("failure", failure))
		return success, fmt.Errorf(failure.Message)
	}
}

// GetServiceByName returns a service by name and namespace name
func (c *Client) GetServiceByName(ctx context.Context, namespace, service string) (eve.Service, error) {
	var success eve.Service
	var failure eveerror.RestError

	r, err := c.sling.New().Get(fmt.Sprintf("namespaces/%s/services/%s", namespace, service)).Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api GetServiceByName request", zap.Error(err))
		return success, err
	}
	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		return success, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return success, nil
	default:
		return success, fmt.Errorf(failure.Message)
	}
}

// SetNamespaceVersion sets the version on the namespace
func (c *Client) SetNamespaceVersion(ctx context.Context, version string, id int) (eve.Namespace, error) {
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
		return success, fmt.Errorf(failure.Message)
	}
}

// SetServiceVersion sets the version on the service
func (c *Client) SetServiceVersion(ctx context.Context, version string, id int) (eve.Service, error) {
	var success eve.Service
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
		return success, fmt.Errorf(failure.Message)
	}
}

// DeleteServiceMetadata deletes a metadata key on a service
func (c *Client) DeleteServiceMetadata(ctx context.Context, m string, id int) (params.MetadataMap, error) {
	var success params.MetadataMap
	var failure eveerror.RestError

	// Guard against the user sending key=value
	// we only want to send the key to the API
	metadatakey := m
	if strings.Contains(m, "=") {
		metadatakey = strings.Split(m, "=")[0]
	}

	if strings.Contains(metadatakey, "/") {
		return nil, fmt.Errorf("invalid metadata key: %s", metadatakey)
	}

	r, err := c.sling.New().Delete(fmt.Sprintf("services/%v/metadata/%s", id, metadatakey)).Request()
	if err != nil {
		log.Logger.Error("error preparing eve-api DeleteServiceMetadata request", zap.Error(err))
		return nil, err
	}

	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		log.Logger.Error("error calling eve-api DeleteServiceMetadata", zap.Error(err))
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusPartialContent:
		return success, nil
	default:
		return nil, fmt.Errorf(failure.Message)
	}
}

// GetServiceByID returns a service by an ID
func (c *Client) GetServiceByID(ctx context.Context, id int) (eve.Service, error) {
	var success eve.Service
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
		return success, fmt.Errorf(failure.Message)
	}
}

// GetServicesByNamespace returns all of the services for a given namespace
func (c *Client) GetServicesByNamespace(ctx context.Context, namespace string) ([]eve.Service, error) {
	var success []eve.Service
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
		return nil, fmt.Errorf(failure.Message)
	}
}

// GetEnvironmentByID returns an environment by ID
func (c *Client) GetEnvironmentByID(ctx context.Context, id string) (*eve.Environment, error) {
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
		return nil, fmt.Errorf(failure.Message)
	}
}

// GetEnvironments returns all of the environments
func (c *Client) GetEnvironments(ctx context.Context) ([]eve.Environment, error) {
	var success []eve.Environment
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
		return nil, fmt.Errorf(failure.Message)
	}
}

// GetNamespacesByEnvironment returns all of the namespaces for an environment
func (c *Client) GetNamespacesByEnvironment(ctx context.Context, environmentName string) ([]eve.Namespace, error) {
	var success []eve.Namespace
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
		return nil, fmt.Errorf(failure.Message)
	}
}

// Deploy calls the eve api to deploy resources
func (c *Client) Deploy(ctx context.Context, dp eve.DeploymentPlanOptions, user, channel, ts string) (*eve.DeploymentPlanOptions, error) {
	var success eve.DeploymentPlanOptions
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

	resp, err := c.sling.Do(r.WithContext(ctx), &success, &failure)
	if err != nil {
		log.Logger.Error("error calling eve-api Deploy", zap.Error(err))
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusPartialContent:
		return &success, nil
	default:
		return nil, fmt.Errorf(failure.Message)
	}
}

// GetNamespaceByID returns the namespace by an ID
func (c *Client) GetNamespaceByID(ctx context.Context, id int) (eve.Namespace, error) {
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
		return success, fmt.Errorf(failure.Message)
	}
}
