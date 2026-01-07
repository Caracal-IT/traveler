package main

import (
	"fmt"
	"runtime"
	"sync"
)

// RWMutex

func main() {
	fmt.Println("Start Race condition")

	counter := 0
	const gs = 100
	var wg sync.WaitGroup
	wg.Add(gs)

	var mu sync.Mutex

	for i := 0; i < gs; i++ {
		go func() {
			mu.Lock()
			v := counter
			runtime.Gosched() // Yield
			v++
			counter = v
			mu.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println("Counter: ", counter)
	fmt.Println("End Race condition")
}
