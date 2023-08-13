package blogtemplate

import (
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/aghorui/burlough/blog"
	"github.com/aghorui/burlough/constants"
	"github.com/aghorui/burlough/static"
	"github.com/aghorui/burlough/util"
	"github.com/otiai10/copy"
)

// Contains all required template files
type BlogTemplate struct {
	TemplateFS *fs.FS            // Asset Directory Filesystem
	FrontPage *template.Template // Front Page Template
	IndexPage *template.Template // Index Page Template
	BlogPage *template.Template  // Blog Page Template
}

const IndexPageTemplateFileName = "blog_list.html"
const FrontPageTemplateFileName = "front_page.html"
const BlogPageTemplateFileName  = "blog_page.html"

// Copies asset files of the template to the desired folder.
func (b BlogTemplate) CopyAssetsToFolder(dest string) error {

	if b.TemplateFS == nil {
		// Nothing to copy.
		return nil
	}

	finalDest := filepath.Join(dest, "assets")

	err := os.MkdirAll(finalDest, 0755)
	if err != nil {
		return util.Error(err)
	}

	// This is a weird thing. I have to explicitly set the permissions of the
	// embed.FS files to get the actually correct permissions ORed with the
	// supposed umask. 0644 seems to get the job done.
	err = copy.Copy("assets", finalDest, copy.Options{
		FS: *b.TemplateFS,
		PermissionControl: copy.AddPermission(0644),
	})

	if err != nil {
		return util.Error(err)
	}

	return nil
}

// Loads a template into a struct
func LoadTemplate(folder fs.FS) (BlogTemplate, error) {
	var t BlogTemplate
	var err error

	t.TemplateFS = &folder

	if err != nil {
		return t, util.Error(err)
	}

	t.FrontPage, err = template.New(FrontPageTemplateFileName).ParseFS(folder, FrontPageTemplateFileName)
	if err != nil {
		return t, util.Error(err)
	}

	t.BlogPage, err = template.New(BlogPageTemplateFileName).ParseFS(folder, BlogPageTemplateFileName)
	if err != nil {
		return t, util.Error(err)
	}

	t.IndexPage, err = template.New(IndexPageTemplateFileName).ParseFS(folder, IndexPageTemplateFileName)
	if err != nil {
		return t, util.Error(err)
	}

	return t, nil
}

// Gets the directory listing of default_export_template
func GetDefaultExportTemplateFiles() []fs.DirEntry {
	files, err := static.DefaultExportTemplate.ReadDir("default_export_template")

	if err != nil {
		util.LogErr(err)
		panic(err)
	}

	return files
}

// Gets the FS for default_export_template because we can't do subdirs directly.
func GetDefaultExportTemplateFS() fs.FS {
	dir, err := fs.Sub(static.DefaultExportTemplate, "default_export_template")

	if err != nil {
		util.LogErr(err)
		panic(err)
	}

	return dir
}

var DefaultBlogTemplate BlogTemplate = func() BlogTemplate {
	dir := GetDefaultExportTemplateFS()

	return BlogTemplate{
		TemplateFS:  &dir,
		FrontPage: template.Must(template.New("front_page.html").ParseFS(GetDefaultExportTemplateFS(), "front_page.html")),
		BlogPage:  template.Must(template.New("blog_page.html").ParseFS(GetDefaultExportTemplateFS(), "blog_page.html")),
		IndexPage: template.Must(template.New("blog_list.html").ParseFS(GetDefaultExportTemplateFS(), "blog_list.html")),
	}
}()

func DumpDefaultExportTemplate(dest string) error {
	finalDest := filepath.Join(dest, constants.AppName + "_default_export_template")

	err := os.MkdirAll(finalDest, 0755)
	if err != nil {
		return util.Error(err)
	}

	// This is a weird thing. I have to explicitly set the permissions of the
	// embed.FS files to get the actually correct permissions ORed with the
	// supposed umask. 0644 seems to get the job done.
	err = copy.Copy("default_export_template", finalDest, copy.Options{
		FS: static.DefaultExportTemplate,
		PermissionControl: copy.AddPermission(0644),
	})

	if err != nil {
		return util.Error(err)
	}

	return nil
}

// Input given to templates for generating the final HTML.
type BlogTemplateEntry struct {
	Title string
	Desc string
	GlobalDesc string
	Tags blog.Tags
	GlobalTags blog.Tags
	Created string
	Updated string
	URL string
	Content template.HTML
}

func PrepareBlogTemplateEntry(b blog.BlogFile, finalPath string, globalDesc string, globalTags blog.Tags) BlogTemplateEntry {
	return BlogTemplateEntry{
		Title: b.Title,
		Desc: b.Desc,
		GlobalDesc: globalDesc,
		Tags: b.Tags,
		GlobalTags: globalTags,
		Created: util.GetStandardTimestampString(b.Created),
		Updated: util.GetStandardTimestampString(b.Updated),
		URL: filepath.Join("./", finalPath),
		Content: b.Content,
	}
}