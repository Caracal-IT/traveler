package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	fmt.Println("Welcome to channels - Fan Out")

	c1 := make(chan int)
	c2 := make(chan int)

	go populate2(c1)
	go fanOutIn2(c1, c2)

	for v := range c2 {
		fmt.Println(v)
	}

	fmt.Println("Done")
}

func populate2(c chan int) {
	for i := 0; i < 100; i++ {
		c <- i
	}
	close(c)
}

func fanOutIn2(c1, c2 chan int) {
	var wg sync.WaitGroup
	const goroutines = 10
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			for n := range c1 {
				c2 <- timeToConsumingWork2(n)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	close(c2)
}

func timeToConsumingWork2(n int) int {
	time.Sleep(time.Microsecond * time.Duration(rand.Intn(500)))
	return n + rand.Intn(1000)
}
