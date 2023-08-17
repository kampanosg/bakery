package main

import (
	"fmt"
	"os"
)

var (
	DefaultBakefiles = []string{
		"./Bakefile",
		"./Bakefile.yaml",
		"./Bakefile.yml",
	}
)

func LoadDefaultBakefile() (*os.File, error) {
	for _, b := range DefaultBakefiles {
		file, err := os.Open(b)
		if err != nil {
			continue
		}
		return file, nil
	}
	return nil, fmt.Errorf("unable to find Bakefile")
}
