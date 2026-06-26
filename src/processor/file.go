package processor

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
	Path        string      `json:"path"`        // Full path (root + location)
	Destination string      `json:"destination"` // Full destination path
	Root        string      `json:"root"`        // Root location
	Location    string      `json:"location"`    // Location from root
	Folder      string      `json:"folder"`      // Location folder
	File        string      `json:"file"`        // File name with extension
	Name        string      `json:"name"`        // File name
	Extension   string      `json:"extension"`   // File Extension
	Content     string      `json:"content"`     // File content
	Permission  fs.FileMode `json:"permission"`  // File permissions
	Exists      bool        `json:"exists"`      // File exists flag
	Checksum    []string    `json:"-"`           // Checksum history
	Related     []Related   `json:"-"`           // Related items
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
