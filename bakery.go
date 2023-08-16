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

	fmt.Printf("%v\n", recipe.Version)
}
