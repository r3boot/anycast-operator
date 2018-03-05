package utils

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"
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

func AddCIDR(addr string) string {
	if !strings.Contains(addr, "/") {
		if strings.Contains(addr, ":") {
			// IPv6
			addr += "/128"
		} else {
			// IPv4
			addr += "/32"
		}
	}

	return addr
}
