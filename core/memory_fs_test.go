package core

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type InMemoryFileSystem struct {
	ModTime      time.Time
	Files        map[string]*MemoryFile
	ErrReadFile  map[string]error
	ErrWriteFile map[string]error
	ErrMkdirAll  map[string]error
	ErrWalkFunc  map[string]error
}

func NewInMemoryFileSystem() *InMemoryFileSystem {
	return &InMemoryFileSystem{
		ModTime:      time.Now(),
		Files:        make(map[string]*MemoryFile),
		ErrReadFile:  make(map[string]error),
		ErrWriteFile: make(map[string]error),
		ErrMkdirAll:  make(map[string]error),
		ErrWalkFunc:  make(map[string]error),
	}
}

func (this *InMemoryFileSystem) ReadFile(path string) ([]byte, error) {
	err := this.ErrReadFile[path]
	if err != nil {
		return nil, err
	}
	return this.Files[path].content, nil
}

func (this *InMemoryFileSystem) WriteFile(path string, content []byte, perm os.FileMode) error {
	err := this.ErrWriteFile[path]
	if err != nil {
		return err
	}
	this.Files[path] = &MemoryFile{
		content: content,
		name:    filepath.Base(path),
		mode:    perm,
		modTime: this.ModTime,
	}
	return nil
}

func (this *InMemoryFileSystem) MkdirAll(path string, perm os.FileMode) error {
	err := this.ErrMkdirAll[path]
	if err != nil {
		return err
	}
	components := strings.Split(path, string(os.PathSeparator))
	for x := 1; x < len(components)+1; x++ {
		p := strings.Join(components[:x], string(os.PathSeparator))
		this.Files[p] = &MemoryFile{
			name:    filepath.Base(p),
			mode:    perm,
			modTime: this.ModTime,
			isDir:   true,
		}
	}
	return nil
}

func (this *InMemoryFileSystem) Walk(root string, walk filepath.WalkFunc) error {
	skips := make(map[string]struct{})
	root = filepath.Clean(root)
	var paths []string
	for path := range this.Files {
		if path == root || strings.HasPrefix(path, root+"/") {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	for _, path := range paths {
		shouldSkip := false
		for skip := range skips {
			if strings.HasPrefix(path, skip) {
				shouldSkip = true
				break
			}
		}
		if shouldSkip {
			continue
		}
		file := this.Files[path]
		err := walk(path, file, this.ErrWalkFunc[path])
		if err == filepath.SkipDir {
			skips[path+string(filepath.Separator)] = struct{}{}
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}

///////////////////////////////////////////////////////////////////

type MemoryFile struct {
	content []byte
	name    string
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (this *MemoryFile) Content() string    { return string(this.content) }
func (this *MemoryFile) Name() string       { return this.name }
func (this *MemoryFile) Size() int64        { return int64(len(this.content)) }
func (this *MemoryFile) Mode() os.FileMode  { return this.mode }
func (this *MemoryFile) ModTime() time.Time { return this.modTime }
func (this *MemoryFile) IsDir() bool        { return this.isDir }
func (this *MemoryFile) Sys() any           { panic("NOT NEEDED") }
