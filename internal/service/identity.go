package service

import (
	"context"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

func (p *Provider) AuthCodeURL(state string) string {
	return p.oauth.config.AuthCodeURL(state)
}

func (p *Provider) Verify(ctx context.Context, input string) (*oidc.IDToken, error) {
	return p.oauth.verifier.Verify(ctx, input)
}

func (p *Provider) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return p.oauth.config.Exchange(ctx, code, opts...)
}