package toolchain

import (
	"os"
	"path/filepath"
	"runtime"
)

func GoPath() string {
	base, _ := os.Getwd()
	if runtime.GOOS == "windows" {
		return filepath.Clean(filepath.Join(base, ".cache", "go", "bin", "go.exe"))

	} else {
		return filepath.Clean(filepath.Join(base, ".cache", "go", "bin", "go"))

	}
}

func Exists() bool {

	_, err := os.Stat(GoPath())
	return err == nil
}
