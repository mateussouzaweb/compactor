package compactor

import "io/fs"

// Item struct
type Item struct {
	Path         string
	Folder       string
	File         string
	Name         string
	Extension    string
	Content      string
	Permission   fs.FileMode
	Exists       bool
	Checksum     string
	Previous     string
	Dependencies []*Item
}
