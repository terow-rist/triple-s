package internal

import (
	"errors"
	"os"
	"regexp"
)

func validateBucketName(name string) error {
	validBucketName := regexp.MustCompile(`^[a-z0-9.](?:[a-z0-9.-]*[a-z0-9.])?$`)
	ipFormat := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)

	if len(name) < 3 || len(name) > 63 {
		return errors.New("Error: Bucket name must be between 3 and 63 characters long.")
	}

	if !validBucketName.MatchString(name) {
		return errors.New("Error: Bucket name contains invalid characters.")
	}

	if ipFormat.MatchString(name) {
		return errors.New("Error: Bucket name must not be formatted like an IP address.")
	}

	if checkConsecutive(name) {
		return errors.New("Error: Bucket must not contain two consecutive periods or dashes.")
	}

	return nil
}

func checkConsecutive(str string) bool {
	for i := 0; i < len(str)-1; i++ {
		if str[i] == '.' && str[i+1] == '.' {
			return true
		}
		if str[i] == '-' && str[i+1] == '-' {
			return true
		}
	}
	return false
}

func isBucketEmpty(path string) (bool, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	return len(dirEntries) == 0, nil
}

func isStandardPackage(packageName string) bool {
	return packageName == "cmd" || packageName == "config" || packageName == "internal"
}
