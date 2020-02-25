package main

import (
	"log"
	"os"

	"github.com/mdwhatcott/huguinho/core"
	"github.com/mdwhatcott/huguinho/fs"
)

func main() {
	log.SetFlags(0)
	args := os.Args[1:]
	disk := fs.NewDisk()
	runner := core.NewPipelineRunner(args, disk)
	os.Exit(runner.Run())
}
