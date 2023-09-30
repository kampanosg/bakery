package parser

import (
	"fmt"
	"io/ioutil"
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
	content, err := ioutil.ReadAll(f)
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
	for k, recipe := range b.Recipes {
		steps := make([]string, len(recipe.Steps))
		for i, step := range recipe.Steps {
			vars := getVariables(step)
			if len(vars) == 0 {
				continue
			}
			for _, v := range vars {
				if step, err = extrapolate(&b, step, v); err != nil {
					return nil, fmt.Errorf("cannot parse step %s, %w", step, err)
				}
			}
			steps[i] = step
		}

		recipe.Steps = steps
		recipes[k] = recipe
	}

	b.Recipes = recipes

	return &b, nil
}

// getVariables returns a slice of strings that match the regex
// it's in its own function to make testing easier
func getVariables(s string) []string {
	return r.FindAllString(s, -1)
}

// extrapolate replaces the variable with its value
func extrapolate(b *models.Bakery, step, v string) (string, error) {
	vv := v[1 : len(v)-1] // remove the colons
	val, ok := b.Variables[vv]
	if !ok {
		return "", fmt.Errorf("variable %s not found", vv)
	}
	return strings.ReplaceAll(step, v, val), nil
}
