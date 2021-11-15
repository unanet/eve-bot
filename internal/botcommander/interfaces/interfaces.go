package interfaces

import (
	"context"

	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/eve-bot/internal/chatservice/chatmodels"
	"github.com/unanet/eve/pkg/eve"
)

// ChatProvider interface used to interface with a Chat Provider (i.e. Slack)
// TODO: clean up this interface with msg type switch (remove duplicate signatures)
type ChatProvider interface {
	GetChannelInfo(ctx context.Context, channelID string) (chatmodels.Channel, error)
	PostMessage(ctx context.Context, msg, channel string) (timestamp string)
	PostMessageThread(ctx context.Context, msg, channel, ts string) (timestamp string)
	ErrorNotification(ctx context.Context, user, channel string, err error)
	ErrorNotificationThread(ctx context.Context, user, channel, ts string, err error)
	UserNotificationThread(ctx context.Context, msg, user, channel, ts string)
	DeploymentNotificationThread(ctx context.Context, msg, user, channel, ts string)
	GetUser(ctx context.Context, user string) (*chatmodels.ChatUser, error)
	PostLinkMessageThread(ctx context.Context, msg string, user string, channel string, ts string)
	ShowResultsMessageThread(ctx context.Context, msg, user, channel, ts string)
	ReleaseResultsMessageThread(ctx context.Context, msg, user, channel, ts string)
	PostPrivateMessage(ctx context.Context, msg string, user string)
}

// EveAPI interface used to interface with eve/pipeline API
// TODO: clean up this interface with more generic calls (GET,PUT,POST,DELETE,PATCH with interfaces{})
type EveAPI interface {
	Deploy(ctx context.Context, dp eve.DeploymentPlanOptions, slackUser, slackChannel, ts string) (*eve.DeploymentPlanOptions, error)
	GetEnvironmentByID(ctx context.Context, id string) (*eve.Environment, error)
	GetEnvironments(ctx context.Context) ([]eve.Environment, error)
	GetNamespacesByEnvironment(ctx context.Context, environmentName string) ([]eve.Namespace, error)
	GetServicesByNamespace(ctx context.Context, namespace string) ([]eve.Service, error)
	GetServiceByName(ctx context.Context, namespace, service string) (eve.Service, error)
	GetServiceByID(ctx context.Context, id int) (eve.Service, error)
	DeleteServiceMetadata(ctx context.Context, m string, id int) (params.MetadataMap, error)
	SetServiceVersion(ctx context.Context, version string, id int) (eve.Service, error)
	SetNamespaceVersion(ctx context.Context, version string, id int) (eve.Namespace, error)
	GetNamespaceByID(ctx context.Context, id int) (eve.Namespace, error)
	Release(ctx context.Context, payload eve.Release) ([]eve.Release, error)
	GetMetadata(ctx context.Context, key string) (eve.Metadata, error)
	UpsertMergeMetadata(context.Context, eve.Metadata) (eve.Metadata, error)
	UpsertMetadataServiceMap(context.Context, eve.MetadataServiceMap) (eve.MetadataServiceMap, error)
	DeleteMetadataKey(ctx context.Context, id int, key string) (eve.Metadata, error)
	GetNamespaceJobs(ctx context.Context, ns *eve.Namespace) ([]eve.Job, error)
}

// CommandResolver resolves the input and returns an EvebotCommand (Invalid command instead of an error for error cases)
type CommandResolver interface {
	Resolve(input, channel, user string) commands.EvebotCommand
}

// CommandExecutor interface takes an EvebotCommand and Executes a matching handler
type CommandExecutor interface {
	Execute(ctx context.Context, cmd commands.EvebotCommand, timestamp string)
}
