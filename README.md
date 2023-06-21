# Trash.go

The Trash CLI is a command-line interface for managing a trash system. It allows you to trash, restore, and delete files, as well as view and empty the trash.

## Installation

To install the Trash CLI, you can use the following command:

```shell
$ go get -u github.com/JacobJEdwards/Trash.go/cmd/trash
```

Make sure you have Go installed and configured properly.

## Usage

The Trash CLI supports the following commands:

### View the Trash

To view the contents of the trash, use the `-view` or `-v` flag:

```shell
$ trash -view
```

This will display a list of trashed files along with their original names, paths, and trashed timestamps.

### Empty the Trash

To empty the trash and permanently delete all trashed files, use the `-empty` or `-e` flag:

```shell
$ trash -empty
```

This will remove all files from the trash.

### Restore All Files

To restore all trashed files, use the `-restore-all` or `-ra` flag:

```shell
$ trash -restore-all
```

This will restore all files from the trash to their original locations.

### Restore a File

To restore a specific trashed file, use the `-restore` or `-r` flag followed by the filename:

```shell
$ trash -restore <filename>
```

This will restore the specified file from the trash to its original location.

### Delete a File

To permanently delete a file and move it to the trash, use the `-delete` or `-rm` flag followed by the filename:

```shell
$ trash -delete <filename>
```

This will move the specified file to the trash.

### Help

To display the help message with all available options, use the `-help` or `-h` flag:

```shell
$ trash -help
```

## Configuration

The Trash CLI uses a configuration file to specify the location of the trash and log files. By default, it looks for a `config.yaml` file in the current directory. You can modify the configuration as needed.

## License

This project is licensed under the [MIT License](LICENSE).
```

