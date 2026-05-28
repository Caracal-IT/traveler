package main

import (
	"encoding/json"
	"fmt"
	"os"

	"traveler/pkg/bolt"
)

type Input struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Output struct {
	Greeting string `json:"greeting"`
	Adult    bool   `json:"adult"`
	Tenant   string `json:"tenant"`
}

func main() {
	templateDir := detectTemplateDir()

	request := bolt.Request[Input]{
		TemplateName: "profile.json.tmpl",
		TemplateDir:  templateDir,
		Payload:      Input{Name: "Bolt", Age: 28},
		Model: map[string]any{
			"tenant": "caracal",
		},
	}

	output, err := bolt.RenderAs[Output](request)
	if err != nil {
		exitErr(err)
	}

	payloadJSON, err := json.Marshal(request.Payload)
	if err != nil {
		exitErr(err)
	}
	outputJSON, err := json.Marshal(output)
	if err != nil {
		exitErr(err)
	}

	fmt.Println("Template directory:")
	fmt.Println(templateDir)
	fmt.Println()
	fmt.Println("Template name:")
	fmt.Println(request.TemplateName)
	fmt.Println()
	fmt.Println("Payload object:")
	fmt.Println(string(payloadJSON))
	fmt.Println()
	fmt.Println("Output object:")
	fmt.Println(string(outputJSON))
}

func exitErr(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "demo error: %v\n", err)
	os.Exit(1)
}

func detectTemplateDir() string {
	candidates := []string{
		"templates",
		"cmd/bolt-demo-request/templates",
	}
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
	}
	return "templates"
}
