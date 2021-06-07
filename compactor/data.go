package compactor

import (
	"strings"

	"github.com/mateussouzaweb/compactor/os"
)

//
// To handle data, we have 4 main data lists:
//
// - Plugins: Hold the registered plugins, used for processing.
// - Indexer: The indexer contains information of all source files. In watch mode, this is live updated.
// - Maps: Define the destination for each custom file. If is not in the index, then destination is equal copy.
// - Bundle: The reference model for the creating bundles. Bundle are always temporary structs.
//
// When compactor is executed, it creates bundles by parsing indexer and maps
// Then, each bundle is processed with matching plugins
//

// Plugins list
var _plugins []*Plugin

// Indexer list
var _indexer []*Item

// Maps list
var _maps []*Mapper

// Default bundle model
var _bundle = &Bundle{}

// PLUGIN METHODS

// Register add a new plugin to the index
func Register(plugin *Plugin) {
	_plugins = append(_plugins, plugin)
}

// Unregister removes all plugins that match the given extension
func Unregister(extension string) {

	var list []*Plugin

	for _, _plugin := range _plugins {

		var match bool

		for _, _extension := range _plugin.Extensions {
			if _extension == extension {
				match = true
				break
			}
		}

		if !match {
			list = append(list, _plugin)
		}

	}

	_plugins = list

}

// INDEXER METHODS

// Get the item that match file path on indexer
func Get(path string) *Item {

	for _, item := range _indexer {
		if item.Path == path {
			return item
		}
	}

	return &Item{}
}

// Append file path information to indexer
func Append(path string) error {

	content, checksum, perm := os.Info(path)

	item := Item{
		Path:       path,
		Folder:     os.Dir(path),
		File:       os.File(path),
		Name:       os.Name(path),
		Extension:  os.Extension(path),
		Content:    content,
		Permission: perm,
		Exists:     os.Exist(path),
		Checksum:   checksum,
		Previous:   "",
	}

	_indexer = append(_indexer, &item)

	return nil
}

// Update item information on indexer if match file path
func Update(path string) error {

	for _, item := range _indexer {

		if item.Path != path {
			continue
		}

		exists := os.Exist(item.Path)
		content, checksum, perm := os.Info(item.Path)
		previous := item.Checksum

		item.Content = content
		item.Permission = perm
		item.Exists = exists
		item.Checksum = checksum
		item.Previous = previous

		break
	}

	return nil
}

// Remove item information from indexer if match file path
func Remove(path string) {

	for _, item := range _indexer {

		if item.Path != path {
			continue
		}

		item.Content = ""
		item.Exists = false

		break
	}

}

// Index walks on path and add files to the indexer
func Index(path string) error {

	files, err := os.List(path)

	if err != nil {
		return err
	}

	for _, file := range files {
		_ = Append(file)
	}

	return nil
}

// Matches run callback on indexer and append item if match
func Matches(callback func(item *Item) bool) []*Item {

	var items []*Item

	for _, item := range _indexer {
		if callback(item) {
			items = append(items, item)
		}
	}

	return items
}

// MAPS METHODS

// Map add a new map registration
func Map(files []string, target string) {
	_maps = append(_maps, &Mapper{
		Files:  files,
		Target: target,
	})
}

// BUNDLE METHODS

// DefaultBundle set the default bundle instance
func DefaultBundle(bundle *Bundle) {
	_bundle = bundle
}

// NewBundle create and retrieve a new bundle
func NewBundle() *Bundle {

	bundle := *_bundle
	bundle.Extension = ""
	bundle.Destination.File = ""
	bundle.Destination.Folder = ""
	bundle.Destination.Path = _bundle.Destination.Path
	bundle.Destination.Hashed = _bundle.Destination.Hashed

	return &bundle
}

// GetBundleFromMapper retrieve the final bundle from given mapper
func GetBundleFromMapper(mapper *Mapper) *Bundle {

	bundle := NewBundle()

	// Check if mapper destination is to file or folder
	if strings.HasSuffix(mapper.Target, "/") {
		bundle.Destination.Folder = bundle.CleanPath(mapper.Target)
	} else {
		bundle.Destination.File = bundle.CleanPath(mapper.Target)
	}

	items := Matches(func(item *Item) bool {

		if !bundle.MatchPatterns(item.Path, mapper.Files) {
			return false
		}

		return bundle.ShouldInclude(item.Path)
	})

	bundle.Items = items

	if len(items) != 0 {
		bundle.Extension = items[0].Extension
	}

	return bundle
}

// GetBundles retrieve every possible bundle with current indexer
func GetBundles() []*Bundle {

	var bundles []*Bundle
	used := make(map[string]bool)

	// First create bundle from maps
	for _, mapper := range _maps {

		bundle := GetBundleFromMapper(mapper)

		if len(bundle.Items) == 0 {
			continue
		}

		bundles = append(bundles, bundle)

		for _, item := range bundle.Items {
			used[item.Path] = true
		}

	}

	// Now process only file not included in previous bundles
	for _, item := range _indexer {

		if _, ok := used[item.Path]; ok {
			continue
		}

		bundle := NewBundle()
		bundle.Extension = item.Extension
		bundle.Items = append(bundle.Items, item)
		bundle.Destination.File = bundle.CleanPath(item.Path)

		if bundle.ShouldInclude(item.Path) {
			bundles = append(bundles, bundle)
		}

	}

	return bundles
}

// GetBundleFor retrieve the related bundle for the file
func GetBundleFor(path string) *Bundle {

	// TODO: When a file depends on another, that file should not go to another bundle
	// It should be specially injected in the bundle, because of the dependency

	// First check if file path included in bundle from maps
	for _, mapper := range _maps {

		bundle := GetBundleFromMapper(mapper)

		if bundle.Destination.File == bundle.CleanPath(path) {
			return bundle
		}

		for _, item := range bundle.Items {
			if item.Path == path {
				return bundle
			}
		}

	}

	item := Get(path)

	bundle := NewBundle()
	bundle.Extension = item.Extension
	bundle.Items = append(bundle.Items, item)
	bundle.Destination.File = bundle.CleanPath(item.Path)

	if !bundle.ShouldInclude(item.Path) {
		bundle = NewBundle()
	}

	return bundle
}

// PROCESS METHODS

// Process package by running processors
func Process(bundle *Bundle) error {

	// Make sure folder exists to avoid issues
	destination := bundle.ToDestination(bundle.Destination.File)
	err := os.EnsureDirectory(destination)

	if err != nil {
		return err
	}

	// Determine action based on processable list
	action := "DELETE"

	for _, item := range bundle.Items {
		if item.Exists {
			action = "RUN"
		}
	}

	// Find and run appropriated plugin by detecting extension
	// Generic plugin is always the lastest, so at least one match should happen
	for _, plugin := range _plugins {

		match := false

		for _, _extension := range plugin.Extensions {
			if _extension == bundle.Extension {
				match = true
				break
			}
		}

		if !match && len(plugin.Extensions) != 0 {
			continue
		}

		if action == "RUN" {
			err = plugin.Run(bundle)
		} else {
			err = plugin.Delete(bundle)
		}

		if err != nil {
			return err
		}

		break
	}

	return err
}
