module gitlab.unanet.io/devops/eve-bot

go 1.15

// replace gitlab.unanet.io/devops/eve => ../eve
// replace gitlab.unanet.io/devops/go => ../go

require (
	github.com/aws/aws-sdk-go v1.36.29 // indirect
	github.com/dghubble/sling v1.3.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/golang-migrate/migrate/v4 v4.14.1 // indirect
	github.com/golang/mock v1.4.4
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lib/pq v1.9.0 // indirect
	github.com/prometheus/procfs v0.3.0 // indirect
	github.com/slack-go/slack v0.8.0
	gitlab.unanet.io/devops/eve v0.10.4
	gitlab.unanet.io/devops/go v1.0.5
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
)
