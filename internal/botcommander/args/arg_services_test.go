package args

import (
	"reflect"
	"testing"

	"github.com/unanet/eve/pkg/eve"
)

var fakeServices = Services{Service{Name: "fake-service", Version: "1.0.0"}}

func TestServices_Name(t *testing.T) {
	tests := []struct {
		name string
		svcs Services
		want string
	}{
		{
			name: "happy path",
			svcs: fakeServices,
			want: ServicesName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.svcs.Name(); got != tt.want {
				t.Errorf("Services.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServices_Description(t *testing.T) {
	tests := []struct {
		name string
		svcs Services
		want string
	}{
		{
			name: "happy path",
			svcs: fakeServices,
			want: ServicesDescription,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.svcs.Description(); got != tt.want {
				t.Errorf("Services.Description() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServices_Value(t *testing.T) {
	tests := []struct {
		name string
		svcs Services
		want interface{}
	}{
		{
			name: "happy path",
			svcs: fakeServices,
			want: eve.ArtifactDefinitions{&eve.ArtifactDefinition{Name: fakeServices[0].Name, RequestedVersion: fakeServices[0].Version}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.svcs.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Services.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultServicesArg(t *testing.T) {
	tests := []struct {
		name string
		want Services
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultServicesArg(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultServicesArg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewServicesArg(t *testing.T) {
	type args struct {
		input []string
	}
	tests := []struct {
		name string
		args args
		want Services
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServicesArg(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServicesArg() = %v, want %v", got, tt.want)
			}
		})
	}
}
