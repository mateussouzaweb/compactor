package css

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/sass"
)

// CSS processor. Same as SASS
func Processor(context *compactor.Context) error {
	return sass.Processor(context)
}
