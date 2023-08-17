package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
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

	recipe, err := ParseBakefile(f)
	if err != nil {
		fmt.Printf("unable to parse the Bakefile, %v\n", err)
		return
	}

	args := make([]string, 0)
	for i := 1; i <= len(os.Args[1:]); i++ {
		arg := os.Args[i]
		if strings.HasPrefix(arg, "--") {
			i += 1
			continue
		}
		args = append(args, arg)
	}

	runner := NewDefaultRunner()
	if err := runner.RunCommand(recipe, args); err != nil {
		fmt.Printf("run failed, %v\n", err)
		return
	}
}
