package resources

const (
	// NamespaceName resource key/id
	NamespaceName = "namespaces"
)

// Namespace resource data structure
type Namespace struct {
	baseResource
}

// Name satisfies the resource interface and returns the Namespace Name
func (e Namespace) Name() string {
	return e.name
}

// Description satisfies the resource interface and returns the Namespace Description
func (e Namespace) Description() string {
	return e.description
}

// Value satisfies the resource interface and returns the Namespace Value
func (e Namespace) Value() string {
	return e.value
}

// DefaultNamespace returns the Default Namespace
func DefaultNamespace() Namespace {
	return Namespace{baseResource{
		name:        NamespaceName,
		description: "the namespace resource (i.e. current, prev, prev-1)",
	}}
}
