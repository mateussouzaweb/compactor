package html

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/css"
	"github.com/mateussouzaweb/compactor/pkg/generic"
	"github.com/mateussouzaweb/compactor/pkg/javascript"
)

// HTML minify
func Minify(content string) (string, error) {

	file, err := ioutil.TempFile("", "minify.*.html")

	if err != nil {
		return content, err
	}

	defer os.Delete(file.Name())

	_, err = file.WriteString(content)

	if err != nil {
		return content, err
	}

	_, err = os.Exec(
		"html-minifier",
		"--output", file.Name(),
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
		file.Name(),
	)

	if err != nil {
		return content, err
	}

	content, err = os.Read(file.Name())

	return content, err
}

// ExtractAttribute find the value of the attribute
func ExtractAttribute(html string, attribute string, defaultValue string) string {

	regex := regexp.MustCompile(attribute + `=[\"\']([^"']*)[\"\']`)
	match := regex.FindStringSubmatch(html)
	value := ""

	if match != nil {
		value = match[1]
	}

	if value == "" {
		value = defaultValue
	}

	return value
}

// HTML Format method
func Format(content string) (string, error) {

	var err error
	var file string

	script := `(?i)<script(.+)?>(.+)?</script>`
	regex := regexp.MustCompile(script)
	matches := regex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {

		code := match[0]
		src := ExtractAttribute(code, "src", "")

		// Ignore protocol scripts, only handle relative and absolute paths
		if strings.Contains(src, "://") {
			continue
		}

		file, err = javascript.ResolveProcessor(src)

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
			file, err = css.ResolveProcessor(href)
		} else if extension == ".sass" || extension == ".scss" {
			file, err = css.ResolveProcessor(href)
		} else if extension == ".js" {
			file, err = javascript.ResolveProcessor(href)
		} else if extension == ".ts" {
			file, err = javascript.ResolveProcessor(href)
		}

		if err != nil {
			return content, err
		}

		content = strings.Replace(content, href, file, 1)

	}

	return content, nil
}

// HTML processor
func RunProcessor(bundle *compactor.Bundle) error {

	if bundle.ShouldOutputToMany() {

		for _, item := range bundle.Items {

			if !item.Exists {
				continue
			}

			content, err := Format(item.Content)

			if err != nil {
				return err
			}

			if bundle.ShouldCompress(item.Path) {
				content, err = Minify(content)
				if err != nil {
					return err
				}
			}

			destination := bundle.ToDestination(item.Path)
			err = os.Write(destination, content, item.Permission)

			if err != nil {
				return err
			}

			bundle.Processed(item.Path)

		}

		return nil
	}

	content := ""
	for _, item := range bundle.Items {
		if item.Exists {
			content += item.Content
		}
	}

	content, err := Format(content)

	if err != nil {
		return err
	}

	destination := bundle.ToDestination(bundle.Destination.File)

	if bundle.ShouldCompress(destination) {
		content, err = Minify(content)
		if err != nil {
			return err
		}
	}

	perm := bundle.Items[0].Permission
	err = os.Write(destination, content, perm)

	if err == nil {
		bundle.Written(destination)
	}

	return err
}

func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Extensions: []string{".html", ".htm"},
		Run:        RunProcessor,
		Delete:     generic.DeleteProcessor,
		Resolve:    generic.ResolveProcessor,
	}
}
