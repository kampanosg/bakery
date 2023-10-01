package parser

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/kampanosg/bakery/internal/models"
	"gopkg.in/yaml.v2"
)

var (
	r, _ = regexp.Compile(`(:{1}\w*:{1})`)
)

func ParseBakefile(f *os.File) (*models.Bakery, error) {
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read file, %w", err)
	}

	var b models.Bakery
	if err := yaml.Unmarshal(content, &b); err != nil {
		return nil, fmt.Errorf("cannot unmarshal file, %w", err)
	}

	if err := b.Valid(); err != nil {
		return nil, fmt.Errorf("invalid Bakefile, %w", err)
	}

	recipes := make(map[string]models.Recipe, len(b.Recipes))
	for k, r := range b.Recipes {
		steps := make([]string, len(r.Steps))
		for i, s := range r.Steps {
			vars := parseVars(s)
			if len(vars) > 0 {
				for _, v := range vars {
					if s, err = extrapolate(&b, s, v); err != nil {
						return nil, fmt.Errorf("cannot extrapolate vars, %w", err)
					}
				}
			}
			steps[i] = s
		}

		r.Steps = steps
		recipes[k] = r
	}

	b.Recipes = recipes

	return &b, nil
}

// parseVars returns a slice of strings that match the regex
// e.g. i am :name:, :age: years old -> [:name:, :age:]
func parseVars(s string) []string {
	return r.FindAllString(s, -1)
}

// extrapolate replaces the variable with its value
func extrapolate(b *models.Bakery, step, v string) (string, error) {
	vv := v[1 : len(v)-1] // remove the colons
	val, ok := b.Variables[vv]
	if !ok {
		return "", fmt.Errorf("variable :%s: not found", vv)
	}
	return strings.ReplaceAll(step, v, val), nil
}
