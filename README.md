<h1 align="center">git sw</h1>

<p align="center">
  <img src="https://github.com/user-attachments/assets/e08c8a1d-c00e-4aa8-b31f-338dea710bf7" alt="git-sw demo"/>
</p>

A CLI tool to switch between multiple Git profiles/configs. A package to parse and create a .gitconfig file is also available [here](https://github.com/thansetan/git-sw/tree/main/pkg/gitconfig).

## Installation

Binary releases are available on the [releases page](https://github.com/thansetan/git-sw/releases).

**Go**
```sh
go install github.com/thansetan/git-sw@latest
```

## Usage
```text
usage: git-sw [options] command
Available commands: 
  use       Select a profile to use.
  create    Create a new profile.
  edit      Edit an existing profile in text editor.
  delete    Delete an existing profile.
  list      List all available profiles.

Available options: 
  -g        Run the command globally (can only be used with the 'use', 'edit', and 'delete' commands).
```
