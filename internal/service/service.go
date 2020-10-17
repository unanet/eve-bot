package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"go.uber.org/zap"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

// HandleSlackInteraction handles the interactive callbacks (buttons, dropdowns, etc.)
func (p *Provider) HandleSlackInteraction(req *http.Request) error {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(req.FormValue("payload")), &payload)
	if err != nil {
		return &errors.RestError{Code: http.StatusBadRequest, Message: "failed to parse interactive slack message payload", OriginalError: err}
	}
	log.Logger.Info(fmt.Sprintf("Message button pressed by user %s with value %s", payload.User.Name, payload.Value))
	return nil
}

// HandleSlackAppMentionEvent takes slackevents.AppMentionEvent, resolves an EvebotCommand, and handles/executes it...
func (p *Provider) HandleSlackAppMentionEvent(ctx context.Context, ev *slackevents.AppMentionEvent) {
	// Resolve the input and return a Command object
	cmd := p.CommandResolver.Resolve(ev.Text, ev.Channel, ev.User)

	// SlackAuthEnabled is like a "feature flag"
	// set to true and we will check auth
	// set to false and we will skip the auth check
	if p.Cfg.SlackAuthEnabled {
		if cmd.IsAuthorized(p.allowedChannelMap, p.ChatService.GetChannelInfo) == false {
			_ = p.ChatService.PostMessageThread(ctx, "You are not authorized to perform this action", cmd.Info().Channel, ev.ThreadTimeStamp)
			return
		}
	}

	// Hydrate the Acknowledgement Message and whether or not we should continue...
	ackMsg, cont := cmd.AckMsg()
	// Send the AckMsgFn and get the Timestamp back so we can thread it later on...
	timeStamp := p.ChatService.PostMessageThread(ctx, ackMsg, cmd.Info().Channel, ev.ThreadTimeStamp)
	// If the AckMessage needs to continue (no errors)...
	if cont {
		log.Logger.Debug("execute command handler", zap.Any("cmd_type", reflect.TypeOf(cmd)))
		// Asynchronous Command Handler
		go p.CommandExecutor.Execute(context.TODO(), cmd, timeStamp)
	}
}