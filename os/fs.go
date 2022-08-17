package os

import (
	"io/fs"
	"os"
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

	content, err := os.ReadFile(file)

	if err != nil {
		return "", err
	}

	return string(content), nil
}

// Write content on file
func Write(file string, content string, perm fs.FileMode) error {

	err := os.WriteFile(file, []byte(content), perm)

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

// Rename a file path. Overwrite if already exists
func Rename(origin string, destination string) error {
	return os.Rename(origin, destination)
}

// Chmod apply permissions to file
func Chmod(file string, perm fs.FileMode) error {
	return os.Chmod(file, perm)
}

// Chown apply user and group ownership to file
func Chown(file string, user int, group int) error {
	return os.Chown(file, user, group)
}

// Clean return the cleaned relative path from root
func Clean(path string, root string) string {
	return strings.Replace(path, root, "", 1)
}

// Relative return the relative path from root to the target path
func Relative(base string, target string) string {
	relative, _ := filepath.Rel(base, target)
	return relative
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

// Resolve will check paths with possible extensions until file is detected
func Resolve(file string, extensions []string, path string) string {

	if Exist(filepath.Join(path, file)) {
		return filepath.Join(path, file)
	}

	for _, extension := range extensions {
		if Exist(filepath.Join(path, file+extension)) {
			return filepath.Join(path, file+extension)
		}
	}

	if len(path) <= 1 {
		return file
	}

	return Resolve(file, extensions, Dir(path))
}

// List walks on path and return every found file
func List(root string) ([]string, error) {

	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, err
}
