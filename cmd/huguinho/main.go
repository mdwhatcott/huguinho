package main

import (
	"os"
	"time"
)

func main() {
	os.Exit(NewProgram(time.Now()).Run())
}
