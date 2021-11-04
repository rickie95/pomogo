package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
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

func main() {

	var task_duration_minutes int
	var break_duration_minutes int
	var session_break_duration_minutes int
	var task_per_session int

	flag.IntVar(&task_duration_minutes, "d", DEFAULT_TASK_MINUTES, "Sets task duration")
	flag.IntVar(&break_duration_minutes, "b", DEFAULT_BREAK_MINUTES, "Sets break duration")
	flag.IntVar(&session_break_duration_minutes, "s", DEFAULT_LONG_BREAK_MINUTES, "Sets session break duration")
	flag.IntVar(&task_per_session, "t", DEFAULT_TASK_FOR_SESSION, "Sets task per session")

	flag.Parse()

	task_duration := time.Duration(DEFAULT_TASK_MINUTES) * time.Minute
	break_duration := time.Duration(DEFAULT_BREAK_MINUTES) * time.Minute

	// Creates a ticker and a channel
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	// Creates a channel for OS signals
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c // waits for channel
		ticker.Stop()
		done <- true
		os.Exit(0)
	}()
	tasks := 1
	for {
		println("\nTask")
		go tick(ticker, task_duration, done)
		<-done
		done <- false

		// After 4 tasks the break must be longer
		if tasks > 4 {
			tasks = 0
			break_duration = time.Duration(DEFAULT_LONG_BREAK_MINUTES) * time.Minute
		}

		println("\nBreak")
		go tick(ticker, break_duration, done)
		<-done
		done <- false
		// restore break duration
		break_duration = time.Duration(DEFAULT_BREAK_MINUTES) * time.Minute
	}
}
