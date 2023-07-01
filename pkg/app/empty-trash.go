package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/JacobJEdwards/Trash.go/pkg/config"
	"github.com/JacobJEdwards/Trash.go/pkg/utils"
)

func EmptyTrash(c *config.Config) error {
	trashDir := c.TrashDir
	logFilepath := c.Logfile

	if _, err := os.Stat(logFilepath); os.IsNotExist(err) {
		return fmt.Errorf("log file does not exist: %s", logFilepath)
	}

	if _, err := os.Stat(trashDir); os.IsNotExist(err) {
		return fmt.Errorf("trash directory does not exist: %s", trashDir)
	}

	proceed := utils.ProceedTask("Are you sure you want to empty the trash? [y/N] ")

	if !proceed {
		return nil
	}
	files, err := os.ReadDir(trashDir)

	if err != nil {
		return fmt.Errorf("error reading trash directory: %v", err)
	}

	for _, file := range files {
		filePath := filepath.Join(trashDir, file.Name())
		err := os.RemoveAll(filePath)

		if err != nil {
			return fmt.Errorf("error removing file %s: %v", file.Name(), err)
		}
	}

	err = os.Truncate(logFilepath, 0)

	if err != nil {
		return fmt.Errorf("error writing to log file: %v", err)
	}

	return nil
}
