package utils

import (
	"fmt"
	"os/user"
	"path/filepath"
)

func ExpandTilde(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("ExpandTilde: user.Current: %v", err)
	}

	return filepath.Join(usr.HomeDir, path[1:]), nil
}
