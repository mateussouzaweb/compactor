package compactor

import (
	"github.com/mateussouzaweb/compactor/os"
)

// Plugins index list
// Hold the registered plugins, used for processing
var _plugins []*Plugin

// File index list
// The index contains information of all source files
var _files []*File

// Packages list
// The packages index contains information of all files that resolves to a destination
var _packages []*File

// ======== PLUGINS ===== //

// AddPlugin add a new plugin to the index
func AddPlugin(plugin *Plugin) {
	_plugins = append(_plugins, plugin)
}

// RemovePlugin removes all plugins from index that match the given namespace
func RemovePlugin(namespace string) {

	var list []*Plugin

	for _, _plugin := range _plugins {
		if namespace != _plugin.Namespace {
			list = append(list, _plugin)
		}
	}

	_plugins = list

}

// GetPlugin retrieves the first found plugin for the given extension
func GetPlugin(extension string) *Plugin {

	for _, plugin := range _plugins {

		// Extension plugin
		for _, _extension := range plugin.Extensions {
			if _extension == extension {
				return plugin
			}
		}

		// Generic plugin
		// Generic plugin is always the lastest, so at least one match should happen
		if len(plugin.Extensions) == 0 {
			return plugin
		}

	}

	return &Plugin{}
}

// ======== FILES ===== //

// GetFiles retrieve the indexed files
func GetFiles() []*File {
	return _files
}

// GetFile retrieves the file that match path on index
func GetFile(path string) *File {

	for _, file := range _files {
		if file.Path == path {
			return file
		}
	}

	return &File{}
}

// AppendFile appends file information to index from its path
func AppendFile(path string, root string) error {

	location := os.Clean(path, root)
	content, checksum, perm := os.Info(path)

	file := File{
		Path:        path,
		Destination: "",
		Root:        root,
		Location:    location,
		Folder:      os.Dir(location),
		File:        os.File(location),
		Name:        os.Name(location),
		Extension:   os.Extension(location),
		Content:     content,
		Permission:  perm,
		Exists:      os.Exist(path),
		Checksum:    checksum,
		Previous:    "",
	}

	_files = append(_files, &file)

	return nil
}

// UpdateFile updates file information on index if matches path
func UpdateFile(path string) error {

	for _, file := range _files {

		if file.Path != path {
			continue
		}

		exists := os.Exist(file.Path)
		content, checksum, perm := os.Info(file.Path)
		current := file.Checksum

		file.Content = content
		file.Permission = perm
		file.Exists = exists
		file.Checksum = checksum

		if file.Checksum != current {
			file.Previous = current
		}

		break
	}

	return nil
}

// RemoveFile removes the file information from index if match path
func RemoveFile(path string) {

	for _, file := range _files {

		if file.Path != path {
			continue
		}

		file.Content = ""
		file.Exists = false

		break
	}

}

// IndexFile index the root path files to the index, resolve related and determine destination
func IndexFiles(options *Options, root string) error {

	// First walks on path and add files to the index
	paths, err := os.List(root)

	if err != nil {
		return err
	}

	for _, path := range paths {
		if GetFile(path).Path == "" {
			AppendFile(path, root)
		} else {
			UpdateFile(path)
		}
	}

	// With the updated index, resolve each file to discovery the final destination path
	// We also detect the list of related files that the file have
	for _, file := range _files {

		plugin := GetPlugin(file.Extension)
		destination, err := plugin.Resolve(options, file)

		if err != nil {
			return err
		}

		related, err := plugin.Related(options, file)

		if err != nil {
			return err
		}

		file.Destination = destination
		file.Related = related

	}

	return nil
}

// ======== PACKAGES ===== //

// FindPackages retrieves the full list of detected packages
func FindPackages(options *Options) []*File {

	var packages []*File

	// Create and retrieve every possible package with current index
	// First create a list of ignored files
	// These files cannot have an exclusive package because they are dependencies
	ignore := make(map[string]bool)

	for _, file := range _files {
		for _, related := range file.Related {
			if related.Dependency {
				ignore[related.File.Path] = true
			}
		}
	}

	for _, file := range _files {

		// Prevent if should be ignored
		if _, ok := ignore[file.Path]; ok {
			continue
		}

		if options.ShouldInclude(file.Path) {
			packages = append(packages, file)
		}

	}

	// Replace current index
	_packages = packages

	return _packages
}

// FindPackage retrieves the related package from given path
func FindPackage(options *Options, path string) *File {

	for _, file := range _packages {

		source := options.ToSource(path)
		destination := options.ToDestination(path)

		if file.Path == source || file.Destination == destination {
			return file
		}

		relatedFiles := file.FindRelated(true)
		for _, related := range relatedFiles {
			if related.File.Path == source || related.File.Destination == destination {
				return file
			}
		}

	}

	return &File{}
}

// Process execute file packaging by running plugin methods
func Process(options *Options, file *File) error {

	// Make sure folder exists to avoid issues
	err := os.EnsureDirectory(file.Destination)

	if err != nil {
		return err
	}

	// Find the appropriated plugin by detecting extension
	plugin := GetPlugin(file.Extension)

	// Init action
	err = plugin.Init(options)

	if err != nil {
		return err
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

	toDelete := []string{}

	// File name
	destination := file.Destination
	checksum := options.ToHashed(destination, file.Checksum)
	previous := options.ToHashed(destination, file.Previous)
	clean := options.ToNonHashed(destination, file.Checksum)
	cleanPrevious := options.ToNonHashed(destination, file.Previous)

	toDelete = append(
		toDelete,
		destination,
		checksum,
		previous,
		clean,
		cleanPrevious,
	)

	// Related auto generated dependencies
	for _, related := range file.Related {
		if related.Dependency && related.Source == "" {

			// Variations from related file and checksum
			destination := related.File.Destination
			hashed := options.ToHashed(destination, related.File.Checksum)
			previous := options.ToHashed(destination, related.File.Previous)
			clean := options.ToNonHashed(destination, related.File.Checksum)
			cleanPrevious := options.ToNonHashed(destination, related.File.Previous)

			// Variations from main file checksum
			hashedFromMain := options.ToHashed(destination, file.Checksum)
			previousFromMain := options.ToHashed(destination, file.Previous)
			cleanFromMain := options.ToNonHashed(destination, file.Checksum)
			cleanPreviousFromMain := options.ToNonHashed(destination, file.Previous)

			toDelete = append(
				toDelete,
				destination,
				hashed,
				previous,
				clean,
				cleanPrevious,
				hashedFromMain,
				previousFromMain,
				cleanFromMain,
				cleanPreviousFromMain,
			)

		}
	}

	for _, file := range toDelete {

		if !os.Exist(file) {
			continue
		}

		err := os.Delete(file)

		if err != nil {
			return err
		}

	}

	return nil
}
