package main

import "fmt"

func main() {
	fmt.Println("Welcome to channels - select")

	even := make(chan int)
	odd := make(chan int)
	quit := make(chan bool)

	go send2(even, odd, quit)

	receive2(even, odd, quit)
}

func send2(e, o chan<- int, q chan<- bool) {
	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			e <- i
		} else {
			o <- i
		}
	}

	close(q)
}

func receive2(e, o <-chan int, q <-chan bool) {
	for {
		select {
		case v := <-e:
			fmt.Println("Even:", v)
		case v := <-o:
			fmt.Println("Odd:", v)
		case i, ok := <-q:
			if !ok {
				fmt.Println("From comma ok", i, ok)
				return
			} else {
				fmt.Println("From comma ok", i)
			}
		}
	}
}
