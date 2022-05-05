package compactor

import (
	"io/fs"
)

// Related struct
type Related struct {
	Type       string
	Dependency bool
	Source     string
	Path       string
	Item       *Item
}

// Item struct
type Item struct {
	Path       string      // Full path (root + location)
	Root       string      // Root location
	Location   string      // Location from root
	Folder     string      // Location folder
	File       string      // File name with extension
	Name       string      // File name
	Extension  string      // File Extension
	Content    string      // File content
	Permission fs.FileMode // File permissions
	Exists     bool        // File exists flag
	Checksum   string      // Current Checksum
	Previous   string      // Previous Checksum
	Related    []Related   // Related items
}

// GetRelatedPaths retrieve the related paths of the item
func (i *Item) GetRelatedPaths(onlyDependencies bool) []string {

	var paths []string
	for _, related := range i.Related {
		if onlyDependencies && !related.Dependency {
			continue
		}

		paths = append(paths, related.Item.Path)

		if len(related.Item.Related) > 0 {
			paths = append(paths, related.Item.GetRelatedPaths(onlyDependencies)...)
		}
	}

	return paths
}
