package render

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aghorui/burlough/blog"
	"github.com/aghorui/burlough/parse"
	"github.com/aghorui/burlough/util"
)

type RenderPageInput struct {
	Title string
	Desc string
	Tags blog.Tags
	Entries []blog.BlogTemplateEntry
}

func renderIndexPage(t *blog.BlogTemplate, params blog.ConfigFileParams, entries []blog.BlogTemplateEntry) ([]byte, error) {
	var buf bytes.Buffer
	err := t.IndexPage.Execute(&buf, RenderPageInput{
		Title: params.Title,
		Desc: params.Desc,
		Tags: params.Tags,
		Entries: entries,
	})

	if err != nil {
		return buf.Bytes(), util.Error(err)
	}

	return buf.Bytes(), nil
}

func renderFrontPage(t *blog.BlogTemplate, params blog.ConfigFileParams, entries []blog.BlogTemplateEntry) ([]byte, error) {
	var buf bytes.Buffer
	err := t.FrontPage.Execute(&buf, RenderPageInput{
		Title: params.Title,
		Desc: params.Desc,
		Tags: params.Tags,
		Entries: entries,
	})

	if err != nil {
		return buf.Bytes(), util.Error(err)
	}

	return buf.Bytes(), nil
}

func renderBlogPage(t *blog.BlogTemplate, page blog.BlogTemplateEntry) ([]byte, error) {
	var buf bytes.Buffer
	err := t.BlogPage.Execute(&buf, page)

	if err != nil {
		return buf.Bytes(), util.Error(err)
	}

	return buf.Bytes(), nil
}

func Render(basePath string, tmpl *blog.BlogTemplate, params blog.ConfigFileParams, renderOverride string) error {
	entries := make([]blog.BlogTemplateEntry, 0, len(params.Files))

	var renderPath string

	if renderOverride != "" {
		renderPath = renderOverride
	} else {
		renderPath = params.RenderPath
	}

	err := tmpl.CopyAssetsToFolder(renderPath)
	if err != nil {
		return util.Error(err)
	}

	err = os.MkdirAll(renderPath, 0755)
	if err != nil {
		return util.Error(err)
	}

	// Prepare all articles
	for index, file := range params.Files {
		fmt.Fprintf(os.Stderr, "Processing %v (%v/%v)\n", file.Path, index + 1, len(params.Files))

		data, err := os.ReadFile(filepath.Join(basePath, file.Path))
		if err != nil {
			return util.Error(err)
		}

		finalPath := util.ExtractFilename(file.Path) + ".html"

		page, noMetadata, err := parse.ParseBlogFile(data)

		if err != nil {
			return fmt.Errorf("Error encountered while parsing %v: %w", file.Path, err)
		}

		if noMetadata {
			fmt.Fprintf(os.Stderr, "Warning: file %v has no metadata.\n", file.Path)
		}

		te := blog.PrepareBlogTemplateEntry(blog.BlogFile{
			BlogMetadata: file,
			BlogFileContents: page,
		}, finalPath, params.Desc, params.Tags)

		entries = append(entries, te)

		renderedPage, err := renderBlogPage(tmpl, te)

		if err != nil {
			return fmt.Errorf("Error encountered while rendering %v: %w", file.Path, err)
		}

		err = os.WriteFile(
			filepath.Join(renderPath, finalPath),
			renderedPage, 0644)

		if err != nil {
			return fmt.Errorf("Error encountered while writing %v: %w", filepath.Join(renderPath, finalPath), err)
		}
	}

	// Prepare blog index
	indexPage, err := renderIndexPage(tmpl, params, entries)
	if err != nil {
		return util.Error(err)
	}

	err = os.WriteFile(
		filepath.Join(renderPath, "blog_index" + ".html"),
		indexPage, 0644)
	if err != nil {
		return fmt.Errorf("Error encountered while blog index file: %w", err)
	}

	// Prepare front page
	frontPage, err := renderFrontPage(tmpl, params, entries)
	if err != nil {
		return util.Error(err)
	}

	err = os.WriteFile(
		filepath.Join(renderPath, "index" + ".html"),
		frontPage, 0644)
	if err != nil {
		return fmt.Errorf("Error encountered while site index file: %w", err)
	}

	return nil
}