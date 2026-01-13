package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Println("Welcome to channels - context")

	ctx, cancel := context.WithCancel(context.Background())

	fmt.Println("Error check 1:", ctx.Err())
	fmt.Println("Num goroutines 1:", runtime.NumGoroutine())

	go func() {
		n := 0

		for {
			select {
			case <-ctx.Done():
				return
			default:
				n++
				time.Sleep(time.Millisecond * 200)
				fmt.Println("Working", n)
			}
		}
	}()

	time.Sleep(time.Second * 2)
	fmt.Println("Error check 2:", ctx.Err())
	fmt.Println("Num goroutines 2:", runtime.NumGoroutine())

	fmt.Println("About to cancel context")
	cancel()

	fmt.Println("Error check 3:", ctx.Err())
	fmt.Println("Num goroutines 3:", runtime.NumGoroutine())
}
