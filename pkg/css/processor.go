package css

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/sass"
)

// CSS processor. Same as SASS
func Processor(bundle *compactor.Bundle, logger *compactor.Logger) error {
	return sass.Processor(bundle, logger)
}
