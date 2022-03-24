package compactor

import (
	"github.com/mateussouzaweb/compactor/os"
)

// IndexItems walks on path and add files to the index
func IndexItems(path string) error {

	files, err := os.List(path)

	if err != nil {
		return err
	}

	for _, file := range files {
		if Get(file).Path == "" {
			Append(file)
		} else {
			Update(file)
		}
	}

	return nil
}

// IndexDependencies detect the dependencies on the index
func IndexDependencies() error {

	for _, item := range _items {

		plugin := GetPlugin(item.Extension)

		if plugin.Namespace == "" {
			continue
		}

		detected, err := plugin.Dependencies(item)

		if err != nil {
			return err
		}

		dependencies := Matches(func(match *Item) bool {
			for _, path := range detected {
				if match.Path == path {
					return true
				}
			}
			return false
		})

		item.Dependencies = append(item.Dependencies, dependencies...)

	}

	return nil
}

// IndexBundles creates and retrieve every possible bundle with current index
func IndexBundles() error {

	// First create a list of ignored dependencies
	// These files cannot have an exclusive bundle
	// Because they are dependency of another main file
	ignore := make(map[string]bool)

	for _, item := range _items {
		for _, dependency := range item.Dependencies {
			ignore[dependency.Path] = true
		}
	}

	// Now create the bundle registry
	for _, item := range _items {

		// Prevent if should be ignored
		if _, ok := ignore[item.Path]; ok {
			continue
		}

		bundle := NewBundle()
		bundle.Extension = item.Extension
		bundle.Destination.File = bundle.CleanPath(item.Path)
		bundle.Items = append(bundle.Items, item)
		bundle.Items = append(bundle.Items, item.Dependencies...)

		if bundle.ShouldInclude(item.Path) {
			AddBundle(bundle)
		}

	}

	return nil
}

// Index add path files to the index, resolve dependencies and create bundles
func Index(path string) error {

	err := IndexItems(path)

	if err != nil {
		return err
	}

	err = IndexDependencies()

	if err != nil {
		return err
	}

	return IndexBundles()
}

// Process execute bundle packaging by running plugin
func Process(bundle *Bundle) error {

	// Make sure folder exists to avoid issues
	destination := bundle.ToDestination(bundle.Destination.File)
	err := os.EnsureDirectory(destination)

	if err != nil {
		return err
	}

	// Determine action based on processable file list
	action := "DELETE"

	for _, item := range bundle.Items {
		if item.Exists {
			action = "EXECUTE"
		}
	}

	// Find the appropriated plugin by detecting extension
	plugin := GetPlugin(bundle.Extension)

	// Init action
	err = plugin.Init(bundle)

	if err != nil {
		return err
	}

	// Delete action
	if action == "DELETE" {
		return plugin.Delete(bundle)
	}

	// Execute action
	err = plugin.Execute(bundle)

	if err != nil {
		return err
	}

	// Optimize
	return plugin.Optimize(bundle)
}
