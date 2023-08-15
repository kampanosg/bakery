package main

import "os"

const (
	DefaultBakefile = "./Bakefile"
)

func LoadBakefile() (*os.File, error) {
	file, err := os.Open(DefaultBakefile)
	if err != nil {
		return nil, err
	}
	return file, nil
}
