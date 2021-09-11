package service

import (
	"context"

	"github.com/coreos/go-oidc"
)

func (p *Provider) AuthCodeURL(state string) string {
	return p.oidc.AuthCodeURL(state)
}

func (p *Provider) Verify(ctx context.Context, input string) (*oidc.IDToken, error) {
	return p.oidc.Verify(ctx, input)
}
