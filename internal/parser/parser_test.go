package parser

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/kampanosg/bakery/internal/models"
)

func TestParseBakefile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Bakery
		wantErr bool
	}{
		{
			name: "error empty bakefile",
			args: args{
				path: "test-files/Bakefile-empty",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error no recipes",
			args: args{
				path: "test-files/Bakefile-no-recipes",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error invalid syntax",
			args: args{
				path: "test-files/Bakefile-invalid",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error wrong variable",
			args: args{
				path: "test-files/Bakefile-wrong-vars",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success with version",
			args: args{
				path: "test-files/Bakefile-version",
			},
			want: &models.Bakery{
				Version: "v1",
				Recipes: map[string]models.Recipe{
					"list": {
						Steps: []string{"ls -al"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "success with metadata",
			args: args{
				path: "test-files/Bakefile-metadata",
			},
			want: &models.Bakery{
				Metadata: map[string]string{
					"author": "Darth Vader",
				},
				Recipes: map[string]models.Recipe{
					"list": {
						Steps: []string{"ls -al"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "success with defaults",
			args: args{
				path: "test-files/Bakefile-defaults",
			},
			want: &models.Bakery{
				Defaults: []string{"list"},
				Recipes: map[string]models.Recipe{
					"list": {
						Steps: []string{"ls -al"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "success with full syntax",
			args: args{
				path: "test-files/Bakefile-full",
			},
			want: &models.Bakery{
				Version: "v1",
				Metadata: map[string]string{
					"author": "Darth Vader",
				},
				Variables: map[string]string{
					"path": "/home/darth-vader",
				},
				Defaults: []string{"list"},
				Recipes: map[string]models.Recipe{
					"create-user": {
						Steps: []string{
							"ls -al",
							"mkdir /home/darth-vader",
						},
						Description: "creates the user dir",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "success with multiple vars",
			args: args{
				path: "test-files/Bakefile-multiple-vars",
			},
			want: &models.Bakery{
				Variables: map[string]string{
					"path": "/home/darth-vader",
					"file": "passwords.txt",
				},
				Recipes: map[string]models.Recipe{
					"create-pwords": {
						Steps: []string{
							"touch /home/darth-vader/passwords.txt",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := tt

			t.Parallel()

			wd, _ := os.Getwd()
			f, err := os.Open(fmt.Sprintf("%s/%s", wd, tc.args.path))
			if err != nil {
				t.Errorf("ParseBakefile() unable to open file in args, %v", err)
				return
			}
			defer f.Close()

			got, err := ParseBakefile(f)
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseBakefile() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ParseBakefile() = %v, want %v", got, tc.want)
			}
		})
	}
}
