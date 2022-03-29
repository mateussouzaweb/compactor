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

// IndexRelated detect the related assets on the index
func IndexRelated() error {

	for _, item := range _items {

		plugin := GetPlugin(item.Extension)

		if plugin.Namespace == "" {
			continue
		}

		detected, err := plugin.Related(item)

		if err != nil {
			return err
		}

		item.Related = detected

	}

	return nil
}

// IndexBundles creates and retrieve every possible bundle with current index
func IndexBundles() error {

	// First create a list of ignored files
	// These files cannot have an exclusive bundle
	// Because they are dependency of another main file
	ignore := make(map[string]bool)

	for _, item := range _items {
		for _, related := range item.Related {
			if related.Dependency {
				ignore[related.Item.Path] = true
			}
		}
	}

	// Create a new list of bundles
	var bundles []*Bundle

	// Now create the bundle registry
	for _, item := range _items {

		// Prevent if should be ignored
		if _, ok := ignore[item.Path]; ok {
			continue
		}

		bundle := *Default
		bundle.Extension = item.Extension
		bundle.Item = item

		if bundle.ShouldInclude(item.Path) {
			bundles = append(bundles, &bundle)
		}

	}

	// Update the bundles index
	SetBundles(bundles)

	return nil
}

// Index add path files to the index, resolve related and create bundles
func Index(path string) error {

	err := IndexItems(path)

	if err != nil {
		return err
	}

	err = IndexRelated()

	if err != nil {
		return err
	}

	return IndexBundles()
}

// Process execute bundle packaging by running plugin
func Process(bundle *Bundle) error {

	// Make sure folder exists to avoid issues
	destination := bundle.ToDestination(bundle.Item.Path)
	err := os.EnsureDirectory(destination)

	if err != nil {
		return err
	}

	// Determine action based on processable file
	action := "DELETE"

	if bundle.Item.Exists {
		action = "EXECUTE"
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
