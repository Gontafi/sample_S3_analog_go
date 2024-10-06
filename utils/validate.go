package utils

import (
	"regexp"
	"strings"
)

const validNamePattern = `^[a-z0-9]([a-z0-9\-]{0,61}[a-z0-9]|(?:[a-z0-9]{1,61}(\.[a-z0-9]{1,61})*))[a-z0-9]$`

var ErrInvalidBucketName = "invalid match"

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
