package io

import (
	"os"
	"path/filepath"
)

type Disk struct{}

func NewDisk() *Disk {
	return new(Disk)
}

func (this *Disk) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (this *Disk) WriteFile(path string, content []byte, perm os.FileMode) error {
	return os.WriteFile(path, content, perm)
}

func (this *Disk) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (this *Disk) Walk(root string, walk filepath.WalkFunc) error {
	return filepath.Walk(root, walk)
}
