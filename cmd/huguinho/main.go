package main

import (
	"log"
	"os"

	"github.com/mdwhatcott/huguinho/core"
	"github.com/mdwhatcott/huguinho/io"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	os.Exit(core.NewPipelineRunner(os.Args[1:], io.NewDisk()).Run())
}
