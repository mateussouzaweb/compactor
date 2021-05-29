package compactor

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

// WalkFilesCallback type
type WalkFilesCallback func(path string) error

// WalkFiles find files in source path and process callback for every result
func WalkFiles(root string, callback WalkFilesCallback) error {

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return callback(path)
	})

	return err
}

// ListFiles walks on source path and return every found file
func ListFiles(root string) ([]string, error) {

	var files []string

	err := WalkFiles(root, func(path string) error {
		files = append(files, path)
		return nil
	})

	return files, err
}

// FindFiles retrieve files from source path that exactly match names
func FindFiles(root string, names []string) ([]string, error) {

	var files []string

	err := WalkFiles(root, func(path string) error {

		name := filepath.Base(path)

		for _, item := range names {
			if name == item {
				files = append(files, path)
				break
			}
		}

		return nil
	})

	return files, err
}

// FindFilesMatch retrieve files from source path that match patterns
func FindFilesMatch(root string, patterns []string) ([]string, error) {

	var files []string

	err := WalkFiles(root, func(thePath string) error {

		for _, pattern := range patterns {

			file := strings.Replace(thePath, root, "", 1)
			file = strings.TrimLeft(file, "/")

			matched, err := path.Match(pattern, file)

			if err != nil {
				return err
			}

			if matched {
				files = append(files, thePath)
				break
			}

		}

		return nil
	})

	return files, err
}

// ExistDirectory check if directory exists
func ExistDirectory(path string) bool {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

// EnsureDirectory makes sure directory exists from file path
func EnsureDirectory(file string) error {

	path := filepath.Dir(file)

	if !ExistDirectory(path) {

		err := os.MkdirAll(path, 0775)

		if err != nil {
			return err
		}

	}

	return nil
}
