package runner

import (
	"fmt"
	"os"
	"os/exec"
)

type OSAgent struct{}

func (e *OSAgent) Run(cmd string) error {
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return fmt.Errorf("err: %w", err)
	}

	return nil
}
