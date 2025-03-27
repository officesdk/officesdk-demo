package main

import (
	"fmt"
	"os"

	"office-demo/cmd"
	_ "office-demo/cmd"
	_ "office-demo/cmd/server"
)

func main() {
	if err := cmd.RootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return
}
