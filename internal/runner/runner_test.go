package runner

import (
	"errors"
	"testing"

	"github.com/kampanosg/bakery/internal/models"
)

type testCommandAgent struct {
	executorHandler func(cmd string) error
}

func (e *testCommandAgent) Execute(cmd string) error {
	return e.executorHandler(cmd)
}

func TestRunner_RunCommand(t *testing.T) {
	type fields struct {
		b    *models.Bakery
		args []string
	}
	tests := []struct {
		name     string
		executor *testCommandAgent
		fields   fields
		wantErr  bool
	}{
		{
			name: "error nil bakery",
			fields: fields{
				args: []string{"list"},
			},
			wantErr: true,
		},
		{
			name: "error undefined recipe",
			fields: fields{
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
			fields: fields{
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
			fields: fields{
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
			fields: fields{
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
			name: "error when recipe is private",
			fields: fields{
				b: &models.Bakery{
					Recipes: map[string]models.Recipe{
						"build": {
							Private: true,
							Steps:   []string{"go build -o app ./..."},
						},
					},
				},
				args: []string{"build"},
			},
			wantErr: true,
		},
		{
			name: "error when recipe is private and called from default",
			fields: fields{
				b: &models.Bakery{
					Defaults: []string{"build"},
					Recipes: map[string]models.Recipe{
						"build": {
							Private: true,
							Steps:   []string{"go build -o app ./..."},
						},
					},
				},
				args: []string{},
			},
			wantErr: true,
		},
		{
			name: "error when recipe is private and called from multiple args",
			fields: fields{
				b: &models.Bakery{
					Recipes: map[string]models.Recipe{
						"build": {
							Private: true,
							Steps:   []string{"go build -o app ./..."},
						},
						"list": {
							Steps: []string{"ls -al"},
						},
					},
				},
				args: []string{"list", "build"},
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
			fields: fields{
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
			fields: fields{
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
			fields: fields{
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
			fields: fields{
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
			fields: fields{
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
			fields: fields{
				b: &models.Bakery{
					Version: "0.6.9",
				},
				args: []string{"version"},
			},
			wantErr: false,
		},
		{
			name: "success, ignore failure",
			fields: fields{
				b: &models.Bakery{
					Recipes: map[string]models.Recipe{
						"list": {
							Steps: []string{"^lss -al"},
						},
					},
				},
				args: []string{"list"},
			},
			executor: &testCommandAgent{
				executorHandler: func(cmd string) error {
					return errors.New("syntaxt error, lss is not valid")
				},
			},
			wantErr: false,
		},
		{
			name: "success, ignore whitespace",
			fields: fields{
				b: &models.Bakery{
					Recipes: map[string]models.Recipe{
						"clean": {
							Steps: []string{"rm app"},
						},
						"build": {
							Steps: []string{
								"^ clean",
								"go build -o app ./...",
							},
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
			name: "success, multiple recipes",
			fields: fields{
				b: &models.Bakery{
					Recipes: map[string]models.Recipe{
						"clean": {
							Steps: []string{"rm app"},
						},
						"build": {
							Steps: []string{
								"^ clean",
								"go build -o app ./...",
							},
						},
					},
				},
				args: []string{"build", "clean"},
			},
			executor: &testCommandAgent{
				executorHandler: func(cmd string) error {
					return nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := tt

			t.Parallel()

			vbs := false
			r := NewRunner(tc.executor, vbs)

			if err := r.Run(tc.fields.b, tc.fields.args); (err != nil) != tc.wantErr {
				t.Errorf("Runner.RunCommand() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestRunner_GetPrintableVersion(t *testing.T) {
	type fields struct {
		agent CommandAgent
	}
	type args struct {
		b *models.Bakery
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "empty when nil version",
			fields: fields{
				agent: &testCommandAgent{},
			},
			args: args{
				b: &models.Bakery{},
			},
			want: "Bakefile Version: \n",
		},
		{
			name: "empty when empty version",
			fields: fields{
				agent: &testCommandAgent{},
			},
			args: args{
				b: &models.Bakery{
					Version: "",
				},
			},
			want: "Bakefile Version: \n",
		},
		{
			name: "success",
			fields: fields{
				agent: &testCommandAgent{},
			},
			args: args{
				b: &models.Bakery{
					Version: "v0.6.9",
				},
			},
			want: "Bakefile Version: v0.6.9\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := tt

			t.Parallel()

			r := &Runner{
				agent: tc.fields.agent,
			}
			if got := r.GetPrintableVersion(tc.args.b); got != tc.want {
				t.Errorf("Runner.GetPrintableVersion() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRunner_GetPrintableHelp(t *testing.T) {
	type fields struct {
		agent CommandAgent
	}
	type args struct {
		b *models.Bakery
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "empty when nil recipes",
			fields: fields{
				agent: &testCommandAgent{},
			},
			args: args{
				b: &models.Bakery{},
			},
			want: "Available Recipes in Bakefile:\n",
		},
		{
			name: "empty when no recipes",
			fields: fields{
				agent: &testCommandAgent{},
			},
			args: args{
				b: &models.Bakery{
					Recipes: map[string]models.Recipe{},
				},
			},
			want: "Available Recipes in Bakefile:\n",
		},
		{
			name: "success",
			fields: fields{
				agent: &testCommandAgent{},
			},
			args: args{
				b: &models.Bakery{
					Recipes: map[string]models.Recipe{
						"list": {
							Description: "lists files in dir",
						},
					},
				},
			},
			want: "Available Recipes in Bakefile:\n- list: lists files in dir\n",
		},
		{
			name: "success multiple recipes",
			fields: fields{
				agent: &testCommandAgent{},
			},
			args: args{
				b: &models.Bakery{
					Recipes: map[string]models.Recipe{
						"list": {
							Description: "lists files in dir",
						},
						"build": {
							Description: "builds the app",
						},
					},
				},
			},
			want: "Available Recipes in Bakefile:\n- list: lists files in dir\n- build: builds the app\n",
		},
		{
			name: "success when help is defined",
			fields: fields{
				agent: &testCommandAgent{},
			},
			args: args{
				b: &models.Bakery{
					Help: "This is the help",
				},
			},
			want: "This is the help",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := tt

			t.Parallel()

			r := &Runner{
				agent: tc.fields.agent,
			}
			if got := r.GetPrintableHelp(tc.args.b); got != tc.want {
				t.Errorf("Runner.GetPrintableHelp() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRunner_GetPrintableAuthor(t *testing.T) {
	type fields struct {
		agent CommandAgent
	}
	type args struct {
		b *models.Bakery
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "empty when nil author",
			fields: fields{
				agent: &testCommandAgent{},
			},
			args: args{
				b: &models.Bakery{},
			},
			want: "Bakefile Author: \n",
		},
		{
			name: "empty when empty author",
			fields: fields{
				agent: &testCommandAgent{},
			},
			args: args{
				b: &models.Bakery{
					Metadata: map[string]string{
						"author": "",
					},
				},
			},
			want: "Bakefile Author: \n",
		},
		{
			name: "success",
			fields: fields{
				agent: &testCommandAgent{},
			},
			args: args{
				b: &models.Bakery{
					Metadata: map[string]string{
						"author": "Post Malone",
					},
				},
			},
			want: "Bakefile Author: Post Malone\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Runner{
				agent: tt.fields.agent,
			}
			if got := r.GetPrintableAuthor(tt.args.b); got != tt.want {
				t.Errorf("Runner.GetPrintableAuthor() = %v, want %v", got, tt.want)
			}
		})
	}
}
