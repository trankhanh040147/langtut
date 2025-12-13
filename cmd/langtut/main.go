package main

import (
	"fmt"
	"os"

	"github.com/trankhanh040147/langtut/internal/cli"
)

func main() {
	cli.InitRoot()
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
