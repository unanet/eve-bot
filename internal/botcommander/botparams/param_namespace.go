package botparams

type Namespace struct {
	baseParam
}

func (e Namespace) Name() string {
	return e.name
}

func (e Namespace) Description() string {
	return e.description
}

func DefaultNamespace() Namespace {
	return Namespace{baseParam{
		name:        "namespace",
		description: "the namespace where the resource lives (i.e. k8s namespaces)",
	}}
}

func NewNamespaceParam(val string) Namespace {
	p := DefaultNamespace()
	p.value = val
	return p
}
