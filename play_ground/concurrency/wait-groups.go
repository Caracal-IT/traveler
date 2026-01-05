package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		fmt.Println("Hello from goroutine 1")
		time.Sleep(1 * time.Second)
		wg.Done()
	}()

	go func() {
		fmt.Println("Hello from goroutine 2")
		time.Sleep(1 * time.Second)
		wg.Done()
	}()

	wg.Wait()
}
