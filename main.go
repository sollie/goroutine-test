package main

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
	"time"
)

type Worker struct {
	Function string
	Args     []interface{}
	Sleep    int32
}

var (
	functionMap      = make(map[string]reflect.Value)
	functionMapMutex sync.Mutex
)

func registerFunction(name string, function interface{}) {
	functionMapMutex.Lock()
	defer functionMapMutex.Unlock()

	functionMap[name] = reflect.ValueOf(function)
}

func callFunctionByName(w Worker) interface{} {
	functionMapMutex.Lock()
	defer functionMapMutex.Unlock()

	functionInfo, found := functionMap[w.Function]
	if !found {
		return nil
	}

	args := make([]reflect.Value, len(w.Args))
	for i, arg := range w.Args {
		args[i] = reflect.ValueOf(arg)
	}

	result := functionInfo.Call(args)

	if len(result) > 0 {
		return result[0].Interface()
	}

	return nil
}

func main() {
	workers := []Worker{
		{
			"reverse",
			[]interface{}{"Hello World"},
			2,
		},
		{
			"uppercase",
			[]interface{}{"Me gustan los tacos"},
			4,
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
			fmt.Printf("Goroutine %d is running\n", id)

			result := callFunctionByName(w)
			fmt.Printf("Goroutine %d: %s(%s) = %s\n", id, w.Function, w.Args, result)

			fmt.Printf("Sleeping for %d seconds\n", w.Sleep)
			time.Sleep(time.Duration(w.Sleep) * time.Second)
		}
	}
}
