package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"traveler/pkg/roretry"
)

func main() {
	configPath := detectConfigPath()

	runInlineDemo()
	runConfigMapDemo(configPath)
	runConfigMethodDemo(configPath)
	runConfigInitObjectDemo()
}

func runInlineDemo() {
	var primaryAttempts int
	var secondaryAttempts int

	request := roretry.Request{
		Start: "",
		Commands: []roretry.Command{
			{
				Name: "primary",
				Run: func(ctx context.Context) error {
					primaryAttempts++
					if primaryAttempts >= 3 {
						return nil
					}
					return fmt.Errorf("primary temporary failure attempt=%d", primaryAttempts)
				},
			},
			{
				Name:       "secondary",
				MaxRetries: 2,
				Backoff:    200 * time.Millisecond,
				Run: func(ctx context.Context) error {
					secondaryAttempts++
					if secondaryAttempts >= 2 {
						return nil
					}
					return fmt.Errorf("secondary temporary failure attempt=%d", secondaryAttempts)
				},
			},
		},
		OverallCycles: 2,
	}

	status, err := roretry.Run(context.Background(), request)
	printStatus("inline command list status", status, err)
	fmt.Printf("primary attempts: %d\n", primaryAttempts)
	fmt.Printf("secondary attempts: %d\n", secondaryAttempts)
}

func runConfigMapDemo(configPath string) {
	var primaryAttempts int
	var secondaryAttempts int

	commandMap := map[string]roretry.CommandFunc{
		"primary": func(ctx context.Context) error {
			primaryAttempts++
			if primaryAttempts >= 3 {
				return nil
			}
			return fmt.Errorf("primary temporary failure attempt=%d", primaryAttempts)
		},
		"secondary": func(ctx context.Context) error {
			secondaryAttempts++
			if secondaryAttempts >= 2 {
				return nil
			}
			return fmt.Errorf("secondary temporary failure attempt=%d", secondaryAttempts)
		},
	}

	request := roretry.Request{
		ConfigFile:    configPath,
		CommandMap:    commandMap,
		OverallCycles: 0, // use config
	}

	status, err := roretry.Run(context.Background(), request)
	printStatus("config + command-map status", status, err)
	fmt.Printf("primary attempts (config): %d\n", primaryAttempts)
	fmt.Printf("secondary attempts (config): %d\n", secondaryAttempts)
}

func runConfigMethodDemo(configPath string) {
	target := &demoCommandTarget{}

	request := roretry.Request{
		ConfigFile:    configPath,
		CommandTarget: target,
		OverallCycles: 0, // use config
	}

	status, err := roretry.Run(context.Background(), request)
	printStatus("config + method-name status", status, err)
	fmt.Printf("primary attempts (methods): %d\n", target.primaryAttempts)
	fmt.Printf("secondary attempts (methods): %d\n", target.secondaryAttempts)
}

func runConfigInitObjectDemo() {
	target := &demoCommandTarget{}
	request := roretry.Request{
		Config: &roretry.Config{
			Start:         "primary",
			OverallCycles: 4,
			Groups: []roretry.GroupConfig{
				{
					Name:          "primary",
					Method:        "RunPrimary",
					MaxRetries:    3,
					BackoffMillis: 100,
				},
				{
					Name:          "secondary",
					Method:        "RunSecondary",
					MaxRetries:    2,
					BackoffMillis: 100,
				},
			},
		},
		CommandTarget: target,
	}

	status, err := roretry.Run(context.Background(), request)
	printStatus("config init-object + method-name status", status, err)
	fmt.Printf("primary attempts (init object): %d\n", target.primaryAttempts)
	fmt.Printf("secondary attempts (init object): %d\n", target.secondaryAttempts)
}

func printStatus(title string, status roretry.Status, runErr error) {
	if runErr != nil {
		_, _ = fmt.Fprintf(os.Stderr, "roretry demo error (%s): %v\n", title, runErr)
	}

	statusJSON, marshalErr := json.MarshalIndent(status, "", "  ")
	if marshalErr != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to marshal status: %v\n", marshalErr)
		os.Exit(1)
	}

	fmt.Printf("%s:\n", title)
	fmt.Println(string(statusJSON))
}

type demoCommandTarget struct {
	primaryAttempts   int
	secondaryAttempts int
}

func (d *demoCommandTarget) RunPrimary(ctx context.Context) error {
	d.primaryAttempts++
	if d.primaryAttempts >= 3 {
		return nil
	}
	return fmt.Errorf("primary temporary failure attempt=%d", d.primaryAttempts)
}

func (d *demoCommandTarget) RunSecondary(ctx context.Context) error {
	d.secondaryAttempts++
	if d.secondaryAttempts >= 2 {
		return nil
	}
	return fmt.Errorf("secondary temporary failure attempt=%d", d.secondaryAttempts)
}

func detectConfigPath() string {
	candidates := []string{
		"cmd/roretry-demo/roretry.yaml",
		"roretry.yaml",
	}
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate
		}
	}
	return ""
}
