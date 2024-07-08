package main

import (
	"fmt"
	"os/exec"
)

type GPGFormat string

const (
	OPENPGP GPGFormat = "openpgp"
	SSH     GPGFormat = "ssh"
	X509    GPGFormat = "x509"
)

var gpgFormat = []GPGFormat{OPENPGP, SSH, X509}

func isGitDirectory() bool {
	cmd := exec.Command("git", "rev-parse")
	err := cmd.Run()
	if err != nil && cmd.ProcessState.ExitCode() != 128 {
		panic(err)
	}
	return cmd.ProcessState.ExitCode() == 0
}

func unsetConfig(pattern string) error {
	var cmd *exec.Cmd
	if isGitDirectory() {
		cmd = exec.Command("git", "config", "--unset", "include.path", pattern)
	} else {
		cmd = exec.Command("git", "config", "--global", "--unset", "include.path", pattern)
	}
	gitOutput, err := cmd.CombinedOutput()
	if err != nil && cmd.ProcessState.ExitCode() != 5 { // try to unset an option that does not exist will give exit 5
		fmt.Printf("git: %s", string(gitOutput))
		return err
	}
	return nil
}

func applyConfig(configPath string, isGlobal bool) error {
	var cmd *exec.Cmd
	if isGlobal {
		cmd = exec.Command("git", "config", "--global", "--replace-all", "include.path", configPath, fmt.Sprintf("%s.*gitconfig$", saveDirName))
	} else {
		cmd = exec.Command("git", "config", "--replace-all", "include.path", configPath, fmt.Sprintf("%s.*gitconfig$", saveDirName))
	}
	gitOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("git: %s", string(gitOutput))
		return err
	}
	return nil
}

func getCurrentConfig() (string, error) {
	var cmd *exec.Cmd
	if !isGlobal && isGitDirectory() {
		cmd = exec.Command("git", "config", "--get", "include.path", fmt.Sprintf("%s.*gitconfig$", saveDirName))
	} else {
		cmd = exec.Command("git", "config", "--global", "--get", "include.path", fmt.Sprintf("%s.*gitconfig$", saveDirName))
	}
	gitOutput, err := cmd.CombinedOutput()
	if err != nil && cmd.ProcessState.ExitCode() != 1 {
		fmt.Println(string(gitOutput))
		fmt.Println(cmd.String())
		return "", err
	}
	return string(gitOutput), nil
}
