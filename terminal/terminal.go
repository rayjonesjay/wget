package terminal

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// termios related constants
const (
	TCGETS = 0x5401
	TCSETS = 0x5402
)

// TIOCGWINSZ IOCTL command for getting cursor position
const TIOCGWINSZ = 0x5413

// termios struct matching Linux x86_64
type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Line   uint8
	Cc     [19]uint8
	Pad    uint8
	Ispeed uint32
	Ospeed uint32
}

// winsize struct for getting terminal dimensions
type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// Terminal represents a virtual terminal screen.
// Use Init to create an instance.
// Ideally, only one instance should be created;
// having more instances creates a race condition as there is only one os.Stdout
type Terminal struct {
	originalTermios termios
	rows            int
	cols            int
	originalRow     int
	originalCol     int
}

// getTermios gets the current terminal attributes
func getTermios(fd uintptr) (*termios, error) {
	var term termios
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		uintptr(TCGETS),
		uintptr(unsafe.Pointer(&term)),
	)
	if errno != 0 {
		return nil, errno
	}
	return &term, nil
}

// setTermios sets the terminal attributes
func setTermios(fd uintptr, term *termios) error {
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		uintptr(TCSETS),
		uintptr(unsafe.Pointer(term)),
	)
	if errno != 0 {
		return errno
	}
	return nil
}

// getWinsize gets terminal dimensions
func getWinsize(fd uintptr) (*winsize, error) {
	var ws winsize
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		uintptr(TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)),
	)
	if errno != 0 {
		return nil, errno
	}
	return &ws, nil
}

// Size gets the terminal size, width and height, respectively, by making a
// system call. If the terminal size cannot be polled for from the system, it
// returns a default terminal size
func Size() (width, height int) {
	width, height = 35, 150
	ws, err := getWinsize(os.Stdout.Fd())
	if err != nil {
		return
	}
	width, height = int(ws.Col), int(ws.Row)-1
	if width <= 0 {
		width = 35
	}

	if height <= 0 {
		height = 160
	}
	return
}

// Init initializes a new virtual screen while preserving terminal history
func Init() (*Terminal, error) {
	// Get file descriptor for stdin
	fd := os.Stdin.Fd()

	// Get and save original terminal attributes
	originalTerm, err := getTermios(fd)
	if err != nil {
		return nil, fmt.Errorf("failed to get terminal attributes: %v", err)
	}

	// Create a copy for modification
	rawTerm := *originalTerm

	// Modify terminal attributes for raw mode while preserving original flags
	// Clear specific flags we want to modify
	rawTerm.Iflag &^= syscall.IXON | syscall.IXOFF | syscall.ICRNL | syscall.INLCR | syscall.IGNCR | syscall.IUCLC
	rawTerm.Lflag &^= syscall.ICANON | syscall.ECHO | syscall.ISIG | syscall.IEXTEN
	rawTerm.Oflag &^= syscall.OPOST

	// Set desired flags
	rawTerm.Iflag |= syscall.ICRNL

	// Apply the modified attributes
	if err := setTermios(fd, &rawTerm); err != nil {
		return nil, fmt.Errorf("failed to set terminal attributes: %v", err)
	}

	// Get terminal size
	ws, err := getWinsize(fd)
	if err != nil {
		return nil, fmt.Errorf("failed to get terminal size: %v", err)
	}

	// Save current cursor position
	fmt.Print("\033[6n")
	var originalRow, originalCol int
	_, _ = fmt.Scanf("\033[%d;%dR", &originalRow, &originalCol)

	// Switch to alternate screen buffer
	fmt.Print("\033[?1049h")

	// Clear the alternate screen
	fmt.Print("\033[2J")

	// Move cursor to top-left
	fmt.Print("\033[H")

	return &Terminal{
		originalTermios: *originalTerm,
		rows:            int(ws.Row),
		cols:            int(ws.Col),
		originalRow:     originalRow,
		originalCol:     originalCol,
	}, nil
}

// EndWin restores the original terminal state and screen content
func (s *Terminal) EndWin() error {
	// Switch back to main screen buffer
	fmt.Print("\033[?1049l")

	// Restore original terminal attributes
	if err := setTermios(os.Stdin.Fd(), &s.originalTermios); err != nil {
		return fmt.Errorf("failed to restore terminal attributes: %v", err)
	}

	// Restore cursor position
	fmt.Printf("\033[%d;%dH", s.originalRow, s.originalCol)

	return nil
}

// PrintAt prints text at specific coordinates
// Rows and columns are zero-based
func (s *Terminal) PrintAt(row, col int, text string) {
	if row >= 0 && row < s.rows && col >= 0 && col < s.cols {
		fmt.Printf("\033[%d;%dH%s", row+1, col+1, text)
	}
}

// Print prints text at specific row in the terminal.
// Rows and columns are zero-based
func (s *Terminal) Print(row int, text string) {
	s.PrintAt(row, 0, text)
}

// Clear clears the current screen buffer
func (s *Terminal) Clear() {
	fmt.Print("\033[2J\033[H")
}
