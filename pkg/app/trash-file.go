package app

import (
	"fmt"
	"github.com/JacobJEdwards/Trash.go/pkg/config"
	"os"
	"path/filepath"
)

func TrashFiles(files []*os.File, c *config.Config) error {
	for _, file := range files {
		// Get the path to the trash directory
		trashDir := filepath.Join(c.TrashDir, file.Name())

		// Move the file to the trash directory
		err := os.Rename(file.Name(), trashDir)

		if err != nil {
			return fmt.Errorf("Error moving file %s to trash directory %s: %s", file.Name(), trashDir, err)
		}

		err = SetLog(file, c)

		if err != nil {
			return fmt.Errorf("Error setting log for file %s: %s", file.Name(), err)
		}
	}

	return nil
}

func TrashFile(file *os.File, c *config.Config) error {
	// Get the path to the trash directory
	trashDir := filepath.Join(c.TrashDir, file.Name())

	// Move the file to the trash directory
	err := os.Rename(file.Name(), trashDir)
	if err != nil {
		return fmt.Errorf("Error moving file %s to trash directory %s: %s", file.Name(), trashDir, err)
	}

	err = SetLog(file, c)

	if err != nil {
		return fmt.Errorf("Error setting log for file %s: %s", file.Name(), err)
	}

	return nil
}
