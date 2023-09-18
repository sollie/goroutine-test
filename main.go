package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Worker struct {
	Function string
	Args     []interface{}
	Sleep    int32
}

func main() {
	workers := []Worker{
		{
			Function: "reverse",
			Args:     []interface{}{"Hello World"},
			Sleep:    2,
		},
		{
			Function: "uppercase",
			Args:     []interface{}{"Me gustan los tacos"},
			Sleep:    4,
		},
		{
			Function: "caesar",
			Args: []interface{}{
				"Lol Caesar",
				13,
			},
			Sleep: 6,
		},
	}

	shutdownChan := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < len(workers); i++ {
		wg.Add(1)
		go doWork(i, workers[i], shutdownChan, &wg)
	}

	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	// Wait for either a SIGINT signal or a timeout
	select {
	case <-sigChan:
		fmt.Println("Received SIGINT. Shutting down...")
	case <-time.After(120 * time.Second):
		fmt.Println("Timeout. Shutting down...")
	}

	// Signal the goroutines to shut down
	close(shutdownChan)

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("All goroutines have shut down")
}

func doWork(id int, w Worker, shutdownChan chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-shutdownChan:
			fmt.Printf("Goroutine %d received shutdown signal\n", id)
			return
		default:
			result := callFunctionByName(w)
			fmt.Printf("Goroutine %d: %s(%s) = %s \t - sleeping for %d seconds\n", id, w.Function, w.Args, result, w.Sleep)

			time.Sleep(time.Duration(w.Sleep) * time.Second)
		}
	}
}
