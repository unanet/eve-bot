package botcommander

import (
	"reflect"
	"strings"
	"testing"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botcommands"
)

func TestNewResolver(t *testing.T) {
	tests := []struct {
		name string
		want Resolver
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResolver(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResolver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvebotResolver_Resolve(t *testing.T) {

	user := "dummy"
	channel := "chanelID"
	ts := "2342341233.221"
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

	type args struct {
		input, channel, user, ts string
	}
	tests := []struct {
		name string
		ebr  *EvebotResolver
		args args
		want botcommands.EvebotCommand
	}{
		{
			name: "valid basic deploy command",
			ebr:  &EvebotResolver{},
			args: args{input: validDeployCmd, channel: channel, user: user, ts: ts},
			want: botcommands.NewDeployCommand(strings.Fields(validDeployCmd)[1:], channel, user, ts),
		},
		{
			name: "valid full deploy command",
			ebr:  &EvebotResolver{},
			args: args{input: validDeployCmdFull, channel: channel, user: user, ts: ts},
			want: botcommands.NewDeployCommand(strings.Fields(validDeployCmdFull)[1:], channel, user, ts),
		},
		{
			name: "valid services deploy command",
			ebr:  &EvebotResolver{},
			args: args{input: validDeployCmdService, channel: channel, user: user, ts: ts},
			want: botcommands.NewDeployCommand(strings.Fields(validDeployCmdService)[1:], channel, user, ts),
		},
		{
			name: "valid dryrun deploy command",
			ebr:  &EvebotResolver{},
			args: args{input: validDeployCmdDryrun, channel: channel, user: user, ts: ts},
			want: botcommands.NewDeployCommand(strings.Fields(validDeployCmdDryrun)[1:], channel, user, ts),
		},
		{
			name: "valid force deploy command",
			ebr:  &EvebotResolver{},
			args: args{input: validDeployCmdForce, channel: channel, user: user, ts: ts},
			want: botcommands.NewDeployCommand(strings.Fields(validDeployCmdForce)[1:], channel, user, ts),
		},
		{
			name: "invalid deploy command",
			ebr:  &EvebotResolver{},
			args: args{input: invalidDeployCmd, channel: channel, user: user, ts: ts},
			want: botcommands.NewInvalidCommand(strings.Fields(invalidDeployCmd)[1:], channel, user, ts),
		},
		{
			name: "invalid deploy command length",
			ebr:  &EvebotResolver{},
			args: args{input: invalidDeployCmdLen, channel: channel, user: user, ts: ts},
			want: botcommands.NewDeployCommand(strings.Fields(invalidDeployCmdLen)[1:], channel, user, ts),
		},
		{
			name: "invalid deploy command length 2",
			ebr:  &EvebotResolver{},
			args: args{input: invalidDeployCmdLen2, channel: channel, user: user, ts: ts},
			want: botcommands.NewDeployCommand(strings.Fields(invalidDeployCmdLen2)[1:], channel, user, ts),
		},
		{
			name: "deploy command help",
			ebr:  &EvebotResolver{},
			args: args{input: deployCmdHelp, channel: channel, user: user, ts: ts},
			want: botcommands.NewDeployCommand(strings.Fields(deployCmdHelp)[1:], channel, user, ts),
		},
		{
			name: "deploy command help 2",
			ebr:  &EvebotResolver{},
			args: args{input: deployCmdHelp2, channel: channel, user: user, ts: ts},
			want: botcommands.NewDeployCommand(strings.Fields(deployCmdHelp2)[1:], channel, user, ts),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ebr := &EvebotResolver{}

			resolvedCmd := ebr.Resolve(tt.args.input, tt.args.channel, tt.args.user, tt.args.ts)

			if !reflect.DeepEqual(resolvedCmd, tt.want) {
				t.Errorf("got = %v, want %v", resolvedCmd, tt.want)
			}
		})
	}
}
