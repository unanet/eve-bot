package commander

type EvebotCommandArgs []EvebotArg

func (eba EvebotCommandArgs) String() string {
	var msg string
	for _, v := range eba {
		msg = msg + v.Name() + " - " + v.Description() + "\n"
	}
	return msg
}

type EvebotArg interface {
	Name() string
	Description() string
}

/*
	ARGUMENT: Dryrun
*/

type EvebotArgDryrun bool

func (ebad EvebotArgDryrun) Name() string {
	return "dryrun"
}

func (ebad EvebotArgDryrun) Description() string {
	return "generates plan only doesn't actually change any state"
}

func NewDryrunArg() EvebotArgDryrun {
	return false
}

/*
	ARGUMENT: Force
*/

type EvebotArgForce bool

func (ebaf EvebotArgForce) Name() string {
	return "force"
}

func (ebaf EvebotArgForce) Description() string {
	return "forces the command to execute"
}

func NewForceArg() EvebotArgForce {
	return false
}

/*
	ARGUMENT: Services
*/

type EvebotArgServices []string

func (ebas EvebotArgServices) Name() string {
	return "services"
}

func (ebas EvebotArgServices) Description() string {
	return "comma separated list of services with name:version syntax"
}

func NewServicesArg() EvebotArgServices {
	return EvebotArgServices{}
}

/*
	ARGUMENT: Databases
*/

type EvebotArgDatabases []string

func (ebad EvebotArgDatabases) Name() string {
	return "databases"
}

func (ebad EvebotArgDatabases) Description() string {
	return "comma separated list of databases"
}

func NewDatabasesArg() EvebotArgDatabases {
	return EvebotArgDatabases{}
}
