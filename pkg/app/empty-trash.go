package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/JacobJEdwards/Trash.go/pkg/config"
)

func EmptyTrash(c *config.Config) error {
	trashDir := c.TrashDir
	logFile := c.Logfile
	proceed := proceedEmptyTrash()

	if !proceed {
		return nil
	}

	log, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer log.Close()

	files, err := os.ReadDir(trashDir)

	if err != nil {
		return fmt.Errorf("Error reading trash directory: %v", err)
	}

	for _, file := range files {
		filePath := filepath.Join(trashDir, file.Name())
		err := os.RemoveAll(filePath)

		if err != nil {
			return fmt.Errorf("Error removing file %s: %v", file.Name(), err)
		}
	}

	err = os.Truncate(logFile, 0)

	if err != nil {
		return fmt.Errorf("Error writing to log file: %v", err)
	}

	return nil
}

func proceedEmptyTrash() bool {
	var proceed string

	fmt.Print("Are you sure you want to empty the trash? [y/N] ")
	_, err := fmt.Scanln(&proceed)

	if err != nil {
		return false
	}

	switch proceed {
	case "y", "Y", "yes", "Yes":
		return true
	default:
		return false
	}
}
