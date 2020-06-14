package resources

const (
	ServiceName = "services"
)

type Service struct {
	baseResource
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
	return Service{baseResource{
		name:        ServiceName,
		description: "the service resource (i.e. platform, unanetbi, auto, subcontractor, unanet-analytics)",
	}}
}
