package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var msg string
var wg sync.WaitGroup

func updateMessage(newMsg string, m *sync.Mutex) {
	defer wg.Done()

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

	m.Lock()
	msg = newMsg
	m.Unlock()
}

func main() {
	fmt.Println("Mutex example")

	msg = "Hello World"
	var mu sync.Mutex

	wg.Add(2)
	go updateMessage("Hallo Universe", &mu)
	go updateMessage("Hello Gopher", &mu)
	wg.Wait()

	fmt.Println("Final message:", msg)
}
