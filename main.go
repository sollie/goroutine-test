package main

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode"
)

type Worker struct {
	Function string
	Input    []reflect.Value
	Sleep    int32
}

var (
	functionMap = map[string]interface{}{
		"reverse":   reverse,
		"uppercase": uppercase,
		"caesar":    caesar,
	}
)

func main() {
	workers := []Worker{
		{
			"reverse",
			[]reflect.Value{reflect.ValueOf("Hello World")},
			2,
		},
		{
			"uppercase",
			[]reflect.Value{reflect.ValueOf("Me gustan los tacos")},
			4,
		},
		{
			"caesar",
			[]reflect.Value{
				reflect.ValueOf("Lol Caesar"),
				reflect.ValueOf(13),
			},
			6,
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
			fmt.Printf("Goroutine %d is running\n", id)
			fn, ok := functionMap[w.Function]
			if !ok {
				fmt.Printf("Function %s not found\n", w.Function)
				return
			}
			result := reflect.ValueOf(fn).Call(w.Input)
			fmt.Printf("Goroutine %d: %s(%s) = %s\n", id, w.Function, w.Input, result)
			fmt.Printf("Sleeping for %d seconds\n", w.Sleep)
			time.Sleep(time.Duration(w.Sleep) * time.Second)
		}
	}
}

func reverse(s string) string {
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

func caesar(input string, shift int) string {
	runes := []rune(input)
	shifted := make([]rune, len(runes))

	for i, char := range runes {
		if unicode.IsLetter(char) {
			var base rune
			if unicode.IsUpper(char) {
				base = 'A'
			} else {
				base = 'a'
			}
			shifted[i] = (char-base+rune(shift))%26 + base
		} else {
			shifted[i] = char
		}
	}

	return string(shifted)
}
