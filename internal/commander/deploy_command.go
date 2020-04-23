package commander

import "strings"

type EvebotDeployCommand struct {
}

func NewEvebotDeployCommand() *EvebotDeployCommand {
	return &EvebotDeployCommand{}
}

func (edc *EvebotDeployCommand) Name() string {
	return "deploy"
}

func (edc *EvebotDeployCommand) Examples() EvebotCommandExamples {
	return EvebotCommandExamples{
		"- deploy {{ namespace }} in {{ environment }}",
		"- deploy {{ namespace }} in {{ environment }} services={{ artifact_name:artifact_version }}",
		"- deploy {{ namespace }} in {{ environment }} services={{ artifact_name:artifact_version }} dryrun={{ true }}",
		"- deploy {{ namespace }} in {{ environment }} services={{ artifact_name:artifact_version }} dryrun={{ true }} redeploy={{ true }}",
		"\n",
		"`Examples:`",
		"- deploy current in qa",
		"- deploy current in qa services=infocus-cloud-client:2020.1 dryrun=true",
		"- deploy current in qa services=infocus-cloud-client:2020.1,infocus-proxy:2020.1 dryrun=true redeploy=true",
	}
}

func (edc *EvebotDeployCommand) IsHelpRequest(input []string) bool {
	for _, s := range input {
		if strings.ToLower(s) == "help" {
			return true
		}
	}
	return false
}
