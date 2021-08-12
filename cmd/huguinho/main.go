package main

import (
	"log"
	"os"
	"time"

	"github.com/mdwhatcott/huguinho/core"
	"github.com/mdwhatcott/huguinho/io"
)

var Version = "dev"

func main() {
	runner := core.NewPipelineRunner(
		Version,
		os.Args[1:],
		io.Disk{},
		time.Now,
		log.New(os.Stderr, "", log.Lshortfile),
	)
	os.Exit(runner.Run())
}
