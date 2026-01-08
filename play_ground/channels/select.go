package main

import "fmt"

func main() {
	fmt.Println("Welcome to channels - select")

	even := make(chan int)
	odd := make(chan int)
	quit := make(chan int)

	go send(even, odd, quit)

	receive(even, odd, quit)
}

func send(e, o chan<- int, q chan<- int) {
	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			e <- i
		} else {
			o <- i
		}
	}

	q <- 0

	close(e)
	close(o)
	close(q)
}

func receive(e, o <-chan int, q <-chan int) {
	for {
		select {
		case v := <-e:
			fmt.Println("Even:", v)
		case v := <-o:
			fmt.Println("Odd:", v)
		case <-q:
			fmt.Println("Quitting...")
			return
		}
	}

}
