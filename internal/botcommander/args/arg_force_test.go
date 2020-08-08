package args

import (
	"reflect"
	"testing"
)

var fakeForceTrue = Force(true)

var fakeForceFalse = Force(false)

func TestForce_Name(t *testing.T) {
	tests := []struct {
		name string
		a    Force
		want string
	}{
		{
			name: "happy path - true force",
			a:    fakeForceTrue,
			want: ForceDeployName,
		},
		{
			name: "happy path - false force",
			a:    fakeForceFalse,
			want: ForceDeployName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Name(); got != tt.want {
				t.Errorf("Force.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForce_Value(t *testing.T) {
	tests := []struct {
		name string
		a    Force
		want interface{}
	}{
		{
			name: "happy path - true force",
			a:    fakeForceTrue,
			want: true,
		},
		{
			name: "happy path - false force",
			a:    fakeForceFalse,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Force.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForce_Description(t *testing.T) {
	tests := []struct {
		name string
		a    Force
		want string
	}{
		{
			name: "happy path - true force",
			a:    fakeForceTrue,
			want: ForceDeployDescription,
		},
		{
			name: "happy path - false force",
			a:    fakeForceFalse,
			want: ForceDeployDescription,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Description(); got != tt.want {
				t.Errorf("Force.Description() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultForceArg(t *testing.T) {
	tests := []struct {
		name string
		want Force
	}{
		{
			name: "happy path - true force",
			want: DefaultForceArg(),
		},
		{
			name: "happy path - false force",
			want: DefaultForceArg(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultForceArg(); got != tt.want {
				t.Errorf("DefaultForceArg() = %v, want %v", got, tt.want)
			}
		})
	}
}
