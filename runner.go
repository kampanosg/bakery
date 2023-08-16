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

func (r *Runner) RunCommand(b *Bakery, recipe string) error {
	rcp, ok := b.Recipes[recipe]
	if !ok {
		return fmt.Errorf("recipe not found, %s", recipe)
	}

	for i, step := range rcp.Steps {
		fmt.Printf("[%d/%d] - %s\n", i, len(rcp.Steps), step)
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
