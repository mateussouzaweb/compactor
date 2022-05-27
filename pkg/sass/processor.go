package sass

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/css"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:  "sass",
		Extensions: []string{".sass", ".scss", ".css"},
		Init:       css.Init,
		Resolve:    css.Resolve,
		Related:    css.Related,
		Transform:  css.Transform,
		Optimize:   generic.Optimize,
	}
}
