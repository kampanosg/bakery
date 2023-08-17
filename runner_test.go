package main

import (
	"errors"
	"testing"
)

type testCommandExecutor struct {
	executorHandler func(cmd string) error
}

func (e *testCommandExecutor) Run(cmd string) error {
	return e.executorHandler(cmd)
}

func TestRunner_RunCommand(t *testing.T) {
	type args struct {
		b    *Bakery
		args []string
	}
	tests := []struct {
		name     string
		executor *testCommandExecutor
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
			name: "error no args",
			fields: args{
				b: &Bakery{},
			},
			wantErr: true,
		},
		{
			name: "error undefined recipe",
			fields: args{
				b: &Bakery{
					Recipes: map[string]Recipe{
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
				b: &Bakery{
					Recipes: map[string]Recipe{
						"build": {
							Steps: []string{"go bild *.go"},
						},
					},
				},
				args: []string{"build"},
			},
			executor: &testCommandExecutor{
				executorHandler: func(cmd string) error {
					return errors.New("invalid command bild")
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: args{
				b: &Bakery{
					Recipes: map[string]Recipe{
						"build": {
							Steps: []string{"go build *.go"},
						},
					},
				},
				args: []string{"build"},
			},
			executor: &testCommandExecutor{
				executorHandler: func(cmd string) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "success calling other recipes",
			fields: args{
				b: &Bakery{
					Recipes: map[string]Recipe{
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
			executor: &testCommandExecutor{
				executorHandler: func(cmd string) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "success, print help",
			fields: args{
				b: &Bakery{
					Recipes: map[string]Recipe{
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
				b: &Bakery{
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
