package main

import (
	"reflect"
	"sync"
)

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
