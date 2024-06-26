package main

import (
	"fmt"
	"os"
	"spock/term"
	"spock/menu"
)

func main() {
	// Enable raw mode
	term.EnableRawMode()
	defer term.DisableRawMode()

	err := menu.RunMenu()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
