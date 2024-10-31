package main

import (
	"fmt"
	"time"
	"wget/terminal"
)

func main() {
	term, err := terminal.Init()
	if err != nil {
		fmt.Printf("Failed to initialize terminal screen: %v\n", err)
		return
	}

	defer func(term *terminal.Terminal) {
		err := term.EndWin()
		if err != nil {
			fmt.Printf("Failed to restore terminal screen: %v\n", err)
			return
		}
	}(term)

	exitSignal := make(chan any, 1)
	go func() {
		time.Sleep(5 * time.Second)
		exitSignal <- 3
	}()

	term.PrintAt(1, 0, "1. One")
	term.PrintAt(2, 0, "2. Two")
	term.PrintAt(3, 0, "3. Three")
	term.PrintAt(0, 0, "0. Zero")

	<-exitSignal
}
