package typescript

import (
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
)

// TSRename struct
type TSRename struct {
	FromPath string
	FromName string
	ToPath   string
	ToName   string
}

// RenameDestination fix names of destinations files
func RenameDestination(bundle *compactor.Bundle) error {

	var changes []TSRename
	var changed bool

	// Since we have the dependency graph of the bundle
	// We can predit the files that typescript will generate to create the list of changes
	from := bundle.ToDestination(bundle.Item.Path)
	from = bundle.ToExtension(from, ".js")
	to := bundle.ToHashed(from, bundle.Item.Checksum)

	changes = append(changes, TSRename{
		FromPath: from,
		FromName: os.File(from),
		ToPath:   to,
		ToName:   os.File(to),
	})

	for _, related := range bundle.Item.Related {
		if related.Item.Exists && related.Type == "import" {

			from := bundle.ToDestination(related.Item.Path)
			from = bundle.ToExtension(from, ".js")
			to := bundle.ToHashed(from, related.Item.Checksum)

			changes = append(changes, TSRename{
				FromPath: from,
				FromName: os.File(from),
				ToPath:   to,
				ToName:   os.File(to),
			})

		}
	}

	// With the list of changes, we then rename the files to the hashed version
	for _, change := range changes {

		// Make basic verifications
		if change.FromPath == change.ToPath {
			continue
		}
		if !os.Exist(change.FromPath) {
			continue
		}

		// Rename the main file
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

	// And finally we fix import references on the main file
	main := changes[0]
	for _, related := range bundle.Item.Related {
		if related.Item.Exists && related.Type == "import" {

			err := os.Replace(main.ToPath,
				related.Path,
				strings.Replace(
					related.Path,
					os.Name(main.FromName),
					os.Name(main.ToName),
					1,
				),
			)

			if err != nil {
				return err
			}

		}
	}

	return nil
}
