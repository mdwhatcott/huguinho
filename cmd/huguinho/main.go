package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/mdwhatcott/huguinho/core"
)

func main() {
	log.SetFlags(0)
	args := os.Args[1:]
	disk := NewDisk()
	runner := core.NewPipelineRunner(args, disk)
	os.Exit(runner.Run())
}

////////////////////////////////////////////////////////////

type Disk struct{}

func NewDisk() *Disk {
	return new(Disk)
}

func (this *Disk) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (this *Disk) WriteFile(path string, content []byte, perm os.FileMode) error {
	return ioutil.WriteFile(path, content, perm)
}

func (this *Disk) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (this *Disk) Walk(root string, walk filepath.WalkFunc) error {
	return filepath.Walk(root, walk)
}
