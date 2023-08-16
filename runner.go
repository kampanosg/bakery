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

func (r *Runner) RunCommand(b *Bakery, cmd string) error {
	switch cmd {
	case HelpCommand:
		fmt.Printf("help me\n")
	default:
		fmt.Printf("not found\n")
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
