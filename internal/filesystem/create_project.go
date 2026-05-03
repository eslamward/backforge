package filesystem

import (
	"os"
	"path/filepath"
	"strings"
)

func WriteProjectFile(baseFile, filePath, content string) error {

	fullPath := filepath.Join("output", baseFile, filePath)

	listOFFolder := strings.Split(fullPath, "/")

	allDir := strings.Join(listOFFolder[:len(listOFFolder)], "/")

	err := os.MkdirAll(filepath.Dir(allDir), 0755)
	if err != nil {
		return err
	}

	return writeFile(fullPath, content)

}

func CreatFolder(folderName string) error {
	fullPath := filepath.Join("output", folderName)

	err := os.MkdirAll(filepath.Dir(fullPath), 0755)
	if err != nil {
		return err
	}
	return nil
}
