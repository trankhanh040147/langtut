package main

import (
	"fmt"
	"os"

	"github.com/trankhanh040147/langtut/internal/cli"
)

func init() {
	cli.InitRoot()
}

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
