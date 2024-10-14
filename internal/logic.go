package internal

import (
	"errors"
	"os"
	"regexp"
)

// understand regexp and finish this function!!!
func validateBucketName(name string) error {
	// Regex for valid bucket name format
	validBucketName := regexp.MustCompile(`^[a-z0-9]+(?:[a-z0-9-]*[a-z0-9])?$`)
	// Regex to prevent IP address format
	ipFormat := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$`)

	// Length check
	if len(name) < 3 || len(name) > 63 {
		return errors.New("Error: Bucket name must be between 3 and 63 characters long.")
	}

	// Validate characters
	if !validBucketName.MatchString(name) {
		return errors.New("Error: Bucket name contains invalid characters.")
	}

	// Check for IP address-like format
	if ipFormat.MatchString(name) {
		return errors.New("Error: Bucket name must not be formatted like an IP address.")
	}

	return nil
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
