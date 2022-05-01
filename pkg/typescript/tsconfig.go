package typescript

import (
	"encoding/json"
	"path/filepath"

	"github.com/mateussouzaweb/compactor/os"
)

// TSConfig struct
type TSConfig struct {
	CompilerOptions map[string]interface{} `json:"compilerOptions,omitempty"`
	WatchOptions    map[string]interface{} `json:"watchOptions,omitempty"`
	TypeAcquisition map[string]interface{} `json:"typeAcquisition,omitempty"`
	Exclude         []string               `json:"exclude,omitempty"`
	Extends         string                 `json:"extends,omitempty"`
	Files           []string               `json:"files,omitempty"`
	Include         []string               `json:"include,omitempty"`
	References      []string               `json:"references,omitempty"`
}

// FindConfig locate the user defined TypeScript config file
func FindConfig(path string) string {

	if os.Exist(filepath.Join(path, "jsconfig.json")) {
		return filepath.Join(path, "jsconfig.json")
	}
	if os.Exist(filepath.Join(path, "tsconfig.json")) {
		return filepath.Join(path, "tsconfig.json")
	}
	if len(path) <= 1 {
		return ""
	}

	return FindConfig(os.Dir(path))
}

// ReadConfig data from config file if exists
func ReadConfig(path string) (*TSConfig, error) {

	config := TSConfig{}

	if path == "" {
		return &config, nil
	}

	content, err := os.Read(path)

	if err != nil {
		return &config, err
	}

	err = json.Unmarshal([]byte(content), &config)

	if err != nil {
		return &config, err
	}

	return &config, nil
}
