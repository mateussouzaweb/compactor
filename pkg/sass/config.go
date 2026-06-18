package sass

// SassConfig struct
type SassConfig struct {
	SourceMap               bool   `json:"sourceMap,omitempty"`
	SourceMapIncludeSources bool   `json:"watchOptions,omitempty"`
	Style                   string `json:"style,omitempty"`
}
