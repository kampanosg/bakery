package main

type Recipe struct {
	Version string `yaml:"version"`
}

func (r *Recipe) Valid() error {
	return nil
}
