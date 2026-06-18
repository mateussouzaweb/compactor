package css

import (
	"github.com/mateussouzaweb/compactor/src/plugins/generic"
	"github.com/mateussouzaweb/compactor/src/plugins/sass"
	"github.com/mateussouzaweb/compactor/src/processor"
)

// Plugin return the compactor plugin instance
func Plugin() *processor.Plugin {
	return &processor.Plugin{
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
