package eveapi

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve/pkg/eve"
)

func headerMsg(val string) string {
	return fmt.Sprintf("\n*%s*", strings.Title(strings.ToLower(val)))
}

func availableLabel(svc *eve.DeployService) string {
	if svc.ArtifactName == svc.ServiceName {
		return fmt.Sprintf("\n%s:%s", svc.ServiceName, svc.AvailableVersion)
	} else {
		return fmt.Sprintf("\n%s (%s):%s", svc.ServiceName, svc.ArtifactName, svc.AvailableVersion)
	}
}

func artifactResultBlock(svcs eve.DeployServices) string {
	result := ""

	if svcs == nil || len(svcs) == 0 {
		return ""
	}

	for _, svc := range svcs {
		if len(result) == 0 {
			result = availableLabel(svc)
		} else {
			result = result + availableLabel(svc)
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

func environmentNamespaceMsg(deploymentResponsePayload *eve.NSDeploymentPlan) string {
	return fmt.Sprintf("```Namespace: %s\nEnvironment: %s\nCluster: %s```", deploymentResponsePayload.Namespace.Alias, deploymentResponsePayload.EnvironmentName, deploymentResponsePayload.Namespace.ClusterName)
}
