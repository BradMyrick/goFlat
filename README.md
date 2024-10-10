# goFlat

goFlat is a powerful and versatile command-line tool written in Go, designed to flatten the structure of a given directory into a single document. It supports generating both PDF and text files, focusing on code and readable content while automatically excluding binary and image files.

## Features

- **Flexible Output**: Generate flattened representations of directories in PDF or text format.
- **Intelligent Content Selection**: Automatically excludes binary and image files, focusing on textual content.
- **Efficient Processing**: Utilizes concurrent file processing for improved performance.
- **Command-Line Interface**: Easy to integrate into scripts or workflows for automation.
- **Customizable**: Supports various flags for customization, including verbose logging and specifying the directory to flatten.
- **Progress Indication**: Provides visual feedback during the flattening process.

## Installation

Ensure you have Go 1.21 or later installed on your system, then run:

```bash
go install github.com/BradMyrick/goFlat@latest
```

## Usage

To flatten a directory into a PDF:

```bash
goFlat pdf --folder /path/to/folder
```

To generate a text file instead:

```bash
goFlat txt --folder /path/to/folder
```

Additional options:

- Use `--verbose` or `-v` for detailed logging.
- Use `--log-level` or `-l` to set the log verbosity level (debug, info, warn, error, fatal, panic).

Both commands will produce an output file named after the root folder with the appropriate extension in the current directory.

## Contributing

Contributions are welcome! If you have suggestions for improvements or encounter any issues, please file them on the [project's issues page](https://github.com/BradMyrick/goFlat/issues).

Before contributing:
1. Ensure your code adheres to the Go community coding standards.
2. Include appropriate tests for your changes.
3. For significant changes, please open an issue first to discuss what you would like to change.

## License

goFlat is released under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Roadmap

- [x] Add gitignore file
- [x] Implement concurrent file processing
- [x] Add progress indication
- [ ] Add support for custom output file names and paths
- [ ] Implement support for additional file types (e.g., Markdown, HTML)
- [ ] Add support for customizing the output format (e.g., font size, color)
- [ ] Implement support for flattening nested directories with customizable depth
- [ ] Add option to include/exclude specific file types
- [ ] Implement a web interface for online file flattening

## Acknowledgements

Special thanks to all contributors and the community for their invaluable feedback and support.