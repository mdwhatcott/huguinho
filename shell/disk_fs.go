package shell

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mdwhatcott/huguinho/contracts"
)

type Disk struct{}

func NewDisk() *Disk {
	return new(Disk)
}
func (this *Disk) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}
func (this *Disk) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}
func (this *Disk) Create(path string) (io.WriteCloser, error) {
	return os.Create(path)
}
func (this *Disk) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
func (this *Disk) WriteFile(path string, content []byte, perm os.FileMode) error {
	return ioutil.WriteFile(path, content, perm)
}
func (this *Disk) Mkdir(path string, perm os.FileMode) error {
	return os.Mkdir(path, perm)
}
func (this *Disk) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
func (this *Disk) Remove(path string) error {
	return os.Remove(path)
}
func (this *Disk) RemoveAll(path string) error {
	return os.RemoveAll(path)
}
func (this *Disk) Walk(root string, walk filepath.WalkFunc) error {
	return filepath.Walk(root, walk)
}

// TEST (integration)?
func (this *Disk) CopyFile(source, destination string) error {
	_, err := this.Stat(source)
	if os.IsNotExist(err) {
		return err
	}

	err = this.MkdirAll(filepath.Dir(destination), 0755)
	if err != nil {
		return err
	}

	src, err := this.Open(source)
	if err != nil {
		return contracts.StackTraceError(err)
	}

	defer func() { _ = src.Close() }()

	dst, err := this.Create(destination)
	if err != nil {
		return contracts.StackTraceError(err)
	}

	defer func() { _ = dst.Close() }()

	_, err = io.Copy(dst, src)
	return err
}
