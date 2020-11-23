package commands

import "testing"

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
			name: "test regex input",
			args: args{input: "'//(.+).<http://unanet.io/*|unanet.io/*>'"},
			want: "'//(.+).unanet.io/*'",
		},
		{
			name: "test regex input",
			args: args{input: "\"//(.+).<http://unanet.io/*|unanet.io/*>\""},
			want: "\"//(.+).unanet.io/*\"",
		},
		{
			name: "test regex input",
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
