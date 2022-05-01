package html

import (
	"regexp"
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/css"
	"github.com/mateussouzaweb/compactor/pkg/generic"
	"github.com/mateussouzaweb/compactor/pkg/javascript"
	"github.com/mateussouzaweb/compactor/pkg/typescript"
)

// Init processor
func Init(bundle *compactor.Bundle) error {
	return os.NodeRequire("html-minifier", "html-minifier")
}

// ExtractAttribute find the value of the attribute
func ExtractAttribute(html string, attribute string, defaultValue string) string {

	regex := regexp.MustCompile(attribute + `=("([^"]*)"|'([^']*)')`)
	match := regex.FindStringSubmatch(html)

	if match != nil {
		return strings.Trim(match[1], `'"`)
	}

	return defaultValue
}

// Related processor
func Related(item *compactor.Item) ([]compactor.Related, error) {

	var related []compactor.Related

	// Detect imports
	regex := regexp.MustCompile(`<!-- @import ?("(.+)"|'(.+)') -->`)
	matches := regex.FindAllStringSubmatch(item.Content, -1)

	for _, match := range matches {
		source := match[0]
		path := strings.Trim(match[1], `'"`)

		file := os.Resolve(path, os.Dir(item.Path))
		related = append(related, compactor.Related{
			Type:       "partial",
			Dependency: true,
			Source:     source,
			Path:       path,
			Item:       compactor.Get(file),
		})
	}

	// Detect scripts
	regex = regexp.MustCompile(`<script(.+)?>(.+)?</script>`)
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

		file := os.Resolve(src, os.Dir(item.Path))
		related = append(related, compactor.Related{
			Type:       "other",
			Dependency: false,
			Source:     code,
			Path:       src,
			Item:       compactor.Get(file),
		})

	}

	// Detect links
	regex = regexp.MustCompile(`<link(.+)?\/?>`)
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

		file := os.Resolve(href, os.Dir(item.Path))
		related = append(related, compactor.Related{
			Type:       "other",
			Dependency: false,
			Source:     code,
			Path:       href,
			Item:       compactor.Get(file),
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

		if related.Dependency {
			continue
		}

		var err error
		var file string

		extension := os.Extension(related.Path)
		path := related.Path

		if extension == ".css" {
			file, err = css.Resolve(path, bundle.Item)
		} else if extension == ".sass" || extension == ".scss" {
			file, err = css.Resolve(path, bundle.Item)
		} else if extension == ".js" || extension == ".mjs" {
			file, err = javascript.Resolve(path, bundle.Item)
		} else if extension == ".ts" || extension == ".mts" {
			file, err = typescript.Resolve(path, bundle.Item)
		} else if extension == ".tsx" || extension == ".jsx" {
			file, err = typescript.Resolve(path, bundle.Item)
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
		Resolve:    generic.Resolve,
		Execute:    Execute,
		Optimize:   Optimize,
		Delete:     generic.Delete,
	}
}
