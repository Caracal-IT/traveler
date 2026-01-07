package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

// RWMutex

func main() {
	fmt.Println("Start Race condition")

	var counter int64

	const gs = 100
	var wg sync.WaitGroup
	wg.Add(gs)

	for i := 0; i < gs; i++ {
		go func() {
			atomic.AddInt64(&counter, 1)
			runtime.Gosched() // Yield
			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println("Counter: ", counter)
	fmt.Println("End Race condition")
}
