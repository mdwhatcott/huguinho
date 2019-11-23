package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mdwhatcott/static/contracts"
)

func LoadFiles(folder string) map[contracts.Path]contracts.File {
	content := make(map[contracts.Path]contracts.File)
	_ = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		data, _ := ioutil.ReadFile(path)
		content[contracts.Path(strings.TrimPrefix(path, folder))] = contracts.File(data)
		return nil
	})
	return content
}

func WriteFile(path string, data []byte) {
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	_ = ioutil.WriteFile(path, data, 0644)
}
