package app

import (
	"errors"
	"fmt"
	"sync"

	"github.com/JacobJEdwards/Trash.go/pkg/config"
)

func RestoreAll(c *config.Config) ([]LogEntry, error) {
	logEntries, err := GetLog(c)

	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(len(logEntries))

	results := make(chan *LogEntry, len(logEntries))
	restoreErrs := make(chan error, len(logEntries))

	for _, logEntry := range logEntries {
        defer wg.Done()

		go func(logEntry *LogEntry) {
			result, err := RestoreFile(logEntry.OriginalName, c)
			if err != nil {
				restoreErrs <- errors.New(fmt.Sprintf("Error restoring %s: %s\n", logEntry.OriginalName, err))
			} else {
				results <- result
			}
		}(&logEntry)
	}

	wg.Wait()
	close(results)
	close(restoreErrs)

	var restoredEntries []LogEntry
	var restoreErrors []error

	for result := range results {
		restoredEntries = append(restoredEntries, *result)
	}

	for err := range restoreErrs {
		restoreErrors = append(restoreErrors, err)
	}

	if len(restoreErrors) > 0 {
		return restoredEntries, errors.New(fmt.Sprintf("Errors restoring files: %s", restoreErrors))
	}

	return restoredEntries, nil
}
