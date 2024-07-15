package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/manifoldco/promptui"
)

func hash(s string) (string, error) {
	h := md5.New()
	_, err := fmt.Fprint(h, s)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func formatError(err error) string {
	label := promptui.Styler(promptui.BGRed, promptui.FGBlack)("ERROR")
	errMsg := promptui.Styler(promptui.FGRed)(err)
	return fmt.Sprintf("%s %s", label, errMsg)
}

func errorAndExit(err error) {
	if errors.Is(err, promptui.ErrInterrupt) {
		return
	}
	fmt.Println(formatError(err))
	os.Exit(1)
}

func openTextEditor(filePath string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("notepad")
	default:
		cmd = exec.Command("vim")
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Args = append(cmd.Args, filePath)
	return cmd.Run()
}

func successMessage(profileName string, action Action) string {
	label := promptui.Styler(promptui.BGGreen, promptui.FGWhite)("SUCCESS")
	text := promptui.Styler(promptui.FGGreen)(fmt.Sprintf("%s profile \"%s\"", action, profileName))

	return fmt.Sprintf("%s %s", label, text)
}
