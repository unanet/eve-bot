module github.com/unanet/eve-bot

go 1.16

// replace (
// 	github.com/unanet/go => ../go
// )

require (
	github.com/aws/aws-sdk-go v1.40.27
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/dghubble/sling v1.3.0
	github.com/go-chi/chi v4.1.0+incompatible
	github.com/go-chi/jwtauth v4.0.4+incompatible // indirect
	github.com/go-chi/render v1.0.1
	github.com/golang/mock v1.6.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/slack-go/slack v0.9.3
	github.com/unanet/eve v0.0.0-20210727233527-18eb447b118d
	github.com/unanet/go v0.0.0-20210911201142-860a725048b6
	go.uber.org/zap v1.18.1
)
