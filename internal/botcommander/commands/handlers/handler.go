package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

var (
	errInvalidAPIResp = errors.New("invalid api response")
)

// CommandHandler is the interface that Handles EvebotCommands
type CommandHandler interface {
	Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string)
}

func mapToEveService(s eve.Service) eveapimodels.EveService {
	return eveapimodels.EveService{
		ID:              s.ID,
		NamespaceID:     s.NamespaceID,
		NamespaceName:   s.NamespaceName,
		ArtifactID:      s.ArtifactID,
		ArtifactName:    s.ArtifactName,
		OverrideVersion: s.OverrideVersion,
		DeployedVersion: s.DeployedVersion,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
		Name:            s.Name,
		StickySessions:  s.StickySessions,
		Count:           s.Count,
	}
}

func metaDataServiceKey(service, namespace string) string {
	return fmt.Sprintf("eve-bot:%s:%s", service, namespace)
}

func resolveNamespace(ctx context.Context, api eveapi.Client, cmd commands.EvebotCommand) (eve.Namespace, error) {
	var nv eve.Namespace

	dynamicOpts := cmd.Options()

	// Gotta get the namespaces first, since we are working with the Alias, and not the Name/ID
	namespaces, err := api.GetNamespacesByEnvironment(ctx, dynamicOpts[params.EnvironmentName].(string))
	if err != nil {
		return nv, err
	}

	for _, v := range namespaces {
		if strings.ToLower(v.Alias) == strings.ToLower(dynamicOpts[params.NamespaceName].(string)) {
			nv = v
			break
		}
	}

	if nv.ID == 0 {
		return nv, fmt.Errorf("invalid namespace")
	}
	return nv, nil
}
