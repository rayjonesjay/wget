package main

import (
	"fmt"
	"os"
	"runtime"

	"wget/args"
	"wget/downloader"
	"wget/help"

	"wget/syscheck"
	"wget/xerr"
)

func main() {
	operatingSys := runtime.GOOS
	if err := syscheck.CheckOperatingSystem(operatingSys); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
		return
	}

	arguments := os.Args[1:]
	if len(arguments) == 0 {
		// no arguments were passed, so return
		xerr.WriteError(help.UsageMessage, 1, true)
		return
	}

	ctx := args.DownloadContext(arguments)
	downloader.Get(ctx)
}
