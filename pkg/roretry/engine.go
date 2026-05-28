package roretry

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"
)

const (
	defaultGroupRetries = 3
	defaultMaxCycles    = 1
)

// CommandFunc executes a single command attempt.
type CommandFunc func(ctx context.Context) error

// Command defines a retry group in the execution chain.
type Command struct {
	Name       string
	MaxRetries int
	Backoff    time.Duration
	Run        CommandFunc
}

// Request defines the execution request for the retry engine.
type Request struct {
	Commands       []Command
	CommandMap     map[string]CommandFunc
	CommandTarget  any
	Config         *Config
	Start          string
	OverallCycles  int
	OverallBackoff time.Duration
	ConfigFile     string
}

// Status reports the result of running the retry engine.
type Status struct {
	Success        bool      `json:"success"`
	Group          string    `json:"group"`
	Cycle          int       `json:"cycle"`
	GroupAttempt   int       `json:"group_attempt"`
	TotalAttempts  int       `json:"total_attempts"`
	Start          string    `json:"start"`
	StartedAt      time.Time `json:"started_at"`
	FinishedAt     time.Time `json:"finished_at"`
	LastError      string    `json:"last_error,omitempty"`
	CompletedCycle int       `json:"completed_cycle"`
}

// Run executes retry groups in sequence until one command succeeds or limits are reached.
func Run(ctx context.Context, request Request) (Status, error) {
	plan, err := resolvePlan(request)
	if err != nil {
		return Status{}, err
	}

	startedAt := time.Now().UTC()
	totalAttempts := 0
	var lastErr error

	for cycle := 1; plan.maxCycles < 0 || cycle <= plan.maxCycles; cycle++ {
		startIndex := 0
		if cycle == 1 {
			startIndex = plan.startIndex
		}
		order := plan.commands[startIndex:]

		for _, cmd := range order {
			for attempt := 1; attempt <= cmd.MaxRetries; attempt++ {
				if err := ctx.Err(); err != nil {
					return Status{
						Success:        false,
						Group:          cmd.Name,
						Cycle:          cycle,
						GroupAttempt:   attempt,
						TotalAttempts:  totalAttempts,
						Start:          plan.start,
						StartedAt:      startedAt,
						FinishedAt:     time.Now().UTC(),
						LastError:      err.Error(),
						CompletedCycle: cycle - 1,
					}, err
				}

				totalAttempts++
				err := cmd.Run(ctx)
				if err == nil {
					return Status{
						Success:        true,
						Group:          cmd.Name,
						Cycle:          cycle,
						GroupAttempt:   attempt,
						TotalAttempts:  totalAttempts,
						Start:          plan.start,
						StartedAt:      startedAt,
						FinishedAt:     time.Now().UTC(),
						CompletedCycle: cycle,
					}, nil
				}

				lastErr = err
				if attempt < cmd.MaxRetries && cmd.Backoff > 0 {
					if err := waitWithContext(ctx, cmd.Backoff); err != nil {
						return Status{
							Success:        false,
							Group:          cmd.Name,
							Cycle:          cycle,
							GroupAttempt:   attempt,
							TotalAttempts:  totalAttempts,
							Start:          plan.start,
							StartedAt:      startedAt,
							FinishedAt:     time.Now().UTC(),
							LastError:      err.Error(),
							CompletedCycle: cycle - 1,
						}, err
					}
				}
			}
		}

		shouldContinue := plan.maxCycles < 0 || cycle < plan.maxCycles
		if shouldContinue && plan.overallBackoff > 0 {
			if err := waitWithContext(ctx, plan.overallBackoff); err != nil {
				return Status{
					Success:        false,
					Cycle:          cycle,
					TotalAttempts:  totalAttempts,
					Start:          plan.start,
					StartedAt:      startedAt,
					FinishedAt:     time.Now().UTC(),
					LastError:      err.Error(),
					CompletedCycle: cycle - 1,
				}, err
			}
		}
	}

	status := Status{
		Success:        false,
		Cycle:          plan.maxCycles,
		TotalAttempts:  totalAttempts,
		Start:          plan.start,
		StartedAt:      startedAt,
		FinishedAt:     time.Now().UTC(),
		CompletedCycle: plan.maxCycles,
	}
	if lastErr != nil {
		status.LastError = lastErr.Error()
	}
	if status.LastError == "" {
		status.LastError = "no command succeeded"
	}

	return status, fmt.Errorf("roretry: %s", status.LastError)
}

type resolvedPlan struct {
	commands       []Command
	startIndex     int
	start          string
	maxCycles      int
	overallBackoff time.Duration
}

func (r Request) startIndex(commands []Command, start string) int {
	if start == "" {
		return 0
	}
	for i, cmd := range commands {
		if cmd.Name == start {
			return i
		}
	}
	return -1
}

func resolvePlan(request Request) (resolvedPlan, error) {
	cfg, err := resolveConfig(request)
	if err != nil {
		return resolvedPlan{}, err
	}

	groupByName := map[string]GroupConfig{}
	for _, group := range cfg.Groups {
		groupByName[group.Name] = group
	}

	commands, err := resolveCommands(request, cfg.Groups, groupByName)
	if err != nil {
		return resolvedPlan{}, err
	}

	start := request.Start
	if start == "" {
		start = cfg.Start
	}
	if start == "" {
		start = cfg.StartGroup
	}
	startIndex := request.startIndex(commands, start)
	if startIndex < 0 {
		return resolvedPlan{}, fmt.Errorf("start group %q not found in commands", start)
	}

	maxCycles := request.OverallCycles
	if maxCycles == 0 {
		maxCycles = cfg.OverallCycles
	}
	if maxCycles == 0 {
		maxCycles = defaultMaxCycles
	}
	if maxCycles < -1 {
		return resolvedPlan{}, errors.New("overall cycles must be -1 or greater")
	}

	overallBackoff := request.OverallBackoff
	if overallBackoff <= 0 && cfg.OverallBackoff > 0 {
		overallBackoff = cfg.OverallBackoff
	}
	if overallBackoff < 0 {
		return resolvedPlan{}, errors.New("overall backoff must be zero or greater")
	}

	return resolvedPlan{
		commands:       commands,
		startIndex:     startIndex,
		start:          start,
		maxCycles:      maxCycles,
		overallBackoff: overallBackoff,
	}, nil
}

func resolveConfig(request Request) (Config, error) {
	var cfg Config
	if request.ConfigFile != "" {
		loaded, err := loadConfig(request.ConfigFile)
		if err != nil {
			return Config{}, err
		}
		cfg = loaded
	}

	if request.Config != nil {
		cfg = mergeConfig(cfg, *request.Config)
	}

	normalizeConfig(&cfg)
	return cfg, nil
}

func mergeConfig(base Config, override Config) Config {
	merged := base
	if override.Start != "" {
		merged.Start = override.Start
	}
	if override.StartGroup != "" {
		merged.StartGroup = override.StartGroup
	}
	if override.OverallCycles != 0 {
		merged.OverallCycles = override.OverallCycles
	}
	if override.OverallBackoff != 0 {
		merged.OverallBackoff = override.OverallBackoff
	}
	if override.OverallBackoffMillis != 0 {
		merged.OverallBackoffMillis = override.OverallBackoffMillis
	}
	if len(override.Groups) > 0 {
		merged.Groups = override.Groups
	}
	return merged
}

func resolveCommands(request Request, groupConfigs []GroupConfig, groupByName map[string]GroupConfig) ([]Command, error) {
	if len(request.Commands) > 0 {
		commands := make([]Command, 0, len(request.Commands))
		seenNames := map[string]struct{}{}
		for _, command := range request.Commands {
			if command.Name == "" {
				return nil, errors.New("command name is required")
			}
			if command.Run == nil {
				return nil, fmt.Errorf("command %q has nil runner", command.Name)
			}
			if _, exists := seenNames[command.Name]; exists {
				return nil, fmt.Errorf("duplicate command name %q", command.Name)
			}
			seenNames[command.Name] = struct{}{}

			applyGroupDefaults(&command, groupByName[command.Name])
			commands = append(commands, command)
		}
		return commands, nil
	}

	if len(groupConfigs) == 0 {
		return nil, errors.New("at least one command is required")
	}

	commands := make([]Command, 0, len(groupConfigs))
	seenNames := map[string]struct{}{}
	for _, group := range groupConfigs {
		if group.Name == "" {
			return nil, errors.New("config group name is required")
		}
		if _, exists := seenNames[group.Name]; exists {
			return nil, fmt.Errorf("duplicate command name %q", group.Name)
		}
		seenNames[group.Name] = struct{}{}

		runner := request.CommandMap[group.Name]
		if runner == nil {
			var err error
			runner, err = resolveMethodRunner(request.CommandTarget, group.Method)
			if err != nil {
				return nil, fmt.Errorf("resolve command %q: %w", group.Name, err)
			}
		}

		command := Command{
			Name:       group.Name,
			MaxRetries: group.MaxRetries,
			Backoff:    group.Backoff,
			Run:        runner,
		}
		applyGroupDefaults(&command, groupByName[group.Name])
		commands = append(commands, command)
	}

	return commands, nil
}

func resolveMethodRunner(target any, methodName string) (CommandFunc, error) {
	if methodName == "" {
		return nil, errors.New("missing runner: provide command map entry or config method")
	}
	if target == nil {
		return nil, fmt.Errorf("command target is required for method %q", methodName)
	}

	targetValue := reflect.ValueOf(target)
	method := targetValue.MethodByName(methodName)
	if !method.IsValid() {
		return nil, fmt.Errorf("method %q not found on command target", methodName)
	}

	contextType := reflect.TypeOf((*context.Context)(nil)).Elem()
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	methodType := method.Type()

	if methodType.NumIn() != 1 || !methodType.In(0).Implements(contextType) {
		return nil, fmt.Errorf("method %q must have signature func(context.Context) error", methodName)
	}
	if methodType.NumOut() != 1 || !methodType.Out(0).Implements(errorType) {
		return nil, fmt.Errorf("method %q must have signature func(context.Context) error", methodName)
	}

	return func(ctx context.Context) error {
		result := method.Call([]reflect.Value{reflect.ValueOf(ctx)})[0]
		if result.IsNil() {
			return nil
		}
		return result.Interface().(error)
	}, nil
}

func applyGroupDefaults(command *Command, override GroupConfig) {
	if command.MaxRetries <= 0 {
		if override.MaxRetries > 0 {
			command.MaxRetries = override.MaxRetries
		} else {
			command.MaxRetries = defaultGroupRetries
		}
	}
	if command.Backoff <= 0 && override.Backoff > 0 {
		command.Backoff = override.Backoff
	}
}

func waitWithContext(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
