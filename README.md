**Warning: Initial testing is still incomplete.**

<img src="./doc/logo.svg" alt="Burlough" />

<!---------------------------------
 ___          _               _
| _ )_  _ _ _| |___ _  _ __ _| |_
| _ \ || | '_| / _ \ || / _` | ' \
|___/\_,_|_| |_\___/\_,_\__, |_||_|
                        |___/
---------------------------------->

Burlough is yet another static blog generator that processes Markdown files
into HTML pages. Rendered blogs can be put as-is into static webpage hosting
services like Github Pages.

## Motivation

* I didn't want to learn another static blog generator system.
* I wanted to build a static blog generator based on my idea of how a blog
  generator should work (clarity of mind.)
* I didn't want to mess with editing a config file. I want the cli application
  itself to do that through commands (sort of like Git.)
* I wanted something that could also act like a note-taking application.

## Building

You need to have Go `1.20` (or above) on your system. GNU Make is recommended
but not required.

First, clone the repository and get into the directory.

```
git clone https://github.com/aghorui/burlough
cd burlough
```

Then run make.

```
make build-release
```

This will build the program and put it in a newly created directory called
`build.`

You may copy the executable created into $PATH to have it accessible from
everywhere.

If you do not have make, chdir to the `src` directory, and then run the
following command:

```
go build -o <path to built executable>
```

The Makefile is a tool for convenience rather than a necessity.

## Basic Usage

```
Usage: brlo <subcommand> [arguments]

The subcommands are:

	init      Initialize a new project in the current directory
	new       Add a new blog file to the project (untracked by project)
	config    Show or change a project configuration value
	scan      Scan and update the project file
	list      List all tracked files in project
	edit      Edit a given file
	render    Render the project into a finished blog

The following arguments are also supported:

Usage:
  -dump_template
    	Dump the default template to the current directory.
  -version
    	Print version information.
```

### Creating a New Blog

Burlough blogs are simply a folder of markdown files with a single metadata
file. To start a blog, open a terminal in an empty or nonempty subfolder and
type:

```
brlo init -title="This Is My Cool Blog" -desc="You better read it."
```

This will create a metadata file called `burlough.json` in the current folder.


### Creating Blog Files

All blog files are markdown files and end with the `*.md` extension. Burlough
can generate you a pre-filled markdown file for you in the project directory if
you use the `new` command:

```
brlo new -title="My New Blog Post"
```

Note that this does not actually add the file into the blog project.


### Blog Metadata

Metadata for blog files are a set of parameters for the document and are written
in either in YAML or TOML. You can currently have the following metadata options
in a blog file:

* `title`: The title for your document.
* `tags`: The tags for your document.
* `desc`: A short description for your document.

To add metadata to a blog file, you can add a frontmatter section as follows at
the top of the document. TOML and YAML have different delimiters for the
frontmatter blocks:

TOML:
```
+++
title = "title goes here"
tags = [ "tags", "go", "here" ]
desc = """
	description goes here
"""
+++

(Actual content of the document goes here)
```

YAML:
```
---
title: title goes here
tags: tags, go, here
desc: description goes here
---

(Actual content of the document goes here)
```

You can select the default type of frontmatter to generate using the new command
by setting the `metadata_type` option to either `toml` or `yaml`. See the
Configuration section below for details on configuring your project.


### Adding Files to the Blog

Burlough will scan for blog files in the project directory and add them to the
project on using the `scan` command:

```
brlo scan
```

To see the files that are currently being tracked by the project, use the
`list` command:

```
brlo list
```

### Rendering/Exporting the Blog

To render the blog into a set of HTML pages, use the `render` command:

```
brlo render
```

By default, this will generate the files in a separate directory called `output`
as a sibling to the project directory. This can be changed by setting the
`renderpath` parameter. See the Configuration section below for details on
configuring your project.


## Configuration

```
Usage: brlo config <subcommand> [arguments]

The subcommands are:

	get [key]             Gets a config key
	set -[key]=[value]    Sets a config key
	list                  Lists all keys and their values
```

You can configure the parameters of the project using the `config` command.
You can get or set them using the `get` or `set` commands respectively. To list
all of the available parameters and their current values, use the `list`
command.


## Editing and Using Custom Templates

You can change the look and feel of your blog by providing a custom template to
your blog. You can get the default template used by the program to modify by
running the following command:

```
brlo -dump_template
```

This will create a directory in the current location with all of the default
files for a template. The directory has the following file structure:

```
burlough_template/
│
├── assets                 -> The assets directory. On render, all the contents
│   │                         of this folder are copied over.
│   │
│   ├── template_icon.svg  -> The website icon. Defined in Template HTML.
│   │
│   └── template_main.css  -> The website CSS. Defined in Template HTML.
│
├── blog_list.html         -> The Index page. Lists all the blog files.
│
├── blog_page.html         -> The template page for an individual blog.
│
└── front_page.html        -> The front page of the blog.

```

The HTML template files use Go's `html/template` or `text/template` template
syntax.

You can specify this directory as the template by setting the `templatepath`
config parameter:

```
brlo config set -templatepath="path/to/template"
```


## Todo

* Tests are incomplete. Rewrite tests and write new tests.
* Return values of different subcommands may be incorrect.
* Render should implicitly run scan for new files?
* Complete Documentation

## Acknowledgements

This program makes direct use of the following Go projects:

```
Goldmark              https://pkg.go.dev/github.com/yuin/goldmark
Goldmark-Highlighting https://pkg.go.dev/github.com/yuin/goldmark-highlighting
Goldmark-Frontmatter  https://pkg.go.dev/go.abhg.dev/goldmark/frontmatter
Copy                  https://pkg.go.dev/github.com/otiai10/copy
```

## License

This software is available under the MIT license. Please see
[LICENSE](./LICENSE) for details.