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

	// vars := make(map[string]string, len(b.Variables))

	recipes := make(map[string]models.Recipe, len(b.Recipes))
	for k, recipe := range b.Recipes {
		steps := make([]string, len(recipe.Steps))
		for i, step := range recipe.Steps {
			fmt.Printf("step: %s\n", step)
			vars := r.FindAllString(step, -1)
			fmt.Printf("vars: %v\n", vars)
			if len(vars) == 0 {
				continue
			}
			for _, v := range vars {
				vv := v[1 : len(v)-1]
				val, ok := b.Variables[vv]
				if !ok {
					return nil, fmt.Errorf("variable %s not found", vv)
				}
				step = strings.ReplaceAll(step, v, val)
				fmt.Printf("replaced: %s\n", step)
			}
			steps[i] = step
		}
		recipe.Steps = steps
		recipes[k] = recipe
	}
	b.Recipes = recipes
	fmt.Printf("b: %+v\n", b)
	panic("stop")

	return &b, nil
}
