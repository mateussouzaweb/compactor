package compactor

import (
	"path/filepath"
	"strings"
)

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
	Source      string
	Destination string
	Target      string
	Compress    Compress
	SourceMap   SourceMap
	Progressive Progressive
	Files       []string
	Include     []string
	Exclude     []string
	Ignore      []string
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

	file = strings.Replace(file, b.Source, "", 1)
	file = strings.Replace(file, b.Destination, "", 1)
	file = strings.TrimLeft(file, "/")
	file = strings.TrimLeft(file, "\\")

	return file
}

// Return the full source path for file
func (b *Bundle) SourcePath(file string) string {
	return filepath.Join(b.Source, b.CleanPath(file))
}

// Return the full destination path for file
func (b *Bundle) DestinationPath(file string) string {
	return filepath.Join(b.Destination, b.CleanPath(file))
}

// MatchPatterns return if file match one of the given patterns
func (b *Bundle) MatchPatterns(file string, patterns []string) bool {

	file = b.CleanPath(file)

	for _, pattern := range patterns {

		result, err := filepath.Match(pattern, file)

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

	if b.MatchPatterns(file, b.Include) {
		return false
	}

	if b.MatchPatterns(file, b.Exclude) {
		return true
	}

	return false
}

// ShouldIgnore return if processing should be ignored for given file
func (b *Bundle) ShouldIgnore(file string) bool {
	return b.MatchPatterns(file, b.Ignore)
}

// ShouldCompress return if compress is enabled for given file
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

// AddFile add file to bundle file list
func (b *Bundle) AddFile(file string) bool {

	file = b.CleanPath(file)

	if b.ShouldIgnore(file) || b.ShouldSkip(file) {
		return false
	}

	for _, existing := range b.Files {
		if existing == file {
			return true
		}
	}

	b.Files = append(b.Files, file)

	return true
}

// RemoveFile remove file from bundle file list
func (b *Bundle) RemoveFile(file string) {

	file = b.CleanPath(file)

	for index, existing := range b.Files {
		if existing == file {
			b.Files = append(b.Files[:index], b.Files[index+1:]...)
			return
		}
	}

}

// Retrieve files from bundle
func (b *Bundle) GetFiles() []string {
	// TODO: read from pattern?
	return b.Files
}
