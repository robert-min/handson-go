package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/handson-go/chap2/sub-cmd-arch/cmd"
)

var errInvalidSubCommand = errors.New("Invalid sub-command specified")

// Print usage
func printUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: mync [http|grpc] -h\n")
	cmd.HandleHttp(w, []string{"-h"})
	cmd.HandleGrpc(w, []string{"-h"})
}

// Handling command
func handleCommand(w io.Writer, args []string) error {
	var err error

	if len(args) < 1 {
		err = errInvalidSubCommand
	} else {
		switch args[0] {
		case "http":
			err = cmd.HandleHttp(w, args[1:])
		case "grpc":
			err = cmd.HandleGrpc(w, args[1:])
		case "-h":
			printUsage(w)
		case "-help":
			printUsage(w)
		default:
			err = errInvalidSubCommand
		}
	}

	return err

}

func main() {
	err := handleCommand(os.Stdout, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
}
