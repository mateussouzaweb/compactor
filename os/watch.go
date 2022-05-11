package os

import (
	"log"
	"time"

	"github.com/radovskyb/watcher"
)

// WatchCallback type
type WatchCallback func(path string) error

// Watch check for changes on path files and trigger events on change or on delete file
func Watch(path string, interval time.Duration, onChange WatchCallback, onDelete WatchCallback) {

	w := watcher.New()
	// TODO: learn from https://github.com/Shopify/themekit/blob/master/src/file/watcher.go

	go func() {
		for {
			select {
			case event := <-w.Event:

				if !event.IsDir() {

					if event.Op&watcher.Create == watcher.Create {
						onChange(event.Path)
					} else if event.Op&watcher.Write == watcher.Write {
						onChange(event.Path)
					} else if event.Op&watcher.Chmod == watcher.Chmod {
						onChange(event.Path)
					} else if event.Op&watcher.Rename == watcher.Rename {
						// TODO: not processing, maybe is a bug on the extension?
						onDelete(event.OldPath)
						onChange(event.Path)
					} else if event.Op&watcher.Move == watcher.Move {
						// TODO: not processing, maybe is a bug on the extension?
						onDelete(event.OldPath)
						onChange(event.Path)
					} else if event.Op&watcher.Remove == watcher.Remove {
						onDelete(event.Path)
					}

				}

			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	err := w.AddRecursive(path)
	if err != nil {
		log.Fatalln(err)
	}

	err = w.Start(interval)
	if err != nil {
		log.Fatalln(err)
	}

}
