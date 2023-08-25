package runner

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/kampanosg/bakery/internal/models"
)

const (
	HelpCmd    = "help"
	VersionCmd = "version"
	AuthorCmd  = "author"

	IgnoreFailureToken = "^"
)

type (
	CommandAgent interface {
		Execute(cmd string) error
	}

	Runner struct {
		agent CommandAgent
	}
)

func NewRunner(e CommandAgent) *Runner {
	return &Runner{
		agent: e,
	}
}

func (r *Runner) RunCommand(b *models.Bakery, args []string) error {
	if b == nil {
		return fmt.Errorf("nil bakery")
	}

	if len(args) == 0 {
		return r.runDefaults(b)
	}

	var print string
	input := args[0]

	switch input {
	case HelpCmd:
		print = r.GetPrintableHelp(b)
	case VersionCmd:
		print = r.GetPrintableVersion(b)
	case AuthorCmd:
		print = r.GetPrintableAuthor(b)
	default:
		return r.run(b, input)
	}

	fmt.Printf(print)

	return nil
}

func (r *Runner) run(b *models.Bakery, input string) error {
	rcp, ok := b.Recipes[input]
	if !ok {
		return fmt.Errorf("undefined recipe, %s", input)
	}
	return r.runSteps(b, rcp.Steps)
}

func (r *Runner) runSteps(b *models.Bakery, steps []string) error {
	for _, step := range steps {
		ignoreFail := shouldIgnoreFailure(step)
		if ignoreFail {
			step = step[1:]
		}

		fmt.Printf("%s\n", step)

		recipe, ok := b.Recipes[step]
		if ok {
			if err := r.runSteps(b, recipe.Steps); err != nil {
				return fmt.Errorf("unable to run steps, %w", err)
			}
			continue
		}

		if err := r.agent.Execute(step); err != nil {
			if ignoreFail {
				continue
			}
			return fmt.Errorf("unable to run step %s, %w", step, err)
		}
	}
	return nil
}

func (r *Runner) runDefaults(b *models.Bakery) error {
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

func (r *Runner) GetPrintableHelp(b *models.Bakery) string {
	var buffer bytes.Buffer
	buffer.WriteString("Available Recipes in Bakefile:\n")
	for k, r := range b.Recipes {
		buffer.WriteString(fmt.Sprintf("- %s: %s\n", k, r.Description))
	}
	return buffer.String()
}

func (r *Runner) GetPrintableVersion(b *models.Bakery) string {
	return fmt.Sprintf("Bakefile Version: %s\n", b.Version)
}

func (r *Runner) GetPrintableAuthor(b *models.Bakery) string {
	var buffer bytes.Buffer
	buffer.WriteString("Bakefile Author: ")
	if author, ok := b.Metadata["author"]; ok {
		buffer.WriteString(fmt.Sprintf("%s", author))
	}
	buffer.WriteString("\n")
	return buffer.String()
}

func shouldIgnoreFailure(step string) bool {
	return strings.HasPrefix(step, IgnoreFailureToken)
}
