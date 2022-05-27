package os

import (
	"net/http"
	"path/filepath"
	"strings"
)

// RequestCallback type
type RequestCallback func(uri string) error

// Server start a file server with given path and port
func Server(root string, port string, onRequest RequestCallback) error {

	// Make sure root folder exists
	err := EnsureDirectory(root)

	if err != nil {
		return err
	}

	// Attach server handle
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {

		uri := filepath.Clean(request.URL.Path)
		uri = strings.TrimSuffix(uri, "/")

		// Make sure we are requesting a file when trying to get an unknown uri
		if Extension(uri) == "" {
			uri += "/index.html"
		}

		// Run callback on uri
		err := onRequest(uri)

		// In case of error, return error 500
		if err != nil {
			http.Error(response, http.StatusText(500), 500)
			return
		}

		path := filepath.Join(root, strings.TrimPrefix(uri, "/"))
		indexFile := filepath.Join(root, "index.html")

		// If not exists requested path, then try to reply with root index.html file
		if !Exist(path) {
			if !Exist(indexFile) {
				http.Error(response, http.StatusText(500), 500)
			} else {
				http.ServeFile(response, request, indexFile)
			}
			return
		}

		// If everything ok, just serve the file
		http.ServeFile(response, request, path)

	})

	return http.ListenAndServe(":"+port, mux)
}
