package params

// Param is the parameter interface
type Param interface {
	Name() string
	Description() string
	Value() string
}

// Params is a slice of params
type Params []Param

// String satisfies the interface
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
