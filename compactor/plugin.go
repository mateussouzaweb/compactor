package compactor

// Init function
// Used to start up the plugin and check for their dependencies
type InitFunc = func(options *Options) error

// Resolve function
// Used to transform the source path to the final destination path
type ResolveFunc = func(options *Options, file *File) (string, error)

// Related function
// Detects and generates the list of related files to the given item, being a dependency or not
type RelatedFunc = func(options *Options, file *File) ([]Related, error)

// Transform function
// Used to transform the file content to the final format by applying compilation if necessary
type TransformFunc = func(options *Options, file *File) error

// Optimize function
// Used to run optimizations algorithms to compress file
type OptimizeFunc = func(options *Options, file *File) error

// Plugin struct
type Plugin struct {
	Namespace  string
	Extensions []string
	Init       InitFunc
	Resolve    ResolveFunc
	Related    RelatedFunc
	Transform  TransformFunc
	Optimize   OptimizeFunc
}
