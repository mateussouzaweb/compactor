package html

import (
	"regexp"
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

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
func Related(options *compactor.Options, file *compactor.File) ([]compactor.Related, error) {

	var related []compactor.Related

	// Detect imports
	regex := regexp.MustCompile(`<!-- @import ?("(.+)"|'(.+)') -->`)
	matches := regex.FindAllStringSubmatch(file.Content, -1)
	extensions := []string{".html", ".htm"}

	for _, match := range matches {
		source := match[0]
		path := strings.Trim(match[1], `'"`)
		filePath := os.Resolve(path, extensions, os.Dir(file.Path))

		if compactor.GetFile(filePath).Path != "" {
			related = append(related, compactor.Related{
				Type:       "partial",
				Dependency: true,
				Source:     source,
				Path:       path,
				File:       compactor.GetFile(filePath),
			})
		}
	}

	// Detect scripts
	regex = regexp.MustCompile(`<script(.+)?>(.+)?</script>`)
	matches = regex.FindAllStringSubmatch(file.Content, -1)
	extensions = []string{".js", ".mjs", ".jsx", ".ts", ".mts", ".tsx"}

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

		filePath := os.Resolve(src, extensions, os.Dir(file.Path))

		if compactor.GetFile(filePath).Path != "" {
			related = append(related, compactor.Related{
				Type:       "other",
				Dependency: false,
				Source:     code,
				Path:       src,
				File:       compactor.GetFile(filePath),
			})
		}

	}

	// Detect links
	regex = regexp.MustCompile(`<link(.+)?\/?>`)
	matches = regex.FindAllStringSubmatch(file.Content, -1)

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

		filePath := os.Resolve(href, []string{}, os.Dir(file.Path))

		if compactor.GetFile(filePath).Path != "" {
			related = append(related, compactor.Related{
				Type:       "other",
				Dependency: false,
				Source:     code,
				Path:       href,
				File:       compactor.GetFile(filePath),
			})
		}

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
func MergeContent(file *compactor.File) string {

	if !file.Exists {
		return ""
	}

	content := file.Content

	for _, related := range file.Related {
		if related.Type == "partial" && related.File.Exists {

			// Solves recursively
			content = strings.Replace(
				content,
				related.Source,
				MergeContent(related.File),
				1,
			)

		}
	}

	return content
}

// Transform processor
func Transform(options *compactor.Options, file *compactor.File) error {

	content := MergeContent(file)

	for _, related := range file.Related {

		if related.Dependency {
			continue
		}

		path := related.Path
		destination := options.CleanPath(related.File.Destination)

		if strings.HasPrefix(path, "/") {
			destination = "/" + destination
		} else if strings.HasPrefix(path, "./") {
			destination = "./" + destination
		}

		content = strings.Replace(content, path, destination, 1)

	}

	destination := file.Destination
	perm := file.Permission
	err := os.Write(destination, content, perm)

	if err != nil {
		return err
	}

	return nil
}

// Optimize processor
func Optimize(options *compactor.Options, file *compactor.File) error {

	if !options.ShouldCompress(file.Path) {
		return nil
	}

	content, err := os.Read(file.Destination)

	if err != nil {
		return err
	}

	content, err = Minify(content)

	if err != nil {
		return err
	}

	perm := file.Permission
	err = os.Write(file.Destination, content, perm)

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
		Init:       generic.Init,
		Shutdown:   generic.Shutdown,
		Resolve:    generic.Resolve,
		Related:    Related,
		Transform:  Transform,
		Optimize:   Optimize,
	}
}
