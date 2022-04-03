package typescript

import (
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
)

// TSRename struct
type TSRename struct {
	Name     string
	Source   string
	FromPath string
	FromName string
	ToPath   string
	ToName   string
}

// FindRenames will create a list of possible renames based on internal rules
func FindRenames(bundle *compactor.Bundle, item *compactor.Item) []TSRename {

	var changes []TSRename

	from := bundle.ToDestination(item.Path)
	from = bundle.ToExtension(from, ".js")
	to := bundle.ToHashed(from, item.Checksum)

	changes = append(changes, TSRename{
		Name:     item.Name,
		Source:   "",
		FromPath: from,
		FromName: os.File(from),
		ToPath:   to,
		ToName:   os.File(to),
	})

	for _, related := range item.Related {
		if related.Item.Exists && related.Type == "import" {

			from := bundle.ToDestination(related.Item.Path)
			from = bundle.ToExtension(from, ".js")
			to := bundle.ToHashed(from, related.Item.Checksum)

			changes = append(changes, TSRename{
				Name:     related.Item.Name,
				Source:   related.Source,
				FromPath: from,
				FromName: os.File(from),
				ToPath:   to,
				ToName:   os.File(to),
			})

			changes = append(changes, FindRenames(bundle, related.Item)...)

		}
	}

	return changes
}

// RenameDestination fix names of destinations files
func RenameDestination(bundle *compactor.Bundle) error {

	// Since we have the dependency graph of the bundle
	// We can predit the files that typescript will generate to create the list of changes
	changes := FindRenames(bundle, bundle.Item)

	// With the list of changes, we then rename the files to the hashed version
	var changed bool
	for _, change := range changes {

		// Make basic verifications
		if change.FromPath == change.ToPath {
			continue
		}
		if !os.Exist(change.FromPath) {
			continue
		}

		// Rename the file
		err := os.Rename(change.FromPath, change.ToPath)

		if err != nil {
			return err
		}

		// Also rename source-maps
		if os.Exist(change.FromPath + ".map") {
			err = os.Rename(change.FromPath+".map", change.ToPath+".map")

			if err != nil {
				return err
			}
		}

		changed = true
	}

	if !changed {
		return nil
	}

	// Files has been renamed, we need to fix source-map references
	for _, change := range changes {

		err := os.Replace(change.ToPath,
			"sourceMappingURL="+change.FromName+".map",
			"sourceMappingURL="+change.ToName+".map",
		)

		if err != nil {
			return err
		}

		if os.Exist(change.ToPath + ".map") {

			err = os.Replace(change.ToPath+".map",
				"\"file\":\""+change.FromName+"\"",
				"\"file\":\""+change.ToName+"\"",
			)

			if err != nil {
				return err
			}

		}

	}

	// And finally we fix import references on the files
	for _, change := range changes {
		for _, update := range changes {

			if update.Source == "" {
				continue
			}

			newSource := strings.Replace(
				update.Source,
				update.Name,
				os.Name(update.ToName),
				1,
			)

			err := os.Replace(
				change.ToPath,
				update.Source,
				newSource,
			)

			if err != nil {
				return err
			}

		}
	}

	return nil
}
