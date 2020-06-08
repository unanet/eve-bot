package slackservice

import (
	"github.com/slack-go/slack"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice/chatmodels"
)

func mapSlackUser(slackUser *slack.User) *chatmodels.ChatUser {
	return &chatmodels.ChatUser{Name: slackUser.Name}
}
