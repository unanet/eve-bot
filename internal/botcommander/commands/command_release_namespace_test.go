package commands

import (
	"reflect"
	"testing"
)

func Test_ReleaseNamespace_resolveDynamicOptions(t *testing.T) {
	type args struct {
		input []string
	}
	tests := []struct {
		name string
		args args
		want CommandOptions
	}{
		{
			name: "test releasing a namespace from feed",
			args: args{
				input: []string{
					"release-namespace", "current", "dev-int", "from", "foo",
				},
			},
			want: CommandOptions{
				"namespace": "current",
				"environment": "dev-int",
				"from_feed": "foo",
				"to_feed": "",
			},
		},
		{
			name: "test releasing a namespace from feed to feed",
			args: args{
				input: []string{
					"release-namespace", "current", "dev-int", "from", "foo", "to", "bar",
				},
			},
			want: CommandOptions{
				"namespace": "current",
				"environment": "dev-int",
				"from_feed": "foo",
				"to_feed": "bar",
			},
		},
		// Error handling
		{
			name: "test empty",
			args: args{
				input: []string{},
			},
			want: CommandOptions{},
		},
		{
			name: "test releasing an artifact without an environment",
			args: args{
				input: []string{
					"release-namespace", "current",
				},
			},
			want: CommandOptions{},
		},
		{
			name: "test releasing an artifact without a feeds",
			args: args{
				input: []string{
					"release-namespace", "current", "dev-int",
				},
			},
			want: CommandOptions{},
		},
		{
			name: "test releasing an artifact without from feed",
			args: args{
				input: []string{
					"release-namespace", "current", "dev-int", "from",
				},
			},
			want: CommandOptions{},
		},
		{
			name: "test releasing an artifact with missing to feed name",
			args: args{
				input: []string{
					"release-namespace", "current", "dev-int", "from", "foo", "to",
				},
			},
			want: CommandOptions{
				"namespace": "current",
				"environment": "dev-int",
				"from_feed": "foo",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			cmd := NewReleaseNamespaceCommand(tt.args.input, "", "")

			if got := cmd.Options(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v\nwant %v", got, tt.want)
			}
		})
	}
}
