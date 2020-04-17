package okta

// Config contains the okta configuration
type Config struct {
	ClientID     string `split_words:"true" required:"true"`
	ClientSecret string `split_words:"true" required:"true"`
	IssuerURL    string `split_words:"true" required:"true" default:"https://dev-528196-admin.okta.com/oauth2/default"`
}
