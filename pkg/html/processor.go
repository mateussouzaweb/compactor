package html

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/css"
	"github.com/mateussouzaweb/compactor/pkg/generic"
	"github.com/mateussouzaweb/compactor/pkg/javascript"
)

// Init processor
func Init(bundle *compactor.Bundle) error {
	return os.NodeRequire("html-minifier", "html-minifier")
}

// ExtractAttribute find the value of the attribute
func ExtractAttribute(html string, attribute string, defaultValue string) string {

	regex := regexp.MustCompile(attribute + `="([^"]*)"`)
	match := regex.FindStringSubmatch(html)

	if match != nil {
		return match[1]
	}

	regex = regexp.MustCompile(attribute + `='([^']*)'`)
	match = regex.FindStringSubmatch(html)

	if match != nil {
		return match[1]
	}

	return defaultValue
}

// Related processor
func Related(item *compactor.Item) ([]compactor.Related, error) {

	var related []compactor.Related

	// Detect imports
	regex := regexp.MustCompile(`<!-- @import "(.+)" -->`)
	matches := regex.FindAllStringSubmatch(item.Content, -1)

	for _, match := range matches {
		source := match[0]
		path := match[1]
		file := filepath.Join(os.Dir(item.Path), path)

		if os.Extension(file) == "" {
			file += item.Extension
		}

		related = append(related, compactor.Related{
			Type:   "partial",
			Source: source,
			Path:   path,
			Item:   compactor.Get(file),
		})
	}

	// Detect scripts
	regex = regexp.MustCompile(`(?i)<script(.+)?>(.+)?</script>`)
	matches = regex.FindAllStringSubmatch(item.Content, -1)

	for _, match := range matches {

		code := match[0]
		src := ExtractAttribute(code, "src", "")

		// Ignore if is not a src script
		if src == "" {
			continue
		}

		// Ignore protocol scripts, only handle relative and absolute paths
		if strings.Contains(src, "://") {
			continue
		}

		file := filepath.Join(os.Dir(item.Path), src)
		related = append(related, compactor.Related{
			Type:   "other",
			Source: code,
			Path:   src,
			Item:   compactor.Get(file),
		})

	}

	// Detect links
	regex = regexp.MustCompile(`(?i)<link(.+)?\/?>`)
	matches = regex.FindAllStringSubmatch(item.Content, -1)

	for _, match := range matches {

		code := match[0]
		rel := ExtractAttribute(code, "rel", "")
		href := ExtractAttribute(code, "href", "")
		as := ExtractAttribute(code, "as", "")

		if rel == "" || href == "" {
			continue
		}
		if rel != "stylesheet" && (rel == "preload" && as == "") {
			continue
		}

		file := filepath.Join(os.Dir(item.Path), href)
		related = append(related, compactor.Related{
			Type:   "other",
			Source: code,
			Path:   href,
			Item:   compactor.Get(file),
		})

	}

	return related, nil
}

// Minify HTML content
func Minify(content string) (string, error) {

	file := os.TemporaryFile("minify.html")
	defer os.Delete(file)

	err := os.Write(file, content, 0775)

	if err != nil {
		return content, err
	}

	_, err = os.Exec(
		"html-minifier",
		"--output", file,
		"--collapse-whitespace",
		"--conservative-collapse",
		"--remove-comments",
		"--remove-script-type-attributes",
		"--remove-style-link-type-attributes",
		"--use-short-doctype",
		"--minify-urls", "true",
		"--minify-css", "true",
		"--minify-js", "true",
		"--ignore-custom-fragments", "/{{[{]?(.*?)[}]?}}/",
		file,
	)

	if err != nil {
		return content, err
	}

	content, err = os.Read(file)

	return content, err
}

// MergeContent returns the content of the item with the replaced partials dependencies
func MergeContent(item *compactor.Item) string {

	if !item.Exists {
		return ""
	}

	content := item.Content

	for _, related := range item.Related {
		if related.Type == "partial" && related.Item.Exists {

			// Solves recursively
			content = strings.Replace(
				content,
				related.Source,
				MergeContent(related.Item),
				1,
			)

		}
	}

	return content
}

// Execute processor
func Execute(bundle *compactor.Bundle) error {

	content := MergeContent(bundle.Item)

	for _, related := range bundle.Item.Related {

		if related.Type != "other" {
			continue
		}

		var err error
		var file string

		extension := os.Extension(related.Path)
		path := related.Path

		if extension == ".css" {
			file, err = css.Resolve(path)
		} else if extension == ".sass" || extension == ".scss" {
			file, err = css.Resolve(path)
		} else if extension == ".js" {
			file, err = javascript.Resolve(path)
		} else if extension == ".ts" {
			file, err = javascript.Resolve(path)
		} else {
			continue
		}

		if err != nil {
			return err
		}

		content = strings.Replace(content, path, file, 1)

	}

	destination := bundle.ToDestination(bundle.Item.Path)
	perm := bundle.Item.Permission
	err := os.Write(destination, content, perm)

	if err != nil {
		return err
	}

	return nil
}

// Optimize processor
func Optimize(bundle *compactor.Bundle) error {

	if !bundle.ShouldCompress(bundle.Item.Path) {
		return nil
	}

	destination := bundle.ToDestination(bundle.Item.Path)
	content, err := os.Read(destination)

	if err != nil {
		return err
	}

	content, err = Minify(content)

	if err != nil {
		return err
	}

	perm := bundle.Item.Permission
	err = os.Write(destination, content, perm)

	if err != nil {
		return err
	}

	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:  "html",
		Extensions: []string{".html", ".htm"},
		Init:       Init,
		Related:    Related,
		Execute:    Execute,
		Optimize:   Optimize,
		Delete:     generic.Delete,
		Resolve:    generic.Resolve,
	}
}
