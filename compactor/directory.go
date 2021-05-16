package compactor

import (
	"os"
	"path/filepath"
)

// ListFiles walks on filepath and return every found file
func ListFiles(root string) ([]string, error) {

	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, err
}

// FindFiles retrieve files from source path
func FindFiles(root string, names []string) ([]string, error) {

	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		name := filepath.Base(path)
		match := false

		for _, item := range names {
			if name == item {
				match = true
				break
			}
		}

		if !match {
			return nil
		}

		files = append(files, path)

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
