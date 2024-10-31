package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"
	"wget/fileio"
	"wget/temp"
	"wget/terminal"
)

func main() {
	// configure file logging to temporary application logger file
	logger, err := os.OpenFile(path.Join(temp.Dir(), "logger.log"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("failed to setup file logging: logging to stderr instead: %v\n", err)
	}
	log.SetOutput(logger)
	defer fileio.Close(logger)

	log.Println("[OK] Starting...")
	term, err := terminal.Init()
	if err != nil {
		fmt.Printf("Failed to initialize terminal screen: %v\n", err)
		return
	}

	defer func(term *terminal.Terminal) {
		log.Println("Cleaning up...")
		err := term.EndWin()
		if err != nil {
			fmt.Printf("Failed to restore terminal screen: %v\n", err)
			return
		}
		log.Println("Done")
	}(term)
	wg := sync.WaitGroup{}

	log.Println("[OK] Creating Progress Terminal")
	exit := make(chan struct{}, 1)
	p := terminal.New(term, exit)
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("[OK] Executing `Progress.Run` in new Goroutine")
		p.Run()
		log.Println("[OK] `Progress.Run` exiting new Goroutine")
	}()

	log.Println("[OK] Creating Goroutine 3")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			time.Sleep(time.Millisecond * 50)
			text := fmt.Sprintf("%d. hello@%d #3", i, i)
			log.Println(text)
			p.SetLine(i, text)
		}
	}()

	log.Println("[OK] Waiting for all Goroutines to finish.")
	wg.Wait()
	log.Println("[OK] All Goroutines finished. Exiting...")
}
