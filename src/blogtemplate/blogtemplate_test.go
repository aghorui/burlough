package blogtemplate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aghorui/burlough/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateFunctions(t *testing.T) {
	dir := t.TempDir()

	assert.NotPanics(t, func() { DumpDefaultExportTemplate(dir) }, "should not panic")
	assert.NotPanics(t, func() { GetDefaultExportTemplateFiles() }, "should not panic")
	assert.NotPanics(t, func() { GetDefaultExportTemplateFS() }, "should not panic")
}

func TestTemplate(t *testing.T) {
	dir := t.TempDir()

	templatePath := filepath.Join(dir, constants.AppName + "_default_export_template")

	require.NoError(t, DumpDefaultExportTemplate(dir))

	require.DirExists(t, templatePath)

	tmpl, err := LoadTemplate(os.DirFS(templatePath))
	assert.NoError(t, err, "there shouldn't be any error while loading the default template")

	assert.NoError(t, tmpl.CopyAssetsToFolder(filepath.Join(dir, "assets")), "there shouldn't be any error while copying template assets to a folder")
}