package parse

import (
	"bytes"
	"html/template"

	"github.com/aghorui/burlough/util"
	"github.com/aghorui/burlough/blog"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)


func ParseBlogFile(src []byte) (blog.BlogFileContents, bool, error) {
	var dest bytes.Buffer
	var parseResult blog.BlogFileContents
	var noMetadata bool = false

	// TODO put this into the context.
	md := goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{},
			highlighting.NewHighlighting(
				highlighting.WithStyle("tango"),
			),
			extension.GFM,
		),
	)

	pc := parser.NewContext()

	err := md.Convert(src, &dest, parser.WithContext(pc))

	if err != nil {
		util.Error(err)
		return parseResult, noMetadata, err
	}

	metadata := frontmatter.Get(pc)

	if metadata != nil {
		if err := metadata.Decode(&parseResult); err != nil {
			return parseResult, noMetadata, err
		}
	} else {
		noMetadata = true
	}

	parseResult.Content = template.HTML(dest.Bytes())

	if parseResult.Title == "" {
		parseResult.Title = "(No Title)"
	}

	return parseResult, noMetadata, nil
}