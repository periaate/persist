package partdb

import (
	"fmt"
	"os"
)

// EnsureDir checks if a directory exists at the given path.
// If the path exists and is a directory, it returns nil.
// If the path exists but is not a directory, it returns an error.
// If the path does not exist, it creates the directory and returns nil.
func EnsureDir(path string) error {
	fileInfo, err := os.Stat(path)
	if err == nil {
		if fileInfo.IsDir() {
			return nil
		}
		return fmt.Errorf("path exists but is not a directory: %s", path)
	}

	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %s", err)
		}
		return nil
	}

	return fmt.Errorf("unknown error: %s", err)
}
