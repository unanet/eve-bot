package evebotservice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"go.uber.org/zap"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

// HandleInteraction handles the interactive callbacks (buttons, dropdowns, etc.)
func (p *Provider) HandleSlackInteraction(req *http.Request) error {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(req.FormValue("payload")), &payload)
	if err != nil {
		return botError(err, "failed to parse interactive slack message payload", http.StatusInternalServerError)
	}
	log.Logger.Info(fmt.Sprintf("Message button pressed by user %s with value %s", payload.User.Name, payload.Value))
	return nil
}

// HandleSlackAppMentionEvent takes slackevents.AppMentionEvent, resolves an EvebotCommand, and handles/executes it...
func (p *Provider) HandleSlackAppMentionEvent(ctx context.Context, ev *slackevents.AppMentionEvent) error {
	// Resolve the input and return a Command object
	cmd := p.CommandResolver.Resolve(ev.Text, ev.Channel, ev.User)
	// Hydrate the Acknowledgement Message and whether or not we should continue...
	ackMsg, cont := cmd.AckMsg()
	// Send the AckMsg and get the Timestamp back so we can thread it later on...
	timeStamp := p.ChatService.PostMessageThread(ctx, ackMsg, cmd.Channel(), ev.ThreadTimeStamp)
	// If the AckMessage needs to continue (no errors)...
	if cont {
		log.Logger.Debug("execute command handler", zap.Any("cmd_type", reflect.TypeOf(cmd)))
		go p.CommandExecutor.Execute(context.TODO(), cmd, timeStamp) // Asynchronous Command Handler
	}
	// Let's get back to the party and take a few more request...
	return nil
}
