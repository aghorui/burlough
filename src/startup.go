package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/chzyer/readline"

	"github.com/aghorui/burlough/blog"
	"github.com/aghorui/burlough/blogtemplate"
	"github.com/aghorui/burlough/constants"
	"github.com/aghorui/burlough/project"
	"github.com/aghorui/burlough/util"
)

const (
	CommandInit       = "init"
	CommandNewFile    = "new"
	CommandConfig     = "config"
	CommandConfigSet  = "set"
	CommandConfigGet  = "get"
	CommandScan       = "scan"
	CommandList       = "list"
	CommandEdit       = "edit"
	CommandRender     = "render"
)

const usageString =
`Usage: %v <subcommand> [arguments]

The subcommands are:

	init      Initialize a new project in the current directory
	new       Add a new blog file to the project
	config    Show or change a project configuration value
	scan      Scan and update the project file
	list      List all tracked files in project
	edit      Edit a given file
	render    Render the project into a finished blog

The following arguments are also supported:

`

const configUsageString =
`Usage: %v config <subcommand> [arguments]

The subcommands are:

	get [key]             Gets a config key
	set -[key]=[value]    Sets a config key
	list                  Lists all keys and their values

`

var ErrInvalidArguments        = fmt.Errorf("Invalid Arguments.")
var ErrProjectAlreadyExists    = fmt.Errorf("Project file already exists in current folder.")
var ErrProjectDoesNotExist     = fmt.Errorf("Project file does not exist in current folder. Create a project using the 'init' subcommand.")
var ErrInvalidMetadataType     = fmt.Errorf("Invalid header metadata type. Type must be either 'toml' or 'yaml'.")
var ErrMalformedConfigFile     = fmt.Errorf("Config file values seem to be incorrect. Have you modified them?")
var ErrNoBlogFiles             = fmt.Errorf("There are no tracked blog files in the current directory. Please add the files in the directory using the 'scan' subcommand.")
var ErrNoConfigOptionSpecified = fmt.Errorf("No Config Option Specified.")

func LoadConfig(args []string) error {
	var showVersion bool
	var dumpTemplate bool

	defaultFlags := flag.NewFlagSet("", flag.ExitOnError)
	defaultFlags.BoolVar(&showVersion, "version", false, "Print version information.")
	defaultFlags.BoolVar(&dumpTemplate, "dump_template", false, "Dump the default template to the current directory.")

	if len(args) - 1 < 1 {
		fmt.Fprintf(os.Stderr, constants.AppName + " " + constants.AppVersion + " - Static Blog Generator\n\n")
		fmt.Fprintf(os.Stderr, usageString, args[0])
		defaultFlags.Usage()
		return ErrInvalidArguments
	}

	switch args[1] {
	case CommandInit:
		var c blog.ConfigFileParams
		var tags string
		var scan bool
		var wizard bool
		var metadataType string

		initFlags := flag.NewFlagSet("init", flag.ExitOnError)
		initFlags.StringVar(&c.Title, "title", "My Blog", "Name of your blog.")
		initFlags.StringVar(&c.Desc, "description", "", "Short description of your blog.")
		initFlags.StringVar(&tags, "tags", "", "Global Tags for blog.")
		initFlags.StringVar(&c.RenderPath, "renderpath", "../blogfiles", "Output directory for your blog.")
		initFlags.StringVar(&c.TemplatePath, "templatepath", "", "Template for your blog.")
		initFlags.StringVar(&metadataType, "metadata_type", "toml", "Default Header Metadata Type for your files (toml/yaml).")
		initFlags.BoolVar(&c.UseFileTimestampAsCreationDate,  "use_file_timestamp_as_creation_date", true, "Use the file modification time as the creation date.")
		initFlags.BoolVar(&scan, "scan", true, "Scan current directory for blog files immediately.")
		initFlags.BoolVar(&wizard, "wizard", false, "Enter init parameters using a wizard.")

		_ = initFlags.Parse(args[2:])

		// We use the wizard if no flags have been specified by the user or the
		// user specifically asks for the wizard to be run.
		wizard = wizard || !isFlagSet(initFlags)

		switch strings.ToLower(metadataType) {
		case "toml":
			c.MetadataType = blog.TOML
		case "yaml":
			c.MetadataType = blog.YAML
		default:
			return ErrInvalidMetadataType
		}

		c.Tags = util.SplitCommaList(tags)

		err := initProject(c, scan, wizard)
		if err != nil {
			return err
		}


	case CommandNewFile:
		var c blog.BlogFileContents
		var tags string
		var filenameOverride string
		var edit bool
		var wizard bool

		newFlags := flag.NewFlagSet("new", flag.ExitOnError)
		newFlags.StringVar(&c.Title, "title", "My Blog Post", "Name of your blog")
		newFlags.StringVar(&c.Desc, "description", "", "Short description of your blog post")
		newFlags.StringVar(&tags, "tags", "", "Tags for your blog post")
		newFlags.StringVar(&filenameOverride, "filename", "", "Set explicit filename")
		newFlags.BoolVar(&edit, "edit", true, "Edit the file after creation")
		newFlags.BoolVar(&wizard, "wizard", true, "Edit the file after creation")

		// We use the wizard if no flags have been specified by the user or the
		// user specifically asks for the wizard to be run.
		wizard = wizard || !isFlagSet(newFlags)


		c.Tags = util.SplitCommaList(tags)

		_ = newFlags.Parse(args[2:])

		err := newFile(c, filenameOverride, edit, wizard)
		if err != nil {
			return err
		}


	case CommandConfig:
		if len(args) - 1 < 2 {
			fmt.Fprintf(os.Stderr, configUsageString, args[0])
			return ErrInvalidArguments
		}

		switch args[2] {
		case "get":
			if len(args) - 1 < 3 {
				fmt.Fprintf(os.Stderr, configUsageString, args[0])
				return ErrInvalidArguments
			}

			path, err := os.Getwd()

			if err != nil {
				return err
			}

			state, err := project.Load(path)
			if err != nil {
				return err
			}

			switch args[3] {
			case "title":
				fmt.Printf("%v\n", state.Title)

			case "description":
				fmt.Printf("%v\n", state.Desc)

			case "tags":
				fmt.Printf("%v\n", state.Tags)

			case "renderpath":
				fmt.Printf("%v\n", state.RenderPath)

			case "templatepath":
				fmt.Printf("%v\n", state.TemplatePath)

			case "metadata_type":
				if state.MetadataType == blog.TOML {
					fmt.Printf("toml\n");
				} else if state.MetadataType == blog.YAML {
					fmt.Printf("yaml\n");
				} else {
					panic(fmt.Errorf("Invalid Metadata Type found: %v.", state.MetadataType));
				}

			case "use_file_timestamp_as_creation_date":
				fmt.Printf("%v\n", state.UseFileTimestampAsCreationDate)
			}


		case "set":
			path, err := os.Getwd()

			if err != nil {
				return err
			}

			state, err := project.Load(path)
			if err != nil {
				return err
			}

			var tags string
			var metadataType string

			switch state.MetadataType {
			case blog.TOML:
				metadataType = "toml"
			case blog.YAML:
				metadataType = "yaml"
			default:
				return ErrMalformedConfigFile
			}

			cfgFlags := flag.NewFlagSet("config set", flag.ExitOnError)
			cfgFlags.StringVar(&state.Title, "title", state.Title, "Name of your blog.")
			cfgFlags.StringVar(&state.Desc, "description", state.Desc, "Short description of your blog.")
			cfgFlags.StringVar(&tags, "tags", "CHANGEME", "Global Tags for blog.")
			cfgFlags.StringVar(&state.RenderPath, "renderpath", state.RenderPath, "Output directory for your blog.")
			cfgFlags.StringVar(&state.TemplatePath, "templatepath", state.TemplatePath, "Template for your blog.")
			cfgFlags.StringVar(&metadataType, "metadata_type", metadataType, "Default Header Metadata Type for your files (toml/yaml).")
			cfgFlags.BoolVar(&state.UseFileTimestampAsCreationDate,  "use_file_timestamp_as_creation_date", state.UseFileTimestampAsCreationDate, "Use the file modification time as the creation date.")

			_ = cfgFlags.Parse(args[3:])

			if !isFlagSet(cfgFlags) {
				cfgFlags.Usage()
				return ErrNoConfigOptionSpecified
			}

			switch strings.ToLower(metadataType) {
			case "toml":
				state.MetadataType = blog.TOML
			case "yaml":
				state.MetadataType = blog.YAML
			default:
				return ErrInvalidMetadataType
			}

			state.Tags = util.SplitCommaList(tags)

			err = state.WriteConfig()
			if err != nil {
				return err
			}

		case "list":
			path, err := os.Getwd()

			if err != nil {
				return err
			}

			state, err := project.Load(path)
			if err != nil {
				return err
			}

			fmt.Printf("title='%v'\n", state.Title)
			fmt.Printf("description='%v'\n", state.Desc)
			fmt.Printf("tags='%v'\n", state.Tags)
			fmt.Printf("renderpath='%v'\n", state.RenderPath)
			fmt.Printf("templatepath='%v'\n", state.TemplatePath)
			fmt.Printf("metadata_type=")

			if state.MetadataType == blog.TOML {
				fmt.Printf("'toml'\n");
			} else if state.MetadataType == blog.YAML {
				fmt.Printf("'yaml'\n");
			} else {
				panic(fmt.Errorf("Invalid Metadata Type found: %v.", state.MetadataType));
			}

			fmt.Printf("use_file_timestamp_as_creation_date='%v'\n", state.UseFileTimestampAsCreationDate)


		default:
			fmt.Fprintf(os.Stderr, configUsageString, args[0])
			return ErrInvalidArguments
		}

	case CommandScan:
		err := scanProject()
		if err != nil {
			return err
		}

	case CommandList:
		err := listFiles()
		if err != nil {
			return err
		}

	case CommandEdit:
		if len(args) - 1 < 2 {
			fmt.Fprintf(os.Stderr, `Usage: %v edit [filename]\n`, args[0])
			return nil
		}

		err := editFile(args[2])
		if err != nil {
			return err
		}

	case CommandRender:
		var renderOverride string
		renderFlags := flag.NewFlagSet("render", flag.ExitOnError)
		renderFlags.StringVar(&renderOverride, "path", "", "Output directory for your blog. (override)")

		_ = renderFlags.Parse(args[2:])

		err := renderProject(renderOverride)
		if err != nil {
			return err
		}

	default:
		_ = defaultFlags.Parse(args[1:])

		if showVersion {
			fmt.Printf(constants.AppName + " " + constants.AppVersion + "\n")
			return nil

		} else if dumpTemplate {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}

			err = blogtemplate.DumpDefaultExportTemplate(dir)
			if err != nil {
				return err
			}

		} else {
			fmt.Fprintf(os.Stderr, constants.AppName + " " + constants.AppVersion + " - Static Blog Generator\n\n")
			fmt.Fprintf(os.Stderr, usageString, args[0])
			defaultFlags.Usage()
			return ErrInvalidArguments
		}
	}

	return nil
}

func isFlagSet(fset *flag.FlagSet) bool {
	flagIsSet := false
	fset.Visit(func(f *flag.Flag) {
		flagIsSet = true
	})

	return flagIsSet
}

func projectFileExists() bool {
	_, err := os.Stat(project.ProjectConfigFileName)

	if !os.IsNotExist(err) {
		return true
	} else {
		return false
	}
}

func getMetadataType(metadataType string) (blog.MetadataType, error) {
	switch strings.ToLower(metadataType) {
	case "toml":
		return blog.TOML, nil
	case "yaml":
		return blog.YAML, nil
	default:
		return blog.Invalid, ErrInvalidMetadataType
	}
}

func wizardPromptString(rl *readline.Instance, text string, defaultVal string, target *string) error {
	rl.SetPrompt(fmt.Sprintf("%v (default: %v): ", text, defaultVal))
	line, err := rl.Readline()

	if err != nil && err != io.EOF {
		return err
	}

	if line == "" {
		*target = defaultVal
	} else {
		*target = line
	}

	return nil
}

func wizardPromptBool(rl *readline.Instance, text string, defaultVal bool, target *bool) error {
	for {
		defaultValStr := ""

		if defaultVal == true {
			defaultValStr = "y"
		} else {
			defaultValStr = "n"
		}

		rl.SetPrompt(fmt.Sprintf("%v (y/n) (default: %v): ", text, defaultValStr))
		line, err := rl.Readline()

		if err != nil && err != io.EOF {
			return err
		}

		if line == "y" || line == "Y" {
			*target = true
			break
		} else if line == "n" || line == "N" {
			*target = false
			break
		} else if line == "" {
			*target = defaultVal
			break
		} else {
			fmt.Printf("Please enter 'y' or 'n' only.\n")
		}
	}

	return nil
}

func initWizard(defaults blog.ConfigFileParams) (blog.ConfigFileParams, error) {
	var params blog.ConfigFileParams = defaults
	var tagStr = ""
	var metadataTypeStr = ""

	path, err := os.Getwd()
	if err != nil {
		return params, err
	}

	rl, err := readline.New("")
	if err != nil {
		return params, err
	}
	defer rl.Close()


	fmt.Printf("Creating a new %v blog in '%v'...\n", constants.AppName, path)

	if err := wizardPromptString(rl, "Title of blog", params.Title, &params.Title); err != nil {
		return params, err
	}

	if err := wizardPromptString(rl, "Description of blog", params.Desc, &params.Desc); err != nil {
		return params, err
	}

	if err := wizardPromptString(rl, "Tags for blog (comma separated list)", params.Tags.String(), &tagStr); err != nil {
		return params, err
	} else {
		params.Tags = util.SplitCommaList(tagStr)
	}

	if err := wizardPromptString(rl, "Output path for blog", params.RenderPath, &params.RenderPath); err != nil {
		return params, err
	}

	if err := wizardPromptString(rl, "Path to blog template", params.TemplatePath, &params.TemplatePath); err != nil {
		return params, err
	}

	for {
		if err := wizardPromptString(rl, "Default blog metadata type (toml/yaml)", "toml", &metadataTypeStr); err != nil {
			return params, err
		}

		params.MetadataType, err = getMetadataType(metadataTypeStr)

		if err != nil {
			fmt.Printf("Please enter 'toml' or 'yaml' only.\n")
		} else {
			break
		}
	}

	return params, nil
}

func initProject(c blog.ConfigFileParams, scan bool, wizard bool) error {
	if projectFileExists() {
		return ErrProjectAlreadyExists
	}

	path, err := os.Getwd()

	if err != nil {
		return err
	}

	if wizard {
		c, err = initWizard(c)
		if err != nil {
			return err
		}
	}

	state, ul, err := project.Init(path, c, scan)
	if err != nil {
		return err
	}

	printUpdateLog(ul)

	err = state.WriteConfig()

	if err != nil {
		return err
	}

	fmt.Printf("Project Initialized at %v\n", path)
	return nil
}

func newFileWizard(defaults blog.BlogFileContents, edit bool) (blog.BlogFileContents, bool, error) {
	var params blog.BlogFileContents = defaults
	var tagStr = ""

	rl, err := readline.New("")
	if err != nil {
		return params, edit, err
	}
	defer rl.Close()


	fmt.Printf("Creating a new blog entry...\n")

	if err := wizardPromptString(rl, "Title of blog entry", params.Title, &params.Title); err != nil {
		return params, edit, err
	}

	if err := wizardPromptString(rl, "Description of blog entry", params.Desc, &params.Desc); err != nil {
		return params, edit, err
	}

	if err := wizardPromptString(rl, "Tags for blog (comma separated list)", params.Tags.String(), &tagStr); err != nil {
		return params, edit, err
	} else {
		params.Tags = util.SplitCommaList(tagStr)
	}

	if err := wizardPromptBool(rl, "Edit file now?", edit, &edit); err != nil {
		return params, edit, err
	} else {
		params.Tags = util.SplitCommaList(tagStr)
	}

	return params, edit, nil
}


func newFile(b blog.BlogFileContents, filenameOverride string, edit bool, wizard bool) error {
	if !projectFileExists() {
		return ErrProjectDoesNotExist
	}

	path, err := os.Getwd()

	if err != nil {
		return err
	}

	state, err := project.Load(path)
	if err != nil {
		return err
	}

	if wizard {
		b, edit, err = newFileWizard(b, edit)
		if err != nil {
		return err
		}
	}

	path, err = state.NewFile(b, filenameOverride)
	if err != nil {
		return err
	}

	fmt.Printf("Created file: %v\n", path)
	fmt.Printf("Scanning for changes...\n")

	ul, err := state.Scan()
	if err != nil {
		return err
	}

	printUpdateLog(ul)

	err = state.WriteConfig()
	if err != nil {
		return err
	}

	if edit {
		editFile(path)
	}

	return nil
}

func printUpdateLog(u []project.UpdateLog) {
	fmt.Printf("Scanned %v files.\n", len(u))
	for _, l := range u {
		switch l.UpdateMode {
		case project.Created:
			fmt.Printf("Created: ")
		case project.Updated:
			fmt.Printf("Updated: ")
		case project.Deleted:
			fmt.Printf("Deleted: ")
		default:
			panic(util.Error(fmt.Errorf("BUG: Found invalid Update Mode: %v", l.UpdateMode)))
		}

		fmt.Printf("%v\n", l.Path)
	}
}

func scanProject() error {
	if !projectFileExists() {
		return ErrProjectDoesNotExist
	}

	path, err := os.Getwd()

	if err != nil {
		return err
	}

	state, err := project.Load(path)
	if err != nil {
		return err
	}

	ul, err := state.Scan()
	if err != nil {
		return err
	}

	printUpdateLog(ul)

	err = state.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

func listFiles() error {
	if !projectFileExists() {
		return ErrProjectDoesNotExist
	}

	path, err := os.Getwd()

	if err != nil {
		return err
	}

	state, err := project.Load(path)
	if err != nil {
		return err
	}

	for _, r := range state.Files {
		fmt.Printf("%v\n", r.Path)
	}

	return nil
}

func editFile(filename string) error {
	if !projectFileExists() {
		return ErrProjectDoesNotExist
	}

	path, err := os.Getwd()

	if err != nil {
		return err
	}

	state, err := project.Load(path)
	if err != nil {
		return err
	}

	fmt.Printf("Editing file: %v\n", filename)

	editor, ok := os.LookupEnv(constants.AppEnvironmentVarPrefix + "EDITOR")

	if !ok {
		editor, ok = os.LookupEnv("EDITOR")
	}

	if !ok {
		return fmt.Errorf("Error: Neither $EDITOR nor $" + constants.AppEnvironmentVarPrefix + "EDITOR" + " were set.")
	}

	var finalFilename string

	if filename == "latest" {
		if len(state.Files) > 0 {
			finalFilename = state.Files[0].Path
		} else {
			return ErrNoBlogFiles
		}
	} else {
		finalFilename = filename
	}

	cmd := exec.Command(editor, finalFilename)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func renderProject(renderOverride string) error {
	if !projectFileExists() {
		return ErrProjectDoesNotExist
	}

	path, err := os.Getwd()

	if err != nil {
		return err
	}

	state, err := project.Load(path)
	if err != nil {
		return err
	}

	fmt.Printf("Scanning for changes in %v\n", path)

	ul, err := state.Scan()
	if err != nil {
		return err
	}

	printUpdateLog(ul)

	err = state.WriteConfig()
	if err != nil {
		return err
	}

	fmt.Printf("Rendering %v\n", path)

	err = state.Render(renderOverride)
	if err != nil {
		return err
	}

	var renderPath string

	if renderOverride != "" {
		renderPath = renderOverride
	} else {
		renderPath = state.RenderPath
	}

	fmt.Fprintf(os.Stderr, "Rendered to %v.\n", renderPath)
	return nil
}


