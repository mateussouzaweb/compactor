package compactor

// Action struct
type Action struct {
	Type string
}

// Return if action is to process
func (a *Action) IsProcess() bool {
	return a.Type == "PROCESS"
}

// Return if action is to delete
func (a *Action) IsDelete() bool {
	return a.Type == "DELETE"
}
