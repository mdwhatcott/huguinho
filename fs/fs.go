package fs

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mdwhatcott/huguinho/contracts"
)

func LoadFiles(folder string) map[contracts.Path]contracts.File {
	content := make(map[contracts.Path]contracts.File)
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		content[contracts.Path(strings.TrimPrefix(path, folder))] = contracts.File(data)
		return nil
	})
	if err != nil {
		log.Println("Error walking filesystem:", err)
	}
	return content
}

func WriteFile(path string, data []byte) {
	mkdir(filepath.Dir(path))
	writeFile(path, data)
}
func mkdir(dir string) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatal(err)
	}
}
func writeFile(path string, data []byte) {
	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
