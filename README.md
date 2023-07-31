**Warning: Initial testing is still incomplete.**

![Burlough](./doc/logo.svg)

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

You need to have Go `1.20` (or above) and GNU Make installed on your system.

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

## Basic Usage

```
Usage: ./burlough <subcommand> [arguments]

The commands are:

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
    	Print Version Information
```

## Editing Templates


## Todo

* Tests are incomplete. Rewrite tests and write new tests.
* Return values of different subcommands may be incorrect.

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