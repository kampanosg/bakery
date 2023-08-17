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
		file, err := open(b)
		if err != nil {
			continue
		}
		return file, nil
	}
	return nil, fmt.Errorf("unable to find Bakefile")
}

func LoadBakefileFromLocation(path string) (*os.File, error) {
	return open(path)
}

func open(file string) (*os.File, error) {
	return os.Open(file)
}
