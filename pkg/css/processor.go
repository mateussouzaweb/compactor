package css

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/generic"
	"github.com/mateussouzaweb/compactor/pkg/sass"
)

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:  "css",
		Extensions: []string{".css"},
		Init:       sass.Init,
		Shutdown:   sass.Shutdown,
		Resolve:    sass.Resolve,
		Related:    sass.Related,
		Transform:  sass.Transform,
		Optimize:   generic.Optimize,
	}
}
