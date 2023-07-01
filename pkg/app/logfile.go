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

func SetLog(trashedFile *os.File, c *config.Config) error {
	logFile := c.Logfile

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to open log file '%s': %v", logFile, err))
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

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

	logEntryString := GenerateLogFileEntry(&logEntry)

	_, err = file.WriteString(logEntryString)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to write log entry to '%s': %v", logFile, err))
	}

	_, err = GetLog(c)
	if err != nil {
		return err
	}

	return file.Sync()
}

func GetLog(c *config.Config) ([]LogEntry, error) {
	logFile := c.Logfile

	file, err := os.OpenFile(logFile, os.O_RDONLY|os.O_CREATE, 0775)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to open log file '%s': %v", logFile, err))
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

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

func RemoveLog(filename string, c *config.Config) (*LogEntry, error) {

	logEntries, err := GetLog(c)
	var foundEntry LogEntry

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to get log entries: %v", err))
	}

	for i, entry := range logEntries {
		if entry.OriginalName == filename {
			logEntries = append(logEntries[:i], logEntries[i+1:]...)
			foundEntry = entry

			break
		}
	}

	if foundEntry == (LogEntry{}) {
		return nil, errors.New(fmt.Sprintf("Failed to find log entry for '%s'", filename))
	}

	err = WriteAllEntries(logEntries, c)

	if err != nil {
		return &foundEntry, errors.New(fmt.Sprintf("Failed to write log entries: %v", err))
	}

	return &foundEntry, nil
}

func WriteAllEntries(logEntries []LogEntry, c *config.Config) error {
	logFile := c.Logfile

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0775)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to open log file '%s': %v", logFile, err))
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	for _, entry := range logEntries {
		logEntryString := GenerateLogFileEntry(&entry)

		_, err = file.WriteString(logEntryString)

		if err != nil {
			return errors.New(fmt.Sprintf("Failed to write log entry to '%s': %v", logFile, err))
		}
	}

	return file.Sync()
}
