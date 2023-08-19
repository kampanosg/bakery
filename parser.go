package main

import (
	"io/ioutil"
	"os"

	"github.com/kampanosg/bakery/internal/models"
	"gopkg.in/yaml.v2"
)

func ParseBakefile(f *os.File) (*models.Bakery, error) {
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var r models.Bakery
	if err := yaml.Unmarshal(content, &r); err != nil {
		return nil, err
	}

	if err := r.Valid(); err != nil {
		return nil, err
	}

	return &r, nil
}
