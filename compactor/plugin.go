package compactor

// Run function
type RunFunc = func(bundle *Bundle) error

// Delete function
type DeleteFunc = func(bundle *Bundle) error

// Resolve function
type ResolveFunc = func(file string) (string, error)

// Plugin struct
type Plugin struct {
	Extensions []string
	Run        RunFunc
	Delete     DeleteFunc
	Resolve    ResolveFunc
}
