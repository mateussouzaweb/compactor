package compactor

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/mateussouzaweb/compactor/os"
)

// Source struct
type Source struct {
	Path    string
	Include []string
	Exclude []string
}

// Destination struct
type Destination struct {
	Path   string
	Hashed bool
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

// Progressive struct
type Progressive struct {
	Enabled bool
	Include []string
	Exclude []string
}

// Logs struct
type Logs struct {
	Processed []string
	Skipped   []string
	Ignored   []string
	Written   []string
	Deleted   []string
	Optimized []string
}

// Bundle struct
type Bundle struct {
	Extension   string
	Item        *Item
	Source      Source
	Destination Destination
	Compress    Compress
	SourceMap   SourceMap
	Progressive Progressive
	Logs        Logs
}

// CleanPath return the clean file, without source and destination path
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

// ShouldInclude return if file should be included on the bundle
func (b *Bundle) ShouldInclude(file string) bool {

	if len(b.Source.Exclude) != 0 && b.MatchPatterns(file, b.Source.Exclude) {
		return false
	}

	if len(b.Source.Include) != 0 && !b.MatchPatterns(file, b.Source.Include) {
		return false
	}

	return true
}

// ShouldCompress return if compress should be enabled for given file
func (b *Bundle) ShouldCompress(file string) bool {

	if !b.Compress.Enabled {
		return false
	}

	if len(b.Compress.Exclude) != 0 && b.MatchPatterns(file, b.Compress.Exclude) {
		return false
	}
	if len(b.Compress.Include) != 0 && !b.MatchPatterns(file, b.Compress.Include) {
		return false
	}

	return true
}

// ShouldGenerateSourceMap return if source map should be generated for given file
func (b *Bundle) ShouldGenerateSourceMap(file string) bool {

	if !b.SourceMap.Enabled {
		return false
	}

	if len(b.SourceMap.Exclude) != 0 && b.MatchPatterns(file, b.SourceMap.Exclude) {
		return false
	}
	if len(b.SourceMap.Include) != 0 && !b.MatchPatterns(file, b.SourceMap.Include) {
		return false
	}

	return true
}

// ShouldGenerateProgressive return if progressive formats should be generated for given file
func (b *Bundle) ShouldGenerateProgressive(file string) bool {

	if !b.Progressive.Enabled {
		return false
	}

	if len(b.Progressive.Exclude) != 0 && b.MatchPatterns(file, b.Progressive.Exclude) {
		return false
	}
	if len(b.Progressive.Include) != 0 && !b.MatchPatterns(file, b.Progressive.Include) {
		return false
	}

	return true
}

// ToSource transform and return the full source path for file
func (b *Bundle) ToSource(file string) string {
	return filepath.Join(b.Source.Path, b.CleanPath(file))
}

// ToDestination transform and return the full destination path for file
func (b *Bundle) ToDestination(file string) string {
	return filepath.Join(b.Destination.Path, b.CleanPath(file))
}

// ToExtension return a file converted to a specific extension
func (b *Bundle) ToExtension(file string, extension string) string {

	previousExtension := os.Extension(file)
	file = strings.TrimSuffix(file, previousExtension)
	file = file + extension

	return file
}

// ToHashed return a file converted to a hashed name to avoid caching
func (b *Bundle) ToHashed(file string, hash string) string {

	if hash == "" || !b.Destination.Hashed {
		return file
	}

	extension := os.Extension(file)
	file = strings.TrimSuffix(file, extension)
	file = file + "." + hash + extension

	return file
}

// ToNonHashed return a file converted back to a non hashed name
func (b *Bundle) ToNonHashed(file string, hash string) string {

	if hash == "" || !b.Destination.Hashed {
		return file
	}

	extension := os.Extension(file)
	file = strings.TrimSuffix(file, extension)
	file = strings.TrimSuffix(file, "."+hash)
	file = file + extension

	return file
}

// Processed append path to processed list
func (b *Bundle) Processed(path string) {
	b.Logs.Processed = append(b.Logs.Processed, path)
}

// Skipped append path to skipped list
func (b *Bundle) Skipped(path string) {
	b.Logs.Skipped = append(b.Logs.Skipped, path)
}

// Ignored append path to ignored list
func (b *Bundle) Ignored(path string) {
	b.Logs.Ignored = append(b.Logs.Ignored, path)
}

// Written append path to written list
func (b *Bundle) Written(path string) {
	b.Logs.Written = append(b.Logs.Written, path)
}

// Deleted append path to deleted list
func (b *Bundle) Deleted(path string) {
	b.Logs.Deleted = append(b.Logs.Deleted, path)
}

// Optimized append path to optimized list
func (b *Bundle) Optimized(path string) {
	b.Logs.Optimized = append(b.Logs.Optimized, path)
}
