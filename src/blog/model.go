// Meant to keep data structures used within the program

package blog

import (
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aghorui/burlough/util"
	"github.com/otiai10/copy"
)

type MetadataType int

const (
	TOML MetadataType = 0
	YAML MetadataType = 1
)

type FileHash string

type Tags []string

func (t Tags) String() string {
	return strings.Join([]string(t), ", ")
}

// The Blog file's contents after parsing it
type BlogFileContents struct {
	Title string `yaml:"title"`
	Desc string `yaml:"desc"`
	Tags []string `yaml:"tags"`
	Content template.HTML
}

// Data for a Given Blog File
type BlogMetadata struct {
	Path string       // Relative path to file ('a.md', 'a/b.md', etc.)
	Hash FileHash     // current SHA1 sum of the file
	Updated time.Time // Update date of the file (bumped if there is a hash mismatch)
	Created time.Time // Creation date of the file
}

// Complete representation of a blog file (Data + Metadata).
type BlogFile struct {
	BlogFileContents
	BlogMetadata
}

// Parameters for a blog project unmarshalled from a config file.
type ConfigFileParams struct {
	Title string                        // Title of the blog
	Desc string                         // Short description of the blog. Goes in the <meta> tags.
	Tags []string                       // Tags for the blog. Goes in the <meta> tags.
	BlogURLPathPrefix string            // This prefix will be added to all in-site URLs that are generated.
	RenderPath string                   // Path to where the rendered files should be put.
	TemplatePath string                 // Path to template.
	UseFileTimestampAsCreationDate bool // Use File Timestamp As Creation date.
	Files []BlogMetadata                // List of blog markdown files.
	MetadataType MetadataType           // Type of the blog file metadata (TOML/YAML)
}

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

// Input given to templates for generating the final HTML.
type BlogTemplateEntry struct {
	Title string
	Desc string
	GlobalDesc string
	Tags Tags
	GlobalTags Tags
	Created string
	Updated string
	URL string
	Content template.HTML
}

func PrepareBlogTemplateEntry(b BlogFile, finalPath string, globalDesc string, globalTags Tags) BlogTemplateEntry {
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

// sort.Interface Implementation for BlogMetadata.
type BlogMetadataSortCreatedDescending []BlogMetadata

func (b BlogMetadataSortCreatedDescending) Len() int {
	return len(b)
}

func (b BlogMetadataSortCreatedDescending) Less(i, j int) bool {
	c := time.Time.Compare(b[i].Created, b[j].Created)
	if c < 0 {
		return false
	} else {
		return true
	}
}

func (b BlogMetadataSortCreatedDescending) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}