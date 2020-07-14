package resolver

import (
	"reflect"
	"strings"
	"testing"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
)

func TestEvebotResolver_Resolve(t *testing.T) {

	user := "dummy"
	channel := "chanelID"
	validDeployCmd := "@evebot deploy current to int"
	validDeployCmdFull := validDeployCmd + " " + "services=infocus-cloud-client:2020.2.232,infocus-proxy:2020.2.199 dryrun=true force=true"
	validDeployCmdService := validDeployCmd + " " + "services=infocus-cloud-client,infocus-proxy"
	validDeployCmdDryrun := validDeployCmd + " " + "dryrun=true"
	validDeployCmdForce := validDeployCmd + " " + "force=true"
	invalidDeployCmd := "@evebot deployy current to int"
	invalidDeployCmdLen := "@evebot deploy current to"
	invalidDeployCmdLen2 := "@evebot deploy current"
	deployCmdHelp := "@evebot deploy"
	deployCmdHelp2 := "@evebot deploy help"
	setMetaDataCmd := "@evebot set metadata for unaneta in current una-int una-qa unanet_unanet_unanet.org_access.dataManager.people.default_to_all=false"
	setMetaDataCmdInput := "@evebot set metadata for unaneta in current una-int una-qa unanet_unanet_<https://unanet.org|unanet.org>_access.dataManager.people.default_to_all=false"
	setMetaDataCmdCleanInput := "@evebot set metadata for unaneta in current una-int una-qa unanet_unanet_unanet.org_access.dataManager.people.default_to_all=false"

	type args struct {
		input, channel, user string
	}
	tests := []struct {
		name string
		ebr  *EvebotResolver
		args args
		want commands.EvebotCommand
	}{
		{
			name: "valid set metadata clean urls",
			ebr:  &EvebotResolver{},
			args: args{input: setMetaDataCmdInput, channel: channel, user: user},
			want: commands.NewSetCommand(strings.Fields(setMetaDataCmdCleanInput)[1:], channel, user),
		},
		{
			name: "valid set metadata",
			ebr:  &EvebotResolver{},
			args: args{input: setMetaDataCmd, channel: channel, user: user},
			want: commands.NewSetCommand(strings.Fields(setMetaDataCmd)[1:], channel, user),
		},
		{
			name: "valid basic deploy command",
			ebr:  &EvebotResolver{},
			args: args{input: validDeployCmd, channel: channel, user: user},
			want: commands.NewDeployCommand(strings.Fields(validDeployCmd)[1:], channel, user),
		},
		{
			name: "valid full deploy command",
			ebr:  &EvebotResolver{},
			args: args{input: validDeployCmdFull, channel: channel, user: user},
			want: commands.NewDeployCommand(strings.Fields(validDeployCmdFull)[1:], channel, user),
		},
		{
			name: "valid services deploy command",
			ebr:  &EvebotResolver{},
			args: args{input: validDeployCmdService, channel: channel, user: user},
			want: commands.NewDeployCommand(strings.Fields(validDeployCmdService)[1:], channel, user),
		},
		{
			name: "valid dryrun deploy command",
			ebr:  &EvebotResolver{},
			args: args{input: validDeployCmdDryrun, channel: channel, user: user},
			want: commands.NewDeployCommand(strings.Fields(validDeployCmdDryrun)[1:], channel, user),
		},
		{
			name: "valid force deploy command",
			ebr:  &EvebotResolver{},
			args: args{input: validDeployCmdForce, channel: channel, user: user},
			want: commands.NewDeployCommand(strings.Fields(validDeployCmdForce)[1:], channel, user),
		},
		{
			name: "invalid deploy command",
			ebr:  &EvebotResolver{},
			args: args{input: invalidDeployCmd, channel: channel, user: user},
			want: commands.NewInvalidCommand(strings.Fields(invalidDeployCmd)[1:], channel, user),
		},
		{
			name: "invalid deploy command length",
			ebr:  &EvebotResolver{},
			args: args{input: invalidDeployCmdLen, channel: channel, user: user},
			want: commands.NewDeployCommand(strings.Fields(invalidDeployCmdLen)[1:], channel, user),
		},
		{
			name: "invalid deploy command length 2",
			ebr:  &EvebotResolver{},
			args: args{input: invalidDeployCmdLen2, channel: channel, user: user},
			want: commands.NewDeployCommand(strings.Fields(invalidDeployCmdLen2)[1:], channel, user),
		},
		{
			name: "deploy command help",
			ebr:  &EvebotResolver{},
			args: args{input: deployCmdHelp, channel: channel, user: user},
			want: commands.NewDeployCommand(strings.Fields(deployCmdHelp)[1:], channel, user),
		},
		{
			name: "deploy command help 2",
			ebr:  &EvebotResolver{},
			args: args{input: deployCmdHelp2, channel: channel, user: user},
			want: commands.NewDeployCommand(strings.Fields(deployCmdHelp2)[1:], channel, user),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ebr := &EvebotResolver{}

			resolvedCmd := ebr.Resolve(tt.args.input, tt.args.channel, tt.args.user)

			if !reflect.DeepEqual(resolvedCmd, tt.want) {
				t.Errorf("got = %v, want %v", resolvedCmd, tt.want)
			}
		})
	}
}
