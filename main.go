package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"wget/fileio"
	"wget/temp"

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

	// configure file logging to temporary application logger file
	logger, err := os.OpenFile(path.Join(temp.Dir(), "logger.log"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("failed to setup file logging: logging to stderr instead: %v\n", err)
	}
	log.SetOutput(logger)
	defer fileio.Close(logger)

	// check command-line args and download the defined files
	ctx := args.DownloadContext(arguments)
	downloader.Get(ctx)
}
