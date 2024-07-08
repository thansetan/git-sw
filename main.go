package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultConfigName = "default"
	saveDirName       = "git-sw"
)

var (
	userHomeDir, saveDirPath string
	profiles                 []Profile
)

func main() {
	var err error
	parseFlag()
	if len(os.Args) == 1 {
		flag.Usage()
		return
	}
	cmd := flag.Arg(0)
	action := getAction(strings.ToLower(cmd))
	if !action.IsValid() {
		fmt.Println(formatError(fmt.Errorf("invalid command = %s", cmd)))
		flag.Usage()
		os.Exit(1)
	}
	if _, ok := allowedGlobal[action]; isGlobal && !ok {
		errorAndExit(errors.New("flag -g can only be used with the 'use', 'edit', and 'delete' commands"))
	}

	userHomeDir, err = os.UserHomeDir()
	if err != nil {
		errorAndExit(err)
	}
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		errorAndExit(err)
	}
	saveDirPath = filepath.Join(userConfigDir, saveDirName)
	err = os.MkdirAll(saveDirPath, 0744)
	if err != nil {
		errorAndExit(err)
	}
	err = copyDefault()
	if err != nil {
		errorAndExit(err)
	}
	profiles, err = getProfiles(saveDirPath)
	if err != nil {
		errorAndExit(err)
	}

	command, ok := commands[action]
	if !ok {
		errorAndExit(ErrNotImplemented)
	}
	err = command.Func()
	if err != nil {
		errorAndExit(err)
	}
}
