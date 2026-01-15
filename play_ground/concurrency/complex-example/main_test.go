package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func Test_main(t *testing.T) {
	stdOut := os.Stdout

	defer func() {
		os.Stdout = stdOut
	}()

	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	if !strings.Contains(output, "$34320.00") {
		t.Errorf("Expected output to be '%s', got '%s'", "$34320.00", string(result))
	}
}
