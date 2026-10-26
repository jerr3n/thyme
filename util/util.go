package util

import (
	"os"
	"path/filepath"
	"strings"
)

func ExpandHome(in string) (string, error) {
	if in != "~" && !strings.HasPrefix(in, "~/") {
		return in, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, in[1:]), nil
}
