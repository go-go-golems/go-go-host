package doc

import (
	"embed"

	"github.com/go-go-golems/glazed/pkg/help"
)

//go:embed *.md
var docFS embed.FS

// DocFS returns the embedded filesystem containing the agent CLI docs.
func DocFS() embed.FS { return docFS }

func AddDocToHelpSystem(helpSystem *help.HelpSystem) error {
	return helpSystem.LoadSectionsFromFS(docFS, ".")
}
