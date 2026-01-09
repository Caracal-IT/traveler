package main

import "fmt"

func main() {
	fmt.Println("Welcome to channels - select")

	c := make(chan int)

	go func() {
		c <- 42
		close(c)
	}()

	v, ok := <-c
	fmt.Println(v, ok)

	v, ok = <-c
	fmt.Println(v, ok)

	if v, ok = <-c; !ok {
		fmt.Println("Value", v)
	} else {
		fmt.Println("No value")
	}
}
