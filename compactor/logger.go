package compactor

// Logger struct
type Logger struct {
	Processed []string
	Skipped   []string
	Ignored   []string
	Written   []string
	Deleted   []string
}

// AddProcessed append file to processed list
func (l *Logger) AddProcessed(file string) {
	l.Processed = append(l.Processed, file)
}

// AddSkipped append file to skipped list
func (l *Logger) AddSkipped(file string) {
	l.Skipped = append(l.Skipped, file)
}

// AddIgnored append file to ignored list
func (l *Logger) AddIgnored(file string) {
	l.Ignored = append(l.Ignored, file)
}

// AddWritten append file to written list
func (l *Logger) AddWritten(file string) {
	l.Written = append(l.Written, file)
}

// AddDeleted append file to deleted list
func (l *Logger) AddDeleted(file string) {
	l.Deleted = append(l.Deleted, file)
}
