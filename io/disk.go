package io

import (
	"os"
	"path/filepath"
)

type Disk struct{}

func (Disk) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
func (Disk) WriteFile(path string, content []byte, perm os.FileMode) error {
	return os.WriteFile(path, content, perm)
}
func (Disk) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
func (Disk) Walk(root string, walk filepath.WalkFunc) error {
	return filepath.Walk(root, walk)
}
