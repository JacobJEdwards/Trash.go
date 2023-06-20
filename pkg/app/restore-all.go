package app

import (
	"fmt"
	"github.com/JacobJEdwards/Trash.go/pkg/config"
)

func RestoreAll(c *config.Config) ([]LogEntry, error) {
	logEntries, err := GetLog(c)

	if err != nil {
		return nil, err
	}

	for _, logEntry := range logEntries {

		_, err := RestoreFile(logEntry.OriginalName, c)

		if err != nil {
			fmt.Printf("Error restoring %s: %s\n", logEntry.OriginalName, err)
			continue
		}
	}

	return logEntries, nil
}
