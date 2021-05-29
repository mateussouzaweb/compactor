package xml

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/generic"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/xml"
)

// XML minify
func Minify(content string) (string, error) {

	m := minify.New()
	m.AddFunc("generic", xml.Minify)

	content, err := m.String("generic", content)

	return content, err
}

// XML processor
func Processor(action *compactor.Action, bundle *compactor.Bundle, logger *compactor.Logger) error {

	if action.IsDelete() {
		return generic.DeleteProcessor(bundle, logger, []string{})
	}

	files := bundle.GetFiles()

	for _, file := range files {

		content, perm, err := compactor.ReadFileAndPermission(file)

		if err != nil {
			return err
		}

		if bundle.ShouldCompress(file) {
			content, err = Minify(content)
			if err != nil {
				return err
			}
		}

		destination := bundle.ToDestination(file)
		err = compactor.WriteFile(destination, content, perm)

		if err != nil {
			return err
		}

		logger.AddProcessed(file)

	}

	return nil
}
