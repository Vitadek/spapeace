package menu

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"spock/config"
	"spock/jsonhandler"
	"spock/term"
	"strings"
	"syscall"
)

func RunMenu() error {
	options := []string{"Option 1", "Search and Parse JSON", "Configure Program"}
	selected := 0

	// Load the configuration
	err := config.LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

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
			} else if options[selected] == "Configure Program" {
				runConfigSubmenu()
			} else {
				runSubmenu(options[selected])
			}
		case "esc":
			return nil
		}
	}
}

func runConfigSubmenu() {
	cfg := config.GetConfig()
	options := []string{"Ansible Vault Location", "Ansible Vault Password", "Host Limitations"}
	selected := 0

	for {
		// Clear screen and move cursor to top-left
		fmt.Print("\033[2J\033[H")

		// Display current configurations and submenu options
		fmt.Printf("Current Configuration:\n")
		fmt.Printf("Ansible Vault Location: %s\n", cfg.AnsibleVaultLocation)
		fmt.Printf("Ansible Vault Password: %s\n", cfg.AnsibleVaultPassword)
		fmt.Printf("Host Limitations: %s\n\n", cfg.HostLimitations)

		fmt.Println("Options:")
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
		fmt.Println("e: Edit option")
		fmt.Println("Esc: Back to main menu")

		// Move the cursor to the bottom of the screen
		fmt.Print("\033[999B")

		// Read keyboard input
		key := term.GetKey()

		switch key {
		case "k":
			selected = (selected - 1 + len(options)) % len(options)
		case "j":
			selected = (selected + 1) % len(options)
		case "e":
			editConfigOption(options[selected], cfg)
			config.SaveConfig()
		case "esc":
			return
		}
	}
}

func editConfigOption(option string, cfg *config.Config) {
	term.DisableRawMode()      // Temporarily disable raw mode for input
	defer term.EnableRawMode() // Re-enable raw mode after getting the input

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Current %s: ", option)
	switch option {
	case "Ansible Vault Location":
		fmt.Println(cfg.AnsibleVaultLocation)
	case "Ansible Vault Password":
		fmt.Println(cfg.AnsibleVaultPassword)
	case "Host Limitations":
		fmt.Println(cfg.HostLimitations)
	}

	fmt.Printf("Enter new %s: ", option)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input) // Remove surrounding whitespace, including newline character

	switch option {
	case "Ansible Vault Location":
		cfg.AnsibleVaultLocation = input
	case "Ansible Vault Password":
		cfg.AnsibleVaultPassword = input
	case "Host Limitations":
		cfg.HostLimitations = input
	}

	fmt.Println("Value changed, press any key to continue...")
	fmt.Scanln() // Wait for input to resume raw mode
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

