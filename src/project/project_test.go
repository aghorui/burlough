package project

import (
	"path/filepath"
	"testing"

	"github.com/aghorui/burlough/blog"
	"github.com/aghorui/burlough/static"
	"github.com/aghorui/burlough/util"
	"github.com/stretchr/testify/assert"
)

func TestProject(t *testing.T) {
	dir := t.TempDir()
	outDir := t.TempDir()

	var b blog.ConfigFileParams

	b.RenderPath = outDir

	util.GenerateTestMarkdownFiles(dir)

	{
		state, log, err := Init(dir, b, true)

		assert.NoError(t, err, "there shouldn't be any errors during init")
		assert.Greater(t, len(log), 0, "there should be files recorded during scan")
		assert.Equal(t, dir, state.BasePath, "the base path should be the same as the one specified")

		err = state.WriteConfig()

		assert.NoError(t, err, "there shouldn't be any errors during project file write")
		assert.FileExists(t, filepath.Join(dir, ProjectConfigFileName), "project file should be made in dir")

		name, err := state.NewFile(static.GenerateDefaultBlogFileContents(blog.BlogFileContents{}), "")

		assert.NotEmpty(t, name, "file path should not be empty")
		assert.NoError(t, err, "there shouldn't be any errors on adding a new file to project")

		log, err = state.Scan()

		assert.Greater(t, len(log), 0, "there should be files recorded during scan after adding a new file")
		assert.NoError(t, err, "there shouldn't be any errors during project scan")

		err = state.Render(outDir)

		assert.Nil(t, err, "there shouldn't be any errors during project render")
	}
}