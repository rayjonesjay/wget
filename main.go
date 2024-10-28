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
	log.Printf("Args context: %#v\n", ctx)
	// Check for invalid commandline args combinations
	{
		die := func(msg string) {
			xerr.WriteError(
				msg,
				1,
				true)
		}

		if ctx.ConvertLinks && !ctx.Mirror {
			die("bad format: option --convert-links is on but --mirror is off")
			return
		}

		if (len(ctx.Exclude) != 0 || len(ctx.Rejects) != 0) && !ctx.Mirror {
			die("bad format: options [--exclude short hand -X; --reject short hand -R] " +
				"are only valid in --mirror mode")
			return
		}

		if ctx.Mirror && ctx.OutputFile != "" {
			die("bad format: option --mirror with -O specified is ambiguous")
			return
		}

		if len(ctx.Links) > 1 && ctx.OutputFile != "" {
			die("bad format: many URLs to download but -O is specified, this is ambiguous")
			return
		}
	}
	downloader.Get(ctx)
}
