package okta

// Config contains the okta configuration
type Config struct {
	OktaClientID     string `split_words:"true" required:"false"`
	OktaClientSecret string `split_words:"true" required:"false"`
	OktaIssuerURL    string `split_words:"true" required:"false" default:"https://dev-528196-admin.okta.com/oauth2/default"`
}
