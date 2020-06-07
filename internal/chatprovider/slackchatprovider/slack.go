package slackchatprovider

import (
	"context"

	"github.com/slack-go/slack"
)

type slackProvider struct {
	client *slack.Client
}

func New(c *slack.Client) slackProvider {
	return slackProvider{
		client: c,
	}
}

func (sp slackProvider) PostMessage(ctx context.Context) {
	//sp.client.PostMessageContext(ctx)
}
