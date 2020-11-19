package handlers

import (
	"context"
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve/pkg/eve"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
)

func deployHandler(
	ctx context.Context,
	eveAPIClient eveapi.Client,
	chatSvc chatservice.Provider,
	cmd commands.EvebotCommand,
	timestamp string,
	deployOpts eveapimodels.DeploymentPlanOptions) {

	resp, err := eveAPIClient.Deploy(ctx, deployOpts, cmd.Info().User, cmd.Info().Channel, timestamp)
	if err != nil && len(err.Error()) > 0 {
		chatSvc.DeploymentNotificationThread(ctx, err.Error(), cmd.Info().User, cmd.Info().Channel, timestamp)
		return
	}
	if resp == nil {
		chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, errInvalidAPIResp)
		return
	}
	if len(resp.Messages) > 0 {
		chatSvc.UserNotificationThread(ctx, strings.Join(resp.Messages, ","), cmd.Info().User, cmd.Info().Channel, timestamp)
	}
}

func resolveServiceNamespace(
	ctx context.Context,
	eveAPIClient eveapi.Client,
	chatSvc chatservice.Provider,
	cmd commands.EvebotCommand, ts *string) (*eve.Namespace, *eveapimodels.EveService) {

	var ns eve.Namespace
	var svc eveapimodels.EveService
	var svcs eveapimodels.Services
	var err error

	ns, err = resolveNamespace(ctx, eveAPIClient, cmd)
	if err != nil {
		chatSvc.UserNotificationThread(ctx, err.Error(), cmd.Info().User, cmd.Info().Channel, *ts)
		return &ns, &svc
	}
	svcs, err = eveAPIClient.GetServicesByNamespace(ctx, ns.Name)
	if err != nil {
		chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
		return &ns, &svc
	}
	if svcs == nil {
		chatSvc.UserNotificationThread(ctx, "no services", cmd.Info().User, cmd.Info().Channel, *ts)
		return &ns, &svc
	}
	var requestedSvcName string
	var valid bool
	if requestedSvcName, valid = cmd.Options()[params.ServiceName].(string); !valid {
		chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, fmt.Errorf("invalid ServiceName Param"))
		return &ns, &svc
	}
	for _, s := range svcs {
		if strings.ToLower(s.Name) == strings.ToLower(requestedSvcName) {
			svc = mapToEveService(s)
			break
		}
	}
	if svc.ID == 0 {
		chatSvc.UserNotificationThread(ctx, fmt.Sprintf("invalid requested service: %s", requestedSvcName), cmd.Info().User, cmd.Info().Channel, *ts)
		return &ns, &svc
	}
	svc, err = eveAPIClient.GetServiceByID(ctx, svc.ID)
	if err != nil {
		chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
		return &ns, &svc
	}
	return &ns, &svc
}
