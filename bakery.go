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
	fmt.Printf("%v\n", f)

	recipe, err := ParseBakefile()
	if err != nil {
		fmt.Printf("unable to parse the Bafile, %v", err)
		return
	}

	fmt.Printf("%v\n", recipe)
}
