package loader

import "os"

type DefaultFileOpener struct{}

func (f *DefaultFileOpener) Open(path string) (*os.File, error) {
	return os.Open(path)
}
