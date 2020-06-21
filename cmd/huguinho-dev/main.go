package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mdwhatcott/huguinho/core"
	"github.com/mdwhatcott/huguinho/io"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	args := os.Args[1:]
	config, err := core.NewCLIParser(args).Parse()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		err := os.RemoveAll(config.TargetRoot)
		if err != nil {
			http.Error(response, "Failed to remove target directory.", http.StatusInternalServerError)
			return
		}

		runner := core.NewPipelineRunner(args, io.NewDisk())
		errCount := runner.Run()
		if errCount > 0 {
			http.Error(response, "Failed to generate site.", http.StatusInternalServerError)
			return
		}

		http.ServeFile(response, request, filepath.Join(config.TargetRoot, request.URL.Path))
	})

	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
