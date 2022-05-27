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

// Options struct
type Options struct {
	Source      Source
	Destination Destination
	Compress    Compress
	SourceMap   SourceMap
	Progressive Progressive
}

// CleanPath return the clean path, without source and destination path
func (o *Options) CleanPath(path string) string {

	path = strings.Replace(path, o.Source.Path, "", 1)
	path = strings.Replace(path, o.Destination.Path, "", 1)
	path = strings.TrimLeft(path, "/")

	return path
}

// MatchPatterns return if file match one of the given patterns
func (o *Options) MatchPatterns(file string, patterns []string) bool {

	file = o.CleanPath(file)

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

// ShouldInclude return if path should be included on the bundle
func (o *Options) ShouldInclude(path string) bool {

	if len(o.Source.Exclude) != 0 && o.MatchPatterns(path, o.Source.Exclude) {
		return false
	}

	if len(o.Source.Include) != 0 && !o.MatchPatterns(path, o.Source.Include) {
		return false
	}

	return true
}

// ShouldCompress return if compress should be enabled for given path
func (o *Options) ShouldCompress(path string) bool {

	if !o.Compress.Enabled {
		return false
	}

	if len(o.Compress.Exclude) != 0 && o.MatchPatterns(path, o.Compress.Exclude) {
		return false
	}
	if len(o.Compress.Include) != 0 && !o.MatchPatterns(path, o.Compress.Include) {
		return false
	}

	return true
}

// ShouldGenerateSourceMap return if source map should be generated for given path
func (o *Options) ShouldGenerateSourceMap(path string) bool {

	if !o.SourceMap.Enabled {
		return false
	}

	if len(o.SourceMap.Exclude) != 0 && o.MatchPatterns(path, o.SourceMap.Exclude) {
		return false
	}
	if len(o.SourceMap.Include) != 0 && !o.MatchPatterns(path, o.SourceMap.Include) {
		return false
	}

	return true
}

// ShouldGenerateProgressive return if progressive formats should be generated for given path
func (o *Options) ShouldGenerateProgressive(path string) bool {

	if !o.Progressive.Enabled {
		return false
	}

	if len(o.Progressive.Exclude) != 0 && o.MatchPatterns(path, o.Progressive.Exclude) {
		return false
	}
	if len(o.Progressive.Include) != 0 && !o.MatchPatterns(path, o.Progressive.Include) {
		return false
	}

	return true
}

// ToSource transform and return the full source path for given path
func (o *Options) ToSource(path string) string {
	return filepath.Join(o.Source.Path, o.CleanPath(path))
}

// ToDestination transform and return the full destination path for given path
func (o *Options) ToDestination(path string) string {
	return filepath.Join(o.Destination.Path, o.CleanPath(path))
}

// ToExtension return a file path converted to a specific extension
func (o *Options) ToExtension(path string, extension string) string {

	previousExtension := os.Extension(path)
	path = strings.TrimSuffix(path, previousExtension)
	path = path + extension

	return path
}

// ToHashed return a file path converted to a hashed name to avoid caching
func (o *Options) ToHashed(path string, hash string) string {

	if hash == "" || !o.Destination.Hashed {
		return path
	}

	extension := os.Extension(path)
	path = strings.TrimSuffix(path, extension)
	path = path + "." + hash + extension

	return path
}

// ToNonHashed return a file path converted back to a non hashed name
func (o *Options) ToNonHashed(path string, hash string) string {

	if hash == "" || !o.Destination.Hashed {
		return path
	}

	extension := os.Extension(path)
	path = strings.TrimSuffix(path, extension)
	path = strings.TrimSuffix(path, "."+hash)
	path = path + extension

	return path
}
