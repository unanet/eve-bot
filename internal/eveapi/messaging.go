package eveapi

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

func headerMsg(val string) string {
	return fmt.Sprintf("\n*%s*", strings.Title(strings.ToLower(val)))
}

func availableLabel(svc *eve.DeployService) string {
	log.Logger.Debug("available label", zap.Any("deploy_service", *svc))
	return fmt.Sprintf("\n%s:%s", svc.ArtifactName, svc.AvailableVersion)
}


func artifactResultBlock(svcs eve.DeployServices, eveResult eve.DeployArtifactResult) string {
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

func environmentNamespaceMsg(env, ns string) string {
	return fmt.Sprintf("```Namespace: %s\nEnvironment: %s```", ns, env)
}
