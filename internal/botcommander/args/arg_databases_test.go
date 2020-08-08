package args

import (
	"fmt"
	"reflect"
	"testing"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"
)

var db = Database{
	Name:    "fake-db",
	Version: "0.1",
}

var db2 = Database{
	Name:    "fake-db-2",
	Version: "0.1.1",
}

func TestDatabases_Name(t *testing.T) {
	tests := []struct {
		name string
		dbs  Databases
		want string
	}{
		{
			name: "happy path",
			dbs:  Databases{db},
			want: DatabasesName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dbs.Name(); got != tt.want {
				t.Errorf("Databases.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabases_Description(t *testing.T) {
	tests := []struct {
		name string
		dbs  Databases
		want string
	}{
		{
			name: "happy path",
			dbs:  Databases{db},
			want: DatabaseDescription,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dbs.Description(); got != tt.want {
				t.Errorf("Databases.Description() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabases_Value(t *testing.T) {
	tests := []struct {
		name string
		dbs  Databases
		want interface{}
	}{
		{
			name: "happy path",
			dbs:  Databases{db},
			want: eveapimodels.ArtifactDefinitions{&eveapimodels.ArtifactDefinition{
				Name:             db.Name,
				RequestedVersion: db.Version,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dbs.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Databases.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultDatabasesArg(t *testing.T) {
	tests := []struct {
		name string
		want Databases
	}{
		{
			name: "happy path",
			want: Databases{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultDatabasesArg(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultDatabasesArg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDatabasesArg(t *testing.T) {
	type args struct {
		input []string
	}
	tests := []struct {
		name string
		args args
		want Databases
	}{
		{
			name: "happy path",
			args: args{input: []string{db.Name}},
			want: Databases{Database{Name: db.Name}},
		},
		{
			name: "happy path name version",
			args: args{input: []string{fmt.Sprintf("%s:%s", db.Name, db.Version), fmt.Sprintf("%s:%s", db2.Name, db2.Version)}},
			want: Databases{Database{Name: db.Name, Version: db.Version}, Database{Name: db2.Name, Version: db2.Version}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDatabasesArg(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDatabasesArg() = %v, want %v", got, tt.want)
			}
		})
	}
}
