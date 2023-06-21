package app

import (
	"fmt"
	"time"
)

func GenerateLogFileEntry(logEntry *LogEntry) string {
	return fmt.Sprintf("%v %v %v\n", logEntry.TrashTime.Format(time.RFC3339), logEntry.OriginalName, logEntry.OriginalPath)
}

func OutputLogEntries(logEntries []LogEntry) {
	for _, logEntry := range logEntries {
        OutputLogEntry(&logEntry)
	}
}

func OutputLogEntry(logEntry *LogEntry) {
    formattedTime := logEntry.TrashTime.Format("2006-01-02 15:04:05")
    outString := fmt.Sprintf("%s %s %s", formattedTime, logEntry.OriginalName, logEntry.OriginalPath)
    fmt.Println(outString)
}
