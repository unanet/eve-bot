package eveapi

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

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

func headerMsg(val string) string {
	return fmt.Sprintf("\n*%s*", strings.Title(strings.ToLower(val)))
}

func availableLabel(svc *eve.DeployService) string {
	log.Logger.Debug("available label", zap.Any("deploy_service", *svc))
	return fmt.Sprintf("\n%s:%s", svc.ArtifactName, svc.AvailableVersion)
}

func deployedLabel(svc *eve.DeployService) string {
	log.Logger.Debug("deployed label", zap.Any("deploy_service", *svc))
	return fmt.Sprintf("\n%s:%s", svc.ArtifactName, svc.AvailableVersion)
}

func artifactResultBlock(svcs eve.DeployServices, eveResult eve.DeployArtifactResult) string {
	result := ""

	if svcs == nil || len(svcs) == 0 {
		return ""
	}

	for _, svc := range svcs {
		switch eveResult {
		case eve.DeployArtifactResultNoop:
			if len(result) == 0 {
				result = availableLabel(svc)
			} else {
				result = result + availableLabel(svc)
			}
		case eve.DeployArtifactResultSuccess:
			if len(result) == 0 {
				result = deployedLabel(svc)
			} else {
				result = result + deployedLabel(svc)
			}
		case eve.DeployArtifactResultFailed:
			if len(result) == 0 {
				result = availableLabel(svc)
			} else {
				result = result + availableLabel(svc)
			}
		}
	}

	return result
}

func apiMessages(msgs []string) string {
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

func environmentNamespaceMsg(env, ns string) string {
	return fmt.Sprintf("```Namespace: %s\nEnvironment: %s```", ns, env)
}

func (cbs *CallbackState) ToChatMsg() string {

	if cbs == nil {
		log.Logger.Error("invalid callback state")
		return ""
	}

	if cbs.Payload.NothingToDeploy() {
		return fmt.Sprintf("\n<@%s>, we're all caught up! There is nothing to deploy...\n", cbs.User)
	}

	var result string

	switch cbs.Payload.Status {
	case eve.DeploymentPlanStatusComplete:
		result = fmt.Sprintf("\n<@%s>, your deployment is complete...\n\n%s", cbs.User, environmentNamespaceMsg(cbs.Payload.EnvironmentName, cbs.Payload.Namespace.Alias))
	case eve.DeploymentPlanStatusErrors:
		result = fmt.Sprintf("\n<@%s>, we encountered some errors during the deployment...\n\n%s", cbs.User, environmentNamespaceMsg(cbs.Payload.EnvironmentName, cbs.Payload.Namespace.Alias))
	case eve.DeploymentPlanStatusDryrun:
		result = fmt.Sprintf("\n<@%s>, here's your *dryrun* results ...\n\n%s", cbs.User, environmentNamespaceMsg(cbs.Payload.EnvironmentName, cbs.Payload.Namespace.Alias))
	case eve.DeploymentPlanStatusPending:
		result = fmt.Sprintf("\n<@%s>, your deployment is pending. Here's the plan...\n\n%s", cbs.User, environmentNamespaceMsg(cbs.Payload.EnvironmentName, cbs.Payload.Namespace.Alias))
	}

	var deploymentResults string

	if cbs.Payload.Services != nil {
		for svcResult, svcs := range cbs.Payload.Services.TopResultMap() {
			// Let's break out early when this is a pending/dryrun result
			if cbs.Payload.Status == eve.DeploymentPlanStatusPending || cbs.Payload.Status == eve.DeploymentPlanStatusDryrun {
				deploymentResults = "\n```" + artifactResultBlock(svcs, svcResult) + "```"
				break
			}

			if len(deploymentResults) == 0 {
				deploymentResults = headerMsg(svcResult.String()) + "\n```" + artifactResultBlock(svcs, svcResult) + "```"
			} else {
				deploymentResults = deploymentResults + headerMsg(svcResult.String()) + "\n```" + artifactResultBlock(svcs, svcResult) + "```"
			}
		}
		result = result + "\n" + deploymentResults
	}

	if cbs.Payload.Messages == nil || len(cbs.Payload.Messages) == 0 {
		return result
	}

	return result + headerMsg("Messages") + "\n```" + apiMessages(cbs.Payload.Messages) + "```"
}
