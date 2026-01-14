package main

import (
	"sync"
	"testing"
)

// go test -race .
func Test_updateMessage(t *testing.T) {
	var mu sync.Mutex
	msg = "Hello, world!"

	wg.Add(2)
	go updateMessage("x", &mu)
	go updateMessage("Goodbye, cruel world!", &mu)
	wg.Wait()

	if msg != "Goodbye, cruel world!" {
		t.Errorf("Expected msg to be Goodbye, cruel world!, got %s", msg)
	}
}
