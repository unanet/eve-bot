package service

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"net/http/httptest"

	"github.com/slack-go/slack/slackevents"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/executor"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resolver"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice/chatmodels"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

type MockChatService struct {
}

func (mcs MockChatService) GetChannelInfo(ctx context.Context, channelID string) (chatmodels.Channel, error) {
	return chatmodels.Channel{ID: channelID, Name: "somethingcool"}, nil
}

func (mcs MockChatService) PostMessage(ctx context.Context, msg, channel string) (timestamp string) {
	return "2372372323"
}

func (mcs MockChatService) PostMessageThread(ctx context.Context, msg, channel, ts string) (timestamp string) {
	return "2372372323"
}

func (mcs MockChatService) ErrorNotification(ctx context.Context, user, channel string, err error) {

}

func (mcs MockChatService) ErrorNotificationThread(ctx context.Context, user, channel, ts string, err error) {

}

func (mcs MockChatService) UserNotificationThread(ctx context.Context, msg, user, channel, ts string) {

}

func (mcs MockChatService) DeploymentNotificationThread(ctx context.Context, msg, user, channel, ts string) {

}

func (mcs MockChatService) GetUser(ctx context.Context, user string) (*chatmodels.ChatUser, error) {
	return &chatmodels.ChatUser{Name: user}, nil
}

func (mcs MockChatService) PostLinkMessageThread(ctx context.Context, msg string, user string, channel string, ts string) {

}

func (mcs MockChatService) ShowResultsMessageThread(ctx context.Context, msg, user, channel, ts string) {

}

func (mcs MockChatService) ReleaseResultsMessageThread(ctx context.Context, msg, user, channel, ts string) {

}

type MockResolver struct {
}

func (mr MockResolver) Resolve(input, channel, user string) commands.EvebotCommand {
	return commands.NewInvalidCommand(strings.Fields(input), channel, user)
}

type MockExecutor struct {
}

func (me MockExecutor) Execute(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {

}

type MockEveAPIClient struct {
}

func (meac MockEveAPIClient) Deploy(ctx context.Context, dp eveapimodels.DeploymentPlanOptions, slackUser, slackChannel, ts string) (*eveapimodels.DeploymentPlanOptions, error) {
	return &eveapimodels.DeploymentPlanOptions{}, nil
}

func (meac MockEveAPIClient) GetEnvironmentByID(ctx context.Context, id string) (*eve.Environment, error) {
	return &eve.Environment{}, nil
}

func (meac MockEveAPIClient) GetEnvironments(ctx context.Context) (eveapimodels.Environments, error) {
	return eveapimodels.Environments{}, nil
}

func (meac MockEveAPIClient) GetNamespacesByEnvironment(ctx context.Context, environmentName string) (eveapimodels.Namespaces, error) {
	return eveapimodels.Namespaces{}, nil
}

func (meac MockEveAPIClient) GetServicesByNamespace(ctx context.Context, namespace string) (eveapimodels.Services, error) {
	return eveapimodels.Services{}, nil
}

func (meac MockEveAPIClient) GetServiceByName(ctx context.Context, namespace, service string) (eve.Service, error) {
	return eve.Service{}, nil
}

func (meac MockEveAPIClient) GetServiceByID(ctx context.Context, id int) (eveapimodels.EveService, error) {
	return eveapimodels.EveService{}, nil
}

func (meac MockEveAPIClient) SetServiceMetadata(ctx context.Context, metadata params.MetadataMap, id int) (params.MetadataMap, error) {
	return params.MetadataMap{}, nil
}

func (meac MockEveAPIClient) DeleteServiceMetadata(ctx context.Context, m string, id int) (params.MetadataMap, error) {
	return params.MetadataMap{}, nil
}

func (meac MockEveAPIClient) SetServiceVersion(ctx context.Context, version string, id int) (eveapimodels.EveService, error) {
	return eveapimodels.EveService{}, nil
}

func (meac MockEveAPIClient) SetNamespaceVersion(ctx context.Context, version string, id int) (eve.Namespace, error) {
	return eve.Namespace{}, nil
}

func (meac MockEveAPIClient) GetNamespaceByID(ctx context.Context, id int) (eve.Namespace, error) {
	return eve.Namespace{}, nil
}

func (meac MockEveAPIClient) Release(ctx context.Context, payload eve.Release) (eve.Release, error) {
	return eve.Release{}, nil
}

func (meac MockEveAPIClient) GetMetadata(ctx context.Context, key string) (eve.Metadata, error) {
	return eve.Metadata{}, nil
}

func (meac MockEveAPIClient) UpsertMergeMetadata(context.Context, eve.Metadata) (eve.Metadata, error) {
	return eve.Metadata{}, nil
}

func (meac MockEveAPIClient) UpsertMetadataServiceMap(context.Context, eve.MetadataServiceMap) (eve.MetadataServiceMap, error) {
	return eve.MetadataServiceMap{}, nil
}

func (meac MockEveAPIClient) DeleteMetadataKey(ctx context.Context, id int, key string) (eve.Metadata, error) {
	return eve.Metadata{}, nil
}

type MockConfig struct {
}

var (
	mockChatService       = MockChatService{}
	mockResolver          = MockResolver{}
	mockExecutor          = MockExecutor{}
	mockEveAPI            = MockEveAPIClient{}
	cfg                   = &config.Config{}
	mockallowedChannelMap = map[string]interface{}{}
	mockHTTPRequest       = httptest.NewRequest("POST", "/somewhere", nil)
)

func TestProvider_HandleSlackInteraction(t *testing.T) {

	//mockSlackInteractionBody := `{"type": "data","token": "data","callback_id": "data","response_url": "data","trigger_id": "data","action_ts": "data","team": "data","channel": "data","user": "data","original_message": "data","message": "data","name": "data","value": "data","message_ts": "data","attachment_id": "data","actions": "data","view": "data","action_id": "data","api_app_id": "data","block_id": "data","container": "data"}`

	// req := httptest.NewRequest("POST", "/somewhere?payload=", strings.NewReader(mockSlackInteractionBody))

	// req.Form = url.Values{}
	// req.PostForm = url.Values{}

	// req.PostForm.Set("payload", mockSlackInteractionBody)
	// req.Form.Set("payload", mockSlackInteractionBody)

	type fields struct {
		ChatService       chatservice.Provider
		CommandResolver   resolver.Resolver
		CommandExecutor   executor.Executor
		EveAPI            eveapi.Client
		Cfg               *config.Config
		allowedChannelMap map[string]interface{}
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "sad path - invalid post body",
			fields: fields{
				ChatService:       mockChatService,
				CommandResolver:   mockResolver,
				CommandExecutor:   mockExecutor,
				EveAPI:            mockEveAPI,
				Cfg:               cfg,
				allowedChannelMap: mockallowedChannelMap,
			},
			args: args{
				req: mockHTTPRequest,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{
				ChatService:       tt.fields.ChatService,
				CommandResolver:   tt.fields.CommandResolver,
				CommandExecutor:   tt.fields.CommandExecutor,
				EveAPI:            tt.fields.EveAPI,
				Cfg:               tt.fields.Cfg,
				allowedChannelMap: tt.fields.allowedChannelMap,
			}
			if err := p.HandleSlackInteraction(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Provider.HandleSlackInteraction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProvider_HandleSlackAppMentionEvent(t *testing.T) {
	type fields struct {
		ChatService       chatservice.Provider
		CommandResolver   resolver.Resolver
		CommandExecutor   executor.Executor
		EveAPI            eveapi.Client
		Cfg               *config.Config
		allowedChannelMap map[string]interface{}
	}
	type args struct {
		ctx context.Context
		ev  *slackevents.AppMentionEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{
				ChatService:       tt.fields.ChatService,
				CommandResolver:   tt.fields.CommandResolver,
				CommandExecutor:   tt.fields.CommandExecutor,
				EveAPI:            tt.fields.EveAPI,
				Cfg:               tt.fields.Cfg,
				allowedChannelMap: tt.fields.allowedChannelMap,
			}
			p.HandleSlackAppMentionEvent(tt.args.ctx, tt.args.ev)
		})
	}
}
