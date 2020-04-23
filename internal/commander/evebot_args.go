package commander

type EvebotArg interface {
	Name() string
}

type EvebotArgDryrun bool
type EvebotArgForce bool
type EvebotArgServices []string
type EvebotArgDatabases []string
type EvebotArgs []EvebotArg

func (ebad EvebotArgDryrun) Name() string {
	return "dryrun"
}

func (ebaf EvebotArgForce) Name() string {
	return "force"
}

func (ebas EvebotArgServices) Name() string {
	return "services"
}

func (ebad EvebotArgDatabases) Name() string {
	return "databases"
}
