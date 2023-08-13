package parse

import (
	"testing"

	"github.com/aghorui/burlough/util"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	{
		_, noMetadata, err := ParseBlogFile(util.GetTestFile("markdown/standard_toml.md"));
		assert.False(t, noMetadata, "there should be metadata in standard_toml.md")
		assert.NoError(t, err, "there shouldn't be any errors while parsing standard_toml.md")
	}

	{
		_, noMetadata, err := ParseBlogFile(util.GetTestFile("markdown/no_metadata.md"));
		assert.True(t, noMetadata, "there shouldn't be metadata in no_metadata.md")
		assert.NoError(t, err, "there shouldn't be any errors while parsing the file.")
	}
}