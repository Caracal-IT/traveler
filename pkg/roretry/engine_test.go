package roretry

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRun_SucceedsAfterFirstGroupRetries(t *testing.T) {
	attempts := 0
	request := Request{
		Commands: []Command{
			{
				Name:       "primary",
				MaxRetries: 3,
				Run: func(ctx context.Context) error {
					attempts++
					if attempts < 3 {
						return errors.New("temporary")
					}
					return nil
				},
			},
			{
				Name: "secondary",
				Run: func(ctx context.Context) error {
					return nil
				},
			},
		},
		OverallCycles: 1,
	}

	status, err := Run(context.Background(), request)
	require.NoError(t, err)
	require.True(t, status.Success)
	require.Equal(t, "primary", status.Group)
	require.Equal(t, 3, status.GroupAttempt)
	require.Equal(t, 3, status.TotalAttempts)
}

func TestRun_MovesToNextGroupAfterRetries(t *testing.T) {
	request := Request{
		Commands: []Command{
			{
				Name:       "primary",
				MaxRetries: 3,
				Run: func(ctx context.Context) error {
					return errors.New("still failing")
				},
			},
			{
				Name: "secondary",
				Run: func(ctx context.Context) error {
					return nil
				},
			},
		},
		OverallCycles: 1,
	}

	status, err := Run(context.Background(), request)
	require.NoError(t, err)
	require.True(t, status.Success)
	require.Equal(t, "secondary", status.Group)
	require.Equal(t, 1, status.GroupAttempt)
	require.Equal(t, 4, status.TotalAttempts)
}

func TestRun_RestartsFromBeginningAfterCycle(t *testing.T) {
	primaryAttempts := 0
	secondaryAttempts := 0

	request := Request{
		Commands: []Command{
			{
				Name:       "primary",
				MaxRetries: 2,
				Run: func(ctx context.Context) error {
					primaryAttempts++
					if primaryAttempts == 3 {
						return nil
					}
					return errors.New("primary fail")
				},
			},
			{
				Name:       "secondary",
				MaxRetries: 1,
				Run: func(ctx context.Context) error {
					secondaryAttempts++
					return errors.New("secondary fail")
				},
			},
		},
		OverallCycles: 3,
	}

	status, err := Run(context.Background(), request)
	require.NoError(t, err)
	require.True(t, status.Success)
	require.Equal(t, "primary", status.Group)
	require.Equal(t, 2, status.Cycle)
	require.Equal(t, 1, status.GroupAttempt)
	require.Equal(t, 4, status.TotalAttempts)
	require.Equal(t, 3, primaryAttempts)
	require.Equal(t, 1, secondaryAttempts)
}

func TestRun_StartOnlyAffectsFirstCycle(t *testing.T) {
	primaryAttempts := 0

	request := Request{
		Commands: []Command{
			{
				Name:       "primary",
				MaxRetries: 1,
				Run: func(ctx context.Context) error {
					primaryAttempts++
					if primaryAttempts >= 2 {
						return nil
					}
					return errors.New("primary fail")
				},
			},
			{
				Name:       "secondary",
				MaxRetries: 1,
				Run: func(ctx context.Context) error {
					return errors.New("secondary fail")
				},
			},
		},
		Start:         "secondary",
		OverallCycles: 3,
	}

	status, err := Run(context.Background(), request)
	require.NoError(t, err)
	require.True(t, status.Success)
	require.Equal(t, "primary", status.Group)
	require.Equal(t, 3, status.Cycle)
	require.Equal(t, 4, status.TotalAttempts)
}

func TestRun_ConfigAppliesGroupSettings(t *testing.T) {
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, "roretry.yaml")
	config := []byte("overall_cycles: 2\ngroups:\n  - name: primary\n    max_retries: 3\n")
	require.NoError(t, os.WriteFile(configPath, config, 0o644))

	attempts := 0
	request := Request{
		ConfigFile: configPath,
		Commands: []Command{
			{
				Name: "primary",
				Run: func(ctx context.Context) error {
					attempts++
					if attempts >= 3 {
						return nil
					}
					return errors.New("retry")
				},
			},
		},
	}

	status, err := Run(context.Background(), request)
	require.NoError(t, err)
	require.True(t, status.Success)
	require.Equal(t, 3, status.GroupAttempt)
}

func TestRun_BuildsCommandsFromConfigAndMap(t *testing.T) {
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, "roretry.yaml")
	config := []byte(
		"overall_cycles: 2\nstart: primary\ngroups:\n  - name: primary\n    method: RunPrimary\n    max_retries: 3\n  - name: secondary\n    method: RunSecondary\n    max_retries: 2\n",
	)
	require.NoError(t, os.WriteFile(configPath, config, 0o644))

	primaryAttempts := 0
	secondaryAttempts := 0
	request := Request{
		ConfigFile: configPath,
		CommandMap: map[string]CommandFunc{
			"primary": func(ctx context.Context) error {
				primaryAttempts++
				if primaryAttempts >= 3 {
					return nil
				}
				return errors.New("retry")
			},
			"secondary": func(ctx context.Context) error {
				secondaryAttempts++
				return nil
			},
		},
	}

	status, err := Run(context.Background(), request)
	require.NoError(t, err)
	require.True(t, status.Success)
	require.Equal(t, "primary", status.Group)
	require.Equal(t, 3, status.GroupAttempt)
	require.Equal(t, "primary", status.Start)
	require.Equal(t, 3, primaryAttempts)
	require.Equal(t, 0, secondaryAttempts)
}

func TestRun_FailsWhenConfigGroupHasNoRunner(t *testing.T) {
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, "roretry.yaml")
	config := []byte("groups:\n  - name: primary\n")
	require.NoError(t, os.WriteFile(configPath, config, 0o644))

	request := Request{
		ConfigFile: configPath,
		CommandMap: map[string]CommandFunc{},
	}

	_, err := Run(context.Background(), request)
	require.EqualError(t, err, `resolve command "primary": missing runner: provide command map entry or config method`)

	request.CommandMap = map[string]CommandFunc{
		"secondary": func(ctx context.Context) error { return nil },
	}
	_, err = Run(context.Background(), request)
	require.EqualError(t, err, `resolve command "primary": missing runner: provide command map entry or config method`)
}

func TestRun_BuildsCommandsFromConfigAndMethods(t *testing.T) {
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, "roretry.yaml")
	config := []byte(
		"overall_cycles: 2\nstart: primary\ngroups:\n  - name: primary\n    method: RunPrimary\n    max_retries: 3\n  - name: secondary\n    method: RunSecondary\n    max_retries: 2\n",
	)
	require.NoError(t, os.WriteFile(configPath, config, 0o644))

	target := &testCommandTarget{}
	request := Request{
		ConfigFile:    configPath,
		CommandTarget: target,
	}

	status, err := Run(context.Background(), request)
	require.NoError(t, err)
	require.True(t, status.Success)
	require.Equal(t, "primary", status.Group)
	require.Equal(t, 3, status.GroupAttempt)
	require.Equal(t, "primary", status.Start)
	require.Equal(t, 3, target.primaryAttempts)
	require.Equal(t, 0, target.secondaryAttempts)
}

func TestRun_ConfigInitObjectAppliesSettings(t *testing.T) {
	attempts := 0
	request := Request{
		Config: &Config{
			OverallCycles: 2,
			Start:         "primary",
			Groups: []GroupConfig{
				{
					Name:       "primary",
					MaxRetries: 3,
				},
			},
		},
		Commands: []Command{
			{
				Name: "primary",
				Run: func(ctx context.Context) error {
					attempts++
					if attempts >= 3 {
						return nil
					}
					return errors.New("retry")
				},
			},
		},
	}

	status, err := Run(context.Background(), request)
	require.NoError(t, err)
	require.True(t, status.Success)
	require.Equal(t, "primary", status.Group)
	require.Equal(t, 3, status.GroupAttempt)
	require.Equal(t, "primary", status.Start)
}

func TestRun_ConfigInitObjectBuildsCommandsFromMethods(t *testing.T) {
	target := &testCommandTarget{}
	request := Request{
		Config: &Config{
			Start: "primary",
			Groups: []GroupConfig{
				{
					Name:       "primary",
					Method:     "RunPrimary",
					MaxRetries: 3,
				},
				{
					Name:       "secondary",
					Method:     "RunSecondary",
					MaxRetries: 2,
				},
			},
		},
		CommandTarget: target,
	}

	status, err := Run(context.Background(), request)
	require.NoError(t, err)
	require.True(t, status.Success)
	require.Equal(t, "primary", status.Group)
	require.Equal(t, 3, status.GroupAttempt)
	require.Equal(t, "primary", status.Start)
	require.Equal(t, 3, target.primaryAttempts)
	require.Equal(t, 0, target.secondaryAttempts)
}

func TestRun_FailsWhenConfiguredMethodMissingOrInvalid(t *testing.T) {
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, "roretry.yaml")
	config := []byte("groups:\n  - name: primary\n    method: DoesNotExist\n")
	require.NoError(t, os.WriteFile(configPath, config, 0o644))

	request := Request{
		ConfigFile:    configPath,
		CommandTarget: &testCommandTarget{},
	}
	_, err := Run(context.Background(), request)
	require.EqualError(t, err, `resolve command "primary": method "DoesNotExist" not found on command target`)

	config = []byte("groups:\n  - name: primary\n    method: InvalidSignature\n")
	require.NoError(t, os.WriteFile(configPath, config, 0o644))
	_, err = Run(context.Background(), request)
	require.EqualError(t, err, `resolve command "primary": method "InvalidSignature" must have signature func(context.Context) error`)
}

func TestRun_FailsOnOverallCycleLimit(t *testing.T) {
	request := Request{
		Commands: []Command{
			{
				Name:       "primary",
				MaxRetries: 1,
				Run: func(ctx context.Context) error {
					return errors.New("always fails")
				},
			},
		},
		OverallCycles: 2,
	}

	status, err := Run(context.Background(), request)
	require.Error(t, err)
	require.False(t, status.Success)
	require.Equal(t, 2, status.CompletedCycle)
	require.Equal(t, 2, status.TotalAttempts)
}

func TestRun_AppliesOverallBackoffBetweenCycles(t *testing.T) {
	request := Request{
		Commands: []Command{
			{
				Name:       "primary",
				MaxRetries: 1,
				Run: func(ctx context.Context) error {
					return errors.New("always fails")
				},
			},
		},
		OverallCycles:  2,
		OverallBackoff: 25 * time.Millisecond,
	}

	start := time.Now()
	_, err := Run(context.Background(), request)
	elapsed := time.Since(start)
	require.Error(t, err)
	require.GreaterOrEqual(t, elapsed, 20*time.Millisecond)
}

func TestRun_RespectsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	request := Request{
		Commands: []Command{
			{
				Name: "primary",
				Run: func(ctx context.Context) error {
					return nil
				},
			},
		},
	}

	status, err := Run(ctx, request)
	require.Error(t, err)
	require.False(t, status.Success)
}

func TestRun_ConfigAppliesOverallBackoff(t *testing.T) {
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, "roretry.yaml")
	config := []byte("overall_cycles: 2\noverall_backoff_millis: 25\n")
	require.NoError(t, os.WriteFile(configPath, config, 0o644))

	request := Request{
		ConfigFile: configPath,
		Commands: []Command{
			{
				Name:       "primary",
				MaxRetries: 1,
				Run: func(ctx context.Context) error {
					return errors.New("always fails")
				},
			},
		},
	}

	start := time.Now()
	_, err := Run(context.Background(), request)
	elapsed := time.Since(start)
	require.Error(t, err)
	require.GreaterOrEqual(t, elapsed, 20*time.Millisecond)
}

func BenchmarkRun_SingleSuccess(b *testing.B) {
	request := Request{
		Commands: []Command{
			{
				Name:       "primary",
				MaxRetries: 1,
				Backoff:    time.Millisecond,
				Run: func(ctx context.Context) error {
					return nil
				},
			},
		},
		OverallCycles: 1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Run(context.Background(), request)
		if err != nil {
			b.Fatal(err)
		}
	}
}

type testCommandTarget struct {
	primaryAttempts   int
	secondaryAttempts int
}

func (t *testCommandTarget) RunPrimary(ctx context.Context) error {
	t.primaryAttempts++
	if t.primaryAttempts >= 3 {
		return nil
	}
	return errors.New("retry")
}

func (t *testCommandTarget) RunSecondary(ctx context.Context) error {
	t.secondaryAttempts++
	if t.secondaryAttempts >= 2 {
		return nil
	}
	return errors.New("retry")
}

func (t *testCommandTarget) InvalidSignature() error {
	return nil
}
