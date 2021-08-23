package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/unanet/go/pkg/errors"
	"github.com/unanet/go/pkg/log"
)

// HandleSlackInteraction handles the interactive callbacks (buttons, dropdowns, etc.)
func (p *Provider) HandleSlackInteraction(req *http.Request) error {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(req.FormValue("payload")), &payload)
	if err != nil {
		return errors.RestError{Code: http.StatusBadRequest, Message: "failed to parse interactive slack message payload", OriginalError: err}
	}
	log.Logger.Info(fmt.Sprintf("Message button pressed by user %s with value %s", payload.User.Name, payload.Value))
	return nil
}

// HandleSlackAppMentionEvent takes slackevents.AppMentionEvent, resolves an EvebotCommand, and handles/executes it...
func (p *Provider) HandleSlackAppMentionEvent(ctx context.Context, ev *slackevents.AppMentionEvent) {
	// Resolve the input and return an EvebotCommand object
	cmd := p.CommandResolver.Resolve(ev.Text, ev.Channel, ev.User)

	slackUser, err := p.ChatService.GetUser(ctx, cmd.Info().User)
	if err != nil {
		p.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, ev.ThreadTimeStamp, err)
		return
	}

	if cmd.IsAuthenticated(slackUser, p.UserDB) == false {
		p.ChatService.PostPrivateMessage(ctx, p.MgrSvc.AuthCodeURL(slackUser.FullyQualifiedName()), cmd.Info().User)
		_ = p.ChatService.PostMessageThread(ctx, "You need to login. Please Check your Private DM from `evebot` for an auth link", cmd.Info().Channel, ev.ThreadTimeStamp)
		return
	}

	if cmd.IsAuthorized(p.allowedChannelMap, p.ChatService.GetChannelInfo) == false {
		_ = p.ChatService.PostMessageThread(ctx, "You are not authorized to perform this action", cmd.Info().Channel, ev.ThreadTimeStamp)
		return
	}

	// SlackAuthEnabled is like a "feature flag"
	// set to true we will check auth
	// set to false we will skip the auth check
	//if p.Cfg.SlackAuthEnabled {
	//	if cmd.IsAuthorized(p.allowedChannelMap, p.ChatService.GetChannelInfo, p.ChatService.GetUser, p.UserDB) == false {
	//		_ = p.ChatService.PostMessageThread(ctx, "You are not authorized to perform this action", cmd.Info().Channel, ev.ThreadTimeStamp)
	//		return
	//	}
	//}

	// SlackMaintenanceEnabled is like a "feature flag"
	// set to true, and we are in Maintenance Mode
	// Only Channels set to the EVEBOT_SLACK_CHANNELS_MAINTENANCE environment variable are allowed to issue commands
	// ex:  EVEBOT_SLACK_CHANNELS_MAINTENANCE=my-evebot,evebot-tests
	if p.Cfg.SlackMaintenanceEnabled {
		incomingChannel, err := p.ChatService.GetChannelInfo(ctx, cmd.Info().Channel)
		if err != nil {
			// This shouldn't happen, but if it does, we don't want to be locked out from deploying eve
			// so we will show the error (which is logged) and DevOps will take care of the issue (hopefully...)
			p.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, ev.ThreadTimeStamp, err)
		} else {
			// Not coming from an approved Maintenance channel Show the maintenance mode
			if _, ok := p.allowedMaintenanceChannelMap[incomingChannel.Name]; ok == false {
				_ = p.ChatService.PostMessageThread(ctx, ":construction: Sorry, but we are currently in maintenance mode!", cmd.Info().Channel, ev.ThreadTimeStamp)
				return
			}
		}
	}

	// Hydrate the Acknowledgement Message and whether we should continue...
	ackMsg, cont := cmd.AckMsg()
	// Send the AckMsg and get the Timestamp back, so we can thread it later on...
	timeStamp := p.ChatService.PostMessageThread(ctx, ackMsg, cmd.Info().Channel, ev.ThreadTimeStamp)
	// If the AckMessage needs to continue (no errors)...
	if cont {
		// Asynchronous CommandExecutor call
		// which maps an EveBotCommand to a CommandHandler
		go p.CommandExecutor.Execute(context.TODO(), cmd, timeStamp)
	}
}
