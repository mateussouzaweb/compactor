package html

import (
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

// Related processor
func Related(item *compactor.Item) ([]compactor.Related, error) {

	var patterns []generic.FindPattern
	patterns = append(patterns, generic.FindPattern{
		Type:     "partial",
		Regex:    "<!-- @include \"(.+)\" -->",
		SubMatch: 1,
	})

	return generic.FindRelated(item, patterns)
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

// Format HTML method
func Format(content string) (string, error) {

	var err error
	var file string

	script := `(?i)<script(.+)?>(.+)?</script>`
	regex := regexp.MustCompile(script)
	matches := regex.FindAllStringSubmatch(content, -1)

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

		file, err = javascript.Resolve(src)

		if err != nil {
			return content, err
		}

		content = strings.Replace(content, src, file, 1)

	}

	link := `(?i)<link(.+)?\/?>`
	regex = regexp.MustCompile(link)
	matches = regex.FindAllStringSubmatch(content, -1)

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

		file = href
		extension := os.Extension(href)

		if extension == ".css" {
			file, err = css.Resolve(href)
		} else if extension == ".sass" || extension == ".scss" {
			file, err = css.Resolve(href)
		} else if extension == ".js" {
			file, err = javascript.Resolve(href)
		} else if extension == ".ts" {
			file, err = javascript.Resolve(href)
		}

		if err != nil {
			return content, err
		}

		content = strings.Replace(content, href, file, 1)

	}

	return content, nil
}

// Execute processor
func Execute(bundle *compactor.Bundle) error {

	content := bundle.Item.MergedContent()
	content, err := Format(content)

	if err != nil {
		return err
	}

	destination := bundle.ToDestination(bundle.Item.Path)
	perm := bundle.Item.Permission
	err = os.Write(destination, content, perm)

	if err != nil {
		return err
	}

	bundle.Processed(bundle.Item.Path)

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

	bundle.Optimized(bundle.Item.Path)

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
