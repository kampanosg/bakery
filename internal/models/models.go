package models

import "fmt"

type (
	Bakery struct {
		Version   string            `yaml:"version"`
		Metadata  map[string]string `yaml:"metadata"`
		Variables map[string]string `yaml:"variables"`
		Defaults  []string          `yaml:"defaults"`
		Recipes   map[string]Recipe `yaml:"recipes"`
		Help      string            `yaml:"help"`
	}

	Recipe struct {
		Description string   `yaml:"description"`
		Private     bool     `yaml:"private"`
		Steps       []string `yaml:"steps"`
	}
)

func (r *Bakery) Valid() error {
	if r.Recipes == nil {
		return fmt.Errorf("no recipes found")
	}

	return nil
}
