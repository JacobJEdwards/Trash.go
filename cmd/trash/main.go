package main

import (
	"flag"
	"fmt"
)

var (
	empty      = flag.Bool("empty", false, "empty the trash")
	restoreAll = flag.Bool("restore-all", false, "restore all files")
	restore    = flag.String("restore", "", "restore a file")
)

func main() {
	flag.Parse()
	args := flag.Args()

	if *empty {
		fmt.Println("Emptying the trash")
	}

	if *restoreAll {
		fmt.Println("Restoring all files")
	}

	if *restore != "" {
		fmt.Println(*restore)
	}

	for _, arg := range args {
		fmt.Println(arg)
	}
}
