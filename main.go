package main

import (
	"fisherman/infrastructure/io"
	"fisherman/runner"
	"fmt"
	"os"
	"os/user"
)

func main() {
	fileAccessor := io.NewFileAccessor()
	usr, err := user.Current()
	handleError(err)
	r := runner.NewRunner(fileAccessor, usr)
	handleError(r.Run(os.Args[1:]))
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
