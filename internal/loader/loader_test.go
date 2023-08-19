package loader

import (
	"errors"
	"os"
	"testing"
)

type testFileOpener struct {
	openerHandler func(cmd string) (*os.File, error)
}

func (f *testFileOpener) Open(path string) (*os.File, error) {
	return f.openerHandler(path)
}

func TestBakefileLoader_LoadDefaultBakefile(t *testing.T) {
	type fields struct {
		Opener FileOpener
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "error file not found",
			fields: fields{
				Opener: &testFileOpener{
					openerHandler: func(cmd string) (*os.File, error) {
						return nil, errors.New("file not found")
					},
				},
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				Opener: &testFileOpener{
					openerHandler: func(cmd string) (*os.File, error) {
						return nil, nil
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &BakefileLoader{
				Opener: tt.fields.Opener,
			}
			_, err := l.LoadDefaultBakefile()
			if (err != nil) != tt.wantErr {
				t.Errorf("BakefileLoader.LoadDefaultBakefile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestBakefileLoader_LoadBakefileFromLocation(t *testing.T) {
	type fields struct {
		Opener FileOpener
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "error file not found",
			fields: fields{
				Opener: &testFileOpener{
					openerHandler: func(cmd string) (*os.File, error) {
						return nil, errors.New("file not found")
					},
				},
			},
			args: args{
				path: "UnknownBakefile",
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				Opener: &testFileOpener{
					openerHandler: func(cmd string) (*os.File, error) {
						return nil, nil
					},
				},
			},
			args: args{
				path: "Cupcake",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := tt

			t.Parallel()

			l := &BakefileLoader{
				Opener: tc.fields.Opener,
			}
			_, err := l.LoadBakefileFromLocation(tc.args.path)
			if (err != nil) != tc.wantErr {
				t.Errorf("BakefileLoader.LoadBakefileFromLocation() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}
