package compactor

import (
	"github.com/mateussouzaweb/compactor/os"
)

// Index list
// The index contains information of all source files
var _items []*Item

// Plugins list
// Hold the registered plugins, used for processing
var _plugins []*Plugin

// Bundles list
// Contains the reference for created bundles
var _bundles []*Bundle

// Default bundle model
var Default = &Bundle{
	Destination: Destination{
		Hashed: true,
	},
	Compress: Compress{
		Enabled: true,
	},
	SourceMap: SourceMap{
		Enabled: true,
	},
	Progressive: Progressive{
		Enabled: true,
	},
}

// Get the item that match file path on index
func Get(path string) *Item {

	for _, item := range _items {
		if item.Path == path {
			return item
		}
	}

	return &Item{}
}

// Append file path information to index
func Append(path string, root string) error {

	location := os.Relative(path, root)
	content, checksum, perm := os.Info(path)

	item := Item{
		Path:       path,
		Root:       root,
		Location:   location,
		Folder:     os.Dir(location),
		File:       os.File(location),
		Name:       os.Name(location),
		Extension:  os.Extension(location),
		Content:    content,
		Permission: perm,
		Exists:     os.Exist(path),
		Checksum:   checksum,
		Previous:   "",
	}

	_items = append(_items, &item)

	return nil
}

// Update item information on index if match file path
func Update(path string) error {

	for _, item := range _items {

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

// Remove item information from index if match file path
func Remove(path string) {

	for _, item := range _items {

		if item.Path != path {
			continue
		}

		item.Content = ""
		item.Exists = false

		break
	}

}

// Matches run callback on index and append item if match
func Matches(callback func(item *Item) bool) []*Item {

	var items []*Item

	for _, item := range _items {
		if callback(item) {
			items = append(items, item)
		}
	}

	return items
}

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

// SetBundles renew the index of bundles with a new list
func SetBundles(bundles []*Bundle) {
	_bundles = bundles
}

// GetBundle retrieve the related bundle for the file
func GetBundle(path string) *Bundle {

	for _, bundle := range _bundles {

		sourcePath := bundle.ToSource(path)

		if bundle.Item.Path == path || bundle.Item.Path == sourcePath {
			return bundle
		}

		for _, related := range bundle.Item.Related {
			if related.Item.Path == path || related.Item.Path == sourcePath {
				if related.Dependency {
					return bundle
				}
			}
		}

	}

	return &Bundle{}
}

// GetBundles retrieve all bundles in the index
func GetBundles() []*Bundle {
	return _bundles
}
