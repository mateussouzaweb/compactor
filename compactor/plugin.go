package compactor

// Init function
type InitFunc = func(bundle *Bundle) error

// Related function
type RelatedFunc = func(item *Item) ([]Related, error)

// Resolve function
type ResolveFunc = func(file string, item *Item) (string, error)

// Execute function
type ExecuteFunc = func(bundle *Bundle) error

// Optimize function
type OptimizeFunc = func(bundle *Bundle) error

// Delete function
type DeleteFunc = func(bundle *Bundle) error

// Plugin struct
type Plugin struct {
	Namespace  string
	Extensions []string
	Init       InitFunc
	Related    RelatedFunc
	Resolve    ResolveFunc
	Execute    ExecuteFunc
	Optimize   OptimizeFunc
	Delete     DeleteFunc
}
