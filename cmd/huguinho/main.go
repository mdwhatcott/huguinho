package main

import (
	"log"
	"os"
	"time"

	"github.com/mdwhatcott/huguinho/core"
	"github.com/mdwhatcott/huguinho/io"
)

var Version = "dev"

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	runner := core.NewPipelineRunner(
		Version,
		os.Args[1:],
		io.NewDisk(),
		time.Now,
		log.Default(),
	)
	os.Exit(runner.Run())
}
