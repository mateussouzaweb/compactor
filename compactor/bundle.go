package compactor

import (
	"path"
	"path/filepath"
	"strings"
)

// Source struct
type Source struct {
	Path    string
	Files   []string
	Include []string
	Exclude []string
	Ignore  []string
}

// Destination struct
type Destination struct {
	Path string
	File string
}

// Compress struct
type Compress struct {
	Enabled bool
	Include []string
	Exclude []string
}

// SourceMap struct
type SourceMap struct {
	Enabled bool
	Include []string
	Exclude []string
}

// Compress struct
type Progressive struct {
	Enabled bool
	Include []string
	Exclude []string
}

// Bundle struct
type Bundle struct {
	Extension   string
	Source      Source
	Destination Destination
	Compress    Compress
	SourceMap   SourceMap
	Progressive Progressive
}

// Return the clean file name, with extension
func (b *Bundle) CleanName(file string) string {
	return filepath.Base(file)
}

// Return the clean file extension, without dot
func (b *Bundle) CleanExtension(file string) string {
	return strings.TrimLeft(filepath.Ext(file), ".")
}

// Return the clean file, without source and destination path
func (b *Bundle) CleanPath(file string) string {

	file = strings.Replace(file, b.Source.Path, "", 1)
	file = strings.Replace(file, b.Destination.Path, "", 1)
	file = strings.TrimLeft(file, "/")

	return file
}

// MatchPatterns return if file match one of the given patterns
func (b *Bundle) MatchPatterns(file string, patterns []string) bool {

	file = b.CleanPath(file)

	for _, pattern := range patterns {

		result, err := path.Match(pattern, file)

		if err != nil {
			continue
		}
		if result {
			return true
		}

	}

	return false
}

// ShouldSkip return if processing should be skipped for given file
func (b *Bundle) ShouldSkip(file string) bool {

	if b.MatchPatterns(file, b.Source.Include) {
		return false
	}

	if b.MatchPatterns(file, b.Source.Exclude) {
		return true
	}

	return false
}

// ShouldIgnore return if processing should be ignored for given file
func (b *Bundle) ShouldIgnore(file string) bool {
	return b.MatchPatterns(file, b.Source.Ignore)
}

// ShouldCompress return if compress should be enabled for given file
func (b *Bundle) ShouldCompress(file string) bool {

	if !b.Compress.Enabled {
		return false
	}

	if b.MatchPatterns(file, b.Compress.Include) {
		return true
	}
	if b.MatchPatterns(file, b.Compress.Exclude) {
		return false
	}

	return true
}

// ShouldGenerateSourceMap return if source map should be generated for given file
func (b *Bundle) ShouldGenerateSourceMap(file string) bool {

	if !b.SourceMap.Enabled {
		return false
	}

	if b.MatchPatterns(file, b.SourceMap.Include) {
		return true
	}
	if b.MatchPatterns(file, b.SourceMap.Exclude) {
		return false
	}

	return true
}

// ShouldGenerateProgressive return if progressive formats should be generated for given file
func (b *Bundle) ShouldGenerateProgressive(file string) bool {

	if !b.Progressive.Enabled {
		return false
	}

	if b.MatchPatterns(file, b.Progressive.Include) {
		return true
	}
	if b.MatchPatterns(file, b.Progressive.Exclude) {
		return false
	}

	return true
}

// AddFile add file to bundle source files
func (b *Bundle) AddFile(file string) bool {

	file = b.CleanPath(file)

	for _, existing := range b.Source.Files {
		if existing == file {
			return true
		}
	}

	b.Source.Files = append(b.Source.Files, file)

	return true
}

// RemoveFile remove file from bundle source files
func (b *Bundle) RemoveFile(file string) {

	file = b.CleanPath(file)

	for index, existing := range b.Source.Files {
		if existing == file {
			b.Source.Files = append(
				b.Source.Files[:index],
				b.Source.Files[index+1:]...,
			)
			return
		}
	}

}

// Contains check if file is in bundle source files
func (b *Bundle) ContainsFile(file string) bool {

	file = b.CleanPath(file)

	for _, existing := range b.Source.Files {

		if existing == file {
			return true
		}

		match, err := path.Match(existing, file)

		if err != nil {
			continue
		}
		if match {
			return true
		}

	}

	return false
}

// Retrieve fullpath files from bundle source list
func (b *Bundle) GetFiles() []string {

	files := []string{}
	patterns := []string{}

	for _, file := range b.Source.Files {

		if b.ShouldSkip(file) || b.ShouldIgnore(file) {
			continue
		}

		path := filepath.Join(b.Source.Path, file)

		if ExistFile(path) {
			files = append(files, path)
		} else {
			patterns = append(patterns, file)
		}

	}

	foundInPattern, _ := FindFilesMatch(b.Source.Path, patterns)
	files = append(files, foundInPattern...)

	return files
}

// Detect if bundle processing should have multiple destinations
func (b *Bundle) IsToMultipleDestinations() bool {
	return b.Destination.File == ""
}

// Return the final destination file path
func (b *Bundle) GetDestination() string {
	return filepath.Join(b.Destination.Path, b.Destination.File)
}

// Transform and return the full destination path for file
func (b *Bundle) ToDestination(file string) string {
	return filepath.Join(b.Destination.Path, b.CleanPath(file))
}

// Return a file converted to a specific extension
func (b *Bundle) ToExtension(file string, extension string) string {

	previousExtension := b.CleanExtension(file)
	file = strings.TrimRight(file, "."+previousExtension)
	file = file + "." + extension

	return file
}
