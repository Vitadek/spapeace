package menu

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"spock/term"
	"spock/jsonhandler"
)

func RunMenu() error {
	options := []string{"Option 1", "Search and Parse JSON"}
	selected := 0

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		term.DisableRawMode()
		os.Exit(0)
	}()

	for {
		// Clear screen and move cursor to top-left
		fmt.Print("\033[2J\033[H")

		// Display menu options
		for i, option := range options {
			if i == selected {
				fmt.Printf("> %s\n", option)
			} else {
				fmt.Println(option)
			}
		}

		// Display instructions on the right side
		fmt.Println("\nKeys:")
		fmt.Println("j: Move down")
		fmt.Println("k: Move up")
		fmt.Println("e: Select option")
		fmt.Println("Esc: Exit")

		// Move the cursor to the bottom of the screen
		fmt.Print("\033[999B") // Move cursor to the bottom

		// Read keyboard input
		key := term.GetKey()

		switch key {
		case "k":
			selected = (selected - 1 + len(options)) % len(options)
		case "j":
			selected = (selected + 1) % len(options)
		case "e":
			if options[selected] == "Search and Parse JSON" {
				return jsonhandler.ParseJSON()
			} else {
				runSubmenu(options[selected])
			}
		case "esc":
			return nil
		}
	}
}

func runSubmenu(mainOption string) {
	submenuOptions := []string{mainOption + " SubOption 1", mainOption + " SubOption 2"}
	selected := 0

	for {
		// Clear screen and move cursor to top-left
		fmt.Print("\033[2J\033[H")

		// Display submenu options
		fmt.Println(mainOption, "submenu:")
		for i, option := range submenuOptions {
			if i == selected {
				fmt.Printf("> %s\n", option)
			} else {
				fmt.Println(option)
			}
		}

		// Display instructions on the right side
		fmt.Println("\nKeys:")
		fmt.Println("j: Move down")
		fmt.Println("k: Move up")
		fmt.Println("e: Select option")
		fmt.Println("Esc: Back to main menu")

		// Move the cursor to the bottom of the screen
		fmt.Print("\033[999B") // Move cursor to the bottom

		// Read keyboard input
		key := term.GetKey()

		switch key {
		case "k":
			selected = (selected - 1 + len(submenuOptions)) % len(submenuOptions)
		case "j":
			selected = (selected + 1) % len(submenuOptions)
		case "e":
			fmt.Printf("Selected: %s\n", submenuOptions[selected])
			fmt.Print("Press any key to return to the main menu...")
			term.GetKey() // Wait for any key press
			return
		case "esc":
			return
		}
	}
}
