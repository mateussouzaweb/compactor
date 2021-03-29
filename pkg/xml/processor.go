package xml

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/xml"
)

// XML minify
func Minify(content string) (string, error) {

	m := minify.New()
	m.AddFunc("generic", xml.Minify)

	content, err := m.String("generic", content)

	return content, err
}

// XML processor
func Processor(context *compactor.Context, options *compactor.Options) error {

	content, perm, err := compactor.ReadFileAndPermission(context.Source)

	if err != nil {
		return err
	}

	if options.Minify {
		content, err = Minify(content)

		if err != nil {
			return err
		}
	}

	err = compactor.WriteFile(context.Destination, content, perm)

	if err == nil {
		context.Processed = true
	}

	return err
}
