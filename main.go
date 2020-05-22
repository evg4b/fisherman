package main

import (
	"fisherman/runner"
	"fmt"
	"os"
)

func main() {
	if err := runner.Run(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
