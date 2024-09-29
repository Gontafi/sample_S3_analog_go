package utils

import "regexp"

const bucketNamePattern = `^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$`

var (
	ErrInvalidBucketName = "invalid match"
)

func ValidateBucketName(name string) error {
	_, err := regexp.MatchString(bucketNamePattern, name)
	if err != nil {
		return err
	}

	return nil
}
