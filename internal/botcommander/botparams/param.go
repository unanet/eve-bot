package botparams

type Param interface {
	Name() string
	Description() string
	Value() string
}

type Params []Param

func (p Params) String() string {
	var msg string
	for _, v := range p {
		msg = msg + v.Name() + " - " + v.Description() + "\n"
	}
	return msg
}

type baseParam struct {
	name        string
	description string
	value       string
}

func GetEnvironmentValue(params Params) string {
	for _, v := range params {
		if v.Name() == "environment" {
			return v.Value()
		}
	}
	return ""
}

func GetNamespaceValue(params Params) string {
	for _, v := range params {
		if v.Name() == "namespace" {
			return v.Value()
		}
	}
	return ""
}
