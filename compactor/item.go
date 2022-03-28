package compactor

import (
	"io/fs"
	"strings"
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

// MergedContent returns the content of the item with the replaced partials dependencies
func (i *Item) MergedContent() string {

	if !i.Exists {
		return ""
	}

	content := i.Content

	for _, related := range i.Related {
		if related.Type == "partial" && related.Item.Exists {

			// Solves recursively
			content = strings.Replace(
				content,
				related.Source,
				related.Item.MergedContent(),
				1,
			)

		}
	}

	return content
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
