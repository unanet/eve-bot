package evebotservice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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

// HandleEvent takes an http request and handles the Slack API Event
// this is where we do our request signature validation
// ..and switch the incoming api event types
func (p *Provider) HandleSlackAppMentionEvent(ctx context.Context, ev *slackevents.AppMentionEvent) error {
	// Resolve the input and return a Command object
	cmd := p.CommandResolver.Resolve(ev.Text, ev.Channel, ev.User)
	// Hydrate the Acknowledgement Message and whether or not we should continue...
	ackMsg, cont := cmd.AckMsg()
	// Send the AckMsg and get the Timestamp back so we can thread it later on...
	timeStamp := p.ChatService.PostMessageThread(ctx, ackMsg, cmd.Channel(), ev.ThreadTimeStamp)
	// If the AckMessage needs to continue (no errors)...
	if cont {
		log.Logger.Debug("lets continue to execute command handler....")
		go p.CommandHandler.Execute(context.TODO(), cmd, timeStamp) // Asynchronous Command Handler
		//go p.CommandHandler.Handle(context.TODO(), cmd, timeStamp) // Asynchronous Command Handler
	}
	// Let's get back to the party and take a few more request...
	return nil
}
