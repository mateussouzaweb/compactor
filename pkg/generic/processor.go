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
	target := bundle.GetDestination()
	result := ""

	for _, file := range files {

		if bundle.IsToMultipleDestinations() {

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

	if bundle.IsToMultipleDestinations() {
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
