package main

import (
	"fmt"
	"os"
	"wget/args"
	"wget/downloader"
	"wget/help"

	"wget/errorss"
	"wget/syscheck"
)

func main() {
	if err := syscheck.CheckOperatingSystem(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
		return
	}

	arguments := os.Args[1:]
	if len(arguments) == 0 {
		// no arguments were passed, so return
		errorss.WriteError(help.UsageMessage, 1, true)
		return
	}

	ctx := args.DownloadContext(arguments)
	downloader.Get(ctx)
}
