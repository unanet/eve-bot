package botparams

type Param interface {
	Name() string
	Description() string
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
