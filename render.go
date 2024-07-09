package main

import (
	"fmt"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/manifoldco/promptui"
	"github.com/thansetan/git-sw/pkg/gitconfig"
	"golang.org/x/crypto/ssh"
)

func validateNotEmpty(s string) error {
	if len(s) <= 0 {
		return ErrEmptyField
	}
	return nil
}

func displayCreateForm() (Profile, error) {
	var (
		profile Profile
		err     error
	)
	profile.Config = gitconfig.New()

	profileNamePrompt := promptui.Prompt{
		Label: "Name",
		Validate: func(s string) error {
			err := validateNotEmpty(s)
			if err != nil {
				return err
			}
			profilePath, err := getProfilePath(s)
			if err != nil {
				return err
			}
			_, err = os.Stat(profilePath)
			if err == nil {
				return ErrDuplicateProfile
			}
			return nil
		},
	}

	gitNamePrompt := promptui.Prompt{
		Label: "Git Username",
		Validate: func(s string) error {
			err := validateNotEmpty(s)
			if err != nil {
				return err
			}
			err = gitconfig.ValidateValue(s)
			if err != nil {
				return err
			}
			return nil
		},
	}

	gitEmailPrompt := promptui.Prompt{
		Label: "Git Email",
		Validate: func(s string) error {
			err := validateNotEmpty(s)
			if err != nil {
				return err
			}
			_, err = mail.ParseAddress(s)
			if err != nil {
				return ErrInvalidEmail
			}
			err = gitconfig.ValidateValue(s)
			if err != nil {
				return err
			}
			return nil
		},
	}

	gitWithSigningKeyPrompt := promptui.Prompt{
		Label:     "Add Signing Key",
		IsConfirm: true,
	}

	gitGPGFormatSelect := promptui.Select{
		Label:    "Select Key Format",
		Items:    gpgFormat,
		HideHelp: true,
	}

	profile.Name, err = profileNamePrompt.Run()
	if err != nil {
		return Profile{}, err
	}
	name, err := gitNamePrompt.Run()
	if err != nil {
		return Profile{}, err
	}
	err = profile.Config.Set("user.name", name)
	if err != nil {
		return Profile{}, err
	}
	email, err := gitEmailPrompt.Run()
	if err != nil {
		return Profile{}, err
	}
	err = profile.Config.Set("user.email", email)
	if err != nil {
		return Profile{}, err
	}
	_, err = gitWithSigningKeyPrompt.Run()
	if err == nil {
		ix, _, err := gitGPGFormatSelect.Run()
		if err != nil {
			return Profile{}, err
		}
		keyFormat := gpgFormat[ix]
		signingKey, err := getSigningKeyPrompt(keyFormat).Run()
		if err != nil {
			return Profile{}, err
		}
		err = profile.Config.Set("gpg.format", string(keyFormat))
		if err != nil {
			return Profile{}, err
		}
		err = profile.Config.Set("user.signingKey", signingKey)
		if err != nil {
			return Profile{}, err
		}
		err = profile.Config.Set("commit.gpgsign", "true")
		if err != nil {
			return Profile{}, err
		}
	}

	return profile, nil
}

func displayProfileSelector(profiles []Profile) (Profile, error) {
	keys := &promptui.SelectKeys{
		Prev:     promptui.Key{Code: promptui.KeyPrev, Display: promptui.KeyPrevDisplay},
		Next:     promptui.Key{Code: promptui.KeyNext, Display: promptui.KeyNextDisplay},
		PageUp:   promptui.Key{Code: promptui.KeyBackward, Display: promptui.KeyBackwardDisplay},
		PageDown: promptui.Key{Code: promptui.KeyForward, Display: promptui.KeyForwardDisplay},
		Search:   promptui.Key{Code: 3},
	}

	prompt := promptui.Select{
		Label: "Profile",
		Items: profiles,
		Size:  5,
		Keys:  keys,
		Searcher: func(input string, index int) bool {
			profile := strings.ReplaceAll(profiles[index].Name, " ", "")
			input = strings.TrimSpace(strings.ReplaceAll(input, " ", ""))
			return strings.Contains(strings.ToLower(profile), strings.ToLower(input))
		},
		HideHelp: true,
		Templates: &promptui.SelectTemplates{
			Label:    "Please select one of the available profiles",
			Active:   "> {{ .Name | cyan }}\t{{ if .IsActive }}{{ \"active\" | green }} {{ end }}",
			Inactive: "  {{ .Name | blue }}\t{{ if .IsActive }}{{ \"active\" | green }} {{ end }}",
			Selected: "> {{ .Name | cyan }}",
		},
		StartInSearchMode: true,
	}

	ix, _, err := prompt.Run()
	if err != nil {
		return Profile{}, err
	}
	return profiles[ix], nil
}

func displayProfileList(profiles []Profile) error {
	tw := tabwriter.NewWriter(os.Stdout, 4, 4, 0, ' ', 0)
	_, err := tw.Write([]byte("List of available profiles:\n"))
	if err != nil {
		return err
	}
	for i, profile := range profiles {
		fmt.Fprintf(tw, "%d.\tName: %s ", i+1, profile.Name)
		if profile.IsActive {
			fmt.Fprint(tw, promptui.Styler(promptui.FGGreen)("(active)"))
		}
		fmt.Fprint(tw, "\n")
		fmt.Fprintf(tw, "\tPath: %s\n", filepath.Join(saveDirPath, profile.DirName))
	}
	err = tw.Flush()
	if err != nil {
		return err
	}
	return nil
}

func displayDeleteConfirmation() bool {
	deletePrompt := promptui.Prompt{
		Label:     "You're about to delete a GLOBAL config file, do you want to proceed",
		IsConfirm: true,
	}

	_, err := deletePrompt.Run()
	if err != nil {
		return false
	}
	return err == nil
}

func getSigningKeyPrompt(keyFormat GPGFormat) *promptui.Prompt {
	prompt := new(promptui.Prompt)

	switch keyFormat {
	case OPENPGP:
		prompt.Label = "Enter your GPG key"
		prompt.Validate = func(s string) error {
			err := validateNotEmpty(s)
			if err != nil {
				return err
			}

			err = gitconfig.ValidateValue(s)
			if err != nil {
				return err
			}

			return nil
		}
	case SSH:
		prompt.Label = "Enter path to your public key"
		prompt.Validate = func(s string) error {
			err := validateNotEmpty(s)
			if err != nil {
				return err
			}

			err = gitconfig.ValidateValue(s)
			if err != nil {
				return err
			}

			if filepath.Ext(s) != ".pub" {
				return ErrInvalidPublicKeyExt
			}
			content, err := os.ReadFile(s)
			if err != nil {
				return err
			}
			_, _, _, _, err = ssh.ParseAuthorizedKey(content)
			if err != nil {
				return err
			}
			return nil
		}
	case X509:
		prompt.Label = "Enter your certificate ID"
		prompt.Validate = func(s string) error {
			err := validateNotEmpty(s)
			if err != nil {
				return err
			}

			err = gitconfig.ValidateValue(s)
			if err != nil {
				return err
			}

			return nil
		}
	}
	return prompt
}
