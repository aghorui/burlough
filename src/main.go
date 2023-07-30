package main

import (
	"fmt"
	"log"
	"os"
)

var exitCode int = 0

// Throws a message on a panic.
func unexpectedExitHandler() {
	if r := recover(); r != nil {
		exitCode = 2
		log.Printf("Fatal: %v\n", r)
		log.Println("Exiting")
	}
}

func main() {
	// defer os.Exit(exitCode)
	// defer unexpectedExitHandler()
	
	err := LoadConfig(os.Args)

	if err != nil {
		if err == ErrInvalidArguments {
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}


