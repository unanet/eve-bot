package eveapi

import (
	"fmt"

	"gitlab.unanet.io/devops/eve/pkg/eve"
)

type CallbackState struct {
	User    string               `json:"user"`
	Channel string               `json:"channel"`
	Payload eve.NSDeploymentPlan `json:"payload"`
}

//func newBlockMsgOpt(text string) slack.MsgOption {
//	return slack.MsgOptionBlocks(
//		slack.NewSectionBlock(
//			slack.NewTextBlockObject(
//				slack.MarkdownType,
//				text,
//				false,
//				false),
//			nil,
//			nil),
//		slack.NewDividerBlock())
//}

func artifactResultMsg(services eve.DeployServices) string {
	successfulResultsMsg := ""
	successfulResultsHeader := "*Successful:*\n"
	successfulResults := ""
	failedResultsMsg := ""
	failedResultsHeader := "*Failed:*\n"
	failedResults := ""
	noopResultsMsg := ""
	//noopResultsHeader := "*Noop:*\n"
	noopResults := ""
	for _, svc := range services {
		switch svc.Result {
		case eve.DeployArtifactResultFailed:
			if len(failedResults) == 0 {
				failedResults = fmt.Sprintf("\n%s:%s", svc.ArtifactName, svc.AvailableVersion)
			} else {
				failedResults = failedResults + fmt.Sprintf("\n%s:%s", svc.ArtifactName, svc.AvailableVersion)
			}
		case eve.DeployArtifactResultSucceeded:
			if len(successfulResults) == 0 {
				successfulResults = fmt.Sprintf("\n%s:%s", svc.ArtifactName, svc.DeployedVersion)
			} else {
				successfulResults = successfulResults + fmt.Sprintf("\n%s:%s", svc.ArtifactName, svc.DeployedVersion)
			}
		case eve.DeployArtifactResultNoop:
			if len(noopResults) == 0 {
				noopResults = fmt.Sprintf("\n%s:%s", svc.ArtifactName, svc.AvailableVersion)
			} else {
				noopResults = noopResults + fmt.Sprintf("\n%s:%s", svc.ArtifactName, svc.AvailableVersion)
			}
		}
	}

	if len(successfulResults) > 0 {
		successfulResultsMsg = successfulResultsHeader + successfulResults
	}

	if len(failedResults) > 0 {
		failedResultsMsg = failedResultsHeader + failedResults
	}

	if len(noopResults) > 0 {
		noopResultsMsg = noopResults
	}

	return successfulResultsMsg + failedResultsMsg + noopResultsMsg
}

func apiMessages(msgs []string) string {
	infoHeader := "*Info:*\n"
	infoMsgs := ""
	for _, msg := range msgs {
		if len(infoMsgs) == 0 {
			infoMsgs = "```\n- " + msg + "\n"
		} else {
			infoMsgs = infoMsgs + "- " + msg + "\n"
		}
	}
	if len(infoMsgs) == 0 {
		return ""
	}
	return infoHeader + infoMsgs + "```\n"
}

func environmentNamespaceMsg(env, ns string) string {
	return fmt.Sprintf("```Namespace: %s\nEnvironment: %s```\n\n", ns, env)
}

func (cbs *CallbackState) SlackMsgHeader() string {
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

	artifactMsg := artifactResultMsg(cbs.Payload.Services)
	apiMsgs := apiMessages(cbs.Payload.Messages)

	if len(artifactMsg) > 0 {
		artifactMsg = "```" + artifactMsg + "```"
	}

	if len(apiMsgs) > 0 {
		apiMsgs = "```" + apiMsgs + "```"
	}

	return artifactMsg + apiMsgs
}

type DeploymentPlanOptions struct {
	Artifacts   []ArtifactDefinition `json:"artifacts"`
	ForceDeploy bool                 `json:"force_deploy"`
	DryRun      bool                 `json:"dry_run"`
	CallbackURL string               `json:"callback_url"`
	Environment string               `json:"environment"`
	Namespaces  []string             `json:"namespaces,omitempty"`
	Messages    []string             `json:"messages,omitempty"`
	Type        string               `json:"type"`
	User        string               `json:"user"`
}

type ArtifactDefinition struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	RequestedVersion string `json:"requested_version,omitempty"`
	AvailableVersion string `json:"available_version"`
	ArtifactoryFeed  string `json:"artifactory_feed"`
	ArtifactoryPath  string `json:"artifactory_path"`
	Matched          bool   `json:"-"`
}
