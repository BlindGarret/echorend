package externals

import (
	"os"
	"path/filepath"
)

type FileAccess interface {
	Glob(pattern string) ([]string, error)
	ReadFile(filename string) ([]byte, error)
}

type StdFileAccess struct {
}

func (g *StdFileAccess) Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

func (g *StdFileAccess) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
