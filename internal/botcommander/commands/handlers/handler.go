package handlers

import (
	"context"
	"errors"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"
	"gitlab.unanet.io/devops/eve/pkg/eve"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
)

var (
	errInvalidApiResp = errors.New("invalid api response")
)

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
		Metadata:        s.Metadata,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
		Name:            s.Name,
		StickySessions:  s.StickySessions,
		Count:           s.Count,
	}
}
