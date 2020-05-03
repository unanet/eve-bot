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
		input string
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
			args: args{input: validDeployCmd},
			want: botcommands.NewDeployCommand(strings.Fields(validDeployCmd)[1:]),
		},
		{
			name: "valid full deploy command",
			ebr:  &EvebotResolver{},
			args: args{input: validDeployCmdFull},
			want: botcommands.NewDeployCommand(strings.Fields(validDeployCmdFull)[1:]),
		},
		{
			name: "valid services deploy command",
			ebr:  &EvebotResolver{},
			args: args{input: validDeployCmdService},
			want: botcommands.NewDeployCommand(strings.Fields(validDeployCmdService)[1:]),
		},
		{
			name: "valid dryrun deploy command",
			ebr:  &EvebotResolver{},
			args: args{input: validDeployCmdDryrun},
			want: botcommands.NewDeployCommand(strings.Fields(validDeployCmdDryrun)[1:]),
		},
		{
			name: "valid force deploy command",
			ebr:  &EvebotResolver{},
			args: args{input: validDeployCmdForce},
			want: botcommands.NewDeployCommand(strings.Fields(validDeployCmdForce)[1:]),
		},
		{
			name: "invalid deploy command",
			ebr:  &EvebotResolver{},
			args: args{input: invalidDeployCmd},
			want: botcommands.NewInvalidCommand(strings.Fields(invalidDeployCmd)[1:]),
		},
		{
			name: "invalid deploy command length",
			ebr:  &EvebotResolver{},
			args: args{input: invalidDeployCmdLen},
			want: botcommands.NewDeployCommand(strings.Fields(invalidDeployCmdLen)[1:]),
		},
		{
			name: "invalid deploy command length 2",
			ebr:  &EvebotResolver{},
			args: args{input: invalidDeployCmdLen2},
			want: botcommands.NewDeployCommand(strings.Fields(invalidDeployCmdLen2)[1:]),
		},
		{
			name: "deploy command help",
			ebr:  &EvebotResolver{},
			args: args{input: deployCmdHelp},
			want: botcommands.NewDeployCommand(strings.Fields(deployCmdHelp)[1:]),
		},
		{
			name: "deploy command help 2",
			ebr:  &EvebotResolver{},
			args: args{input: deployCmdHelp2},
			want: botcommands.NewDeployCommand(strings.Fields(deployCmdHelp2)[1:]),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ebr := &EvebotResolver{}

			resolvedCmd := ebr.Resolve(tt.args.input)

			if !reflect.DeepEqual(resolvedCmd, tt.want) {
				t.Errorf("got = %v, want %v", resolvedCmd, tt.want)
			}
		})
	}
}
