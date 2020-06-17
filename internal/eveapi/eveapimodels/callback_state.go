package eveapimodels

import (
	"fmt"

	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

type CallbackState struct {
	User    string               `json:"user"`
	Channel string               `json:"channel"`
	TS      string               `json:"ts"`
	Payload eve.NSDeploymentPlan `json:"payload"`
}

func (cbs *CallbackState) cleanUser() {
	if cbs.User == "" {
		return
	}
	if cbs.User == "channel" {
		cbs.User = "!channel"
		return
	} else {
		cbs.User = "@" + cbs.User
		return
	}
}

func (cbs *CallbackState) nothingToDeployResponse() string {
	msg := fmt.Sprintf("\n<%s>, we're all caught up! There is nothing to deploy...\n", cbs.User)
	if len(cbs.Payload.Messages) > 0 {
		return msg + HeaderMsg("Messages") + "\n```" + APIMessages(cbs.Payload.Messages) + "```"
	}
	return msg
}

func (cbs *CallbackState) appendDeployServicesResult(result *string) {
	var deployServicesResults string
	if cbs.Payload.Services != nil {
		for svcResult, svcs := range cbs.Payload.Services.ToResultMap() {
			// Let's break out early when this is a pending/dryrun result
			if cbs.Payload.Status == eve.DeploymentPlanStatusPending || cbs.Payload.Status == eve.DeploymentPlanStatusDryrun {
				deployServicesResults = "\n```" + ServicesResultBlock(svcs) + "```"
				break
			}

			svcResultMessage := HeaderMsg(svcResult.String()) + "\n```" + ServicesResultBlock(svcs) + "```"

			if len(deployServicesResults) == 0 {
				deployServicesResults = svcResultMessage
			} else {
				deployServicesResults = deployServicesResults + svcResultMessage
			}
		}
		*result = *result + "\n" + deployServicesResults
	}
}

func (cbs *CallbackState) appendDeployMigrationsResult(result *string) {
	var deployMigrationsResults string
	if cbs.Payload.Migrations != nil {
		for migResult, migs := range cbs.Payload.Migrations.ToResultMap() {
			// Let's break out early when this is a pending/dryrun result
			if cbs.Payload.Status == eve.DeploymentPlanStatusPending || cbs.Payload.Status == eve.DeploymentPlanStatusDryrun {
				deployMigrationsResults = "\n```" + MigrationResultBlock(migs) + "```"
				break
			}

			svcResultMessage := HeaderMsg(migResult.String()) + "\n```" + MigrationResultBlock(migs) + "```"

			if len(deployMigrationsResults) == 0 {
				deployMigrationsResults = svcResultMessage
			} else {
				deployMigrationsResults = deployMigrationsResults + svcResultMessage
			}
		}
		*result = *result + "\n" + deployMigrationsResults
	}
}

func (cbs *CallbackState) initialResult() string {
	var ackMessage string
	switch cbs.Payload.Status {
	case eve.DeploymentPlanStatusComplete:
		ackMessage = fmt.Sprintf("your %s deployment is complete", cbs.Payload.DeploymentPlanType())
	case eve.DeploymentPlanStatusErrors:
		ackMessage = "we encountered some errors"
	case eve.DeploymentPlanStatusDryrun:
		ackMessage = "here's your *dryrun* results"
	case eve.DeploymentPlanStatusPending:
		ackMessage = fmt.Sprintf("your %s deployment is pending, here's the plan", cbs.Payload.DeploymentPlanType())
	}

	var result string
	if len(cbs.User) > 0 {
		result = fmt.Sprintf("\n<%s>, %s...\n\n%s", cbs.User, ackMessage, EnvironmentNamespaceMsg(&cbs.Payload))
	} else {
		result = fmt.Sprintf("\n%s...\n\n%s", ackMessage, EnvironmentNamespaceMsg(&cbs.Payload))
	}
	return result
}

func (cbs *CallbackState) appendApiMessages(result *string) string {
	return *result + HeaderMsg("Messages") + "\n```" + APIMessages(cbs.Payload.Messages) + "```"
}

// ToChatMsg takes the eve-api callback payload
// and converts it to a Chat Message (string with formatting/proper messaging)
func (cbs *CallbackState) ToChatMsg() string {
	if cbs == nil {
		log.Logger.Error("invalid callback state")
		return ""
	}

	cbs.cleanUser()

	if cbs.Payload.NothingToDeploy() {
		return cbs.nothingToDeployResponse()
	}

	result := cbs.initialResult()

	cbs.appendDeployServicesResult(&result)

	cbs.appendDeployMigrationsResult(&result)

	if cbs.Payload.Messages == nil || len(cbs.Payload.Messages) == 0 {
		return result
	}

	return cbs.appendApiMessages(&result)
}
