package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitlab.unanet.io/devops/eve-bot/internal/chatservice/chatmodels"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"

	"github.com/golang/mock/gomock"

	"github.com/slack-go/slack/slackevents"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/executor"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resolver"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
)

func TestProvider_HandleSlackInteraction(t *testing.T) {

	//mockSlackInteractionBody := `{"type": "data","token": "data","callback_id": "data","response_url": "data","trigger_id": "data","action_ts": "data","team": "data","channel": "data","user": "data","original_message": "data","message": "data","name": "data","value": "data","message_ts": "data","attachment_id": "data","actions": "data","view": "data","action_id": "data","api_app_id": "data","block_id": "data","container": "data"}`

	// req := httptest.NewRequest("POST", "/somewhere?payload=", strings.NewReader(mockSlackInteractionBody))

	// req.Form = url.Values{}
	// req.PostForm = url.Values{}

	// req.PostForm.Set("payload", mockSlackInteractionBody)
	// req.Form.Set("payload", mockSlackInteractionBody)

	mockChatSvc := chatservice.NewMockProvider(gomock.NewController(t))
	mockResolver := resolver.NewMockResolver(gomock.NewController(t))
	mockExecutor := executor.NewMockExecutor(gomock.NewController(t))
	mockAPI := eveapi.NewMockClient(gomock.NewController(t))

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
				ChatService:       mockChatSvc,
				CommandResolver:   mockResolver,
				CommandExecutor:   mockExecutor,
				EveAPI:            mockAPI,
				Cfg:               &config.Config{},
				allowedChannelMap: map[string]interface{}{},
			},
			args: args{
				req: httptest.NewRequest("POST", "/somewhere", nil),
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

type serviceMocks struct {
	mockChat              chatservice.MockProvider
	mockResolver          resolver.MockResolver
	mockExecutor          executor.MockExecutor
	mockAPI               eveapi.MockClient
	mockCfg               *config.Config
	mockAllowedChannelMap map[string]interface{}
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
		req *http.Request
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSlackEvent := &slackevents.AppMentionEvent{
		Type:            "test",
		User:            "someuser",
		Text:            "show environment",
		TimeStamp:       "2423423423",
		ThreadTimeStamp: "2342342343",
		Channel:         "test",
		EventTimeStamp:  "2342342343",
		UserTeam:        "unknown",
		SourceTeam:      "unknown",
		BotID:           "test",
	}

	mockChatInfo := commands.ChatInfo{
		User:        mockSlackEvent.User,
		Channel:     mockSlackEvent.Channel,
		CommandName: strings.Split(mockSlackEvent.Text, " ")[0],
	}

	mockChatOpts := commands.CommandOptions{}

	mockChatChannel := chatmodels.Channel{
		ID:   "test",
		Name: "WTF",
	}

	mockAllowedChannels := make(map[string]interface{})

	mockChatSvc := chatservice.NewMockProvider(ctrl)
	mockChatSvc.EXPECT().GetChannelInfo(context.Background(), mockSlackEvent.Channel).Return(mockChatChannel, nil)

	mockEveCmd := commands.NewMockEvebotCommand(ctrl)
	mockEveCmd.EXPECT().Info().Return(mockChatInfo)
	mockEveCmd.EXPECT().Options().Return(mockChatOpts)
	mockEveCmd.EXPECT().AckMsg().Return("Ohhh yeah", false)
	mockEveCmd.EXPECT().IsAuthorized(mockAllowedChannels, mockChatSvc.GetChannelInfo).Return(true)

	mockResolver := resolver.NewMockResolver(ctrl)
	mockResolver.
		EXPECT().
		Resolve(mockSlackEvent.Text, mockSlackEvent.Channel, mockSlackEvent.User).
		Return(mockEveCmd)

	mockExecutor := executor.NewMockExecutor(ctrl)
	mockAPI := eveapi.NewMockClient(ctrl)

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		setupMocks func(*serviceMocks)
	}{
		{
			name: "sad path - invalid post body",
			setupMocks: func(t *serviceMocks) {

			},
			fields: fields{
				ChatService:       mockChatSvc,
				CommandResolver:   mockResolver,
				CommandExecutor:   mockExecutor,
				EveAPI:            mockAPI,
				Cfg:               &config.Config{},
				allowedChannelMap: map[string]interface{}{},
			},
			args: args{
				ctx: context.Background(),
				ev:  mockSlackEvent,
				req: httptest.NewRequest("POST", "/somewhere", nil),
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
			if err := p.HandleSlackAppMentionEvent(tt.args.ctx, tt.args.ev); (err != nil) != tt.wantErr {
				t.Errorf("Provider.HandleSlackInteraction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
