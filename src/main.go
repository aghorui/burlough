package main

import (
	"fmt"
	"os"
)


func main() {
	err := LoadConfig(os.Args)

	if err != nil {
		if err == ErrInvalidArguments {
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stderr, "[Error] %v\n", err)
			os.Exit(1)
		}
	}
}


