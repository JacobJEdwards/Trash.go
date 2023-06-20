package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/JacobJEdwards/Trash.go/pkg/app"
	"github.com/JacobJEdwards/Trash.go/pkg/config"
)

var (
	view       = flag.Bool("view", false, "view the trash")
	emptyTrash = flag.Bool("empty", false, "empty the trash")
	restoreAll = flag.Bool("restore-all", false, "restore all files")
	restore    = flag.String("restore", "", "restore a file")
)

func main() {
	flag.Parse()
	args := flag.Args()

	if flag.NFlag() > 1 {
		fmt.Println("Only one flag can be passed at a time")
		return
	}
	if flag.NFlag() == 0 && len(args) == 0 {
		fmt.Println("No arguments passed")
		return
	}

	config, err := config.LoadConfig()

	if err != nil {
		fmt.Printf("error loading config: %v", err)
		return
	}

	if *emptyTrash {
		fmt.Println("Emptying the trash")
		err := app.EmptyTrash(config)

		if err != nil {
			fmt.Printf("Error emptying trash: %v", err)
		}

		return
	}

	if *restoreAll {
		fmt.Println("Restoring all files")
        _, err := app.RestoreAll(config)

        if err != nil {
            fmt.Printf("Error restoring all files: %v", err)
        }

		return
	}

	if *restore != "" {
		fmt.Println("Restoring a file")

        file, err := app.RestoreFile(*restore, config)

		if err != nil {
			fmt.Printf("Error restoring file: %v", err)
		}

		formattedTime := file.TrashTime.Format("2006-01-02 15:04:05")
		outString := fmt.Sprintf("%s %s %s", formattedTime, file.OriginalName, file.OriginalPath)
		fmt.Printf("Restored %s\n", outString)

		return
	}

	if *view {
		fmt.Println("Viewing the trash\n")
		logEntries, err := app.GetLog(config)

		if err != nil {
			fmt.Printf("error getting log: %v", err)
		}

		if len(logEntries) == 0 {
			fmt.Println("Trash is empty")
			return
		}

		for _, logEntry := range logEntries {
			formattedTime := logEntry.TrashTime.Format("2006-01-02 15:04:05")
			outString := fmt.Sprintf("%s %s %s", formattedTime, logEntry.OriginalName, logEntry.OriginalPath)
			fmt.Println(outString)
		}

		return
	}

	// If no flags are passed, assume the user wants to trash a file
	for _, arg := range args {
		file, err := os.Open(arg)

		if err != nil {
			fmt.Printf("Error opening file: %v", err)
			return
		}

		defer file.Close()

		err = app.TrashFile(file, config)

		if err != nil {
			fmt.Printf("Error trashing file: %v", err)
			return
		}

		fmt.Printf("Trashed %s\n", arg)
	}

	return
}
