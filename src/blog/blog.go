// Meant to keep data structures used within the program

package blog

import (
	"html/template"
	"strings"
	"time"
)

type MetadataType int

const (
	Invalid MetadataType = -1
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
	Tags Tags `yaml:"tags"`
	Content template.HTML
}

// Data for a Given Blog File
type BlogMetadata struct {
	Path string       `json:"path"`    // Relative path to file ('a.md', 'a/b.md', etc.)
	Hash FileHash     `json:"hash"`    // current SHA1 sum of the file
	Updated time.Time `json:"updated"` // Update date of the file (bumped if there is a hash mismatch)
	Created time.Time `json:"created"` // Creation date of the file
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

// Complete representation of a blog file (Data + Metadata).
type BlogFile struct {
	BlogFileContents
	BlogMetadata
}

// Parameters for a blog project unmarshalled from a config file.
type ConfigFileParams struct {
	Title string                        `json:"title"`                // Title of the blog
	Desc string                         `json:"description"`          // Short description of the blog. Goes in the <meta> tags.
	Tags Tags                           `json:"tags"`                 // Tags for the blog. Goes in the <meta> tags.
	BlogURLPathPrefix string            `json:"blog_url_path_prefix"` // NOT IMPLEMENTED This prefix will be added to all in-site URLs that are generated.
	RenderPath string                   `json:"renderpath"`           // Path to where the rendered files should be put.
	TemplatePath string                 `json:"templatepath"`         // Path to template.
	UseFileTimestampAsCreationDate bool `json:"use_file_timestamp_as_creation_date"` // Use File Timestamp As Creation date.
	MetadataType MetadataType           `json:"metadata_type"`        // Type of the blog file metadata (TOML/YAML)
	Files []BlogMetadata                `json:"files"`                // List of blog markdown files.
}

