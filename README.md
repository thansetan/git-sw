<h1 align="center">git sw</h1>

a cli tool to switch between multiple git profiles/configs.

## Installation

Binary releases are available on the [releases page](/releases).

**Go**
```sh
go install github.com/thansetan/git-sw@latest
```

## Usage
```text
usage: git-sw [options] command
Available commands: 
  use       Select a config file to use.
  create    Create a new config file.
  edit      Edit an existing config file in text editor.
  delete    Delete an existing config file.
  list      List all available config files.

Available options: 
  -g        Run the command globally (can only be used with the 'use', 'edit', and 'delete' commands).
```