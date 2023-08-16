package main

import "fmt"

type (
	Bakery struct {
		Version string   `yaml:"version"`
		Recipes []Recipe `yaml:"recipes"`
	}

	Recipe struct {
		Description string   `yaml:"description"`
		Default     bool     `yaml:"default"`
		Steps       []string `yaml:"steps"`
	}
)

func (r *Bakery) Valid() error {
	if r.Recipes == nil {
		return fmt.Errorf("no recipes found")
	}

	return nil
}
