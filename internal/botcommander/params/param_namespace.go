package params

const (
	// NamespaceName param key/id
	NamespaceName = "namespace"
)

// Namespace param data struct
type Namespace struct {
	baseParam
}

// Name satisfies the param interface and returns the Namespace Name
func (e Namespace) Name() string {
	return e.name
}

// Description satisfies the param interface and returns the Namespace Description
func (e Namespace) Description() string {
	return e.description
}

// Value satisfies the param interface and returns the Namespace Value
func (e Namespace) Value() string {
	return e.value
}

// DefaultNamespace is the default Namespace (used for help/init)
func DefaultNamespace() Namespace {
	return Namespace{baseParam{
		name:        NamespaceName,
		description: "the namespace where the resource lives (i.e. k8s namespaces)",
	}}
}
