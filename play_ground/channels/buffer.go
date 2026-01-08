package main

import "fmt"

func main() {
	fmt.Println("Welcome to channels - buffered")

	c := make(chan int, 2)

	c <- 42
	c <- 43

	fmt.Println(<-c)
	fmt.Println(<-c)
}
