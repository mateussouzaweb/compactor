package compactor

import (
	"io/fs"
)

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

// Related struct
type Related struct {
	Type       string
	Dependency bool
	Source     string
	Path       string
	Item       *Item
}
