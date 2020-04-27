package commander

type Arg interface {
	Name() string
	Description() string
}

type Args []Arg

func (a Args) String() string {
	var msg string
	for _, v := range a {
		msg = msg + v.Name() + " - " + v.Description() + "\n"
	}
	return msg
}
