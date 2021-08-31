package service

func (p *Provider) AuthCodeURL(userFQDN string) string {
	return p.oidc.AuthCodeURL(userFQDN)
}
