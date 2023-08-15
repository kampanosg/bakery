package main

import "fmt"

func main() {
	f, err := LoadBakefile()
	if err != nil {
		fmt.Printf("unable to load the Bakefile, %v", err)
		return
	}
	fmt.Printf("%v\n", f)
}
