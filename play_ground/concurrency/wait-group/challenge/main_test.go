package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_updateMessage(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	updateMessage("Hello World", &wg)
	wg.Wait()

	if msg != "Hello World" {
		t.Errorf("Expected msg to be Hello World, got %s", msg)
	}
}

func Test_printMessage(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	var wg sync.WaitGroup
	wg.Add(1)

	msg = "Hello Test"
	printMessage(&wg)
	wg.Wait()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if output != "Hello Test\n" {
		t.Errorf("Expected output to be Hello Test, got %s", output)
	}
}

func Test_main(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	expectedMessages := []string{
		"Hello, universe!",
		"Hello, cosmos!",
		"Hello, world!",
	}

	for _, msg := range expectedMessages {
		if !strings.Contains(output, msg) {
			t.Errorf("Expected output to contain %s, got %s", msg, output)
		}
	}
}
