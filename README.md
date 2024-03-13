# goFlat

goFlat is a Go program that takes a folder location as input and generates a flattened PDF file with the folder structure and code contents. 

## Features

- Command-line interface using Cobra
- Reads folder path from command-line flag `--folder` or `-f`
- Walks the specified folder and subfolders 
- Adds folder structure to PDF with bold font
- Includes code contents of each file in monospace font below folder
- Ignores hidden files and folders starting with `.`
- Ignores `.git` folder
- Generates PDF with name of root folder 

## Installation

```
go get github.com/BradMyrick/goFlat
```

## Usage

```
goFlat --folder /path/to/folder
```

Or 

```
goFlat -f /path/to/folder
```

This will generate a PDF file named `folder.pdf` in the current directory with the flattened contents.

## Example

For a folder structure like:

```
root/
    cmd/
        cmd.go
    rpc/
        rpc.go
```

The generated PDF would look like:

```
root
    cmd
        cmd.go 
            package main
            import "fmt"
            func main() {
                fmt.Println("cmd.go")
            }
    rpc
        rpc.go
            package rpc 
            func DoRPC() {
            }
```

## Implementation Details

The key implementation aspects are:

- Using Cobra for CLI 
- Using Viper to read folder flag
- Walking folder tree with `filepath.Walk`
- Adding folder paths with `pdf.SetFont("Arial", "B", 14)` 
- Adding file contents with `pdf.SetFont("Courier", "", 10)`
- Generating PDF with gofpdf library
- Ignoring hidden files/folders by checking `strings.HasPrefix(info.Name(), ".")`

created by kodr.eth
```