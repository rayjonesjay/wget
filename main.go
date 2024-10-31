package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
	"wget/fileio"
	"wget/globals"
	"wget/temp"
	"wget/terminal"

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

	log.Println("[MAIN] Initializing Terminal...")
	term, err := terminal.Init()
	if err != nil {
		fmt.Printf("Failed to initialize terminal screen: %v\n", err)
		return
	}
	log.Println("[MAIN] OK")

	defer func(term *terminal.Terminal) {
		log.Println("[MAIN] Cleaning up...")
		err := term.EndWin()
		if err != nil {
			fmt.Printf("Failed to restore terminal screen: %v\n", err)
			return
		}
		log.Println("[MAIN] OK")
	}(term)

	wg := sync.WaitGroup{}

	log.Println("[MAIN] Creating Progress Terminal...")
	exit := make(chan struct{}, 1)
	p := terminal.New(term, exit)
	wg.Add(1)
	go func() {
		defer wg.Done()
		// If this goroutine finishes,
		//the downloader is may still be running, go down with it too
		defer wg.Done()
		log.Println("[OK] Executing `Progress.Run` in new Goroutine")
		p.Run()
		log.Println("[OK] `Progress.Run` exiting new Goroutine")
	}()
	log.Println("[MAIN] OK")

	log.Println("[MAIN] Starting downloader...")
	wg.Add(1)
	go func() {
		defer wg.Done()
		globals.ProgressTerm = p
		downloader.Get(ctx)
		defer wg.Done()
		if !ctx.BackgroundMode {
			time.Sleep(1 * time.Minute)
		}
	}()
	log.Println("[MAIN] OK")

	log.Println("[MAIN] Waiting for all Goroutines to finish.")
	wg.Wait()
	log.Println("[MAIN] All Goroutines finished. Exiting...")
}
