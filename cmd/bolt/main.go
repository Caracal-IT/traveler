package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"traveler/pkg/bolt"
)

func main() {
	var templateInline string
	var templateFile string
	var templateName string
	var templateDir string
	var configFile string
	var dataInline string
	var dataFile string
	var modelFile string
	var outFile string

	flag.StringVar(&templateInline, "template", "", "inline JSON template content")
	flag.StringVar(&templateFile, "template-file", "", "path to JSON template file")
	flag.StringVar(&templateName, "template-name", "", "template file name resolved from -template-dir")
	flag.StringVar(&templateDir, "template-dir", "", "template directory for -template-name resolution")
	flag.StringVar(&configFile, "config-file", "", "path to engine JSON config file (optional)")
	flag.StringVar(&dataInline, "data", "", "inline JSON input payload")
	flag.StringVar(&dataFile, "data-file", "", "path to JSON input payload file")
	flag.StringVar(&modelFile, "model-file", "", "path to JSON model file (optional)")
	flag.StringVar(&outFile, "out-file", "", "write output JSON to file (optional)")
	flag.Parse()

	if templateFile != "" {
		if templateInline != "" || templateName != "" {
			exitErr(errors.New("use only one of -template, -template-file, or -template-name"))
		}
		fileTemplate, err := readTemplateFile(templateFile)
		if err != nil {
			exitErr(err)
		}
		templateInline = fileTemplate
	}

	if templateInline == "" && templateName == "" {
		exitErr(errors.New("one of -template, -template-file, or -template-name is required"))
	}

	dataPayload, err := readPayload(dataInline, dataFile)
	if err != nil {
		exitErr(err)
	}

	modelPayload, err := readPayload("", modelFile)
	if err != nil {
		exitErr(err)
	}

	engineOptions := []bolt.Option{}
	if templateDir != "" {
		engineOptions = append(engineOptions, bolt.WithTemplateDir(templateDir))
	}
	var engine *bolt.Engine
	if configFile != "" {
		engine, err = bolt.NewEngineFromConfigFile(configFile, engineOptions...)
		if err != nil {
			exitErr(err)
		}
	} else {
		engine = bolt.NewEngine(engineOptions...)
	}
	rendered, err := bolt.Render[any, []byte](engine, bolt.Request[any]{
		Template:     templateInline,
		TemplateName: templateName,
		Payload:      dataPayload,
		Model:        modelPayload,
	})
	if err != nil {
		exitErr(err)
	}

	if outFile != "" {
		if err := os.WriteFile(outFile, rendered, 0o644); err != nil {
			exitErr(fmt.Errorf("write output file: %w", err))
		}
		return
	}

	fmt.Println(string(rendered))
}

func readTemplateFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("read template file: %w", err)
	}
	return string(data), nil
}

func readPayload(inline string, file string) (any, error) {
	if inline != "" && file != "" {
		return nil, errors.New("use either inline JSON or file JSON, not both")
	}
	if inline == "" && file == "" {
		return nil, nil
	}
	if inline != "" {
		return inline, nil
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read json file: %w", err)
	}
	return data, nil
}

func exitErr(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "bolt error: %v\n", err)
	os.Exit(1)
}
