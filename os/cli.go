package os

import "fmt"

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

// Print a info to standard output
func Printf(color string, format string, args ...interface{}) {
	fmt.Printf(color+format+Reset, args...)
}
