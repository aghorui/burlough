package util

import (
	"embed"
	"io/fs"

	"github.com/otiai10/copy"
)

//go:embed testing_files/*
var testFileFSBase embed.FS

var TestFileFS fs.FS = func() fs.FS {
	f, err := fs.Sub(testFileFSBase, "testing_files")

	if err != nil {
		panic(err)
	}

	return f
}()

func GetTestFile(path string) []byte {
	data, err := fs.ReadFile(TestFileFS, path)

	if err != nil {
		panic(err)
	}

	return data
}

func WriteTestFiles(path string, dest string) {
	err := copy.Copy(path, dest, copy.Options{
		FS: TestFileFS,
		PermissionControl: copy.AddPermission(0644),
	})

	if err != nil {
		panic(err)
	}
}

func GenerateTestMarkdownFiles(dest string) {
	WriteTestFiles("markdown", dest)
}

func GenerateTestBadTemplate(dest string) {
	WriteTestFiles("template/bad_template", dest)
}

func GenerateTestEmptyTemplate(dest string) {
	WriteTestFiles("template/no_template", dest)
}

func GenerateTestEmptyConfig(dest string) {
	WriteTestFiles("config/empty.json", dest)
}

func GenerateTestBadConfig(dest string) {
	WriteTestFiles("template/config/bad.json", dest)
}