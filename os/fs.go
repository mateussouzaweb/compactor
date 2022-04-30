package os

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Exist check if file or directory exists
func Exist(path string) bool {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

// Permissions retrieve permissions for file or directory
func Permissions(path string) (fs.FileMode, error) {

	var perm fs.FileMode
	info, err := os.Stat(path)

	if err != nil {
		return perm, err
	}

	perm = info.Mode().Perm()

	return perm, nil
}

// Read retrieve content from file
func Read(file string) (string, error) {

	content, err := ioutil.ReadFile(file)

	if err != nil {
		return "", err
	}

	return string(content), nil
}

// ReadMany retrieve merged content from file list
func ReadMany(files []string) (string, error) {

	buf := bytes.NewBuffer(nil)

	for _, filepath := range files {

		content, err := Read(filepath)

		if err != nil {
			return "", err
		}

		buf.WriteString(content)

	}

	return buf.String(), nil
}

// Write content on file
func Write(file string, content string, perm fs.FileMode) error {

	err := ioutil.WriteFile(file, []byte(content), perm)

	if err != nil {
		return err
	}

	return nil
}

// Copy the origin file into destination
func Copy(origin string, destination string) error {

	content, err := Read(origin)

	if err != nil {
		return err
	}

	perm, err := Permissions(origin)

	if err != nil {
		return err
	}

	err = Write(destination, content, perm)

	return err
}

// Replace content inside file
func Replace(file string, search string, replace string) error {

	content, err := Read(file)

	if err != nil {
		return err
	}

	permissions, err := Permissions(file)

	if err != nil {
		return err
	}

	content = strings.ReplaceAll(content, search, replace)
	err = Write(file, content, permissions)

	return err
}

// Delete remove a file
func Delete(file string) error {

	if Exist(file) {
		return os.Remove(file)
	}

	return nil
}

// Move a file to destination
func Move(origin string, destination string) error {
	return os.Rename(origin, destination)
}

// Rename a file name
func Rename(origin string, destination string) error {
	return Move(origin, destination)
}

// Chmod apply permissions to file
func Chmod(file string, perm fs.FileMode) error {
	return os.Chmod(file, perm)
}

// Chown apply user and group ownership to file
func Chown(file string, user int, group int) error {
	return os.Chown(file, user, group)
}

// Relative return the relative path from root
func Relative(path string, root string) string {
	return strings.Replace(path, root, "", 1)
}

// Dir return the clean directory path for file
func Dir(path string) string {
	return filepath.Dir(path)
}

// File return the clean file name for path, with extension
func File(path string) string {
	return filepath.Base(path)
}

// Name return the clean file name for path, without extension
func Name(path string) string {
	name := filepath.Base(path)
	ext := filepath.Ext(path)
	return strings.TrimSuffix(name, ext)
}

// Extension return the clean file extension, with dot
func Extension(file string) string {
	return filepath.Ext(file)
}

// Info read and return file information: content, checksum and permissions
func Info(file string) (string, string, fs.FileMode) {

	content, err := Read(file)

	if err != nil {
		content = ""
	}

	perm, err := Permissions(file)

	if err != nil {
		perm = fs.FileMode(0644)
	}

	checksum, err := Checksum(content)

	if err != nil {
		checksum = ""
	}

	return content, checksum, perm
}

// WalkCallback type
type WalkCallback func(path string) error

// Walk find files in path and process callback for every result
func Walk(root string, callback WalkCallback) error {

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

// List walks on path and return every found file
func List(root string) ([]string, error) {

	var files []string

	err := Walk(root, func(path string) error {
		files = append(files, path)
		return nil
	})

	return files, err
}

// Find retrieve files from path that exactly match names
func Find(root string, names []string) ([]string, error) {

	var files []string

	err := Walk(root, func(path string) error {

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

// FindMatch retrieve files from path that match patterns
func FindMatch(root string, patterns []string) ([]string, error) {

	var files []string

	err := Walk(root, func(thePath string) error {

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

// EnsureDirectory makes sure directory exists from file path
func EnsureDirectory(file string) error {

	path := filepath.Dir(file)

	if !Exist(path) {

		err := os.MkdirAll(path, 0775)

		if err != nil {
			return err
		}

	}

	return nil
}
