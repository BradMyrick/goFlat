package flatten

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/BradMyrick/goFlat/obs"
	"github.com/jung-kurt/gofpdf"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

const (
	pdfExtension = ".pdf"
	txtExtension = ".txt"
)

var ignoreExtensions = []string{".exe", ".pdf", ".bmp", ".gif", ".png", ".jpg", ".zip", ".tar"}

type Flattener interface {
	Flatten(ctx context.Context, folderPath string) error
}

type TextFlattener struct{}
type PDFFlattener struct{}

func (t *TextFlattener) Flatten(ctx context.Context, folderPath string) error {
	logger := obs.Logger(ctx)

	outputPath := filepath.Base(folderPath) + txtExtension
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	bar := progressbar.Default(-1, "Flattening to TXT")
	defer bar.Finish()

	err = processFiles(ctx, folderPath, func(path string, info os.FileInfo) error {
		bar.Add(1)

		_, err := fmt.Fprintln(writer, strings.TrimPrefix(path, folderPath))
		if err != nil {
			return fmt.Errorf("failed to write path: %w", err)
		}

		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			_, err = writer.Write(content)
			if err != nil {
				return fmt.Errorf("failed to write content: %w", err)
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

func (p *PDFFlattener) Flatten(ctx context.Context, folderPath string) error {
	logger := obs.Logger(ctx)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)

	bar := progressbar.Default(-1, "Flattening to PDF")
	defer bar.Finish()

	err := processFiles(ctx, folderPath, func(path string, info os.FileInfo) error {
		bar.Add(1)

		pdf.Ln(5)
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 0, strings.TrimPrefix(path, folderPath))

		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
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

	outputPath := filepath.Base(folderPath) + pdfExtension
	err = pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return fmt.Errorf("failed to save PDF: %w", err)
	}

	logger.Info("PDF saved", zap.String("outputPath", outputPath))
	return nil
}

func processFiles(_ context.Context, folderPath string, processor func(string, os.FileInfo) error) error {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, runtime.NumCPU())
	errChan := make(chan error, 1)

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == folderPath || strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() && path != folderPath {
				return filepath.SkipDir
			}
			return nil
		}

		for _, ext := range ignoreExtensions {
			if strings.HasSuffix(strings.ToLower(path), ext) {
				return nil
			}
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := processor(path, info); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
		}()

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	wg.Wait()
	close(errChan)

	if err, ok := <-errChan; ok {
		return fmt.Errorf("error during file processing: %w", err)
	}

	return nil
}

func TextFlatten(ctx context.Context, folderPath string) error {
	flattener := &TextFlattener{}
	return flattener.Flatten(ctx, folderPath)
}

func PdfFlatten(ctx context.Context, folderPath string) error {
	flattener := &PDFFlattener{}
	return flattener.Flatten(ctx, folderPath)
}
