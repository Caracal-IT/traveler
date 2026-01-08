package main

import "fmt"

func main() {
	fmt.Println("Welcome to channels - directional")

	c := make(chan int)

	go func(wc chan<- int) {
		for i := 0; i < 10; i++ {
			wc <- i
		}

		close(wc)
	}(c)

	for o := range c {
		fmt.Println(o)
	}
}
