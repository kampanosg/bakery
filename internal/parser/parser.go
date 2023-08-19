package parser

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kampanosg/bakery/internal/models"
	"gopkg.in/yaml.v2"
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

	return &b, nil
}
