package botcommands

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botparams"
)

func NewRestartCommand(cmdFields []string) EvebotCommand {
	cmd := defaultRestartCommand()
	cmd.input = cmdFields
	cmd.resolveParams()
	cmd.resolveArgs()
	return cmd
}

type RestartCmd struct {
	baseCommand
}

func defaultRestartCommand() RestartCmd {
	return RestartCmd{baseCommand{
		name:    "restart",
		summary: "The `restart` command is used to restart services in a specific *namespace* and *environment*",
		usage: bothelp.Usage{
			"restart {{ namespace }} in {{ environment }}",
			"restart {{ namespace }} in {{ environment }} services={{ service_name,service_name }}",
		},
		examples: bothelp.Examples{
			"restart current in una-int",
			"restart current in una-int services=unanetbi",
			"restart current in una-int services=unanetbi,unaneta",
		},
		async:               true,
		optionalArgs:        botargs.Args{botargs.DefaultServicesArg()},
		requiredParams:      botparams.Params{botparams.DefaultNamespace(), botparams.DefaultEnvironment()},
		apiOptions:          make(map[string]interface{}),
		requiredInputLength: 4,
	}}
}
