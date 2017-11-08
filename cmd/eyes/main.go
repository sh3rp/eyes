package main

import (
	"fmt"
	"os"

	"github.com/sh3rp/eyes/cmd/eyes/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
