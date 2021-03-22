package compactor

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"os"
)

// GetPermission retrieve file permission from filename
func GetPermission(filename string) (fs.FileMode, error) {

	var perm fs.FileMode
	info, err := os.Stat(filename)

	if err != nil {
		return perm, err
	}

	perm = info.Mode().Perm()

	return perm, nil
}

// ExistFile check if file exists
func ExistFile(filename string) bool {

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	return true
}

// ReadFile retrieve file content from filename
func ReadFile(filename string) (string, error) {

	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return "", err
	}

	return string(content), nil
}

// ReadFileAndPermission retrieve file content and permissions from filename
func ReadFileAndPermission(filename string) (string, fs.FileMode, error) {

	perm, err := GetPermission(filename)

	if err != nil {
		return "", perm, err
	}

	content, err := ReadFile(filename)

	return content, perm, err
}

// ReadFiles retrieve files content from file list
func ReadFiles(files []string) (string, error) {

	buf := bytes.NewBuffer(nil)

	for _, filepath := range files {

		content, err := ReadFile(filepath)

		if err != nil {
			return "", err
		}

		buf.WriteString(content)

	}

	return buf.String(), nil
}

// WriteFile write content on file
func WriteFile(filename string, content string, perm fs.FileMode) error {

	err := ioutil.WriteFile(filename, []byte(content), perm)

	if err != nil {
		return err
	}

	return nil
}

// CopyFile copy the source file into destination
func CopyFile(source string, destination string) error {

	content, perm, err := ReadFileAndPermission(source)

	if err != nil {
		return err
	}

	err = WriteFile(destination, content, perm)

	return err
}

// DeleteFile remove a file
func DeleteFile(filename string) error {

	if ExistFile(filename) {
		return os.Remove(filename)
	}

	return nil
}

// MoveFile move a file to destination
func MoveFile(source string, destination string) error {
	return os.Rename(source, destination)
}

// RenameFile rename a file name
func RenameFile(source string, destination string) error {
	return MoveFile(source, destination)
}

// ChmodFile apply permission to file
func ChmodFile(filename string, perm fs.FileMode) error {
	return os.Chmod(filename, perm)
}

// ChmodFile apply permission to file
func ChownFile(filename string, user int, group int) error {
	return os.Chown(filename, user, group)
}
