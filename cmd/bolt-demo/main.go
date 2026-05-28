package main

import (
	"fmt"
	"os"

	"traveler/pkg/bolt"
)

func main() {
	engine := bolt.NewEngine()

	templateString := `{"message":"Hello {{.name}}","tenant":"{{.model.tenant}}","role":"{{.json.role}}"}`
	modelString := `{"tenant":"caracal"}`
	inputJSONString := `{"name":"Bolt","role":"admin"}`

	outputJSONString, err := bolt.Render[string, string](engine, bolt.Request[string]{
		Template: templateString,
		Model:    modelString,
		Payload:  inputJSONString,
	})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "demo error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Template (string):")
	fmt.Println(templateString)
	fmt.Println()
	fmt.Println("Model (string):")
	fmt.Println(modelString)
	fmt.Println()
	fmt.Println("Input JSON (string):")
	fmt.Println(inputJSONString)
	fmt.Println()
	fmt.Println("Output JSON (string):")
	fmt.Println(outputJSONString)

	textTemplate := `User {{.name}} has role {{.role}}`
	textOut, err := bolt.Render[string, string](engine, bolt.Request[string]{
		Template: textTemplate,
		Format:   bolt.FormatText,
		Payload:  inputJSONString,
	})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "demo error: %v\n", err)
		os.Exit(1)
	}

	htmlTemplate := `<p>{{.name}}</p>`
	htmlOut, err := bolt.Render[string, string](engine, bolt.Request[string]{
		Template: htmlTemplate,
		Format:   bolt.FormatHTML,
		Payload:  `{"name":"<b>Bolt</b>"}`,
	})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "demo error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("Text Template (string):")
	fmt.Println(textTemplate)
	fmt.Println("Output Text (string):")
	fmt.Println(textOut)
	fmt.Println()
	fmt.Println("HTML Template (string):")
	fmt.Println(htmlTemplate)
	fmt.Println("Output HTML (string):")
	fmt.Println(htmlOut)
}
