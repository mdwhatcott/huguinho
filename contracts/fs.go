package contracts

import (
	"io"
	"os"
	"path/filepath"
)

type (
	// Deprecated
	Path string

	// Deprecated
	File string
)

type FileSystem interface {
	Stat
	Open
	Create
	ReadFile
	WriteFile
	Mkdir
	MkdirAll
	Remove
	RemoveAll
	Walk
}

type (
	Stat interface {
		Stat(path string) (os.FileInfo, error)
	}

	Open interface {
		Open(path string) (io.ReadCloser, error)
	}

	Create interface {
		Create(path string) (io.WriteCloser, error)
	}

	ReadFile interface {
		ReadFile(path string) ([]byte, error)
	}

	WriteFile interface {
		WriteFile(path string, content []byte, perm os.FileMode) error
	}

	Mkdir interface {
		Mkdir(path string, perm os.FileMode) error
	}

	MkdirAll interface {
		MkdirAll(path string, perm os.FileMode) error
	}

	Remove interface {
		Remove(path string) error
	}

	RemoveAll interface {
		RemoveAll(path string) error
	}

	Walk interface {
		Walk(root string, walk filepath.WalkFunc) error
	}
)
