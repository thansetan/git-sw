package main

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/thansetan/git-sw/pkg/gitconfig"
)

type Profile struct {
	Config        *gitconfig.GitConfig
	Name, DirName string
	IsActive      bool
}

func getProfilePath(profileName string) (string, error) {
	dirName, err := hash(profileName)
	if err != nil {
		return "", err
	}
	return filepath.Join(saveDirPath, dirName), nil
}

func saveProfile(dirPath string, profileName string, config *gitconfig.GitConfig) (err error) {
	err = os.MkdirAll(dirPath, 0744)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			errDel := os.RemoveAll(dirPath)
			if errDel != nil {
				fmt.Printf("error deleting directory: %s", dirPath)
			}
		}
	}()
	err = config.Save(filepath.Join(dirPath, ".gitconfig"))
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(dirPath, "profile"))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(profileName)
	if err != nil {
		return err
	}
	err = f.Chmod(0444)
	if err != nil {
		return err
	}

	return nil
}

func copyDefault() error {
	var (
		defaultConf *gitconfig.GitConfig
	)
	dirName, err := getProfilePath(defaultConfigName)
	if err != nil {
		return err
	}
	err = os.RemoveAll(dirName)
	if err != nil {
		return err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configContent, err := os.ReadFile(filepath.Join(homeDir, ".gitconfig"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	} else if errors.Is(err, os.ErrNotExist) {
		defaultConf = gitconfig.New()
		goto saveProfile
	}
	defaultConf, err = gitconfig.Parse(configContent)
	if err != nil {
		return err
	}
saveProfile:
	err = saveProfile(dirName, defaultConfigName, defaultConf)
	if err != nil {
		return err
	}

	return nil
}

func getCurrentProfile() (string, error) {
	currentConfig, err := getCurrentConfig()
	if err != nil {
		return "", err
	}
	profileDir := filepath.Dir(currentConfig)
	if profileDir == "." {
		return "default", nil
	}
	profileFile, err := os.Open(filepath.Join(profileDir, "profile"))
	if err != nil {
		return "", err
	}
	defer profileFile.Close()
	profileName, err := io.ReadAll(profileFile)
	if err != nil {
		return "", err
	}

	return string(profileName), nil
}

func getProfiles(configPath string) ([]Profile, error) {
	var profiles []Profile
	currProfile, err := getCurrentProfile()
	if err != nil {
		return nil, err
	}
	err = filepath.WalkDir(configPath, func(path string, d fs.DirEntry, err error) error {
		if d.Name() == "profile" {
			profileFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer profileFile.Close()
			profileName, err := io.ReadAll(profileFile)
			if err != nil {
				return err
			}
			dirName := filepath.Base(filepath.Dir(path))
			h, err := hash(string(profileName))
			if err != nil || h != dirName { // an entity has changed the profile name, should I remove ??? or just skip it ???
				return fs.SkipDir
			}

			profiles = append(profiles, Profile{
				Name:     string(profileName),
				IsActive: string(profileName) == currProfile,
				DirName:  dirName,
			})
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	slices.SortFunc(profiles, func(a, b Profile) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return profiles, nil
}
