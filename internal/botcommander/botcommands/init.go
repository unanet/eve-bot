package botcommands

var (
	CommandInitializerMap = map[string]interface{}{
		"help":    NewHelpCommand,
		"deploy":  NewDeployCommand,
		"migrate": NewMigrateCommand,
	}
)

func nonHelpCmd() []EvebotCommand {
	var cmds []EvebotCommand

	for k, v := range CommandInitializerMap {
		if k != "help" {
			cmds = append(cmds, v.(func([]string) EvebotCommand)([]string{}))
		}
	}

	return cmds
}
