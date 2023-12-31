package models

import "testing"

func TestBakery_Valid(t *testing.T) {
	type fields struct {
		Version string
		Recipes map[string]Recipe
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "error no recipes",
			fields: fields{
				Version: "1",
			},
			wantErr: true,
		},
		{
			name: "success empty recipes",
			fields: fields{
				Version: "1",
				Recipes: map[string]Recipe{},
			},
			wantErr: false,
		},
		{
			name: "success",
			fields: fields{
				Version: "1",
				Recipes: map[string]Recipe{
					"list": {
						Description: "a step to list the filesystem",
						Steps: []string{
							"ls -al",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := tt

			t.Parallel()

			r := &Bakery{
				Version: tc.fields.Version,
				Recipes: tc.fields.Recipes,
			}
			if err := r.Valid(); (err != nil) != tc.wantErr {
				t.Errorf("Bakery.Valid() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
