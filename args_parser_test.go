package main

import (
	"testing"

	match "github.com/kampanosg/go-match-slices"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "only program name",
			args: []string{"bake"},
			want: []string{},
		},
		{
			name: "only one command",
			args: []string{"bake", "version"},
			want: []string{"version"},
		},
		{
			name: "only flag",
			args: []string{"bake", "--file", "Test"},
			want: []string{},
		},
		{
			name: "flag and command",
			args: []string{"bake", "--file", "Test", "version"},
			want: []string{"version"},
		},
		{
			name: "command and flag",
			args: []string{"bake", "version", "--file", "Test"},
			want: []string{"version"},
		},
		{
			name: "multiple flags",
			args: []string{"bake", "--file", "Test", "--version", "--upgrade"},
			want: []string{},
		},
		{
			name: "multiple flags one command",
			args: []string{"bake", "--file", "Test", "--version", "--upgrade", "version"},
			want: []string{"version"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := tt

			t.Parallel()

			got := ParseArgs(tc.args)
			match.MatchExactly(t, tc.want, got)
		})
	}
}
