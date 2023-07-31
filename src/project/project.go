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


func ErrFileAlreadyExists(filePath string) error {
	return fmt.Errorf("File already exists: %v\n", filePath);
}

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
				return blog.FileHash(""), util.Error(err)
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
			return nil, nil, util.Error(err)
		}

		filehash, err := getHash(fh)
		if err != nil {
			return nil, nil, util.Error(err)
		}

		var createdTime time.Time

		if useFileTimestampAsCreationDate {
			info, err := file.Info()

			if err != nil {
				return nil, nil, util.Error(err)
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

		metaMap[file.Name()] = len(projectFiles) - 1
	}

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
			// This is all we need this function for really. Other than this it's all statistics.

			// We consider the update stamp to be the created stamp in this case.
			if new[v].Hash != blogMetadata.Hash {
				updateLog = append(updateLog, UpdateLog{ Updated, blogMetadata.Path })
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
		return nil, nil, util.Error(err)
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
		return ProjectState{}, nil, util.Error(err)
	}

	base, err := os.ReadDir(basePath)

	if err != nil {
		return ProjectState{}, nil, util.Error(err)
	}

	var updateLog []UpdateLog = nil

	if scan {
		var projectFiles []blog.BlogMetadata
		projectFiles, updateLog, err = prepareBlogMetadata(params.Files, base, params.UseFileTimestampAsCreationDate)

		if err != nil {
			return ProjectState{}, nil, util.Error(err)
		}

		// Replace old files with current.
		params.Files = projectFiles
	}

	if err != nil {
		return ProjectState{}, updateLog, util.Error(err)
	}

	var tmpl blog.BlogTemplate

	if params.TemplatePath == "" {
		tmpl = static.DefaultBlogTemplate
	} else {
		f := os.DirFS(params.TemplatePath)

		tmpl, err = blog.LoadTemplate(f)

		if err != nil {
			return ProjectState{}, nil, util.Error(err)
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
		return ProjectState{}, util.Error(err)
	}

	// Read and unmarshal
	data, err := os.ReadFile(filepath.Join(basePath, ProjectConfigFileName))

	if err != nil {
		if os.IsNotExist(err) {
			return ProjectState{}, ErrNoConfigFileFound
		} else {
			return ProjectState{}, util.Error(err)
		}
	}

	var params blog.ConfigFileParams
	err = json.Unmarshal(data, &params);

	if err != nil {
		return ProjectState{}, util.Error(err)
	}


	var tmpl blog.BlogTemplate

	if params.TemplatePath == "" {
		tmpl = static.DefaultBlogTemplate
	} else {
		f := os.DirFS(params.TemplatePath)

		tmpl, err = blog.LoadTemplate(f)

		if err != nil {
			return ProjectState{}, util.Error(err)
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
		return nil, util.Error(err)
	}

	base, err := os.ReadDir(state.BasePath)

	if err != nil {
		return nil, util.Error(err)
	}

	projectFiles, updateLog, err := prepareBlogMetadata(state.Files, base, state.UseFileTimestampAsCreationDate)

	// Replace old files with current.
	state.Files = projectFiles

	if err != nil {
		return nil, util.Error(err)
	}

	return updateLog, util.Error(err)
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
		return util.Error(err)
	}

	data, err := json.MarshalIndent(state.ConfigFileParams, "", "\t")
	if err != nil {
		return util.Error(err)
	}

	err = os.WriteFile(ProjectConfigFileName, data, 0644);
	if err != nil {
		return util.Error(err)
	}

	return nil
}

func sanitizeString(s string) string {
	sl := strings.Join(strings.Fields(s), DefaultWhitespaceReplacement)
	sanitized := strings.ToLower(string(util.SanitizeRegex.ReplaceAll([]byte(sl), []byte(""))))
	timePrefix := time.Now().Format("20060102")
	return timePrefix + DefaultWhitespaceReplacement + sanitized
}

func (state ProjectState) NewFile(b blog.BlogFileContents, filenameOverride string) (string, error) {
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
		return "", util.Error(err)
	}

	var filename string

	if filenameOverride == "" {
		filename = b.Title
	} else {
		filename = filenameOverride
	}

	filePath := filepath.Join(state.BasePath, sanitizeString(filename) + DefaultBlogFileExtension)

	_, err = os.Stat(filePath)

	if !os.IsNotExist(err) {
		return filePath, ErrFileAlreadyExists(filePath)
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil {
		return filePath, util.Error(err)
	}

	err = static.WriteBlogFile(state.MetadataType, file, b)
	if err != nil {
		return filePath, util.Error(err)
	}

	file.Close()
	if err != nil {
		return filePath, util.Error(err)
	}

	return filePath, nil
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
		return util.Error(err)
	}

	err = render.Render(state.BasePath, &state.Template, state.ConfigFileParams, renderOverride)

	if err != nil {
		return err
	}

	return nil
}