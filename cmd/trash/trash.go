package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/JacobJEdwards/Trash.go/pkg/app"
	"github.com/JacobJEdwards/Trash.go/pkg/config"
)

var (
	view       bool
	emptyTrash bool
	restoreAll bool
	restore    string
	rm         string
)

type TrashCLI struct {
	config *config.Config
}

func init() {
	flag.BoolVar(&view, "view", false, "view the trash")
	flag.BoolVar(&emptyTrash, "empty", false, "empty the trash")
	flag.BoolVar(&restoreAll, "restore-all", false, "restore all files")
	flag.StringVar(&restore, "restore", "", "restore a file")
	flag.StringVar(&rm, "rm", "", "trash a file")

}

func NewTrashCLI() (*TrashCLI, error) {
	config, err := config.LoadConfig()

	if err != nil {
		return nil, err
	}

	flag.Parse()
	flag.Usage = printUsage

	return &TrashCLI{
		config: config,
	}, nil
}

func (cli *TrashCLI) Run() {
	if flag.NFlag() == 0 && len(flag.Args()) == 0 {
		fmt.Println("Error: must pass exactly one flag or at least one file")
		flag.Usage()
		return
	}

	if flag.NFlag() > 1 {
		fmt.Println("Error: must pass exactly one flag or at least one file")
		flag.Usage()
		return
	}

	switch {
	case view:
		cli.viewTrash()
	case emptyTrash:
		cli.clearTrash()
	case restoreAll:
		cli.restoreAllFiles()
	case restore != "":
		cli.restoreFile(restore)
	case rm != "":
		cli.deleteFile(rm)
	default:
		cli.trashFiles(flag.Args())
	}
}

func (cli *TrashCLI) viewTrash() {
	logEntries, err := app.GetLog(cli.config)

	if err != nil {
		fmt.Printf("error getting log: %v", err)
	}

	if len(logEntries) == 0 {
		fmt.Println("Trash is empty")
		return
	}

	cli.outputLogEntries(logEntries)

	return
}

func (cli *TrashCLI) clearTrash() {
	err := app.EmptyTrash(cli.config)

	if err != nil {
		fmt.Printf("Error emptying trash: %v", err)
	}

	return
}

func (cli *TrashCLI) restoreAllFiles() {
	logEntries, err := app.RestoreAll(cli.config)

	if err != nil {
		fmt.Printf("Error restoring all files: %v", err)
		return
	}

	if len(logEntries) == 0 {
		fmt.Println("Trash is empty")
		return
	}

	fmt.Println("Restored files:")
	cli.outputLogEntries(logEntries)

	return
}

func (cli *TrashCLI) restoreFile(file string) {
	fileEntry, err := app.RestoreFile(file, cli.config)

	if err != nil {
		fmt.Printf("Error restoring file: %v", err)
	}

	formattedTime := fileEntry.TrashTime.Format("2006-01-02 15:04:05")
	outString := fmt.Sprintf("%s %s %s", formattedTime, fileEntry.OriginalName, fileEntry.OriginalPath)

	fmt.Printf("Restored %s\n", outString)

	return
}

func (cli *TrashCLI) deleteFile(file string) {
	return
}

func (cli *TrashCLI) trashFiles(files []string) {
	for _, fileName := range files {

		_, err := os.Stat(fileName)

		if err != nil {
			fmt.Printf("Error using file: %v", err)
			continue
		}

		file, err := os.Open(fileName)

		if err != nil {
			fmt.Printf("Error using file: %v", err)
			continue
		}

		defer file.Close()

		err = app.TrashFile(file, cli.config)

		if err != nil {
			fmt.Printf("Error trashing file: %v", err)
		}

		fmt.Printf("Trashed %s\n", fileName)
	}
}

func printUsage() {
	fmt.Println("Usage: trash [flags] [files...]")
	fmt.Println("Flags:")
	flag.PrintDefaults()
}

func main() {
	cli, err := NewTrashCLI()

	if err != nil {
		fmt.Printf("Error creating CLI: %v", err)
		return
	}

	cli.Run()
}

func (cli *TrashCLI) outputLogEntries(logEntries []app.LogEntry) {
	for _, logEntry := range logEntries {
		formattedTime := logEntry.TrashTime.Format("2006-01-02 15:04:05")
		outString := fmt.Sprintf("%s %s %s", formattedTime, logEntry.OriginalName, logEntry.OriginalPath)
		fmt.Println(outString)
	}
}
