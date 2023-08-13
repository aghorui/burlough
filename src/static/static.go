package static

import (
	"embed"
	"html/template"
	"io"
	"strconv"
	"strings"

	"github.com/aghorui/burlough/blog"
	"github.com/aghorui/burlough/util"
)


const DefaultTemplateContents template.HTML = `
# This is the title

--------------------

This is a paragraph.

This is another paragraph.

**This is bold** *This is italic* ~~This is strikethrough~~
`

//go:embed default_export_template/*
var DefaultExportTemplate embed.FS

//go:embed file_template/none.md
var BlogTemplateNoneData []byte

//go:embed file_template/toml.md
var BlogTemplateTOMLData []byte

//go:embed file_template/yaml.md
var BlogTemplateYAMLData []byte

// Coverts a blog tagset to a TOML parseable string array
func TagsToStringTOML(t blog.Tags) template.HTML {
	var b strings.Builder
	var l int = 0

	for _, v := range t {
		l += len(v) + 4 // 2 quotes + comma + space
	}

	l += 4 // "[" + " " + " " + "]"

	b.Grow(l)

	b.WriteString("[ ")

	for _, v := range t[:len(t) - 1] {
		b.WriteString(strconv.Quote(v))
		b.WriteString(", ")
	}

	b.WriteString(strconv.Quote(t[len(t) - 1]))

	b.WriteString(" ]")

	return template.HTML(b.String())
}

// Converts a blog tagset to a YAML parseable string array
func TagsToStringYAML(t blog.Tags) template.HTML {
	return "[ " + template.HTML(strings.Join(t, ", ")) + " ]"
}

var BlogTemplateNone = template.Must(template.New("none.md").Parse(string(BlogTemplateNoneData)))

var BlogTemplateTOML = template.Must(
	template.New("toml.md").Funcs(
		template.FuncMap{
			"tagsToString": TagsToStringTOML,
		}).Parse(string(BlogTemplateTOMLData)))

var BlogTemplateYAML = template.Must(template.New("yaml.md").Funcs(
		template.FuncMap{
			"tagsToString": TagsToStringYAML,
		}).Parse(string(BlogTemplateYAMLData)))

// Generates a default blog file for use. Fields can be overriden with `n`
func GenerateDefaultBlogFileContents(n blog.BlogFileContents) blog.BlogFileContents {
	d := blog.BlogFileContents{
		Title: "Enter Your Title Here",
		Tags: []string{"default_tag_1", "default_tag_2"},
		Desc: "Enter a description here",
		Content: DefaultTemplateContents,
	}

	if n.Title != ""   { d.Title = n.Title }
	if n.Tags != nil   { d.Tags = n.Tags }
	if n.Desc != ""    { d.Desc = n.Desc }
	if n.Content != "" { d.Content = n.Content }

	return d
}

func WriteBlogFile(t blog.MetadataType, wr io.Writer, b blog.BlogFileContents) error {
	var err error

	if t == blog.TOML {
		err = BlogTemplateTOML.Execute(wr, b)
	} else if t == blog.YAML {
		err = BlogTemplateYAML.Execute(wr, b)
	}

	if err != nil {
		return util.Error(err)
	}

	return nil
}