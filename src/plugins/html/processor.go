package html

import (
	"regexp"
	"strings"

	"github.com/mateussouzaweb/compactor/src/plugins/generic"
	"github.com/mateussouzaweb/compactor/src/processor"
	"github.com/mateussouzaweb/compactor/src/system"
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
func Related(options *processor.Options, file *processor.File) ([]processor.Related, error) {

	var related []processor.Related

	// Detect imports
	regex := regexp.MustCompile(`<!-- @import ?("(.+)"|'(.+)') -->`)
	matches := regex.FindAllStringSubmatch(file.Content, -1)
	extensions := []string{".html", ".htm"}

	for _, match := range matches {
		source := match[0]
		path := strings.Trim(match[1], `'"`)
		filePath := system.Resolve(path, extensions, system.Dir(file.Path))

		if processor.GetFile(filePath).Path != "" {
			related = append(related, processor.Related{
				Type:       "partial",
				Dependency: true,
				Source:     source,
				Path:       path,
				File:       processor.GetFile(filePath),
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

		filePath := system.Resolve(src, extensions, system.Dir(file.Path))

		if processor.GetFile(filePath).Path != "" {
			related = append(related, processor.Related{
				Type:       "other",
				Dependency: false,
				Source:     code,
				Path:       src,
				File:       processor.GetFile(filePath),
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

		filePath := system.Resolve(href, []string{}, system.Dir(file.Path))

		if processor.GetFile(filePath).Path != "" {
			related = append(related, processor.Related{
				Type:       "other",
				Dependency: false,
				Source:     code,
				Path:       href,
				File:       processor.GetFile(filePath),
			})
		}

	}

	return related, nil
}

// Minify HTML content
func Minify(content string) (string, error) {

	file := system.TemporaryFile("minify.html")
	defer system.Delete(file)

	err := system.Write(file, content, 0775)
	if err != nil {
		return content, err
	}

	_, err = system.Exec(
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

	content, err = system.Read(file)

	return content, err
}

// MergeContent returns the content of the item with the replaced partials dependencies
func MergeContent(file *processor.File) string {

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
func Transform(options *processor.Options, file *processor.File) error {

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

	err := system.Write(destination, content, perm)
	if err != nil {
		return err
	}

	return nil
}

// Optimize processor
func Optimize(options *processor.Options, file *processor.File) error {

	if !options.ShouldCompress(file.Path) {
		return nil
	}

	content, err := system.Read(file.Destination)
	if err != nil {
		return err
	}

	content, err = Minify(content)
	if err != nil {
		return err
	}

	perm := file.Permission
	err = system.Write(file.Destination, content, perm)
	if err != nil {
		return err
	}

	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *processor.Plugin {
	return &processor.Plugin{
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
