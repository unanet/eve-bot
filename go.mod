module gitlab.unanet.io/devops/eve-bot

go 1.15

// replace gitlab.unanet.io/devops/eve => ../eve
// replace gitlab.unanet.io/devops/go => ../go

require (
	github.com/aws/aws-sdk-go v1.36.14 // indirect
	github.com/dghubble/sling v1.3.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/golang/mock v1.4.4
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/slack-go/slack v0.7.4
	github.com/stretchr/testify v1.5.1 // indirect
	gitlab.unanet.io/devops/eve v0.6.1-0.20210107221943-4f5cddfc7a1e
	gitlab.unanet.io/devops/go v0.6.0
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/sys v0.0.0-20201221093633-bc327ba9c2f0 // indirect
	golang.org/x/tools v0.0.0-20200807224323-c05a0f5be48b // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)
