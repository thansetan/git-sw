package main

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrEmptyField          = errors.New("field can't be empty")
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrDuplicateProfile    = errors.New("profile with given name already exists")
	ErrInvalidAction       = errors.New("invalid action")
	ErrNotImplemented      = errors.New("not implemented")
	ErrEditDefaultConfig   = fmt.Errorf("use '%s -g edit' to edit default config", os.Args[0])
	ErrDeleteDefaultConfig = fmt.Errorf("use '%s -g delete' to delete default config", os.Args[0])
	ErrInvalidPublicKeyExt = errors.New("invalid public key file extension")
	ErrNotGitDirectory     = errors.New("not in a git directory")
)
