package compactor

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// GetPermission retrieve file permission from file
func GetPermission(file string) (fs.FileMode, error) {

	var perm fs.FileMode
	info, err := os.Stat(file)

	if err != nil {
		return perm, err
	}

	perm = info.Mode().Perm()

	return perm, nil
}

// ExistFile check if file exists
func ExistFile(file string) bool {

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}

	return true
}

// ReadFile retrieve file content from file
func ReadFile(file string) (string, error) {

	content, err := ioutil.ReadFile(file)

	if err != nil {
		return "", err
	}

	return string(content), nil
}

// ReadFileAndPermission retrieve file content and permissions from file
func ReadFileAndPermission(file string) (string, fs.FileMode, error) {

	perm, err := GetPermission(file)

	if err != nil {
		return "", perm, err
	}

	content, err := ReadFile(file)

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

// ReadFilesAndPermission retrieve files content from file list and the permision of the first file
func ReadFilesAndPermission(files []string) (string, fs.FileMode, error) {

	perm, err := GetPermission(files[0])

	if err != nil {
		return "", perm, err
	}

	content, err := ReadFiles(files)

	return content, perm, err
}

// WriteFile write content on file
func WriteFile(file string, content string, perm fs.FileMode) error {

	err := ioutil.WriteFile(file, []byte(content), perm)

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
func DeleteFile(file string) error {

	if ExistFile(file) {
		return os.Remove(file)
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
func ChmodFile(file string, perm fs.FileMode) error {
	return os.Chmod(file, perm)
}

// ChmodFile apply permission to file
func ChownFile(file string, user int, group int) error {
	return os.Chown(file, user, group)
}

// GetChecksum retrive the checksum for files
func GetChecksum(files []string) (string, error) {

	content, err := ReadFiles(files)

	if err != nil {
		return "", err
	}

	sum := md5.New()
	_, err = io.WriteString(sum, content)

	inBytes := sum.Sum(nil)[:8]
	hash := hex.EncodeToString(inBytes)

	return hash, err
}

// Return the clean file name, with extension
func CleanName(file string) string {
	return filepath.Base(file)
}

// Return the clean file extension, without dot
func CleanExtension(file string) string {
	return strings.TrimLeft(filepath.Ext(file), ".")
}
