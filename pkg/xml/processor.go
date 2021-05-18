package xml

import (
	"github.com/mateussouzaweb/compactor/compactor"
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
func Processor(bundle *compactor.Bundle, logger *compactor.Logger) error {

	files := bundle.GetFiles()
	target, isDir := bundle.GetDestination()
	result := ""

	for _, file := range files {

		content, err := compactor.ReadFile(file)

		if err != nil {
			return err
		}

		if bundle.ShouldCompress(file) {
			content, err = Minify(content)
			if err != nil {
				return err
			}
		}

		if !isDir {
			result += content
			continue
		}

		destination := bundle.ToDestination(file)
		perm, err := compactor.GetPermission(file)

		if err != nil {
			return err
		}

		err = compactor.WriteFile(destination, content, perm)

		if err != nil {
			return err
		}

		logger.AddProcessed(file)

	}

	if isDir {
		return nil
	}

	perm, err := compactor.GetPermission(files[0])

	if err != nil {
		return err
	}

	err = compactor.WriteFile(target, result, perm)

	if err == nil {
		logger.AddWritten(target)
	}

	return err
}
