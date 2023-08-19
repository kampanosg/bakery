package runner

import (
	"errors"
	"testing"

	"github.com/kampanosg/bakery/internal/models"
)

type testCommandAgent struct {
	executorHandler func(cmd string) error
}

func (e *testCommandAgent) Run(cmd string) error {
	return e.executorHandler(cmd)
}

func TestRunner_RunCommand(t *testing.T) {
	type args struct {
		b    *models.Bakery
		args []string
	}
	tests := []struct {
		name     string
		executor *testCommandAgent
		fields   args
		wantErr  bool
	}{
		{
			name: "error nil bakery",
			fields: args{
				args: []string{"list"},
			},
			wantErr: true,
		},
		{
			name: "error undefined recipe",
			fields: args{
				b: &models.Bakery{
					Recipes: map[string]models.Recipe{
						"build": {},
					},
				},
				args: []string{"list"},
			},
			wantErr: true,
		},
		{
			name: "error thrown by executor",
			fields: args{
				b: &models.Bakery{
					Recipes: map[string]models.Recipe{
						"build": {
							Steps: []string{"go bild *.go"},
						},
					},
				},
				args: []string{"build"},
			},
			executor: &testCommandAgent{
				executorHandler: func(cmd string) error {
					return errors.New("invalid command bild")
				},
			},
			wantErr: true,
		},
		{
			name: "error when default doesnt exist",
			fields: args{
				b: &models.Bakery{
					Defaults: []string{"run"},
					Recipes: map[string]models.Recipe{
						"build": {
							Steps: []string{"go build -o app ./..."},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error when default step fails",
			fields: args{
				b: &models.Bakery{
					Defaults: []string{"build"},
					Recipes: map[string]models.Recipe{
						"build": {
							Steps: []string{"go bild -o app ./..."},
						},
					},
				},
			},
			executor: &testCommandAgent{
				executorHandler: func(cmd string) error {
					return errors.New("unable to run command")
				},
			},
			wantErr: true,
		},
		{
			name: "success when no defaults",
			fields: args{
				b: &models.Bakery{
					Defaults: []string{},
					Recipes: map[string]models.Recipe{
						"build": {
							Steps: []string{"go build -o app ./..."},
						},
					},
				},
			},
			executor: &testCommandAgent{
				executorHandler: func(cmd string) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "success defaults",
			fields: args{
				b: &models.Bakery{
					Defaults: []string{"build"},
					Recipes: map[string]models.Recipe{
						"build": {
							Steps: []string{"go build -o app ./..."},
						},
					},
				},
			},
			executor: &testCommandAgent{
				executorHandler: func(cmd string) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "success",
			fields: args{
				b: &models.Bakery{
					Recipes: map[string]models.Recipe{
						"build": {
							Steps: []string{"go build *.go"},
						},
					},
				},
				args: []string{"build"},
			},
			executor: &testCommandAgent{
				executorHandler: func(cmd string) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "success calling other recipes",
			fields: args{
				b: &models.Bakery{
					Recipes: map[string]models.Recipe{
						"build": {
							Steps: []string{"go build *.go"},
						},
						"run": {
							Steps: []string{"build"},
						},
					},
				},
				args: []string{"run"},
			},
			executor: &testCommandAgent{
				executorHandler: func(cmd string) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "success, print help",
			fields: args{
				b: &models.Bakery{
					Recipes: map[string]models.Recipe{
						"build": {
							Steps:       []string{"go build *.go"},
							Description: "builds the project using go",
						},
						"test": {
							Steps:       []string{"go test -v ./..."},
							Description: "",
						},
						"run": {
							Steps: []string{"./run"},
						},
					},
				},
				args: []string{"help"},
			},
			wantErr: false,
		},
		{
			name: "success, print version",
			fields: args{
				b: &models.Bakery{
					Version: "0.6.9",
				},
				args: []string{"version"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := tt

			t.Parallel()

			r := NewRunner(tc.executor)

			if err := r.RunCommand(tc.fields.b, tc.fields.args); (err != nil) != tc.wantErr {
				t.Errorf("Runner.RunCommand() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
