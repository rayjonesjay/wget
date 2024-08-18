package main

import (
	"fmt"
	"os"

	"wget/args"
	"wget/errors"
	"wget/syscheck"
	"wget/types"
)

func main() {
	if err := syscheck.CheckOperatingSystem(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	// download.DownloadUrl(errors.Link1)
	arguments := os.Args[1:]
	if len(arguments) == 0 {
		// no arguments were passed, so return
		errors.WriteError(types.UsageMessage,1)
		return
	}

	// if arguments are passed.
	args.EvalArgs(arguments)
}
