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
	help       bool
	emptyTrash bool
	restoreAll bool
	restore    string
	rm         string
)

type TrashCLI struct {
	config *config.Config
}

func init() {
	flag.BoolVar(&help, "help", false, "print help")
	flag.BoolVar(&help, "h", false, "print help")

	flag.BoolVar(&view, "view", false, "view the trash")
	flag.BoolVar(&view, "v", false, "view the trash")

	flag.BoolVar(&emptyTrash, "empty", false, "empty the trash")
	flag.BoolVar(&emptyTrash, "e", false, "empty the trash")

	flag.BoolVar(&restoreAll, "restore-all", false, "restore all files")
	flag.BoolVar(&restoreAll, "ra", false, "restore all files")

	flag.StringVar(&restore, "restore", "", "restore a file")
	flag.StringVar(&restore, "r", "", "restore a file")

	flag.StringVar(&rm, "delete", "", "trash a file")
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
	case help:
		flag.Usage()
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
		fmt.Printf("error getting log: %v\n", err)
	}

	if len(logEntries) == 0 {
		fmt.Println("Trash is empty")
		return
	}

	app.OutputLogEntries(logEntries)
}

func (cli *TrashCLI) clearTrash() {
	err := app.EmptyTrash(cli.config)

	if err != nil {
		fmt.Printf("Error emptying trash: %v\n", err)
	}
}

func (cli *TrashCLI) restoreAllFiles() {
	logEntries, err := app.RestoreAll(cli.config)

	if err != nil {
		fmt.Printf("Error restoring all files: %v\n", err)
		return
	}

	if len(logEntries) == 0 {
		fmt.Println("Trash is empty")
		return
	}

	fmt.Println("Restored files:")
	app.OutputLogEntries(logEntries)
}

func (cli *TrashCLI) restoreFile(file string) {
	fileEntry, err := app.RestoreFile(file, cli.config)

	if err != nil {
		fmt.Printf("Error restoring file: %v\n", err)
	}

    fmt.Println("Restored:")
    app.OutputLogEntry(fileEntry)
}

func (cli *TrashCLI) deleteFile(file string) {
	fileEntry, err := app.DeleteFile(file, cli.config)

	if err != nil {
		fmt.Printf("Error deleting file: %v\n", err)
	}

    fmt.Println("Deleted:")
    app.OutputLogEntry(fileEntry)
}

func (cli *TrashCLI) trashFiles(files []string) {
	for _, fileName := range files {

		if _, err := os.Stat(fileName); err != nil {
			fmt.Printf("Error using file: %v\n", err)
			continue
		}

		file, err := os.Open(fileName)

		if err != nil {
			fmt.Printf("Error using file: %v\n", err)
			continue
		}

		defer file.Close()

		err = app.TrashFile(file, cli.config)

		if err != nil {
			fmt.Printf("Error trashing file: %v\n", err)
			return
		}

		fmt.Printf("Trashed %s\n", fileName)
	}
}

func printUsage() {
	fmt.Println("Usage: trash [flags] [files...]")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("")
	fmt.Println("  -h, --help")
	fmt.Println("        Print this help message")
	fmt.Println("")
	fmt.Println("  -v, --view")
	fmt.Println("        View the trash")
	fmt.Println("")
	fmt.Println("  -e, --empty")
	fmt.Println("        Empty the trash")
	fmt.Println("")
	fmt.Println("  -ra, --restore-all")
	fmt.Println("        Restore all files")
	fmt.Println("")
	fmt.Println("  -r, --restore")
	fmt.Println("        Restore a file")
	fmt.Println("")
	fmt.Println("  -rm, --delete")
	fmt.Println("        Delete a file")
	fmt.Println("")
}

func main() {
	cli, err := NewTrashCLI()

	if err != nil {
		fmt.Printf("Error creating CLI: %v\n", err)
		return
	}

	cli.Run()
}

