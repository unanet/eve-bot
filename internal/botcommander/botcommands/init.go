package botcommands

func init() {
	// Add all the Evebot commands here on init
	EvebotCommands = []EvebotCommand{
		DefaultHelpCommand(),
		DefaultDeployCommand(),
		DefaultMigrateCommand(),
	}
}

var (
	EvebotCommands []EvebotCommand
)
