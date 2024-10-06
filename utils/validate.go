package utils

import (
	"regexp"
	"strings"
)

const validNamePattern = `^[a-z0-9]([a-z0-9\-]{0,61}[a-z0-9]|(?:[a-z0-9]{1,61}(\.[a-z0-9]{1,61})*))[a-z0-9]$`

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

func IsValidBucketName(name string) bool {
	if len(name) < 3 || len(name) > 63 {
		return false
	}

	re := regexp.MustCompile(validNamePattern)
	if !re.MatchString(name) {
		return false
	}

	parts := strings.Split(name, ".")
	for _, part := range parts {
		if strings.Contains(part, "--") {
			return false
		}
	}

	return true
}
