package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/JacobJEdwards/Trash.go/pkg/config"
)

type LogEntry struct {
	TrashTime    time.Time
	OriginalName string
	OriginalPath string
}

var (
	logMu sync.Mutex
)

func SetLog(trashedFile *os.File, c *config.Config) error {
	logMu.Lock()
	defer logMu.Unlock()

	logFile := c.Logfile

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to open log file '%s': %v", logFile, err))
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	defer func() {
		if err := writer.Flush(); err != nil {
			fmt.Println(err)
		}
	}()

	if _, err := trashedFile.Stat(); err != nil {
		return errors.New(fmt.Sprintf("Failed to get file info for '%s': %v", trashedFile.Name(), err))
	}

	absPath, err := filepath.Abs(trashedFile.Name())

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get absolute path for '%s': %v", trashedFile.Name(), err))
	}

	logEntry := LogEntry{
		TrashTime:    time.Now(),
		OriginalName: trashedFile.Name(),
		OriginalPath: absPath,
	}

	logEntryString := GenerateLogFileEntry(&logEntry)

	_, err = writer.WriteString(logEntryString)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to write log entry to '%s': %v", logFile, err))
	}

	return file.Sync()
}

func GetLog(c *config.Config) ([]LogEntry, error) {
	logFile := c.Logfile

	logMu.Lock()
	defer logMu.Unlock()

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

func RemoveLog(filename string, c *config.Config) (*LogEntry, error) {
	logMu.Lock()
	defer logMu.Unlock()

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
	logMu.Lock()
	defer logMu.Unlock()

	logFile := c.Logfile

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0775)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to open log file '%s': %v", logFile, err))
	}

	defer file.Close()

	file.Truncate(0)

	var newLogEntries string
	var sb strings.Builder

	for _, entry := range logEntries {
		logEntryString := GenerateLogFileEntry(&entry)

		sb.WriteString(logEntryString)
	}

	newLogEntries = sb.String()

	if newLogEntries == "" {
		return nil
	}

	_, err = file.WriteString(newLogEntries)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to write log entry to '%s': %v", logFile, err))
	}

	return file.Sync()
}
