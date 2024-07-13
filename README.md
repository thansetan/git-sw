<h1 align="center">git sw</h1>

<p align="center">
  <img src="https://github.com/thansetan/git-sw/assets/62317096/61820d18-ca8b-4e31-8f71-b5f473d289a0" alt="git-sw demo"/>
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
