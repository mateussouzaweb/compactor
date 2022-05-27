package os

import (
	"time"
)

// Watch changes tracking
var watchTrack = map[string]string{}

// WatchCallback type
type WatchCallback func(path string, action string) error

// WatchCheckFile will check if file has been changed and run callback once necessary
func WatchCheckFile(file string, onChange WatchCallback) error {

	_, checksum, _ := Info(file)

	if current, ok := watchTrack[file]; ok {
		if current != checksum {
			watchTrack[file] = checksum
			onChange(file, "updated")
			return nil
		} else {
			return nil
		}
	}

	watchTrack[file] = checksum
	onChange(file, "added")

	return nil
}

// Watch will check if there is any change in the files inside path
func Watch(root string, onChange WatchCallback) error {

	// Fill initial checksums
	files, err := List(root)

	if err != nil {
		return err
	}

	for _, file := range files {
		_, checksum, _ := Info(file)
		watchTrack[file] = checksum
	}

	// Start tracking
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	go func() error {
		for {
			select {
			case <-done:
				return nil
			case <-ticker.C:
				files, err := List(root)

				if err != nil {
					return err
				}

				for _, file := range files {
					WatchCheckFile(file, onChange)
				}
			}
		}
	}()

	return nil
}
