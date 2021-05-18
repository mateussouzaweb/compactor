package compactor

// Action struct
type Action struct {
	Type     string
	Multiple bool
}

// Return if action is to delete one
func (a *Action) IsSingleDestinationProcess() bool {
	return a.Type == "PROCESS" && !a.Multiple
}

// Return if action is to delete many
func (a *Action) IsManyDestinationProcess() bool {
	return a.Type == "PROCESS" && a.Multiple
}

// Return if action is to delete one
func (a *Action) IsSingleDestinationDelete() bool {
	return a.Type == "DELETE" && !a.Multiple
}

// Return if action is to delete many
func (a *Action) IsManyDestinationDelete() bool {
	return a.Type == "DELETE" && a.Multiple
}
