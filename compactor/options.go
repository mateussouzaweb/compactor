package compactor

import "path/filepath"

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
	Files       []string
	Destination string
}

// Options struct
type Options struct {
	Source      string
	Destination string
	Development bool
	Watch       bool
	Compress    Compress
	SourceMap   SourceMap
	Progressive Progressive
	Bundles     []Bundle
	Include     []string
	Exclude     []string
	Ignore      []string
}

// CanProcess check if item can be processed based on include and exclude list
func (o *Options) CanProcess(item string, include []string, exclude []string) bool {

	for _, v := range exclude {
		if v == item {
			return false
		}
	}

	if len(include) == 0 {
		return true
	}

	for _, v := range include {
		if v == item {
			return true
		}
	}

	return false
}

// ShouldIgnore return if processing should be ignored for given context
func (o *Options) ShouldIgnore(context *Context) bool {

	for _, v := range o.Ignore {
		if v == context.Path {
			return true
		}
	}

	return false
}

// ShouldSkip return if processing should be skipped for given context
func (o *Options) ShouldSkip(context *Context) bool {

	for _, pattern := range o.Include {

		result, err := filepath.Match(pattern, context.Path)
		if err != nil {
			continue
		}

		if result {
			return false
		}

	}

	for _, pattern := range o.Exclude {

		result, err := filepath.Match(pattern, context.Path)
		if err != nil {
			continue
		}

		if result {
			return true
		}

	}

	return false
}

// ShouldCompress return if compress is enabled for given context
func (o *Options) ShouldCompress(context *Context) bool {

	if o.Development || !o.Compress.Enabled {
		return false
	}

	return o.CanProcess(context.Extension, o.Compress.Include, o.Compress.Exclude)
}

// GenerateSourceMap return if source map should be generated for given context
func (o *Options) GenerateSourceMap(context *Context) bool {

	if !o.SourceMap.Enabled {
		return false
	}

	return o.CanProcess(context.Extension, o.SourceMap.Include, o.SourceMap.Exclude)
}

// GenerateProgressive return if progressive formats should be generated for given context
func (o *Options) GenerateProgressive(context *Context) bool {

	if !o.Progressive.Enabled {
		return false
	}

	return o.CanProcess(context.Extension, o.Progressive.Include, o.Progressive.Exclude)
}
