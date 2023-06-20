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
	empty      = flag.Bool("empty", false, "empty the trash")
	restoreAll = flag.Bool("restore-all", false, "restore all files")
	restore    = flag.String("restore", "", "restore a file")
)

func main() {
	flag.Parse()
	args := flag.Args()

	config, err := config.LoadConfig()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *empty {
		fmt.Println("Emptying the trash")
        err := app.EmptyTrash(&config)

        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
	}

	if *restoreAll {
		fmt.Println("Restoring all files")
	}

	if *restore != "" {
		fmt.Println(*restore)
	}

	if *view {
		fmt.Println("Viewing the trash")
		logEntries, err := app.GetLog(&config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, logEntry := range logEntries {
			fmt.Println(logEntry)
		}

		return
	}

	for _, arg := range args {
		file, err := os.Open(arg)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer file.Close()

		err = app.SetLog(file, &config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
