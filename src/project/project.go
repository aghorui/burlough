package project

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/aghorui/burlough/blog"
	"github.com/aghorui/burlough/render"
	"github.com/aghorui/burlough/static"
	"github.com/aghorui/burlough/util"
)

const ProjectConfigFileName = "burlough.json"
const DefaultBlogFileExtension = ".md"
const DefaultWhitespaceReplacement = "-"
const hashingBufferSize = 2048

var ErrNoConfigFileFound = fmt.Errorf("Config file not found in base path.")

// State of a given blog project.
type ProjectState struct {
	BasePath string            // Path to the base folder of the project.
	                           // The program will chdir to this path.
	Template blog.BlogTemplate // Template Struct.
	blog.ConfigFileParams      // Include config file params into struct
}


// Get the hash of the given file.
func getHash(file *os.File) (blog.FileHash, error) {
	hasher := sha1.New()
	buffer := make([]byte, hashingBufferSize)
	for {
		n, err := file.Read(buffer)

		if err != nil {
			if err != io.EOF {
				return blog.FileHash(""), err
			} else {
				break;
			}
		}

		hasher.Write(buffer[0:n])
	}

	return blog.FileHash(hex.EncodeToString(hasher.Sum(nil))), nil
}

// Type used for making a dictionary between blog filename -> Blog Metadata.
// We still want the slice to be preserved.
type MetadataMap map[string]int

// Scans all blog files (*.md) within a folder and returns metadata for them.
func scanBlogFiles(files []fs.DirEntry, useFileTimestampAsCreationDate bool) ([]blog.BlogMetadata, MetadataMap, error) {
	projectFiles := make([]blog.BlogMetadata, 0, 10)
	metaMap := make(MetadataMap)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != DefaultBlogFileExtension {
			continue
		}

		fh, err := os.Open(file.Name())
		if err != nil {
			return nil, nil, err
		}

		filehash, err := getHash(fh)
		if err != nil {
			return nil, nil, err
		}

		var createdTime time.Time

		if useFileTimestampAsCreationDate {
			info, err := file.Info()

			if err != nil {
				return nil, nil, err
			}

			createdTime = info.ModTime().UTC()
		} else {
			createdTime = time.Now().UTC()
		}

		projectFiles = append(projectFiles, blog.BlogMetadata{
			Path: file.Name(),
			Hash: filehash,
			Updated: time.Time{},
			Created: createdTime,
		})

		metaMap[file.Name()] = len(projectFiles) - 1;
	}

	// We're putting the responsibility of having the metadata sorted in this
	// function
	sort.Sort(blog.BlogMetadataSortCreatedDescending(projectFiles))
	return projectFiles, metaMap, nil
}

type UpdateMode int

// Enum for UpdateLog
const (
	Created UpdateMode = 0
	Updated UpdateMode = 1
	Deleted UpdateMode = 2
)

// Store Update/Create/Delete information between old and new
type UpdateLog struct {
	UpdateMode UpdateMode
	Path string
}

// Looks at the old and new metadata values and updates them.
// Kind of looks like reinventing git.
func updateBlogMetadata(old []blog.BlogMetadata, new []blog.BlogMetadata, newMetaMap MetadataMap) []UpdateLog {
	// 3 cases
	// old exists -> new doesn't exist : delete
	// old doesn't exist -> new exists : add
	// old exists -> new exists        : update

	updateLog := make([]UpdateLog, 0, len(new))

	for _, blogMetadata := range old {
		v, ok := newMetaMap[blogMetadata.Path]

		// Deletion Case
		if !ok {
			updateLog = append(updateLog, UpdateLog{ Deleted, blogMetadata.Path })

		// Updation Case
		} else {
			updateLog = append(updateLog, UpdateLog{ Updated, blogMetadata.Path })

			// This is all we need this function for really. Other than this it's all statistics.

			// We consider the update stamp to be the created stamp in this case.
			if new[v].Hash != blogMetadata.Hash {
				new[v].Updated = new[v].Created
			}

			// And we carry over the old created timestamp.
			new[v].Created = blogMetadata.Created

			// We delete whatever we have processed to single out the new entries.
			delete(newMetaMap, blogMetadata.Path)
		}


	}

	// Creation Case
	for k := range newMetaMap {
		updateLog = append(updateLog, UpdateLog{ Created, k })
	}

	return updateLog
}

// Finalize blog metadata to keep in project state.
func finalizeBlogMetadata(data []blog.BlogMetadata) {
	sort.Sort(blog.BlogMetadataSortCreatedDescending(data))
}

// Combining Function.
func prepareBlogMetadata(old []blog.BlogMetadata, files []fs.DirEntry, useFileTimestampAsCreationDate bool) ([]blog.BlogMetadata, []UpdateLog, error) {
	projectFiles, metaMap, err := scanBlogFiles(files, useFileTimestampAsCreationDate)

	if err != nil {
		return nil, nil, err
	}

	updateLog := updateBlogMetadata(old, projectFiles, metaMap)

	finalizeBlogMetadata(projectFiles)

	return projectFiles, updateLog, nil
}

// Initializes a project with a json file at the root directory
func Init(
	basePath string,
	params blog.ConfigFileParams,
	scan bool) (ProjectState, []UpdateLog, error) {

	wd, err := os.Getwd()

	if err != nil {
		util.LogErr(err)
		panic(err)
	}

	// Switch back to wd after we are done
	defer func(wd string) {
		if err := os.Chdir(wd); err != nil {
			util.LogErr(err)
			panic(err)
		}
	}(wd)

	// Chdir to project base
	if err := os.Chdir(basePath); err != nil {
		return ProjectState{}, nil, err
	}

	base, err := os.ReadDir(basePath)

	if err != nil {
		return ProjectState{}, nil, err
	}

	var updateLog []UpdateLog = nil

	if scan {
		var projectFiles []blog.BlogMetadata
		projectFiles, updateLog, err = prepareBlogMetadata(params.Files, base, params.UseFileTimestampAsCreationDate)

		if err != nil {
			return ProjectState{}, nil, err
		}

		// Replace old files with current.
		params.Files = projectFiles
	}

	if err != nil {
		return ProjectState{}, updateLog, err
	}

	var tmpl blog.BlogTemplate

	if params.TemplatePath == "" {
		tmpl = static.DefaultBlogTemplate
	} else {
		f := os.DirFS(params.TemplatePath)

		tmpl, err = blog.LoadTemplate(f)

		if err != nil {
			return ProjectState{}, nil, err
		}
	}


	return ProjectState{
		BasePath: basePath,
		Template: tmpl,
		ConfigFileParams: params,
	}, updateLog, nil
}

// Loads an existing project and returns a projectparams struct for it.
func Load(basePath string) (ProjectState, error) {
	wd, err := os.Getwd()

	if err != nil {
		util.LogErr(err)
		panic(err)
	}

	// Switch back to wd after we are done
	defer func(wd string) {
		if err := os.Chdir(wd); err != nil {
			util.LogErr(err)
			panic(err)
		}
	}(wd)

	// Chdir to project base
	if err := os.Chdir(basePath); err != nil {
		return ProjectState{}, err
	}

	// Read and unmarshal
	data, err := os.ReadFile(filepath.Join(basePath, ProjectConfigFileName))

	if err != nil {
		if os.IsNotExist(err) {
			return ProjectState{}, ErrNoConfigFileFound
		} else {
			return ProjectState{}, err
		}
	}

	var params blog.ConfigFileParams
	err = json.Unmarshal(data, &params);

	if err != nil {
		return ProjectState{}, err
	}


	var tmpl blog.BlogTemplate

	if params.TemplatePath == "" {
		tmpl = static.DefaultBlogTemplate
	} else {
		f := os.DirFS(params.TemplatePath)

		tmpl, err = blog.LoadTemplate(f)

		if err != nil {
			return ProjectState{}, err
		}
	}

	return ProjectState{
		BasePath: basePath,
		Template: tmpl,
		ConfigFileParams: params,
	}, err
}

func (state *ProjectState) Scan() ([]UpdateLog, error) {
	wd, err := os.Getwd()

	if err != nil {
		util.LogErr(err)
		panic(err)
	}

	// Switch back to wd after we are done
	defer func(wd string) {
		if err := os.Chdir(wd); err != nil {
			util.LogErr(err)
			panic(err)
		}
	}(wd)

	// Chdir to project base
	if err := os.Chdir(state.BasePath); err != nil {
		return nil, err
	}

	base, err := os.ReadDir(state.BasePath)

	if err != nil {
		return nil, err
	}

	projectFiles, updateLog, err := prepareBlogMetadata(state.Files, base, state.UseFileTimestampAsCreationDate)

	// Replace old files with current.
	state.Files = projectFiles

	if err != nil {
		return nil, err
	}

	return updateLog, err
}

func (state ProjectState) WriteConfig() error {
	wd, err := os.Getwd()
	if err != nil {
		util.LogErr(err)
		panic(err)
	}

	// Switch back to wd after we are done
	defer func(wd string) {
		if err := os.Chdir(wd); err != nil {
			util.LogErr(err)
			panic(err)
		}
	}(wd)

	// Chdir to project base
	if err := os.Chdir(state.BasePath); err != nil {
		return err
	}

	data, err := json.MarshalIndent(state.ConfigFileParams, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(ProjectConfigFileName, data, 0644);
	if err != nil {
		return err
	}

	return nil
}

func sanitizeString(s string) string {
	sl := strings.Join(strings.Fields(s), DefaultWhitespaceReplacement)
	sanitized := strings.ToLower(string(util.SanitizeRegex.ReplaceAll([]byte(sl), []byte(""))))
	timePrefix := time.Now().Format("20060102")
	return timePrefix + DefaultWhitespaceReplacement + sanitized
}

func (state ProjectState) NewFile(b blog.BlogFileContents, filenameOverride string) error {
	wd, err := os.Getwd()
	if err != nil {
		util.LogErr(err)
		panic(err)
	}

	// Switch back to wd after we are done
	defer func(wd string) {
		if err := os.Chdir(wd); err != nil {
			util.LogErr(err)
			panic(err)
		}
	}(wd)

	// Chdir to project base
	if err := os.Chdir(state.BasePath); err != nil {
		return err
	}

	var filename string

	if filenameOverride == "" {
		filename = b.Title
	} else {
		filename = filenameOverride
	}

	filePath := filepath.Join(state.BasePath, sanitizeString(filename) + "." + DefaultBlogFileExtension)

	file, err := os.OpenFile(filePath, os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	err = static.WriteBlogFile(state.MetadataType, file, b)
	if err != nil {
		return err
	}

	file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (state ProjectState) Render(renderOverride string) error {
	wd, err := os.Getwd()
	if err != nil {
		util.LogErr(err)
		panic(err)
	}

	// Switch back to wd after we are done
	defer func(wd string) {
		if err := os.Chdir(wd); err != nil {
			util.LogErr(err)
			panic(err)
		}
	}(wd)

	// Chdir to project base
	if err := os.Chdir(state.BasePath); err != nil {
		return err
	}

	err = render.Render(state.BasePath, &state.Template, state.ConfigFileParams, renderOverride)

	if err != nil {
		return err
	}

	return nil
}