package static

import (
	"html/template"
	"os"
	"path/filepath"
	"testing"

	"github.com/aghorui/burlough/blog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestTagFunctions(t *testing.T) {
	assert.Equal(t,
		TagsToStringTOML(blog.Tags([]string{ "a", "1", "2", "\"b"})),
		template.HTML("[ \"a\", \"1\", \"2\", \"\\\"b\" ]"),
		"should match expected output")

	assert.Equal(t,
		TagsToStringTOML(blog.Tags([]string{})),
		template.HTML("[]"),
		"should match expected output")

	assert.Equal(t,
		TagsToStringTOML(blog.Tags(nil)),
		template.HTML("[]"),
		"should match expected output")

	assert.Equal(t,
		TagsToStringYAML(blog.Tags([]string{ "a", "1", "2", "\"b"})),
		template.HTML("[ a, 1, 2, \"b ]"),
		"should match expected output")

	assert.Equal(t,
		TagsToStringYAML(blog.Tags([]string{})),
		template.HTML("[]"),
		"should match expected output")

	assert.Equal(t,
		TagsToStringYAML(blog.Tags(nil)),
		template.HTML("[]"),
		"should match expected output")

}

func TestWriteBlogFile(t *testing.T) {
	dir := t.TempDir()
	b := GenerateDefaultBlogFileContents(blog.BlogFileContents{})

	f, err := os.OpenFile(filepath.Join(dir, "test_toml.md"), os.O_CREATE | os.O_WRONLY, 0644)
	require.NoError(t, err, "should be able to create markdown file")

	err = WriteBlogFile(blog.TOML, f, b)
	assert.NoError(t, err, "should be able to write TOML blog file")

	err = WriteBlogFile(blog.YAML, f, b)
	assert.NoError(t, err, "should be able to write YAML blog file")
}