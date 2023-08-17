package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	HelpCmd    = "help"
	VersionCmd = "version"
)

type (
	CommandExecutor interface {
		Run(cmd string) error
	}

	Runner struct {
		executor CommandExecutor
	}

	DefaultExecutor struct{}
)

func NewDefaultRunner() *Runner {
	return &Runner{
		executor: &DefaultExecutor{},
	}
}

func NewRunner(e CommandExecutor) *Runner {
	return &Runner{
		executor: e,
	}
}

func (r *Runner) RunCommand(b *Bakery, args []string) error {
	if b == nil {
		return fmt.Errorf("nil bakery")
	}

	if len(args) == 0 {
		return r.runDefaults(b)
	}

	input := args[0]

	switch input {
	case HelpCmd:
		r.printHelp(b)
	case VersionCmd:
		r.printVersion(b)
	default:
		return r.run(b, input)
	}

	return nil
}

func (r *Runner) run(b *Bakery, input string) error {
	rcp, ok := b.Recipes[input]
	if !ok {
		return fmt.Errorf("undefined recipe, %s", input)
	}
	return r.runSteps(b, rcp.Steps)
}

func (r *Runner) runSteps(b *Bakery, steps []string) error {
	for _, step := range steps {
		fmt.Printf("%s\n", step)

		recipe, ok := b.Recipes[step]
		if ok {
			if err := r.runSteps(b, recipe.Steps); err != nil {
				return fmt.Errorf("unable to run steps, %w", err)
			}
			continue
		}

		if err := r.executor.Run(step); err != nil {
			return fmt.Errorf("unable to run step %s, %w", step, err)
		}
	}
	return nil
}

func (r *Runner) runDefaults(b *Bakery) error {
	for _, d := range b.Defaults {
		recipe, ok := b.Recipes[d]
		if !ok {
			return fmt.Errorf("undefined recipe %s", d)
		}

		if err := r.runSteps(b, recipe.Steps); err != nil {
			return fmt.Errorf("unable to run steps, %w", err)
		}
	}
	return nil
}

func (r *Runner) printHelp(b *Bakery) {
	fmt.Printf("Available Recipes in Bakefile:\n")
	for k, r := range b.Recipes {
		fmt.Printf("- %s: %s\n", k, r.Description)
	}
}

func (r *Runner) printVersion(b *Bakery) {
	fmt.Printf("Bakefile Version: %s\n", b.Version)
}

func (e *DefaultExecutor) Run(cmd string) error {
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return fmt.Errorf("err: %w", err)
	}

	return nil
}
