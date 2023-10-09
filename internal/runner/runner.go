package runner

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/kampanosg/bakery/internal/models"
)

const (
	HelpCmd    = "help"
	VersionCmd = "version"
	AuthorCmd  = "author"

	IgnoreFailureToken = "^"
)

var (
	cyan = color.New(color.FgCyan)
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

	var msg string

	for _, input := range args {
		switch input {
		case HelpCmd:
			msg = r.GetPrintableHelp(b)
		case VersionCmd:
			msg = r.GetPrintableVersion(b)
		case AuthorCmd:
			msg = r.GetPrintableAuthor(b)
		default:
			return r.run(b, input)
		}

		cyan.Printf("%s", msg)
	}
	return nil
}

func (r *Runner) run(b *models.Bakery, input string) error {
	return r.rrr(b, input)
}

func (r *Runner) runDefaults(b *models.Bakery) error {
	for _, d := range b.Defaults {
		if err := r.rrr(b, d); err != nil {
			return fmt.Errorf("unable to run defaults, %w", err)
		}
	}
	return nil
}

func (r *Runner) rrr(b *models.Bakery, recipe string) error {
	if b == nil {
		return fmt.Errorf("nil bakery")
	}

	rcp, ok := b.Recipes[recipe]
	if !ok {
		return fmt.Errorf("undefined recipe, %s", recipe)
	}

	if rcp.Private {
		return fmt.Errorf("recipe %s is private", recipe)
	}

	return r.runSteps(b, rcp.Steps)
}

func (r *Runner) runSteps(b *models.Bakery, steps []string) error {
	for _, step := range steps {
		ignoreFail := ignoreFail(step)
		if ignoreFail {
			step = step[1:]
		}

		step = strings.TrimSpace(step)

		cyan.Printf("-> %s\n", step)

		recipe, ok := b.Recipes[step]
		if ok {
			if err := r.runSteps(b, recipe.Steps); !ignoreFail && err != nil {
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
		buffer.WriteString(author)
	}
	buffer.WriteString("\n")
	return buffer.String()
}

func ignoreFail(step string) bool {
	return strings.HasPrefix(step, IgnoreFailureToken)
}
