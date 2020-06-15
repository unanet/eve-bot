package params

const (
	ServiceName = "service"
)

type Service struct {
	baseParam
}

func (e Service) Name() string {
	return e.name
}

func (e Service) Description() string {
	return e.description
}

func (e Service) Value() string {
	return e.value
}

func DefaultService() Service {
	return Service{baseParam{
		name:        ServiceName,
		description: "the service param (ex: unaneta, unanetb, platform)",
	}}
}
