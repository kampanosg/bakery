package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kampanosg/bakery/internal/parser"
)

func main() {
	file := flag.String("file", "", "custom location for a Bakefile")
	flag.Parse()

	var f *os.File
	var err error

	if *file == "" {
		f, err = LoadDefaultBakefile()
	} else {
		f, err = LoadBakefileFromLocation(*file)
	}

	if err != nil {
		fmt.Printf("unable to load the Bakefile, %v\n", err)
		return
	}
	defer f.Close()

	recipe, err := parser.ParseBakefile(f)
	if err != nil {
		fmt.Printf("unable to parse the Bakefile, %v\n", err)
		return
	}

	args := ParseArgs(os.Args)

	runner := NewDefaultRunner()
	if err := runner.RunCommand(recipe, args); err != nil {
		fmt.Printf("run failed, %v\n", err)
		return
	}
}
