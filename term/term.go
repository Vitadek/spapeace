package term

import (
	"os"
	"syscall"
	"unsafe"
)

var originalTermios syscall.Termios

func EnableRawMode() {
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

func DisableRawMode() {
	// Reset the terminal to the original settings
	file := os.Stdin
	fd := file.Fd()
	syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TCSETS, uintptr(unsafe.Pointer(&originalTermios)))
}

func GetKey() string {
	var buf [1]byte
	_, err := os.Stdin.Read(buf[:])
	if err != nil {
		return ""
	}
	switch buf[0] {
	case 'k', 'K':
		return "k"
	case 'j', 'J':
		return "j"
	case 'e', 'E':
		return "e"
	case 'n', 'N':
		return "n"
	case 'o', 'O':
		return "o"
	case 'r', 'R':
		return "r"
	case 'x', 'X':
		return "x"
	case 's', 'S':
		return "s"
	case 'h', 'H':
		return "h"
	case 'l', 'L':
		return "l"
	case '/':
		return "/"
	case 27: // Escape character
		return "esc"
	default:
		return ""
	}
}
