package utils

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidBucketName = "invalid match"
	Directory            = "data"
	Help                 = `Simple Storage Service.

**Usage:**
    triple-s [-port <N>] [-dir <S>]  
    triple-s --help

**Options:**
- --help     Show this screen.
- --port N   Port number
- --dir S    Path to the directory`
)

const validNamePattern = `^[a-z0-9]([a-z0-9\-]{0,61}[a-z0-9]|(?:[a-z0-9]{1,61}(\.[a-z0-9]{1,61})*))[a-z0-9]$`

var (
	ErrNameLength = errors.New("name length must be 3-63")
	ErrValidate   = errors.New("invalid naming")
)

func IsValidBucketName(name string) error {
	if len(name) < 3 || len(name) > 63 {
		return ErrNameLength
	}

	re := regexp.MustCompile(validNamePattern)
	if !re.MatchString(name) {
		return ErrValidate
	}

	parts := strings.Split(name, ".")
	for _, part := range parts {
		if strings.Contains(part, "--") {
			return ErrValidate
		}
	}

	return nil
}
