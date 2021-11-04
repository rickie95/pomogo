package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const DEFAULT_TASK_MINUTES = 50
const DEFAULT_BREAK_MINUTES = 5
const DEFAULT_TASK_FOR_SESSION = 4
const DEFAULT_LONG_BREAK_MINUTES = 15

func tick(ticker *time.Ticker, duration time.Duration, done chan bool) {
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			duration -= 1 * time.Second
			if duration < 0 {
				done <- true
			}
			fmt.Printf("\rTick %02d:%02d",
				duration/time.Minute,
				duration/time.Second%60+1)
		}
	}
}

func cleanup() {
	fmt.Println("\n\rTimer stopped")
}

func main() {
	task_duration := time.Duration(5) * time.Second
	break_duration := time.Duration(3) * time.Second

	argLeng := len(os.Args[1:])
	if argLeng > 0 {
		// should be in mm:ss form
		minutes, has_errors := strconv.ParseInt(os.Args[1], 10, 0)
		if has_errors != nil {
			task_duration = time.Duration(minutes) * time.Minute
		} else {
			fmt.Println("Error while parsing. Using default value (50 min)")
		}
	}

	// Creates a channel for OS signals
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()

	// Creates a ticker and a channel
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	for {
		// task
		println("\nTask")
		go tick(ticker, task_duration, done)
		<-done
		done <- false
		println("\nBreak")
		go tick(ticker, break_duration, done)
		<-done
		done <- false
	}
	// Executes tick asyncronusly

	ticker.Stop()
	done <- true
}
