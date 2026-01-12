package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("Welcome to channels - Fan In (Rob Pike)")

	c := fanIn(boring("Joe"), boring("Ann"))

	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}

	fmt.Println("You are boring, I'm leaving")
}

func boring(msg string) <-chan string {
	c := make(chan string)

	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()

	return c
}

func fanIn(cs ...<-chan string) <-chan string {
	out := make(chan string)

	output := func(c <-chan string) {
		for {
			out <- <-c
		}
	}

	for _, c := range cs {
		go output(c)
	}

	return out
}
