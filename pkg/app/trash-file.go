package app

import (
	"errors"
	"fmt"
	"github.com/JacobJEdwards/Trash.go/pkg/config"
	"os"
	"path/filepath"
)

func TrashFiles(files []*os.File, c *config.Config) error {
	for _, file := range files {

		err := TrashFile(file, c)

		if err != nil {
			return err
		}
	}

	return nil
}

func TrashFile(file *os.File, c *config.Config) error {
	// Get the path to the trash directory
	fileName := filepath.Base(file.Name())

	if fileName == "." || fileName == ".." || fileName == filepath.Base(c.TrashDir) || fileName == filepath.Base(c.Logfile) {
		return errors.New(fmt.Sprintf("Error trashing file %s: file name is invalid", file.Name()))
	}

	trashDir := filepath.Join(c.TrashDir, fileName)
	absPath, err := filepath.Abs(file.Name())

	if err != nil {
		return errors.New(fmt.Sprintf("Error getting absolute path for file %s: %s", file.Name(), err))
	}

	// Move the file to the trash directory
	err = os.Rename(absPath, trashDir)

	if err != nil {
		return errors.New(fmt.Sprintf("Error moving file %s to trash directory %s: %s", file.Name(), trashDir, err))
	}

	err = SetLog(file, c)

	if err != nil {
		return errors.New(fmt.Sprintf("Error setting log for file %s: %s", file.Name(), err))
	}

	return nil
}
