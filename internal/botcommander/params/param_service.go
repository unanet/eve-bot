package params

const (
	// ServiceName param key/id
	ServiceName = "service"
)

// Service param data struct
type Service struct {
	baseParam
}

// Name satisfies the param interface and returns the Service Name
func (e Service) Name() string {
	return e.name
}

// Description satisfies the param interface and returns the Service Description
func (e Service) Description() string {
	return e.description
}

// Value satisfies the param interface and returns the Service Value
func (e Service) Value() string {
	return e.value
}
