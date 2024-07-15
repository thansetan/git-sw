package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Command struct {
	Func        func() error
	Description string
}

var commands = map[Action]Command{
	CREATE: {
		Description: "Create a new profile.",
		Func: func() error {
			profile, err := displayCreateForm()
			if err != nil {
				return err
			}
			profilePath, err := getProfilePath(profile.Name)
			if err != nil {
				return err
			}
			err = saveProfile(profilePath, profile.Name, profile.Config)
			if err != nil {
				return err
			}
			fmt.Println(successMessage(profile.Name, CREATE))
			return nil
		},
	},
	USE: {
		Description: "Select a profile to use.",
		Func: func() error {
			if !isGlobal && !isGitDirectory() {
				return ErrNotGitDirectory
			}
			selected, err := displayProfileSelector(profiles)
			if err != nil {
				return err
			}
			if selected.Name == defaultConfigName {
				err = unsetConfig(fmt.Sprintf(`%s.*\.gitconfig$`, saveDirName))
				if err != nil {
					return err
				}
				goto successMsg
			}
			err = applyConfig(filepath.Join(saveDirPath, selected.DirName, ".gitconfig"), isGlobal)
			if err != nil {
				return err
			}
		successMsg:
			fmt.Println(successMessage(selected.Name, USE))
			return nil
		},
	},
	LIST: {
		Description: "List all available profiles.",
		Func: func() error {
			err := displayProfileList(profiles)
			if err != nil {
				return err
			}
			return nil
		},
	},
	EDIT: {
		Description: "Edit an existing profile in text editor.",
		Func: func() error {
			var (
				selected Profile
				err      error
			)
			if isGlobal {
				err = openTextEditor(filepath.Join(userHomeDir, ".gitconfig"))
				if err != nil {
					return err
				}
				selected.Name = ".gitconfig"
				goto successMsg
			}
			selected, err = displayProfileSelector(profiles)
			if err != nil {
				return err
			}
			if selected.Name == "default" {
				return ErrEditDefaultConfig
			}
			err = openTextEditor(filepath.Join(saveDirPath, selected.DirName, ".gitconfig"))
			if err != nil {
				return err
			}
		successMsg:
			fmt.Println(successMessage(selected.Name, EDIT))
			return nil
		},
	},
	DELETE: {
		Description: "Delete an existing profile.",
		Func: func() error {
			var (
				selected     Profile
				err          error
				deleteGlobal bool
			)
			if isGlobal {
				selected.Name = ".gitconfig"
				selected.DirName, err = hash(defaultConfigName)
				if err != nil {
					return err
				}
				deleteGlobal = displayDeleteConfirmation()
				if deleteGlobal {
					goto deleteConfig
				} else {
					return nil
				}
			}
			selected, err = displayProfileSelector(profiles)
			if err != nil {
				return err
			}
			if selected.Name == defaultConfigName {
				return ErrDeleteDefaultConfig
			}
		deleteConfig:
			err = unsetConfig(fmt.Sprintf(`%s.*%s.\.gitconfig$`, saveDirName, selected.DirName))
			if err != nil {
				return err
			}
			err = os.RemoveAll(filepath.Join(saveDirPath, selected.DirName))
			if err != nil {
				return err
			}
			if deleteGlobal {
				err = os.Remove(filepath.Join(userHomeDir, ".gitconfig"))
				if err != nil {
					return err
				}
			}
			fmt.Println(successMessage(selected.Name, DELETE))
			return nil
		},
	},
}
