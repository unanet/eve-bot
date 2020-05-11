package eveapi

import (
	"fmt"

	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

type CallbackState struct {
	User    string               `json:"user"`
	Channel string               `json:"channel"`
	Payload eve.NSDeploymentPlan `json:"payload"`
}

type ArtifactDefinitions []*ArtifactDefinition

type StringList []string

type DeploymentPlanType string

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

// ToChatMsg takes the eve-api callback payload
// and converts it to a Chat Message (string with formatting/proper messaging)
func (cbs *CallbackState) ToChatMsg() string {
	if cbs == nil {
		log.Logger.Error("invalid callback state")
		return ""
	}

	if cbs.Payload.NothingToDeploy() {
		return fmt.Sprintf("\n<@%s>, we're all caught up! There is nothing to deploy...\n", cbs.User)
	}

	var ackMessage string
	switch cbs.Payload.Status {
	case eve.DeploymentPlanStatusComplete:
		ackMessage = "your deployment is complete"
	case eve.DeploymentPlanStatusErrors:
		ackMessage = "we encountered some errors during the deployment"
	case eve.DeploymentPlanStatusDryrun:
		ackMessage = "here's your *dryrun* results"
	case eve.DeploymentPlanStatusPending:
		ackMessage = "your deployment is pending, here's the plan"
	}

	result := fmt.Sprintf("\n<@%s>, %s...\n\n%s", cbs.User, ackMessage, environmentNamespaceMsg(cbs.Payload.EnvironmentName, cbs.Payload.Namespace.Alias))

	var deploymentResults string

	if cbs.Payload.Services != nil {
		for svcResult, svcs := range cbs.Payload.Services.TopResultMap() {
			// Let's break out early when this is a pending/dryrun result
			if cbs.Payload.Status == eve.DeploymentPlanStatusPending || cbs.Payload.Status == eve.DeploymentPlanStatusDryrun {
				deploymentResults = "\n```" + artifactResultBlock(svcs, svcResult) + "```"
				break
			}

			svcResultMessage := headerMsg(svcResult.String()) + "\n```" + artifactResultBlock(svcs, svcResult) + "```"

			if len(deploymentResults) == 0 {
				deploymentResults = svcResultMessage
			} else {
				deploymentResults = deploymentResults + svcResultMessage
			}
		}
		result = result + "\n" + deploymentResults
	}

	if cbs.Payload.Messages == nil || len(cbs.Payload.Messages) == 0 {
		return result
	}

	return result + headerMsg("Messages") + "\n```" + apiMessages(cbs.Payload.Messages) + "```"
}
