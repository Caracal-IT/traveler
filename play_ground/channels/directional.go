package main

import "fmt"

func main() {
	fmt.Println("Welcome to channels - directional")

	c := make(chan int)

	output := func(rc <-chan int) {
		fmt.Println(<-rc)
	}

	go func(wc chan<- int) {
		for i := 0; i < 10; i++ {
			wc <- i
		}
	}(c)

	for i := 0; i < 10; i++ {
		output(c)
	}
}
