package eveapi

import (
	"fmt"

	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

type CallbackState struct {
	User    string               `json:"user"`
	Channel string               `json:"channel"`
	TS      string               `json:"ts"`
	Action  string               `json:"action"`
	Payload eve.NSDeploymentPlan `json:"payload"`
}

type ArtifactDefinitions []*ArtifactDefinition

type StringList []string

type DeploymentPlanType string

func (dpt DeploymentPlanType) String() string {
	return string(dpt)
}

type DeploymentPlanOptions struct {
	Artifacts        ArtifactDefinitions `json:"artifacts"`
	ForceDeploy      bool                `json:"force_deploy"`
	User             string              `json:"user"`
	DryRun           bool                `json:"dry_run"`
	CallbackURL      string              `json:"callback_url"`
	Environment      string              `json:"environment"`
	NamespaceAliases StringList          `json:"namespaces,omitempty"`
	Messages         []string            `json:"messages,omitempty"`
	Type             DeploymentPlanType  `json:"type"`
}

type ArtifactDefinition struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	RequestedVersion string `json:"requested_version,omitempty"`
	AvailableVersion string `json:"available_version"`
	ArtifactoryFeed  string `json:"artifactory_feed"`
	ArtifactoryPath  string `json:"artifactory_path"`
	FunctionPointer  string `json:"function_pointer"`
	FeedType         string `json:"feed_type"`
	Matched          bool   `json:"-"`
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

func (cbs *CallbackState) cleanAction() {
	if cbs.Action == "" {
		cbs.Action = "job"
		return
	}
	if cbs.Action == "application" {
		cbs.Action = "deployment"
		return
	}
}

func (cbs *CallbackState) nothingToDeployResponse() string {
	msg := fmt.Sprintf("\n<%s>, we're all caught up! There is nothing to deploy...\n", cbs.User)
	if len(cbs.Payload.Messages) > 0 {
		return msg + headerMsg("Messages") + "\n```" + apiMessages(cbs.Payload.Messages) + "```"
	}
	return msg
}

func (cbs *CallbackState) appendDeployServicesResult(result *string) {
	var deployServicesResults string
	if cbs.Payload.Services != nil {
		for svcResult, svcs := range cbs.Payload.Services.ToResultMap() {
			// Let's break out early when this is a pending/dryrun result
			if cbs.Payload.Status == eve.DeploymentPlanStatusPending || cbs.Payload.Status == eve.DeploymentPlanStatusDryrun {
				deployServicesResults = "\n```" + servicesResultBlock(svcs) + "```"
				break
			}

			svcResultMessage := headerMsg(svcResult.String()) + "\n```" + servicesResultBlock(svcs) + "```"

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
				deployMigrationsResults = "\n```" + migrationResultBlock(migs) + "```"
				break
			}

			svcResultMessage := headerMsg(migResult.String()) + "\n```" + migrationResultBlock(migs) + "```"

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
		ackMessage = fmt.Sprintf("your %s is complete", cbs.Action)
	case eve.DeploymentPlanStatusErrors:
		ackMessage = "we encountered some errors"
	case eve.DeploymentPlanStatusDryrun:
		ackMessage = "here's your *dryrun* results"
	case eve.DeploymentPlanStatusPending:
		ackMessage = fmt.Sprintf("your %s is pending, here's the plan", cbs.Action)
	}

	var result string
	if len(cbs.User) > 0 {
		result = fmt.Sprintf("\n<%s>, %s...\n\n%s", cbs.User, ackMessage, environmentNamespaceMsg(&cbs.Payload))
	} else {
		result = fmt.Sprintf("\n%s...\n\n%s", ackMessage, environmentNamespaceMsg(&cbs.Payload))
	}
	return result
}

func (cbs *CallbackState) appendApiMessages(result *string) string {
	return *result + headerMsg("Messages") + "\n```" + apiMessages(cbs.Payload.Messages) + "```"
}

// ToChatMsg takes the eve-api callback payload
// and converts it to a Chat Message (string with formatting/proper messaging)
func (cbs *CallbackState) ToChatMsg() string {
	if cbs == nil {
		log.Logger.Error("invalid callback state")
		return ""
	}

	cbs.cleanUser()

	cbs.cleanAction()

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
