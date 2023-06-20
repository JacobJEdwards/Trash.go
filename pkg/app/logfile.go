package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/JacobJEdwards/Trash.go/pkg/config"
)

type LogEntry struct {
	TrashTime    time.Time
	OriginalName string
	OriginalPath string
}

func generateLogFileEntry(logEntry *LogEntry) string {
	return fmt.Sprintf("%v %v %v\n", logEntry.TrashTime.Format(time.RFC3339), logEntry.OriginalName, logEntry.OriginalPath)
}

func SetLog(trashedFile *os.File, c *config.Config) error {
	logFile := c.Logfile

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to open log file '%s': %v", logFile, err))
	}

	defer file.Close()

	fileInfo, err := trashedFile.Stat()

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get file info for '%s': %v", trashedFile.Name(), err))
	}

	absPath, err := filepath.Abs(trashedFile.Name())

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get absolute path for '%s': %v", trashedFile.Name(), err))
	}

	logEntry := LogEntry{
		TrashTime:    time.Now(),
		OriginalName: fileInfo.Name(),
		OriginalPath: absPath,
	}

	logEntryString := generateLogFileEntry(&logEntry)

	_, err = file.WriteString(logEntryString)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to write log entry to '%s': %v", logFile, err))
	}

	GetLog(c)

	return file.Sync()
}

func GetLog(c *config.Config) ([]LogEntry, error) {
	logFile := c.Logfile

	file, err := os.OpenFile(logFile, os.O_RDONLY|os.O_CREATE, 0775)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to open log file '%s': %v", logFile, err))
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var logEntries []LogEntry

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		logEntry := strings.Split(line, " ")
		trashTime, err := time.Parse(time.RFC3339, logEntry[0])

		if err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to parse trash time '%s': %v", logEntry[0], err))
		}

		logEntries = append(logEntries, LogEntry{
			TrashTime:    trashTime,
			OriginalName: logEntry[1],
			OriginalPath: logEntry[2],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to scan log file '%s': %v", logFile, err))
	}

	return logEntries, nil
}

func RemoveLog(logEntry *LogEntry, c *config.Config) error {
	logEntries, err := GetLog(c)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get log entries: %v", err))
	}

	for i, entry := range logEntries {
		if entry.OriginalPath == logEntry.OriginalPath {
			logEntries = append(logEntries[:i], logEntries[i+1:]...)

			break
		}
	}

	err = WriteAllEntries(logEntries, c)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to write log entries: %v", err))
	}

	return nil
}

func WriteAllEntries(logEntries []LogEntry, c *config.Config) error {
	logFile := c.Logfile

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0775)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to open log file '%s': %v", logFile, err))
	}

	defer file.Close()

	file.Truncate(0)

	for _, entry := range logEntries {
		logEntryString := generateLogFileEntry(&entry)

		_, err = file.WriteString(logEntryString)

		if err != nil {
			return errors.New(fmt.Sprintf("Failed to write log entry to '%s': %v", logFile, err))
		}
	}

	return file.Sync()
}
