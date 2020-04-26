package commander

func init() {
	// Add all the Evebot commands here on init
	evebotCommands = []EvebotCommand{
		NewEvebotHelpCommand(),
		NewEvebotDeployCommand(),
	}
}

var (
	evebotCommands []EvebotCommand
)
