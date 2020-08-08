package resources

const (
	// ServiceName resource key/id
	ServiceName = "services"
)

// Service resource data structure
type Service struct {
	baseResource
}

// Name satisfies the resource interface and returns the Service Name
func (e Service) Name() string {
	return e.name
}

// Description satisfies the resource interface and returns the Service Description
func (e Service) Description() string {
	return e.description
}

// Value satisfies the resource interface and returns the Service Value
func (e Service) Value() string {
	return e.value
}

// DefaultService returns the Default Service
func DefaultService() Service {
	return Service{baseResource{
		name:        ServiceName,
		description: "the service resource (i.e. platform, unanetbi, auto, subcontractor, unanet-analytics)",
	}}
}
