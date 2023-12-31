package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/JacobJEdwards/Trash.go/pkg/config"
)

func RestoreFile(fileName string, c *config.Config) (*LogEntry, error) {
	trashDir := c.TrashDir

	trashedFiles, err := os.ReadDir(trashDir)

	if err != nil {
		return nil, err
	}

	for _, trashedFile := range trashedFiles {

		if trashedFile.Name() == fileName {

			trashPath := filepath.Join(trashDir, trashedFile.Name())

			entry, err := RemoveLog(fileName, c)

			if err == nil {
				err = os.Rename(trashPath, entry.OriginalPath)

				if err != nil {
					return nil, err
				}

				return entry, nil
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
