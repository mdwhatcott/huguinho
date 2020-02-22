package contracts

import (
	"os"
	"path/filepath"
)

type FileSystem interface {
	ReadFile
	WriteFile
	MkdirAll
	Walk
}

type (
	ReadFile interface {
		ReadFile(path string) ([]byte, error)
	}

	WriteFile interface {
		WriteFile(path string, content []byte, perm os.FileMode) error
	}

	MkdirAll interface {
		MkdirAll(path string, perm os.FileMode) error
	}

	Walk interface {
		Walk(root string, walk filepath.WalkFunc) error
	}
)
