package commands

import (
	"reflect"
	"testing"
)

func Test_Release_resolveDynamicOptions(t *testing.T) {
	type args struct {
		input []string
	}
	tests := []struct {
		name string
		args args
		want CommandOptions
	}{
		{
			name: "test releasing an artifact from feed",
			args: args{
				input: []string{
					"release", "artifact", "foo", "from", "foo",
				},
			},
			want: CommandOptions{
				"commandType": "artifact",
				"artifact": "foo",
				"version": "",
				"from_feed": "foo",
				"to_feed": "",
			},
		},
		{
			name: "test releasing an artifact with version from feed",
			args: args{
				input: []string{
					"release", "artifact", "foo:1.0.0", "from", "foo",
				},
			},
			want: CommandOptions{
				"commandType": "artifact",
				"artifact": "foo",
				"version": "1.0.0",
				"from_feed": "foo",
				"to_feed": "",
			},
		},
		{
			name: "test releasing an artifact from feed to feed",
			args: args{
				input: []string{
					"release", "artifact", "foo", "from", "foo", "to", "bar",
				},
			},
			want: CommandOptions{
				"commandType": "artifact",
				"artifact": "foo",
				"version": "",
				"from_feed": "foo",
				"to_feed": "bar",
			},
		},
		{
			name: "test releasing an artifact with version from feed to feed",
			args: args{
				input: []string{
					"release", "artifact", "foo:1.0.0", "from", "foo", "to", "bar",
				},
			},
			want: CommandOptions{
				"commandType": "artifact",
				"artifact": "foo",
				"version": "1.0.0",
				"from_feed": "foo",
				"to_feed": "bar",
			},
		},
		// Error handling
		{
			name: "test empty",
			args: args{
				input: []string{"release"},
			},
			want: CommandOptions{},
		},
		{
			name: "test releasing an artifact with invalid to feed",
			args: args{
				input: []string{
					"release", "artifact", "foo:1.0.0", "from", "foo", "to",
				},
			},
			want: CommandOptions{
				"commandType": "artifact",
				"artifact": "foo",
				"version": "1.0.0",
				"from_feed": "foo",
			},
		},
		{
			name: "test releasing an artifact without feeds",
			args: args{
				input: []string{
					"release", "artifact", "foo:1.0.0",
				},
			},
			want: CommandOptions{},
		},
		{
			name: "test releasing an artifact without from feed",
			args: args{
				input: []string{
					"release", "artifact", "foo:1.0.0", "from",
				},
			},
			want: CommandOptions{},
		},

		// Namespace
		{
			name: "test releasing a namespace from feed",
			args: args{
				input: []string{
					"release", "namespace", "current", "dev-int", "from", "foo",
				},
			},
			want: CommandOptions{
				"commandType": "namespace",
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
					"release", "namespace", "current", "dev-int", "from", "foo", "to", "bar",
				},
			},
			want: CommandOptions{
				"commandType": "namespace",
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
				input: []string{"release"},
			},
			want: CommandOptions{},
		},
		{
			name: "test releasing a namespace without an environment",
			args: args{
				input: []string{
					"release", "namespace", "current",
				},
			},
			want: CommandOptions{},
		},
		{
			name: "test releasing a namespace without a feeds",
			args: args{
				input: []string{
					"release", "namespace", "current", "dev-int",
				},
			},
			want: CommandOptions{},
		},
		{
			name: "test releasing a namespace without from feed",
			args: args{
				input: []string{
					"release", "namespace", "current", "dev-int", "from",
				},
			},
			want: CommandOptions{},
		},
		{
			name: "test releasing a namespace with missing to feed name",
			args: args{
				input: []string{
					"release", "namespace", "current", "dev-int", "from", "foo", "to",
				},
			},
			want: CommandOptions{
				"commandType": "namespace",
				"namespace": "current",
				"environment": "dev-int",
				"from_feed": "foo",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			cmd := NewReleaseCommand(tt.args.input, "", "")

			if got := cmd.Options(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v\nwant %v", got, tt.want)
			}
		})
	}
}
