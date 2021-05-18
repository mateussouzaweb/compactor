package svg

import (
	"github.com/mateussouzaweb/compactor/compactor"
)

// SVG minify
func Minify(content string) (string, error) {

	// TODO: Viewbox removal causing bugs
	// _, err = compactor.ExecCommand(
	// 	"svgo",
	// 	"--quiet",
	// 	"--input", target,
	// 	"--output", target,
	// )

	return content, nil
}

// SVG processor
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
