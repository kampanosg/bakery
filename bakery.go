package main

import (
	"fmt"
)

func main() {
	f, err := LoadBakefile()
	if err != nil {
		fmt.Printf("unable to load the Bakefile, %v", err)
		return
	}
	defer f.Close()

	recipe, err := ParseBakefile(f)
	if err != nil {
		fmt.Printf("unable to parse the Bakefile, %v", err)
		return
	}

	cmd := "help"
	runner := NewDefaultRunner()
	if err := runner.RunCommand(recipe, cmd); err != nil {
		fmt.Printf("command failed, %v", err)
		return
	}

	fmt.Printf("done 🧁\n")
}
