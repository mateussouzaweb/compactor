package processor

import (
	"github.com/mateussouzaweb/compactor/src/system"
)

// Process execute file packaging by running plugin methods
func Process(options *Options, file *File) error {

	// Make sure folder exists to avoid issues
	err := system.EnsureDirectory(file.Destination)
	if err != nil {
		return err
	}

	// Find the appropriated plugin by detecting extension
	plugin := GetPlugin(file.Extension)

	// Init action
	if !plugin.Initialized {
		err = plugin.Init(options)
		plugin.Initialized = true

		if err != nil {
			return err
		}
	}

	// Determine action based on processable file
	// If not exists, the run delete action
	if !file.Exists {
		return Delete(options, file)
	}

	// Check if should execute the action
	// Useful to detect if file is in fact updated
	// if !plugin.ShouldExecute(options, file) {
	// 	return nil
	// }

	// Transform action
	err = plugin.Transform(options, file)
	if err != nil {
		return err
	}

	// Optimize action
	return plugin.Optimize(options, file)
}

// Delete removes the destination file(s) for given file
func Delete(options *Options, file *File) error {

	destination := file.Destination
	toDelete := []string{destination}

	// File checksum history
	for _, checksum := range file.Checksum {
		path := options.ToHashed(destination, checksum)
		clean := options.ToNonHashed(destination, checksum)
		toDelete = append(toDelete, path, clean)
	}

	// Related auto generated dependencies
	for _, related := range file.Related {
		if related.Dependency && related.Source == "" {

			destination := related.File.Destination
			toDelete = append(toDelete, destination)

			// Variations from related file checksum
			for _, checksum := range related.File.Checksum {
				path := options.ToHashed(destination, checksum)
				clean := options.ToNonHashed(destination, checksum)
				toDelete = append(toDelete, path, clean)
			}

			// Variations from main file checksum
			for _, checksum := range file.Checksum {
				path := options.ToHashed(destination, checksum)
				clean := options.ToNonHashed(destination, checksum)
				toDelete = append(toDelete, path, clean)
			}

		}
	}

	for _, file := range toDelete {

		if !system.Exist(file) {
			continue
		}

		err := system.Delete(file)
		if err != nil {
			return err
		}

	}

	return nil
}

// Shutdown make sure every plugin has properly shutdown
func Shutdown(options *Options) error {

	for _, plugin := range _plugins {
		if plugin.Initialized {
			err := plugin.Shutdown(options)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
