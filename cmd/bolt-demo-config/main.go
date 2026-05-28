package main

import (
	"encoding/json"
	"fmt"
	"os"

	"traveler/pkg/bolt"
)

type Input struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type Output struct {
	Message string `json:"message"`
	Tenant  string `json:"tenant"`
}

func main() {
	engine, err := bolt.NewEngineFromConfigFile("engine.local.json")
	if err != nil {
		exitErr(err)
	}

	request := bolt.Request[Input]{
		TemplateName: "welcome.json.tmpl",
		Payload:      Input{Name: "Bolt", Role: "admin"},
		Model: map[string]any{
			"tenant": "caracal",
		},
	}

	output, err := bolt.Render[Input, Output](engine, request)
	if err != nil {
		exitErr(err)
	}

	inputJSON, err := json.Marshal(request.Payload)
	if err != nil {
		exitErr(err)
	}
	outputJSON, err := json.Marshal(output)
	if err != nil {
		exitErr(err)
	}

	fmt.Println()
	fmt.Println("Template name:")
	fmt.Println(request.TemplateName)
	fmt.Println()
	fmt.Println("Input object:")
	fmt.Println(string(inputJSON))
	fmt.Println()
	fmt.Println("Output object:")
	fmt.Println(string(outputJSON))
}

func exitErr(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "demo error: %v\n", err)
	os.Exit(1)
}
