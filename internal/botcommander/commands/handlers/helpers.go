package handlers

import (
	"context"
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

func deployHandler(
	ctx context.Context,
	eveAPIClient eveapi.Client,
	chatSvc chatservice.Provider,
	cmd commands.EvebotCommand,
	timestamp string,
	deployOpts eve.DeploymentPlanOptions) {

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
	cmd commands.EvebotCommand, ts *string) (*eve.Namespace, *eve.Service) {

	var ns eve.Namespace
	var svc eve.Service
	var svcs []eve.Service
	var err error

	ns, err = resolveNamespace(ctx, eveAPIClient, cmd)
	if err != nil {
		chatSvc.UserNotificationThread(ctx, err.Error(), cmd.Info().User, cmd.Info().Channel, *ts)
		return nil, nil
	}
	svcs, err = eveAPIClient.GetServicesByNamespace(ctx, ns.Name)
	if err != nil {
		chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
		return nil, nil
	}
	if svcs == nil {
		chatSvc.UserNotificationThread(ctx, "no services", cmd.Info().User, cmd.Info().Channel, *ts)
		return nil, nil
	}
	var requestedSvcName string
	var valid bool
	if requestedSvcName, valid = cmd.Options()[params.ServiceName].(string); !valid {
		chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, fmt.Errorf("invalid ServiceName Param"))
		return nil, nil
	}
	for _, s := range svcs {
		if strings.ToLower(s.Name) == strings.ToLower(requestedSvcName) {
			svc = s
			break
		}
	}
	if svc.ID == 0 {
		chatSvc.UserNotificationThread(ctx, fmt.Sprintf("invalid requested service: %s", requestedSvcName), cmd.Info().User, cmd.Info().Channel, *ts)
		return nil, nil
	}
	return &ns, &svc
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
