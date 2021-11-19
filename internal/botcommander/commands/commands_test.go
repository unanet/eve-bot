package commands

import (
	"fmt"
	"testing"
)

func Test_ValidInputLength(t *testing.T) {
	type args struct {
		input []string
		bounds InputLengthBounds
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test valid",
			args: args{
				input: []string{"1", "2", "3"},
				bounds: InputLengthBounds{
					Min: 2,
					Max: 3,
				},
			},
			want: true,
		},
		{
			name: "test exact bounds",
			args: args{
				input: []string{"1", "2"},
				bounds: InputLengthBounds{
					Min: 2,
					Max: 2,
				},
			},
			want: true,
		},
		{
			name: "test below min",
			args: args{
				input: []string{"1"},
				bounds: InputLengthBounds{
					Min: 2,
					Max: 3,
				},
			},
			want: false,
		},
		{
			name: "test above min",
			args: args{
				input: []string{"1", "2", "3", "4"},
				bounds: InputLengthBounds{
					Min: 2,
					Max: 3,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {

		bc := baseCommand{
			input:      tt.args.input,
			bounds: 	tt.args.bounds,
		}

		t.Run(tt.name, func(t *testing.T) {

			if got := bc.ValidInputLength(); got != tt.want {
				t.Errorf("got = %v\nwant %v", got, tt.want)
			}
		})
	}
}

func Test_Ack(t *testing.T) {
	type args struct {
		input []string
		bounds InputLengthBounds
		errs []error
		info ChatInfo
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test for successful ack",
			args: args {
				input: []string{"foo", "bar"},
				info: ChatInfo {
					CommandName: "foo",
				},
				bounds: InputLengthBounds{
					Max: 2,
					Min: 2,
				},
			},
			want: true,
		},
		{
			name: "test empty",
			args: args{
				input: []string{},
			},
			want: false,
		},
		{
			name: "test invalid command",
			args: args{
				input: []string{"test-ack"},
			},
			want: false,
		},
		{
			name: "test help",
			args: args{
				input: []string{"help"},
			},
			want: false,
		},
		{
			name: "test help from command",
			args: args{
				input: []string{"run", "help"},
			},
			want: false,
		},
		{
			name: "test invalid length",
			args: args{
				input: []string{"run", "help"},
				bounds: InputLengthBounds{
					Max: 1,
					Min: 1,
				},
			},
			want: false,
		},
		{
			name: "test handling errs",
			args: args {
				input: []string{"release", "namespace", "current", "dev-int", "from", "foo"},
				info: ChatInfo {
					CommandName: "foo",
				},
				errs: []error{
					fmt.Errorf("test err"),
				},
				bounds: InputLengthBounds{
					Max: 7,
					Min: 5,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {

		bc := baseCommand{
			input:      tt.args.input,
			bounds: 	tt.args.bounds,
			errs: 		tt.args.errs,
			info: 		tt.args.info,
		}

		t.Run(tt.name, func(t *testing.T) {
			if _, ack := bc.BaseAckMsg(""); ack != tt.want {
				t.Errorf("got = %v\nwant %v", ack, tt.want)
			}
		})
	}
}

func Test_CleanUrls(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test encoded input extended",
			args: args{input: "REPORTING_customer.url.regex=&lt;blah&gt;"},
			want: "REPORTING_customer.url.regex=<blah>",
		},
		{
			name: "test <blah> input",
			args: args{input: "<blah>"},
			want: "blah",
		},
		{
			name: "test encoded input simple",
			args: args{input: "&lt;blah&gt;"},
			want: "<blah>",
		},
		{
			name: "test regex url input",
			args: args{input: "'//(.+).<http://unanet.io/*|unanet.io/*>'"},
			want: "'//(.+).unanet.io/*'",
		},
		{
			name: "test regex url input 2",
			args: args{input: "\"//(.+).<http://unanet.io/*|unanet.io/*>\""},
			want: "\"//(.+).unanet.io/*\"",
		},
		{
			name: "test regex url input 3",
			args: args{input: "//(.+).<http://unanet.io/*|unanet.io/*>"},
			want: "//(.+).unanet.io/*",
		},
		{
			name: "test postgres url with equal 2",
			args: args{input: "jdbc:<postgresql://unanet-aurora-db.app-nonprod.unanet.io:5432/tsgc_qa_current?escapeSyntaxCallMode=callIfNoReturn>"},
			want: "jdbc:postgresql://unanet-aurora-db.app-nonprod.unanet.io:5432/tsgc_qa_current?escapeSyntaxCallMode=callIfNoReturn",
		},
		{
			name: "test postgres url with equal",
			args: args{input: "jdbc:postgresql://unanet-aurora-db.app-nonprod.unanet.io:5432/tsgc_qa_current?escapeSyntaxCallMode=callIfNoReturn"},
			want: "jdbc:postgresql://unanet-aurora-db.app-nonprod.unanet.io:5432/tsgc_qa_current?escapeSyntaxCallMode=callIfNoReturn",
		},
		{
			name: "test postgres url",
			args: args{input: "jdbc:<postgresql://unanet-aurora-db.app-nonprod.unanet.io:5432/aurorab_int_current?escapeSyntaxCallMode=callIfNoReturn>"},
			want: "jdbc:postgresql://unanet-aurora-db.app-nonprod.unanet.io:5432/aurorab_int_current?escapeSyntaxCallMode=callIfNoReturn",
		},
		{
			name: "single url- no pipe",
			args: args{input: "<https://unaneta.qa-latest.unanet.io/platform>"},
			want: "https://unaneta.qa-latest.unanet.io/platform",
		},
		{
			name: "single url",
			args: args{input: "https://unaneta.qa-latest.unanet.io/platform"},
			want: "https://unaneta.qa-latest.unanet.io/platform",
		},
		{
			name: "single linked url",
			args: args{input: "<https://unaneta.qa-latest.unanet.io/platform|unaneta.qa-latest.unanet.io/platform>"},
			want: "unaneta.qa-latest.unanet.io/platform",
		},
		{
			name: "full complex url parse",
			args: args{input: "troy_sampson_<ftp://wtfftp.com|wtfftp.com>_http://thisisclean.com_<https://h-ello.there.com|hellothere.com>are_we_there<https://hello.com|hello.com>something_else-goes__-here<ftp://asdf.com|asdf.com>_wtf_are_wedoing"},
			want: "troy_sampson_wtfftp.com_http://thisisclean.com_hellothere.comare_we_therehello.comsomething_else-goes__-hereasdf.com_wtf_are_wedoing",
		},
		{
			name: "simple single url parse",
			args: args{input: "<https://www.google.com|www.google.com>"},
			want: "www.google.com",
		},
		{
			name: "simple single url",
			args: args{input: "https://www.google.com"},
			want: "https://www.google.com",
		},
		{
			name: "simple single string",
			args: args{input: "wtf_are_we_doing"},
			want: "wtf_are_we_doing",
		},
		{
			name: "ftp url parse",
			args: args{input: "<ftp://somehost|somehost>"},
			want: "somehost",
		},
		{
			name: "https and http url parse",
			args: args{input: "<http://www.somehost.com|somehost.com>_<https://someotherhost.com|www.someotherhost.com>"},
			want: "somehost.com_www.someotherhost.com",
		},
		{
			name: "https and http and ftp url parse",
			args: args{input: "here_we_go<ftp://someftphost|someftphost>_<http://www.somehost.com|somehost.com>_<https://someotherhost.com|www.someotherhost.com>"},
			want: "here_we_gosomeftphost_somehost.com_www.someotherhost.com",
		},
		{
			name: "any url parse",
			args: args{input: "postgresql://unanet-aurora-db.app-nonprod.unanet.io:5432/auroraa_int_current_here_we_go<ftp://someftphost|someftphost>_<http://www.somehost.com|somehost.com>_<https://someotherhost.com|www.someotherhost.com>"},
			want: "postgresql://unanet-aurora-db.app-nonprod.unanet.io:5432/auroraa_int_current_here_we_gosomeftphost_somehost.com_www.someotherhost.com",
		},
		{
			name: "postgres url parse",
			args: args{input: "<postgresql://unanet-aurora-db.app-nonprod.unanet.io:5432/auroraa_int_current|unanet-aurora-db.app-nonprod.unanet.io:5432/auroraa_int_current>"},
			want: "unanet-aurora-db.app-nonprod.unanet.io:5432/auroraa_int_current",
		},
		{
			name: "db url parse",
			args: args{input: "jdbc:postgresql://unanet-aurora-db.app-nonprod.unanet.io:5432/aurorab_int_current?escapeSyntaxCallMode=callIfNoReturn"},
			want: "jdbc:postgresql://unanet-aurora-db.app-nonprod.unanet.io:5432/aurorab_int_current?escapeSyntaxCallMode=callIfNoReturn",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CleanUrls(tt.args.input); got != tt.want {
				t.Errorf("got = %v\nwant %v", got, tt.want)
			}
		})
	}
}
