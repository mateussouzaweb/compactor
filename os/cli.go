package os

import (
	"fmt"
	"os"
)

// Colors
var (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Purple  = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
	Notice  = Cyan
	Warn    = Yellow
	Fatal   = Red
	Success = Green
)

// NoColor check if should avoid color output on console
func NoColor() bool {
	return os.Getenv("NO_COLOR") != "" || os.Getenv("CLICOLOR") == "0"
}

// Printf displays a info to standard output
func Printf(color string, format string, args ...interface{}) {
	if NoColor() {
		fmt.Printf(format, args...)
	} else {
		fmt.Printf(color+format+Reset, args...)
	}
}
