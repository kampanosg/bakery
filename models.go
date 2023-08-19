package main

import "fmt"

type (
	Bakery struct {
		Version  string            `yaml:"version"`
		Metadata map[string]string `yaml:"metadata"`
		Defaults []string          `yaml:"defaults"`
		Recipes  map[string]Recipe `yaml:"recipes"`
	}

	Recipe struct {
		Description string   `yaml:"description"`
		Steps       []string `yaml:"steps"`
	}
)

func (r *Bakery) Valid() error {
	if r.Recipes == nil {
		return fmt.Errorf("no recipes found")
	}

	return nil
}
