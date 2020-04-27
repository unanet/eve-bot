package botcommands

func init() {
	// Add all the Evebot commands here on init
	EvebotCommands = []EvebotCommand{
		NewEvebotHelpCommand(),
		NewEvebotDeployCommand(),
		NewEvebotMigrateCommand(),
	}
}

var (
	EvebotCommands []EvebotCommand
)
