package processor

import "slices"

// Plugins index list
// Hold the registered plugins, used for processing
var _plugins []*Plugin

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
		if slices.Contains(plugin.Extensions, extension) {
			return plugin
		}

		// Generic plugin
		// Generic plugin is always the lastest, so at least one match should happen
		if len(plugin.Extensions) == 0 {
			return plugin
		}

	}

	return &Plugin{}
}
