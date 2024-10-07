package utils

import (
	"bytes"
	"errors"
	"io"
	"regexp"
	"strings"
	"time"
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
	ErrNotFound   = errors.New("not found")
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

func LineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func CurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
