package generic

import "github.com/mateussouzaweb/compactor/compactor"

func Processor(bundle *compactor.Bundle, logger *compactor.Logger) error {

	files := bundle.GetFiles()
	target, isDir := bundle.GetDestination()
	result := ""

	for _, file := range files {

		if isDir {

			destination := bundle.ToDestination(file)
			err := compactor.CopyFile(file, destination)

			if err != nil {
				return err
			}

			logger.AddProcessed(file)
			continue

		}

		content, err := compactor.ReadFile(file)

		if err != nil {
			return err
		}

		result += content

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
