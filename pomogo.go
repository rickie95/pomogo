package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func tick(ticker *time.Ticker, duration time.Duration, done <-chan bool) {
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			duration -= 1 * time.Second
			fmt.Printf("\rTick %d", duration/time.Second)
		}
	}
}

func main() {
	duration := time.Duration(50) * time.Minute

	argLeng := len(os.Args[1:])
	if argLeng > 0 {
		// should be in mm:ss form
		minutes, has_errors := strconv.ParseInt(os.Args[1], 10, 0)
		if has_errors != nil {
			duration = time.Duration(minutes) * time.Minute
		} else {
			fmt.Println("Error while parsing. Using default value (50 min)")
		}

	}

	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	// Executes tick asyncronusly
	go tick(ticker, duration, done)

	time.Sleep(5 * time.Second)
	ticker.Stop()
	done <- true
}
