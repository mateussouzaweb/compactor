package json

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/json"
)

// Json minify
func Minify(content string) (string, error) {

	m := minify.New()
	m.AddFunc("generic", json.Minify)

	content, err := m.String("generic", content)

	return content, err
}

// Json processor
func Processor(context *compactor.Context, options *compactor.Options) error {

	content, perm, err := compactor.ReadFileAndPermission(context.Source)

	if err != nil {
		return err
	}

	content, err = Minify(content)

	if err != nil {
		return err
	}

	err = compactor.WriteFile(context.Destination, content, perm)

	if err == nil {
		context.Processed = true
	}

	return err
}
