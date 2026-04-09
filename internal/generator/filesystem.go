package generator

import (
	"os"
	"path/filepath"
)

func createDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func join(paths ...string) string {
	return filepath.Join(paths...)
}
