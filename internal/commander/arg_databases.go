package commander

/*
	ARGUMENT: Databases
*/

type ArgDatabases []string

func (ebad ArgDatabases) Name() string {
	return "databases"
}

func (ebad ArgDatabases) Description() string {
	return "comma separated list of databases"
}

func NewDatabasesArg() ArgDatabases {
	return ArgDatabases{}
}
