package compactor

// Init function
type InitFunc = func(bundle *Bundle) error

// Related function
type RelatedFunc = func(item *Item) ([]Related, error)

// Execute function
type ExecuteFunc = func(bundle *Bundle) error

// Optimize function
type OptimizeFunc = func(bundle *Bundle) error

// Delete function
type DeleteFunc = func(bundle *Bundle) error

// Resolve function
type ResolveFunc = func(file string) (string, error)

// Plugin struct
type Plugin struct {
	Namespace  string
	Extensions []string
	Init       InitFunc
	Related    RelatedFunc
	Execute    ExecuteFunc
	Optimize   OptimizeFunc
	Delete     DeleteFunc
	Resolve    ResolveFunc
}
