# goFlat

goFlat is a versatile command-line tool written in Go, designed to flatten the structure of a given directory into a single document. It supports generating both PDF and text files, excluding binary and image files to focus on code and readable content.

## Features

- **Flexible Output**: Generates flattened representations of directories in PDF or text format.
- **Selective Flattening**: Automatically excludes binary and image files, focusing on textual content.
- **Command-Line Interface**: Easy to integrate into scripts or workflows for automation.
- **Customizable**: Supports various flags for customization, including specifying the directory to flatten.

## Installation

Ensure you have Go installed on your system, then run:

```bash
go install github.com/BradMyrick/goFlat
```

## Usage

To flatten a directory into a PDF, excluding binary and image files:

```bash
goFlat pdf --folder /path/to/folder
```

To generate a text file instead:

```bash
goFlat txt --folder /path/to/folder
```

Both commands will produce an output file named after the root folder with the appropriate extension in the current directory.

## Contributing

Contributions are welcome! If you have suggestions for improvements or encounter any issues, please file them on the project's issues page.

Before contributing, please ensure your code adheres to the Go community coding standards and includes appropriate tests. For significant changes, please open an issue first to discuss what you would like to change.

## License

goFlat is released under the MIT License. See the LICENSE file for more details.

## TODO's

- [x] Add gitignore file. *complete*
- [ ] Add support for custom output file names and paths.
- [ ] Implement support for additional file types, such as Markdown and HTML.
- [ ] Add support for customizing the output format, such as font size and color.
- [ ] Implement support for flattening nested directories.

```
