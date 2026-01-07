package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	fmt.Println("Start Race condition")
	printDetails()

	counter := 0
	const gs = 100
	var wg sync.WaitGroup
	wg.Add(gs)

	for i := 0; i < gs; i++ {
		go func() {
			v := counter
			v++
			counter = v
			wg.Done()
		}()
	}

	wg.Wait()

	printDetails()
	fmt.Println("Counter: ", counter)
	fmt.Println("End Race condition")
}

func printDetails() {
	fmt.Println("======================")
	fmt.Println("CPUs: ", runtime.NumCPU())
	fmt.Println("Go Routines: ", runtime.NumGoroutine())
}
