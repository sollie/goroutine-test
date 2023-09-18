package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode"
)

type Worker struct {
	Function string
	Input    string
	Sleep    int32
}

func main() {
	workers := []Worker{
		{"reverse", "hello", 2},
		{"uppercase", "hello", 4},
		{"uppercaseNext", "hello", 6},
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
			fmt.Printf("Goroutine %d is running\n", id)
			fmt.Printf("Sleeping for %d seconds\n", w.Sleep)
			time.Sleep(time.Duration(w.Sleep) * time.Second)
		}
	}
}

func reverse(s string) string {
	// Convert string to rune slice.
	// ... This method works on the level of runes, not bytes.
	data := []rune(s)
	result := []rune{}
	for i := len(data) - 1; i >= 0; i-- {
		result = append(result, data[i])
	}
	return string(result)
}

func uppercase(s string) string {
	return strings.ToUpper(s)
}

func uppercaseNext(s string) string {
	data := []rune(s)
	result := []rune{}
	for i := 0; i < len(data); i++ {
		if unicode.IsUpper(data[i]) {
			result = append(result, unicode.ToLower(data[i]))
			result = append(result, unicode.ToUpper(data[i+1]))
			i++
		} else {
			result = append(result, unicode.ToUpper(data[i]))
		}
	}
	return string(result)
}
