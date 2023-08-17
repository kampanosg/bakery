package main

import (
	"fmt"
	"os"
)

const (
	DefaultBakefile = "./Bakefile"
)

var (
	DefaultBakefiles = []string{
		"./Bakefile",
		"./Bakefile.yaml",
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
