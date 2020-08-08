package args

import (
	"reflect"
	"testing"
)

var fakeDryrunTrue = Dryrun(true)

var fakeDryrunFalse = Dryrun(false)

func TestDryrun_Name(t *testing.T) {
	tests := []struct {
		name string
		a    Dryrun
		want string
	}{
		{
			name: "happy path - true dryrun",
			a:    fakeDryrunTrue,
			want: DryrunName,
		},
		{
			name: "happy path - false dryrun",
			a:    fakeDryrunFalse,
			want: DryrunName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Name(); got != tt.want {
				t.Errorf("Dryrun.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDryrun_Value(t *testing.T) {
	tests := []struct {
		name string
		a    Dryrun
		want interface{}
	}{
		{
			name: "happy path - true dryrun",
			a:    fakeDryrunTrue,
			want: true,
		},
		{
			name: "happy path - false dryrun",
			a:    fakeDryrunFalse,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dryrun.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDryrun_Description(t *testing.T) {
	tests := []struct {
		name string
		a    Dryrun
		want string
	}{
		{
			name: "happy path - true dryrun",
			a:    fakeDryrunTrue,
			want: DryrunDescription,
		},
		{
			name: "happy path - false dryrun",
			a:    fakeDryrunFalse,
			want: DryrunDescription,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Description(); got != tt.want {
				t.Errorf("Dryrun.Description() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultDryrunArg(t *testing.T) {
	tests := []struct {
		name string
		want Dryrun
	}{
		{
			name: "happy path - true dryrun",
			want: DefaultDryrunArg(),
		},
		{
			name: "happy path - false dryrun",
			want: DefaultDryrunArg(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultDryrunArg(); got != tt.want {
				t.Errorf("DefaultDryrunArg() = %v, want %v", got, tt.want)
			}
		})
	}
}
