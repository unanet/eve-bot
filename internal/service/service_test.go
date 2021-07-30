package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/unanet/eve-bot/internal/botcommander/interfaces"

	"github.com/unanet/eve-bot/internal/chatservice/chatmodels"

	"github.com/unanet/eve-bot/internal/botcommander/commands"

	"github.com/golang/mock/gomock"

	"github.com/slack-go/slack/slackevents"
	"github.com/unanet/eve-bot/internal/botcommander/executor"
	"github.com/unanet/eve-bot/internal/botcommander/resolver"
	"github.com/unanet/eve-bot/internal/chatservice"
	"github.com/unanet/eve-bot/internal/config"
	"github.com/unanet/eve-bot/internal/eveapi"
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
		ChatService       interfaces.ChatProvider
		CommandResolver   interfaces.CommandResolver
		CommandExecutor   interfaces.CommandExecutor
		EveAPI            interfaces.EveAPI
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
	mockChat              *chatservice.MockProvider
	mockCmd               *commands.MockEvebotCommand
	mockResolver          *resolver.MockResolver
	mockExecutor          *executor.MockExecutor
	mockAPI               *eveapi.MockClient
	mockCfg               *config.Config
	mockAllowedChannelMap map[string]interface{}
}

func newServiceMocks(ctrl *gomock.Controller) *serviceMocks {
	return &serviceMocks{
		mockChat:              chatservice.NewMockProvider(ctrl),
		mockCmd:               commands.NewMockEvebotCommand(ctrl),
		mockResolver:          resolver.NewMockResolver(ctrl),
		mockExecutor:          executor.NewMockExecutor(ctrl),
		mockAPI:               eveapi.NewMockClient(ctrl),
		mockCfg:               &config.Config{},
		mockAllowedChannelMap: make(map[string]interface{}),
	}
}

var mockSlackEvent = &slackevents.AppMentionEvent{
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

var mockChatInfo = commands.ChatInfo{
	User:        mockSlackEvent.User,
	Channel:     mockSlackEvent.Channel,
	CommandName: strings.Split(mockSlackEvent.Text, " ")[0],
}

var mockChatChannel = chatmodels.Channel{
	ID:   "some id",
	Name: "some name",
}

func TestProvider_HandleSlackAppMentionEvent(t *testing.T) {

	type args struct {
		ctx context.Context
		ev  *slackevents.AppMentionEvent
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		args       args
		setupMocks func(*serviceMocks)
		ctrl       *gomock.Controller
	}{
		{
			name: "happy path",
			ctrl: ctrl,
			setupMocks: func(t *serviceMocks) {
				t.mockChat.EXPECT().PostMessageThread(context.Background(), "Ohhh yeah", mockSlackEvent.Channel, mockSlackEvent.ThreadTimeStamp).Return("2342342342")
				t.mockCmd.EXPECT().Info().Return(mockChatInfo)
				t.mockCmd.EXPECT().AckMsg().Return("Ohhh yeah", false)
				t.mockResolver.EXPECT().Resolve(mockSlackEvent.Text, mockSlackEvent.Channel, mockSlackEvent.User).Return(t.mockCmd)
			},
			args: args{
				ctx: context.Background(),
				ev:  mockSlackEvent,
			},
		},
		{
			name: "slack auth disabled",
			ctrl: ctrl,
			setupMocks: func(t *serviceMocks) {
				t.mockChat.EXPECT().PostMessageThread(context.Background(), "Ohhh yeah", mockSlackEvent.Channel, mockSlackEvent.ThreadTimeStamp).Return("2342342342")
				t.mockCmd.EXPECT().Info().Return(mockChatInfo)
				t.mockCmd.EXPECT().AckMsg().Return("Ohhh yeah", false)
				t.mockResolver.EXPECT().Resolve(mockSlackEvent.Text, mockSlackEvent.Channel, mockSlackEvent.User).Return(t.mockCmd)
			},
			args: args{
				ctx: context.Background(),
				ev:  mockSlackEvent,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newServiceMocks(tt.ctrl)
			//tt.setupMocks(m)
			if tt.setupMocks != nil {
				tt.setupMocks(m)
			}
			New(m.mockCfg, m.mockResolver, m.mockAPI, m.mockChat, m.mockExecutor).HandleSlackAppMentionEvent(tt.args.ctx, tt.args.ev)
		})
	}
}
