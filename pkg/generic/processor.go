package generic

import "github.com/mateussouzaweb/compactor/compactor"

func DeleteProcessor(bundle *compactor.Bundle, logger *compactor.Logger, extraFormats []string) error {

	toDelete := []string{}

	if bundle.IsToMultipleDestinations() {
		for _, file := range bundle.GetFiles() {
			destination := bundle.ToDestination(file)
			toDelete = append(toDelete, destination)
		}
	} else {
		destination := bundle.GetDestination()
		toDelete = append(toDelete, destination)
	}

	for _, file := range toDelete {

		if !compactor.ExistFile(file) {
			continue
		}

		err := compactor.DeleteFile(file)
		if err != nil {
			return err
		}

		logger.AddDeleted(file)

	}

	for _, file := range toDelete {
		for _, format := range extraFormats {

			extra := bundle.ToExtension(file, format)

			if !compactor.ExistFile(extra) {
				continue
			}

			err := compactor.DeleteFile(extra)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func Processor(action *compactor.Action, bundle *compactor.Bundle, logger *compactor.Logger) error {

	if action.IsDelete() {
		return DeleteProcessor(bundle, logger, []string{})
	}

	files := bundle.GetFiles()

	if bundle.IsToMultipleDestinations() {

		for _, file := range files {

			destination := bundle.ToDestination(file)
			err := compactor.CopyFile(file, destination)

			if err != nil {
				return err
			}

			logger.AddProcessed(file)

		}

		return nil
	}

	content, perm, err := compactor.ReadFilesAndPermission(files)

	if err != nil {
		return err
	}

	destination := bundle.GetDestination()
	err = compactor.WriteFile(destination, content, perm)

	if err == nil {
		logger.AddWritten(destination)
	}

	return err
}
