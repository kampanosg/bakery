package main

import (
	"flag"
	"os"

	"github.com/fatih/color"
	"github.com/kampanosg/bakery/internal/loader"
	"github.com/kampanosg/bakery/internal/parser"
	"github.com/kampanosg/bakery/internal/runner"
)

var (
	red   = color.New(color.FgRed)
	green = color.New(color.FgGreen)
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
		red.Printf("unable to load the Bakefile, %v\n", err)
		return
	}
	defer f.Close()

	recipe, err := parser.ParseBakefile(f)
	if err != nil {
		red.Printf("unable to parse the Bakefile, %v\n", err)
		return
	}

	args := parser.ParseArgs(os.Args)

	r := runner.NewRunner(&runner.OSAgent{})
	if err := r.Run(recipe, args); err != nil {
		red.Printf("run failed, %v\n", err)
		return
	}

	green.Println("done!")
}
