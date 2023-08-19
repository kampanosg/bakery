package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kampanosg/bakery/internal/loader"
	"github.com/kampanosg/bakery/internal/parser"
	"github.com/kampanosg/bakery/internal/runner"
)

func main() {
	file := flag.String("file", "", "custom location for a Bakefile")
	flag.Parse()

	var f *os.File
	var err error

	l := loader.NewBakefileLoader(&loader.DefaultFileOpener{})

	if *file == "" {
		f, err = l.LoadDefaultBakefile()
	} else {
		f, err = l.LoadBakefileFromLocation(*file)
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

	runner := runner.NewRunner(&runner.OSAgent{})
	if err := runner.RunCommand(recipe, args); err != nil {
		fmt.Printf("run failed, %v\n", err)
		return
	}
}
