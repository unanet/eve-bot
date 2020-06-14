package resources

const (
	NamespaceName = "namespaces"
)

type Namespace struct {
	baseResource
}

func (e Namespace) Name() string {
	return e.name
}

func (e Namespace) Description() string {
	return e.description
}

func (e Namespace) Value() string {
	return e.value
}

func DefaultNamespace() Namespace {
	return Namespace{baseResource{
		name:        NamespaceName,
		description: "the namespace resource (i.e. current, prev, prev-1)",
	}}
}
