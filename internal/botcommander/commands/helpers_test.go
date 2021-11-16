package commands

import (
	"reflect"
	"testing"

	"github.com/unanet/eve-bot/internal/botcommander/params"
)

func Test_Metadata(t *testing.T) {
	type args struct {
		input []string
	}
	tests := []struct {
		name string
		args args
		want params.MetadataMap
	}{
		{
			name: "test handling empty metadata",
			args: args{
				input: []string{},
			},
			want: nil,
		},
		{
			name: "test malformed metadata",
			args: args{
				input: []string{"FOO"},
			},
			want: params.MetadataMap{},
		},
		{
			name: "test handling single key value parameter",
			args: args{
				input: []string{
					"FOO=BAR",
				},
			},
			want: params.MetadataMap{
				"FOO": "BAR",
			},
		},
		{
			name: "test handling multiple key value parameters",
			args: args{
				input: []string{
					"FOO=BAR",
					"FIZZ=BUZZ",
				},
			},
			want: params.MetadataMap{
				"FOO":  "BAR",
				"FIZZ": "BUZZ",
			},
		},
		{
			name: "test handling params with spaces",
			args: args{
				input: []string{
					"FOO=",
					"BAR",
					"FIZZ=BUZZ",
				},
			},
			want: params.MetadataMap{
				"FOO":  "BAR",
				"FIZZ": "BUZZ",
			},
		},
		{
			name: "test handling multiple params with spaces",
			args: args{
				input: []string{
					"FOO=",
					"BAR",
					"FIZZ=BUZZ",
					"TACO=",
					"SHOP",
				},
			},
			want: params.MetadataMap{
				"FOO":  "BAR",
				"FIZZ": "BUZZ",
				"TACO": "SHOP",
			},
		},
		{
			name: "test handling params with space between params",
			args: args{
				input: []string{
					"FOO=BAR TACOS BURRITO",
					"FIZZ=BUZZ",
				},
			},
			want: params.MetadataMap{
				"FOO":  "BAR TACOS BURRITO",
				"FIZZ": "BUZZ",
			},
		},
		{
			name: "test handling params with multiple spaces between params at the beginning",
			args: args{
				input: []string{
					"FOO=BAR",
					"TACOS",
					"BURRITO",
					"FIZZ=BUZZ",
				},
			},
			want: params.MetadataMap{
				"FOO":  "BAR TACOS BURRITO",
				"FIZZ": "BUZZ",
			},
		},
		{
			name: "test handling params with multiple spaces between params at the end",
			args: args{
				input: []string{
					"FOO=BAR",
					"FIZZ=BUZZ",
					"TACOS",
					"BURRITO",
				},
			},
			want: params.MetadataMap{
				"FOO":  "BAR",
				"FIZZ": "BUZZ TACOS BURRITO",
			},
		},
		{
			name: "test handling params with hyphens",
			args: args{
				input: []string{
					"FOO=bar-fizz-buzz",
				},
			},
			want: params.MetadataMap{
				"FOO": "bar-fizz-buzz",
			},
		},
		{
			name: "test handling params with special characters",
			args: args{
				input: []string{
					"FOO=!@#$%^&*()_+{}[],.<>/?`~\\|",
				},
			},
			want: params.MetadataMap{
				"FOO": "!@#$%^&*()_+{}[],.<>/?`~\\|",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hydrateMetadataMap(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v\nwant %v", got, tt.want)
			}
		})
	}
}
