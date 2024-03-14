package flatten

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BradMyrick/goFlat/obs"
	"github.com/jung-kurt/gofpdf"
	"go.uber.org/zap"
)

func TextFlatten(ctx context.Context, folderPath string) (err error) {
	// Get the logger
	logger := obs.Logger(ctx)

	// Create output text file
	outputPath := filepath.Base(folderPath) + ".txt"
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Define a slice with the extensions to ignore
	ignoreExtensions := []string{".exe", ".pdf", ".bmp", ".gif", ".png", ".jpg", ".zip", ".tar"}

	// Walk the folder tree
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
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

		// Check if the file has an extension that should be ignored
		for _, ext := range ignoreExtensions {
			if strings.HasSuffix(strings.ToLower(path), ext) {
				return nil // Skip this file
			}
		}

		// Add folder structure
		_, err = fmt.Fprintln(file, strings.TrimPrefix(path, folderPath))
		if err != nil {
			return err
		}

		// Add file contents
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = file.Write(content)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	logger.Info("TXT saved", zap.String("outputPath", outputPath))
	return nil
}

func PdfFlatten(ctx context.Context, folderPath string) (err error) {
	// Get the logger
	logger := obs.Logger(ctx)

	// Create new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)

	// Define a slice with the extensions to ignore
	ignoreExtensions := []string{".exe", ".pdf", ".bmp", ".gif", ".png", ".jpg", ".zip", ".tar"}

	// Walk the folder tree
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
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

		// Check if the file has an extension that should be ignored
		for _, ext := range ignoreExtensions {
			if strings.HasSuffix(strings.ToLower(path), ext) {
				return nil // Skip this file
			}
		}

		// Add folder structure
		pdf.Ln(5)
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 0, strings.TrimPrefix(path, folderPath))

		// Add file contents for files
		if !info.IsDir() {
			content, err := os.ReadFile(path)
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
		return err
	}

	// Save PDF file
	outputPath := filepath.Base(folderPath) + ".pdf"
	err = pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return err
	}

	logger.Info("PDF saved", zap.String("outputPath", outputPath))
	return nil
}
