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
	return fmt.Sprintf("\n%s", strings.Title(strings.ToLower(val)))
}

func availableLabel(svc *eve.DeployService) string {
	log.Logger.Debug("available label", zap.Any("deploy_service", *svc))
	return fmt.Sprintf("\n%s:%s", svc.ArtifactName, svc.AvailableVersion)
}

func deployedLabel(svc *eve.DeployService) string {
	log.Logger.Debug("deployed label", zap.Any("deploy_service", *svc))
	return fmt.Sprintf("\n%s:%s", svc.ArtifactName, svc.DeployedVersion)
}

func artifactResultBlock(svcResultMap eve.ArtifactDeployResultMap, eveResult eve.DeployArtifactResult, status eve.DeploymentPlanStatus) string {
	result := ""

	svcs := svcResultMap[eveResult]

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

	// this is for the initial callback when we are telling the user about the plan
	if status == eve.DeploymentPlanStatusPending {
		return headerMsg("plan") + result + "\n"
	}

	return headerMsg(eveResult.String()) + result + "\n"
}

func apiMessages(msgs []string) string {
	infoHeader := "Messages:\n"
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
	return infoHeader + infoMsgs
}

func environmentNamespaceMsg(env, ns string) string {
	return fmt.Sprintf("```Namespace: %s\nEnvironment: %s```\n\n", ns, env)
}

func (cbs *CallbackState) SlackMsgHeader() string {

	if cbs.Payload.NothingToDeploy() {
		return fmt.Sprintf("\n<@%s>, we're all caught up! There is nothing to deploy...\n", cbs.User)
	}

	switch cbs.Payload.Status {
	case eve.DeploymentPlanStatusComplete:
		return fmt.Sprintf("\n<@%s>, your deployment is complete...\n\n%s", cbs.User, environmentNamespaceMsg(cbs.Payload.EnvironmentName, cbs.Payload.Namespace.Alias))
	case eve.DeploymentPlanStatusErrors:
		return fmt.Sprintf("\n<@%s>, we encountered some errors during the deployment...\n\n%s", cbs.User, environmentNamespaceMsg(cbs.Payload.EnvironmentName, cbs.Payload.Namespace.Alias))
	case eve.DeploymentPlanStatusDryrun:
		return fmt.Sprintf("\n<@%s>, here's your *dryrun* results ...\n\n%s", cbs.User, environmentNamespaceMsg(cbs.Payload.EnvironmentName, cbs.Payload.Namespace.Alias))
	case eve.DeploymentPlanStatusPending:
		return fmt.Sprintf("\n<@%s>, your deployment is pending. Here's the plan...\n\n%s", cbs.User, environmentNamespaceMsg(cbs.Payload.EnvironmentName, cbs.Payload.Namespace.Alias))
	default:
		return ""
	}
}

func (cbs *CallbackState) SlackMsgResults() string {
	var artifactMsg, apiMsgs string

	if cbs == nil {
		log.Logger.Error("invalid callback state")
		return ""
	}

	if cbs.Payload.Services != nil {
		svcMap := cbs.Payload.Services.TopResultMap()
		log.Logger.Debug("svcMap", zap.Any("map_val", svcMap))

		artifactMsg = artifactResultBlock(svcMap, eve.DeployArtifactResultFailed, cbs.Payload.Status) +
			artifactResultBlock(svcMap, eve.DeployArtifactResultSuccess, cbs.Payload.Status) +
			artifactResultBlock(svcMap, eve.DeployArtifactResultNoop, cbs.Payload.Status)
	}

	if cbs.Payload.Messages != nil {
		apiMsgs = apiMessages(cbs.Payload.Messages)
	}

	if len(artifactMsg) > 0 {
		artifactMsg = "```\n" + artifactMsg + "\n```"
	}

	if len(apiMsgs) > 0 {
		apiMsgs = "\n```\n" + apiMsgs + "\n```\n"
	}

	return artifactMsg + apiMsgs
}
