package render

import (
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/aghorui/burlough/blog"
	"github.com/aghorui/burlough/static"
)

// Checks the Blog Parsing Function for TOML
func TestRendering(t *testing.T) {
	dir := t.TempDir()

	err := os.WriteFile(path.Join(dir, "a.md"), static.BlogTemplateTOMLData, 0644)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	data, err := os.ReadFile(path.Join(dir, "a.md"))

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	log.Printf("%v", string(data))

	var p project.ProjectState = project.ProjectState{
		BasePath: dir,
		ConfigFileParams: blog.ConfigFileParams{
			RenderPath: dir,
			TemplatePath: "",
			UseFileTimestampAsCreationDate: true,
			Files: []blog.BlogMetadata{
				{
					Path: "a.md",
					Hash: "-",
					Updated: time.Now(),
					Created: time.Now(),
				},{
					Path: "a.md",
					Hash: "-",
					Updated: time.Now(),
					Created: time.Now(),
				},{
					Path: "a.md",
					Hash: "-",
					Updated: time.Now(),
					Created: time.Now(),
				},{
					Path: "a.md",
					Hash: "-",
					Updated: time.Now(),
					Created: time.Now(),
				},{
					Path: "a.md",
					Hash: "-",
					Updated: time.Now(),
					Created: time.Now(),
				},
			},
		},
	}

	err = Render(p)

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	data, err = os.ReadFile(path.Join(dir, "a.html"))

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	log.Printf("[[%v]] %v", len(data), string(data))
}