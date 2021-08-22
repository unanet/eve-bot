package slackservice

import (
	"github.com/slack-go/slack"
	"github.com/unanet/eve-bot/internal/chatservice/chatmodels"
)

func mapSlackUser(slackUser *slack.User) *chatmodels.ChatUser {
	return &chatmodels.ChatUser{
		Name: slackUser.Name,
		ID: slackUser.ID,
	}
}
