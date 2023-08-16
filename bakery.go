package main

import (
	"fmt"
	"os"
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

	args := os.Args[1:]

	runner := NewDefaultRunner()
	if err := runner.RunCommand(recipe, args); err != nil {
		fmt.Printf("run failed, %v", err)
		return
	}

	fmt.Printf("done 🧁\n")
}
