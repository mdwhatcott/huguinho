package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type (
	Path string
	File string
)

func LoadContent(folder string) map[Path]File {
	content := make(map[Path]File)
	_ = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		data, _ := ioutil.ReadFile(path)
		content[Path(path)] = File(data)
		return nil
	})
	return content
}
