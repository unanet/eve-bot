package resolver

import (
	"reflect"
	"strings"
	"testing"

	"github.com/unanet/eve-bot/internal/botcommander/commands"
)

func TestEvebotResolver_Resolve(t *testing.T) {

	resolver := New(commands.NewFactory())
	user := "dummy"
	channel := "chanelID"
	validDeployCmd := "@evebot deploy current to int"
	validDeployCmdFull := validDeployCmd + " services=infocus-cloud-client:2020.2.232,infocus-proxy:2020.2.199 dryrun=true force=true"
	validDeployCmdService := validDeployCmd + " services=infocus-cloud-client,infocus-proxy"
	validDeployCmdDryrun := validDeployCmd + " dryrun=true"
	validDeployCmdForce := validDeployCmd + " force=true"
	invalidDeployCmd := "@evebot deployy current to int"
	invalidDeployCmdLen := "@evebot deploy current to"
	invalidDeployCmdLen2 := "@evebot deploy current"
	deployCmdHelp := "@evebot deploy"
	deployCmdHelp2 := "@evebot deploy help"
	setMetaDataCmd := "@evebot set metadata for unaneta in current una-int una-qa unanet_unanet_unanet.org_access.dataManager.people.default_to_all=false"
	setMetaDataCmdInput := "@evebot set metadata for unaneta in current una-int una-qa unanet_unanet_<https://unanet.org|unanet.org>_access.dataManager.people.default_to_all=false"
	setMetaDataCmdCleanInput := "@evebot set metadata for unaneta in current una-int una-qa unanet_unanet_unanet.org_access.dataManager.people.default_to_all=false"
	setMetaDataCmd2Input := "@evebot set metadata for unaneta in latest una-qa unanet_unanet_unanet.external.platform.url=<https://unaneta.qa-latest.unanet.io/platform|unaneta.qa-latest.unanet.io/platform>"
	setMetaDataCmd2CleanInput := "@evebot set metadata for unaneta in latest una-qa unanet_unanet_unanet.external.platform.url=unaneta.qa-latest.unanet.io/platform"
	deleteMetaDataCmd := "@evebot delete metadata for unaneta in current una-int key key2 key3"
	invalidCmd := "@evebot wtf does this do"
	rootCmd := "@evebot"
	setMetaDataCmdDbURL := "@evebot set metadata for auroraa in current una-int unanet_database_unatime.database.url=jdbc:postgresql://unanet-aurora-db.app-nonprod.unanet.io:5432/aurorab_int_current?escapeSyntaxCallMode=callIfNoReturn"
	setMetaDataCmdEncoded := "@evebot set metadata for platform in current una-dev REPORTING_customer.url.regex=&lt;blah&gt;"

	type args struct {
		input, channel, user string
	}
	tests := []struct {
		name string
		args args
		want commands.EvebotCommand
	}{
		{
			name: "set encoded url",
			args: args{input: setMetaDataCmdEncoded, channel: channel, user: user},
			want: commands.NewSetCommand(strings.Fields("set metadata for platform in current una-dev REPORTING_customer.url.regex=<blah>"), channel, user),
		},
		{
			name: "set metadata db url",
			args: args{input: setMetaDataCmdDbURL, channel: channel, user: user},
			want: commands.NewSetCommand(strings.Fields(setMetaDataCmdDbURL)[1:], channel, user),
		},
		{
			name: "delete command",
			args: args{input: deleteMetaDataCmd, channel: channel, user: user},
			want: commands.NewDeleteCommand(strings.Fields(deleteMetaDataCmd)[1:], channel, user),
		},
		{
			name: "root command",
			args: args{input: rootCmd, channel: channel, user: user},
			want: commands.NewRootCmd([]string{""}, channel, user),
		},
		{
			name: "invalid command",
			args: args{input: invalidCmd, channel: channel, user: user},
			want: commands.NewInvalidCommand(strings.Fields(invalidCmd)[1:], channel, user),
		},
		{
			name: "set metadata clean url",
			args: args{input: setMetaDataCmd2Input, channel: channel, user: user},
			want: commands.NewSetCommand(strings.Fields(setMetaDataCmd2CleanInput)[1:], channel, user),
		},
		{
			name: "valid set metadata clean urls",
			args: args{input: setMetaDataCmdInput, channel: channel, user: user},
			want: commands.NewSetCommand(strings.Fields(setMetaDataCmdCleanInput)[1:], channel, user),
		},
		{
			name: "valid set metadata",
			args: args{input: setMetaDataCmd, channel: channel, user: user},
			want: commands.NewSetCommand(strings.Fields(setMetaDataCmd)[1:], channel, user),
		},
		{
			name: "valid basic deploy command",
			args: args{input: validDeployCmd, channel: channel, user: user},
			want: commands.NewDeployCommand(strings.Fields(validDeployCmd)[1:], channel, user),
		},
		{
			name: "valid full deploy command",
			args: args{input: validDeployCmdFull, channel: channel, user: user},
			want: commands.NewDeployCommand(strings.Fields(validDeployCmdFull)[1:], channel, user),
		},
		{
			name: "valid services deploy command",
			args: args{input: validDeployCmdService, channel: channel, user: user},
			want: commands.NewDeployCommand(strings.Fields(validDeployCmdService)[1:], channel, user),
		},
		{
			name: "valid dryrun deploy command",
			args: args{input: validDeployCmdDryrun, channel: channel, user: user},
			want: commands.NewDeployCommand(strings.Fields(validDeployCmdDryrun)[1:], channel, user),
		},
		{
			name: "valid force deploy command",
			args: args{input: validDeployCmdForce, channel: channel, user: user},
			want: commands.NewDeployCommand(strings.Fields(validDeployCmdForce)[1:], channel, user),
		},
		{
			name: "invalid deploy command",
			args: args{input: invalidDeployCmd, channel: channel, user: user},
			want: commands.NewInvalidCommand(strings.Fields(invalidDeployCmd)[1:], channel, user),
		},
		{
			name: "invalid deploy command length",
			args: args{input: invalidDeployCmdLen, channel: channel, user: user},
			want: commands.NewDeployCommand(strings.Fields(invalidDeployCmdLen)[1:], channel, user),
		},
		{
			name: "invalid deploy command length 2",
			args: args{input: invalidDeployCmdLen2, channel: channel, user: user},
			want: commands.NewDeployCommand(strings.Fields(invalidDeployCmdLen2)[1:], channel, user),
		},
		{
			name: "deploy command help",
			args: args{input: deployCmdHelp, channel: channel, user: user},
			want: commands.NewDeployCommand(strings.Fields(deployCmdHelp)[1:], channel, user),
		},
		{
			name: "deploy command help 2",
			args: args{input: deployCmdHelp2, channel: channel, user: user},
			want: commands.NewDeployCommand(strings.Fields(deployCmdHelp2)[1:], channel, user),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cmd := resolver.Resolve(tt.args.input, tt.args.channel, tt.args.user)

			if !reflect.DeepEqual(cmd, tt.want) {
				t.Errorf("\nA = %v\nB = %v", cmd, tt.want)
			}
		})
	}
}
