package main

import (
	"fisherman/infrastructure/io"
	"fisherman/runner"
	"fmt"
	"os"
)

func main() {
	fileAccessor := io.NewFileAccessor()
	r := runner.NewRunner(fileAccessor)
	if err := r.Run(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
