package eveapi

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/go/pkg/log"
)

const allCaughtUpMsg = "We're all caught up! There is nothing to deploy..."

// CallbackState data structure
type CallbackState struct {
	User    string               `json:"user"`
	Channel string               `json:"channel"`
	TS      string               `json:"ts"`
	Payload eve.NSDeploymentPlan `json:"payload"`
}

// ToChatMsg converts the eve-api callback payload to a Chat Message (string with formatting/proper messaging)
func (cbs *CallbackState) ToChatMsg() string {
	if cbs == nil {
		log.Logger.Error("invalid callback state")
		return ""
	}

	cbs.cleanUser()

	if cbs.Payload.NothingToDeploy() && cbs.Payload.Status != eve.DeploymentPlanStatusMessage {
		return cbs.nothingToDeployResponse()
	}

	result := cbs.initMsg()

	cbs.appendDeployServicesResult(&result)

	cbs.appendDeployJobsResult(&result)

	return cbs.appendAPIMessages(&result)
}

// messages converts a slice of strings into a string message
func messages(msgs []string) string {
	infoMsgs := ""
	for _, msg := range msgs {
		if len(infoMsgs) == 0 {
			infoMsgs = "\n- " + msg
		} else {
			infoMsgs = infoMsgs + "\n- " + msg
		}
	}
	if len(infoMsgs) == 0 {
		return ""
	}
	return infoMsgs
}

// headerMsg formats a header msg
func headerMsg(val string) string {
	return fmt.Sprintf("\n*%s*", strings.Title(strings.ToLower(val)))
}

func (cbs *CallbackState) cleanUser() {
	switch cbs.User {
	case "":
		return
	case "channel":
		cbs.User = "!channel"
		return
	default:
		cbs.User = "@" + cbs.User
		return
	}
}

func (cbs *CallbackState) nothingToDeployResponse() string {
	msg := ""
	if len(cbs.User) > 0 {
		msg = fmt.Sprintf("\n<%s>, %s\n", cbs.User, allCaughtUpMsg)
	} else {
		msg = fmt.Sprintf("\n%s\n", allCaughtUpMsg)
	}

	details := ""
	if (cbs.Payload.Namespace != nil) && len(cbs.Payload.Namespace.Alias) > 0 {
		details = fmt.Sprintf("Namespace: %s\n", cbs.Payload.Namespace.Alias)
	}

	if len(cbs.Payload.EnvironmentAlias) > 0 {
		details = details + fmt.Sprintf("Environment: %s\n", cbs.Payload.EnvironmentName)
	}

	if (cbs.Payload.Namespace != nil) && len(cbs.Payload.Namespace.ClusterName) > 0 {
		details = details + fmt.Sprintf("Cluster: %s\n", cbs.Payload.Namespace.ClusterName)
	}

	if len(details) > 0 {
		msg = msg + "```" + details + "```"
	}

	if len(cbs.Payload.Messages) > 0 {
		return msg + headerMsg("Messages") + "\n```" + messages(cbs.Payload.Messages) + "```"
	}
	return msg
}

func (cbs *CallbackState) appendDeployServicesResult(result *string) {
	var deployServicesResults string
	if cbs.Payload.Services != nil {
		for svcResult, svcs := range cbs.Payload.Services.ToResultMap() {
			// Let's break out early when this is a pending/dryrun result
			if cbs.Payload.Status == eve.DeploymentPlanStatusPending || cbs.Payload.Status == eve.DeploymentPlanStatusDryrun {
				deployServicesResults = "\n```" + ChatMessage(svcs) + "```"
				break
			}

			svcResultMessage := headerMsg(svcResult.String()) + "\n```" + ChatMessage(svcs) + "```"

			if len(deployServicesResults) == 0 {
				deployServicesResults = svcResultMessage
			} else {
				deployServicesResults = deployServicesResults + svcResultMessage
			}
		}
		*result = *result + "\n" + deployServicesResults
	}
}

func (cbs *CallbackState) appendDeployJobsResult(result *string) {
	var deployJobsResults string
	if cbs.Payload.Jobs != nil {
		for jobResult, jobs := range cbs.Payload.Jobs.ToResultMap() {
			// Let's break out early when this is a pending/dryrun result
			if cbs.Payload.Status == eve.DeploymentPlanStatusPending || cbs.Payload.Status == eve.DeploymentPlanStatusDryrun {
				deployJobsResults = "\n```" + ChatMessage(jobs) + "```"
				break
			}

			jobResultMessage := headerMsg(jobResult.String()) + "\n```" + ChatMessage(jobs) + "```"

			if len(deployJobsResults) == 0 {
				deployJobsResults = jobResultMessage
			} else {
				deployJobsResults = deployJobsResults + jobResultMessage
			}
		}
		*result = *result + "\n" + deployJobsResults
	}
}

func (cbs *CallbackState) initMsg() string {
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
	case eve.DeploymentPlanStatusMessage:
		ackMessage = fmt.Sprintf("your %s deployment is complete", cbs.Payload.DeploymentPlanType())
	}

	var result string
	if len(cbs.User) > 0 {
		result = fmt.Sprintf("\n<%s>, %s...\n\n%s", cbs.User, ackMessage, ChatMessage(&cbs.Payload))
	} else {
		result = fmt.Sprintf("\n%s...\n\n%s", ackMessage, ChatMessage(&cbs.Payload))
	}
	return result
}

func (cbs *CallbackState) appendAPIMessages(result *string) string {
	if cbs.Payload.Messages == nil || len(cbs.Payload.Messages) == 0 {
		return *result
	}

	return *result + headerMsg("Messages") + "\n```" + messages(cbs.Payload.Messages) + "```"
}
