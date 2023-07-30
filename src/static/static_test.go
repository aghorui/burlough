package static

import (
	"io/fs"
	"log"
	"testing"
)

var expectedFiles = []string{ "blog_list.html", "blog_page.html", "front_page.html", "search.html", "assets" }

// Checks whether all template files are present as they are supposed to be or not
func TestEmbeddedTemplateIntegrity(t *testing.T) {
	files := GetDefaultExportTemplateFiles()

	fileMap := make(map[string]fs.DirEntry)

	for _, f := range files {
		fileMap[f.Name()] = f
		log.Printf("%v", f.Name())
	}

	for _, f := range expectedFiles {
		if _, ok := fileMap[f]; !ok {
			t.Fatalf("File %v not found in embedded directory.", f);
		}
	}

	if !fileMap["assets"].Type().IsDir() {
		t.Fatalf("'assets' is not a directory.");
	}
}

// Checks whether all template files are present as they are supposed to be or not
func TestEmbeddedTemplateFSIntegrity(t *testing.T) {
	files_fs := GetDefaultExportTemplateFS()
	files, err :=  fs.ReadDir(files_fs, ".")

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	// blog_page, err := fs.ReadFile(files_fs, "blog_page.html")

	// if err != nil {
	// 	t.Fatalf("Error: %v", err)
	// }

	// log.Printf("blog_page: %v", string(blog_page))

	fileMap := make(map[string]fs.DirEntry)

	for _, f := range files {
		fileMap[f.Name()] = f
		log.Printf("%v", f.Name())
	}

	for _, f := range expectedFiles {
		if _, ok := fileMap[f]; !ok {
			t.Fatalf("File %v not found in embedded directory.", f);
		}
	}

	if !fileMap["assets"].Type().IsDir() {
		t.Fatalf("'assets' is not a directory.");
	}
}