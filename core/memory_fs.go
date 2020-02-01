package core

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type InMemoryFileSystem struct {
	ModTime      time.Time
	Files        map[string]*MemoryFile
	ErrStat      map[string]error
	ErrOpen      map[string]error
	ErrCreate    map[string]error
	ErrReadFile  map[string]error
	ErrWriteFile map[string]error
	ErrMkdir     map[string]error
	ErrMkdirAll  map[string]error
	ErrRemove    map[string]error
	ErrRemoveAll map[string]error
	ErrWalkFunc  map[string]error
	ErrWalk      map[string]error
}

func NewInMemoryFileSystem() *InMemoryFileSystem {
	return &InMemoryFileSystem{
		ModTime:      time.Now(),
		Files:        make(map[string]*MemoryFile),
		ErrStat:      make(map[string]error),
		ErrOpen:      make(map[string]error),
		ErrCreate:    make(map[string]error),
		ErrReadFile:  make(map[string]error),
		ErrWriteFile: make(map[string]error),
		ErrMkdir:     make(map[string]error),
		ErrMkdirAll:  make(map[string]error),
		ErrRemove:    make(map[string]error),
		ErrRemoveAll: make(map[string]error),
		ErrWalkFunc:  make(map[string]error),
		ErrWalk:      make(map[string]error),
	}
}

func (this *InMemoryFileSystem) Stat(path string) (os.FileInfo, error) {
	return this.Files[path], this.ErrStat[path]
}

func (this *InMemoryFileSystem) Open(path string) (io.ReadCloser, error) {
	return this.Files[path].Reader(), this.ErrOpen[path]
}

func (this *InMemoryFileSystem) Create(path string) (io.WriteCloser, error) {
	err := this.ErrCreate[path]
	if err != nil {
		return nil, err
	}
	this.Files[path] = &MemoryFile{
		name:    filepath.Base(path),
		mode:    0644,
		modTime: this.ModTime,
	}
	return this.Files[path], nil
}

func (this *InMemoryFileSystem) ReadFile(path string) ([]byte, error) {
	return this.Files[path].content, this.ErrReadFile[path]
}

func (this *InMemoryFileSystem) WriteFile(path string, content []byte, perm os.FileMode) error {
	err := this.ErrWriteFile[path]
	if err != nil {
		return err
	}
	this.Files[path] = &MemoryFile{
		name:    filepath.Base(path),
		content: content,
		mode:    perm,
		modTime: this.ModTime,
	}
	return nil
}

func (this *InMemoryFileSystem) Mkdir(path string, perm os.FileMode) error {
	err := this.ErrMkdir[path]
	if err != nil {
		return err
	}
	this.Files[path] = &MemoryFile{
		name:    filepath.Base(path),
		mode:    perm,
		modTime: this.ModTime,
		isDir:   true,
	}
	return nil
}

func (this *InMemoryFileSystem) MkdirAll(path string, perm os.FileMode) error {
	err := this.ErrMkdirAll[path]
	if err != nil {
		return err
	}
	components := strings.Split(path, string(os.PathSeparator))
	for x := 1; x < len(components); x++ {
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

func (this *InMemoryFileSystem) Remove(path string) error {
	err := this.ErrRemove[path]
	if err != nil {
		return err
	}
	delete(this.Files, path)
	return nil
}

func (this *InMemoryFileSystem) RemoveAll(path string) error {
	err := this.ErrRemoveAll[path]
	if err != nil {
		return err
	}
	for key := range this.Files {
		if strings.HasPrefix(key, path) {
			delete(this.Files, key)
		}
	}
	return nil
}

func (this *InMemoryFileSystem) Walk(root string, walk filepath.WalkFunc) error {
	for path, file := range this.Files {
		err := walk(path, file, this.ErrWalkFunc[path])
		if err != nil {
			return err
		}
	}
	return this.ErrWalk[root]
}

///////////////////////////////////////////////////////////////////

type MemoryFile struct {
	name     string
	content  []byte
	mode     os.FileMode
	modTime  time.Time
	isDir    bool
	readErr  error
	closeErr error
	writeErr error
}

func (this *MemoryFile) Reader() io.ReadCloser {
	reader := bytes.NewReader(this.content)
	return NewErringReadCloser(reader, this.readErr, this.closeErr)
}

func (this *MemoryFile) Write(p []byte) (written int, err error) {
	this.content = append(this.content, p...)
	return len(p), this.writeErr
}

func (this *MemoryFile) Close() error {
	return this.closeErr
}

func (this *MemoryFile) Name() string       { return this.name }
func (this *MemoryFile) Size() int64        { return int64(len(this.content)) }
func (this *MemoryFile) Mode() os.FileMode  { return this.mode }
func (this *MemoryFile) ModTime() time.Time { return this.modTime }
func (this *MemoryFile) IsDir() bool        { return this.isDir }
func (this *MemoryFile) Sys() interface{}   { panic("NOT NEEDED") }

////////////////////////////////////////////////////////////

type ErringReadCloser struct {
	io.Reader
	readErr  error
	closeErr error
}

func NewErringReadCloser(reader io.Reader, readErr, closeErr error) *ErringReadCloser {
	return &ErringReadCloser{
		Reader:   reader,
		readErr:  readErr,
		closeErr: closeErr,
	}
}

func (this *ErringReadCloser) Read(p []byte) (int, error) {
	read, err := this.Reader.Read(p)
	if this.readErr != nil && err == nil {
		err = this.readErr
	}
	return read, err
}

func (this *ErringReadCloser) Close() error {
	return this.closeErr
}
