package commands

import (
	"reflect"
	"testing"
)

func Test_ReleaseArtifact_resolveDynamicOptions(t *testing.T) {
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
					"release-artifact", "foo", "from", "foo",
				},
			},
			want: CommandOptions{
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
					"release-artifact", "foo:1.0.0", "from", "foo",
				},
			},
			want: CommandOptions{
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
					"release-artifact", "foo", "from", "foo", "to", "bar",
				},
			},
			want: CommandOptions{
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
					"release-artifact", "foo:1.0.0", "from", "foo", "to", "bar",
				},
			},
			want: CommandOptions{
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
				input: []string{},
			},
			want: CommandOptions{},
		},
		{
			name: "test releasing an artifact with invalid to feed",
			args: args{
				input: []string{
					"release-artifact", "foo:1.0.0", "from", "foo", "to",
				},
			},
			want: CommandOptions{
				"artifact": "foo",
				"version": "1.0.0",
				"from_feed": "foo",
			},
		},
		{
			name: "test releasing an artifact without feeds",
			args: args{
				input: []string{
					"release-artifact", "foo:1.0.0",
				},
			},
			want: CommandOptions{},
		},
		{
			name: "test releasing an artifact without from feed",
			args: args{
				input: []string{
					"release-artifact", "foo:1.0.0", "from",
				},
			},
			want: CommandOptions{},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			cmd := NewReleaseArtifactCommand(tt.args.input, "", "")

			if got := cmd.Options(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v\nwant %v", got, tt.want)
			}
		})
	}
}
