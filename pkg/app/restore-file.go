package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/JacobJEdwards/Trash.go/pkg/config"
)

func RestoreFile(fileName string, c *config.Config) (*LogEntry, error) {
	fileLogs, err := GetLog(c)
	trashDir := c.TrashDir

	if err != nil {
		return nil, err
	}

	trashedFiles, err := os.ReadDir(trashDir)

	if err != nil {
		return nil, err
	}

	for _, trashedFile := range trashedFiles {
		if trashedFile.Name() == fileName {

			trashPath := filepath.Join(trashDir, trashedFile.Name())

			for _, fileDetails := range fileLogs {
				if fileDetails.OriginalName == fileName {
					err := os.Rename(trashPath, fileDetails.OriginalPath)

					if err != nil {
						return nil, err
					}

					err = RemoveLog(&fileDetails, c)

					if err != nil {
						return nil, err
					}

					return &fileDetails, err
				}
			}

			fmt.Printf("file %s not found in logs, moving to home\n", fileName)

			homeDir, err := os.UserHomeDir()

			if err != nil {
				return nil, err
			}

			homePath := filepath.Join(homeDir, fileName)

			err = os.Rename(trashPath, homePath)

			if err != nil {
				return nil, err
			}

		}

	}

	return nil, errors.New(fmt.Sprintf("file %s not found in trash", fileName))
}

func restoreFile(file os.DirEntry, OriginalPath string) error {
	err := os.Rename(file.Name(), OriginalPath)

	return err
}
