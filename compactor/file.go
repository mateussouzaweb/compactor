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
	File       *File
}

// File struct
type File struct {
	Path        string      // Full path (root + location)
	Destination string      // Full destination path
	Root        string      // Root location
	Location    string      // Location from root
	Folder      string      // Location folder
	File        string      // File name with extension
	Name        string      // File name
	Extension   string      // File Extension
	Content     string      // File content
	Permission  fs.FileMode // File permissions
	Exists      bool        // File exists flag
	Checksum    string      // Current Checksum
	Previous    string      // Previous Checksum
	Related     []Related   // Related items
}

// FindRelated retrieve the related paths of the item recursively
func (f *File) FindRelated(onlyDependencies bool) []Related {

	var found []Related
	var related []Related
	existing := make(map[string]bool)

	for _, related := range f.Related {
		if onlyDependencies && !related.Dependency {
			continue
		}

		found = append(found, related)

		if len(related.File.Related) > 0 {
			fromRelated := related.File.FindRelated(onlyDependencies)
			if len(fromRelated) > 0 {
				found = append(found, fromRelated...)
			}
		}
	}

	for _, file := range found {
		if _, ok := existing[file.Path]; !ok {
			existing[file.Path] = true
			related = append(related, file)
		}
	}

	return related
}
