package slackservice

// Config needed for slack
//		EVEBOT_SLACK_SIGNING_SECRET
//		EVEBOT_SLACK_VERIFICATION_TOKEN
//		EVEBOT_SLACK_USER_OAUTH_ACCESS_TOKEN
//		EVEBOT_SLACK_OAUTH_ACCESS_TOKEN
//		EVEBOT_SLACK_CHANNELS_AUTH
//		EVEBOT_SLACK_AUTH_ENABLED
// 		EVEBOT_SLACK_CHANNELS_MAINTENANCE
//		EVEBOT_SLACK_MAINTENANCE_ENABLED
type Config struct {
	SlackSigningSecret        string `split_words:"true" required:"true"`
	SlackVerificationToken    string `split_words:"true" required:"true"`
	SlackUserOauthAccessToken string `split_words:"true" required:"true"`
	SlackOauthAccessToken     string `split_words:"true" required:"true"`
	SlackChannelsAuth         string `split_words:"true" required:"true"`
	SlackChannelsMaintenance  string `split_words:"true" required:"true"`
	SlackAuthEnabled          bool   `split_words:"true" required:"true"`
	SlackMaintenanceEnabled   bool   `split_words:"true" default:"false"`
}
