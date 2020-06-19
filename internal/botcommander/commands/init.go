package commands

var (
	CommandInitializerMap = map[string]interface{}{
		"help":    NewHelpCommand,
		"deploy":  NewDeployCommand,
		"migrate": NewMigrateCommand,
		"show":    NewShowCommand,
		"set":     NewSetCommand,
	}
)

func nonHelpCmd() []EvebotCommand {
	var cmds []EvebotCommand

	for k, v := range CommandInitializerMap {
		if k != "help" {
			cmds = append(cmds, v.(func([]string, string, string) EvebotCommand)([]string{}, "", ""))
		}
	}
	return cmds

}
