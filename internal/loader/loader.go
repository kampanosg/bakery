package loader

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

type (
	FileOpener interface {
		Open(f string) (*os.File, error)
	}

	BakefileLoader struct {
		Opener FileOpener
	}
)

func NewBakefileLoader(o FileOpener) *BakefileLoader {
	return &BakefileLoader{
		Opener: o,
	}
}

func (l *BakefileLoader) LoadDefaultBakefile() (*os.File, error) {
	for _, b := range DefaultBakefiles {
		file, err := l.Opener.Open(b)
		if err != nil {
			continue
		}
		return file, nil
	}
	return nil, fmt.Errorf("Bakefile not found")
}

func (l *BakefileLoader) LoadBakefileFromLocation(path string) (*os.File, error) {
	return l.Opener.Open(path)
}
