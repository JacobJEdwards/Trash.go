package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/JacobJEdwards/Trash.go/pkg/config"
	"github.com/JacobJEdwards/Trash.go/pkg/utils"
)

func DeleteFile(filename string, c *config.Config) (*LogEntry, error) {
	trashDir := c.TrashDir
	proceed := utils.ProceedTask(fmt.Sprintf("Are you sure you want to delete %s? [Y/N]", filename))

	if !proceed {
		return nil, errors.New("user aborted delete")
	}

	trashedFiles, err := os.ReadDir(trashDir)

	if err != nil {
		return nil, err
	}

	for _, file := range trashedFiles {
		if file.Name() == filename {
			trashedFilepath := filepath.Join(trashDir, file.Name())

			err := os.Remove(trashedFilepath)

			if err != nil {
				return nil, err
			}

			entry, err := RemoveLog(filename, c)

			if err != nil {
				return nil, err
			}

			return entry, nil

		}
	}

	return nil, errors.New("file not found in trash")

}
