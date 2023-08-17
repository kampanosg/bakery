package main

import (
	"os"
	"strings"
)

func ParseArgs(args []string) []string {
	parsed := make([]string, 0)
	for i := 1; i <= len(args[1:]); i++ {
		arg := os.Args[i]
		if strings.HasPrefix(arg, "--") {
			i += 1
			continue
		}
		parsed = append(parsed, arg)
	}

	return parsed
}
