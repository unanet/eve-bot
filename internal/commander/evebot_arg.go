package commander

type EvebotCommandArg interface {
	Name() string
	Description() string
}

type EvebotCommandArgs []EvebotCommandArg

func (eba EvebotCommandArgs) String() string {
	var msg string
	for _, v := range eba {
		msg = msg + v.Name() + " - " + v.Description() + "\n"
	}
	return msg
}
