module gitlab.unanet.io/devops/eve-bot

go 1.15

// replace gitlab.unanet.io/devops/eve => ../eve

require (
	github.com/aws/aws-sdk-go v1.32.4 // indirect
	github.com/dghubble/sling v1.3.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/golang/mock v1.4.4
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/slack-go/slack v0.6.5
	gitlab.unanet.io/devops/eve v0.5.1-0.20201217234615-52ce60bce87a
	gitlab.unanet.io/devops/go v0.2.0
	go.uber.org/zap v1.16.0
	golang.org/x/tools v0.0.0-20200807224323-c05a0f5be48b // indirect
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/asaskevich/govalidator.v9 v9.0.0-20180315120708-ccb8e960c48f // indirect
)
