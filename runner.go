package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	HelpCommand = "help"
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
		return fmt.Errorf("not enough args provided")
	}
	recipe := args[0]

	rcp, ok := b.Recipes[recipe]
	if !ok {
		return fmt.Errorf("undefined recipe, %s", recipe)
	}

	for i, step := range rcp.Steps {
		fmt.Printf("[%d/%d] - %s\n", i+1, len(rcp.Steps), step)
		if err := r.executor.Run(step); err != nil {
			return err
		}
	}

	return nil
}

func (e *DefaultExecutor) Run(cmd string) error {
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return err
	}

	return nil
}
