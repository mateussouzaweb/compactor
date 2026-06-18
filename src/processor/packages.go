package processor

// Packages list
// The packages index contains information of all files that resolves to a destination
var _packages []*File

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
