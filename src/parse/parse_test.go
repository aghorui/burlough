package parse

import (
	"log"
	"testing"

	"github.com/aghorui/burlough/static"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

var expectedFiles = []string{ "blog_list.html", "blog_page.html", "front_page.html", "search.html", "assets" }

// Prints an AST for reference
func TestPHONYPrintAST(t *testing.T) {
	k := goldmark.DefaultParser()
	root := k.Parse(text.NewReader(static.BlogTemplateNone))
	cursor := root

	cursor.Dump([]byte(static.BlogTemplateNone), 0)

	log.Printf("Next Sibling of Root: %v", cursor.NextSibling())
	log.Printf("Num Children of Root: %v", cursor.ChildCount())
}


// Checks the Blog Parsing Function for TOML
func TestParsingTOML(t *testing.T) {
	p, err := ParseBlogFile(static.BlogTemplateTOML)

	if err != nil {
		t.Fatalf("ParseBlogFile failed: %v", err);
		return
	}

	log.Printf("Function call succeeded. Returned: %v", p)

	if p.Title != "TITLE"      { t.Fail() }
	if p.Desc != "DESCRIPTION" { t.Fail() }
	if len(p.Tags) != 2        { t.Fail() }
	if p.Tags[0] != "TAG1"     { t.Fail() }
	if p.Tags[1] != "TAG2"     { t.Fail() }
}

// Checks the Blog Parsing Function for YAML
func TestParsingYAML(t *testing.T) {
	p, err := ParseBlogFile(static.BlogTemplateYAML)

	if err != nil {
		t.Fatalf("ParseBlogFile failed: %v", err);
		return
	}

	log.Printf("Function call succeeded. Returned: %v", p)

	if p.Title != "TITLE"      { t.Fail() }
	if p.Desc != "DESCRIPTION" { t.Fail() }
	if len(p.Tags) != 2        { t.Fail() }
	if p.Tags[0] != "TAG1"     { t.Fail() }
	if p.Tags[1] != "TAG2"     { t.Fail() }
}