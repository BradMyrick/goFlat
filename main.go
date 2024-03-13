package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/spf13/cobra"
)

var folderPath string

func main() {
	cmd := &cobra.Command{
		Use:   "goFlat",
		Short: "Converts a folder structure to a flattened PDF",
		Long: "goFlat is a tool to convert a folder structure to a flattened PDF. " +
			"It walks the folder tree and adds the folder structure and file contents to the PDF.",
		Run:   run,
	}

	cmd.Flags().StringVarP(&folderPath, "folder", "f", "", "Folder path to convert (required)")
	cmd.MarkFlagRequired("folder")

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	// Create new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)

	// Walk the folder tree
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip root folder
		if path == folderPath {
			return nil
		}

		// Skip hidden files and folders
		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Add folder structure 
		pdf.Ln(5)
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 0, strings.TrimPrefix(path, folderPath))

		// Add file contents for files
		if !info.IsDir() {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			pdf.Ln(5)
			pdf.SetFont("Courier", "", 10)
			pdf.MultiCell(0, 5, string(content), "", "L", false)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	// Save PDF file
	outputPath := filepath.Base(folderPath) + ".pdf"
	err = pdf.OutputFileAndClose(outputPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("PDF generated at %s\n", outputPath)
}
