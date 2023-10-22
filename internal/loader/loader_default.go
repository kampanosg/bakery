package loader

import (
	"os"
	"path/filepath"
)

type DefaultFileOpener struct{}

func (f *DefaultFileOpener) Open(path string) (*os.File, error) {
	p := filepath.Clean(path)
	return os.Open(p)
}
