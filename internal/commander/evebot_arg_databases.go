package commander

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
