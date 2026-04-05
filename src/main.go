package main

import (
	"os"

	"github.com/attila/focalytics/cmd"
)

var exitFunc = os.Exit

func main() {
	exitFunc(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin *os.File, stdout *os.File, stderr *os.File) int {
	return cmd.Execute(args, stdin, stdout, stderr)
}
