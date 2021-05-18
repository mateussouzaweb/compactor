package css

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/sass"
)

// CSS processor. Same as SASS
func Processor(action *compactor.Action, bundle *compactor.Bundle, logger *compactor.Logger) error {
	return sass.Processor(action, bundle, logger)
}
