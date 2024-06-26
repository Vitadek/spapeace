package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

var originalTermios syscall.Termios

func main() {
	// Enable raw mode
	enableRawMode()
	defer disableRawMode()

	options := []string{"Option 1", "Option 2", "Option 3"}
	selected := 0

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		disableRawMode()
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
		key := getKey()

		switch key {
		case "k": // Up
			selected = (selected - 1 + len(options)) % len(options)
		case "j": // Down
			selected = (selected + 1) % len(options)
		case "e": // Enter
			runSubmenu(options[selected])
		case "esc":
			return
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
		key := getKey()

		switch key {
		case "k": // Up
			selected = (selected - 1 + len(submenuOptions)) % len(submenuOptions)
		case "j": // Down
			selected = (selected + 1) % len(submenuOptions)
		case "e": // Enter
			fmt.Printf("Selected: %s\n", submenuOptions[selected])
			fmt.Print("Press any key to return to the main menu...")
			getKey() // Wait for any key press
			return
		case "esc":
			return
		}
	}
}

func enableRawMode() {
	// Get the current terminal settings
	file := os.Stdin
	fd := file.Fd()
	syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TCGETS, uintptr(unsafe.Pointer(&originalTermios)))

	// Make a copy of the original termios and modify it
	raw := originalTermios
	raw.Iflag &^= syscall.ICRNL | syscall.IXON
	raw.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.ISIG
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0

	// Apply the new terminal settings
	syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TCSETS, uintptr(unsafe.Pointer(&raw)))
}

func disableRawMode() {
	// Reset the terminal to the original settings
	file := os.Stdin
	fd := file.Fd()
	syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TCSETS, uintptr(unsafe.Pointer(&originalTermios)))
}

func getKey() string {
	var buf [3]byte
	_, err := os.Stdin.Read(buf[:])
	if err != nil {
		return ""
	}
	if buf[0] == 'k' || buf[0] == 'K' {
		return "k"
	} else if buf[0] == 'j' || buf[0] == 'J' {
		return "j"
	} else if buf[0] == 'e' || buf[0] == 'E' {
		return "e"
	} else if buf[0] == 27 {
		return "esc"
	}
	return ""
}
