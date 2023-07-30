package static

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/aghorui/burlough/blog"
	"github.com/aghorui/burlough/constants"
	"github.com/aghorui/burlough/util"

	"github.com/otiai10/copy"
)

const DefaultTemplateContents template.HTML = `
# This is the title

--------------------

This is a paragraph.

This is another paragraph.

**This is bold** *This is italic* ~~This is strikethrough~~
`

//go:embed default_export_template/*
var defaultExportTemplate embed.FS

//go:embed file_template/none.md
var BlogTemplateNoneData []byte

//go:embed file_template/toml.md
var BlogTemplateTOMLData []byte

//go:embed file_template/yaml.md
var BlogTemplateYAMLData []byte

var BlogTemplateNone = template.Must(template.New("none.md").Parse(string(BlogTemplateNoneData)))
var BlogTemplateTOML = template.Must(template.New("toml.md").Parse(string(BlogTemplateTOMLData)))
var BlogTemplateYAML = template.Must(template.New("yaml.md").Parse(string(BlogTemplateYAMLData)))

var DefaultBlogTemplate blog.BlogTemplate = func() blog.BlogTemplate {
	dir := GetDefaultExportTemplateFS()

	return blog.BlogTemplate{
		TemplateFS:  &dir,
		FrontPage: template.Must(template.New("front_page.html").ParseFS(GetDefaultExportTemplateFS(), "front_page.html")),
		BlogPage:  template.Must(template.New("blog_page.html").ParseFS(GetDefaultExportTemplateFS(), "blog_page.html")),
		IndexPage: template.Must(template.New("blog_list.html").ParseFS(GetDefaultExportTemplateFS(), "blog_list.html")),
	}
}()

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
		return err
	}

	return nil
}

// Gets the directory listing of default_export_template
func GetDefaultExportTemplateFiles() []fs.DirEntry {
	files, err := defaultExportTemplate.ReadDir("default_export_template")

	if err != nil {
		util.LogErr(err)
		panic(err)
	}

	return files
}

// Gets the FS for default_export_template because we can't do subdirs directly.
func GetDefaultExportTemplateFS() fs.FS {
	dir, err := fs.Sub(defaultExportTemplate, "default_export_template")

	if err != nil {
		util.LogErr(err)
		panic(err)
	}

	return dir
}


func DumpDefaultExportTemplate(dest string) error {
	finalDest := filepath.Join(dest, constants.AppName + "_default_export_template")

	err := os.MkdirAll(finalDest, 0755)
	if err != nil {
		return err
	}

	// This is a weird thing. I have to explicitly set the permissions of the
	// embed.FS files to get the actually correct permissions ORed with the
	// supposed umask. 0644 seems to get the job done.
	err = copy.Copy("default_export_template", finalDest, copy.Options{
		FS: defaultExportTemplate,
		PermissionControl: copy.AddPermission(0644),
	})

	if err != nil {
		return err
	}

	return nil
}