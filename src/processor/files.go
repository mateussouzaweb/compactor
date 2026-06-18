package processor

import (
	"slices"

	"github.com/mateussouzaweb/compactor/src/system"
)

// File index list
// The index contains information of all source files
var _files []*File

// GetFiles retrieve the indexed files
func GetFiles() []*File {
	return _files
}

// GetFile retrieves the file that match path on index
func GetFile(path string) *File {

	for _, file := range _files {
		if file.Path == path {
			return file
		}
	}

	return &File{}
}

// AppendFile appends file information to index from its path
func AppendFile(path string, root string) error {

	location := system.Clean(path, root)
	content, checksum, perm := system.Info(path)

	file := File{
		Path:        path,
		Destination: "",
		Root:        root,
		Location:    location,
		Folder:      system.Dir(location),
		File:        system.File(location),
		Name:        system.Name(location),
		Extension:   system.Extension(location),
		Content:     content,
		Permission:  perm,
		Exists:      system.Exist(path),
		Checksum:    []string{checksum},
	}

	_files = append(_files, &file)

	return nil
}

// UpdateFile updates file information on index if matches path
func UpdateFile(path string) error {

	for _, file := range _files {

		if file.Path != path {
			continue
		}

		exists := system.Exist(file.Path)
		content, checksum, perm := system.Info(file.Path)

		file.Content = content
		file.Permission = perm
		file.Exists = exists

		checksumExists := slices.Contains(file.Checksum, checksum)

		if !checksumExists {
			file.Checksum = append(file.Checksum, checksum)
		}

		break
	}

	return nil
}

// RemoveFile removes the file information from index if match path
func RemoveFile(path string) {

	for _, file := range _files {

		if file.Path != path {
			continue
		}

		file.Content = ""
		file.Exists = false

		break
	}

}

// IndexFile index the root path files to the index, resolve related and determine destination
func IndexFiles(options *Options, root string) error {

	// First walks on path and add files to the index
	paths, err := system.List(root)
	if err != nil {
		return err
	}

	for _, path := range paths {
		if GetFile(path).Path == "" {
			AppendFile(path, root)
		} else {
			UpdateFile(path)
		}
	}

	// With the updated index, resolve each file to discovery the final destination path
	// We also detect the list of related files that the file have
	for _, file := range _files {

		plugin := GetPlugin(file.Extension)
		destination, err := plugin.Resolve(options, file)
		if err != nil {
			return err
		}

		related, err := plugin.Related(options, file)
		if err != nil {
			return err
		}

		file.Destination = destination
		file.Related = related

	}

	return nil
}
