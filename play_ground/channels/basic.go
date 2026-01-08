package main

import "fmt"

func main() {
	fmt.Println("Welcome to channels")

	c := make(chan int)

	go func() {
		c <- 42
	}()

	fmt.Println(<-c)
}
