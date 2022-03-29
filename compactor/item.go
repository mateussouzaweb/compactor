package compactor

import (
	"io/fs"
)

// Item struct
type Item struct {
	Path       string
	Folder     string
	File       string
	Name       string
	Extension  string
	Content    string
	Permission fs.FileMode
	Exists     bool
	Checksum   string
	Previous   string
	Related    []Related
}

// Related struct. Valid types:
// import, export, require, partial, source-map, declaration, link, alternative and other
type Related struct {
	Type   string
	Source string
	Path   string
	Item   *Item
}

// IsDependency determines if the related asset is a dependency.
// Dependencies should be carried with the main file
func (r *Related) IsDependency() bool {
	if r.Type == "import" {
		return true
	}
	if r.Type == "partial" {
		return true
	}
	if r.Type == "source-map" {
		return true
	}
	if r.Type == "declaration" {
		return true
	}
	if r.Type == "alternative" {
		return true
	}
	return false
}
