package flatten

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/BradMyrick/goFlat/obs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestTextFlattener(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "goflat_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	createTestFiles(t, tempDir)

	logger, _ := zap.NewDevelopment()
	ctx := obs.WithLogger(context.Background(), logger)

	flattener := &TextFlattener{}
	err = flattener.Flatten(ctx, tempDir)
	assert.NoError(t, err)

	outputPath := filepath.Base(tempDir) + txtExtension
	_, err = os.Stat(outputPath)
	assert.NoError(t, err)

	content, err := os.ReadFile(outputPath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "test1.txt")
	assert.Contains(t, string(content), "This is test file 1")
	assert.Contains(t, string(content), "test2.txt")
	assert.Contains(t, string(content), "This is test file 2")
	assert.Contains(t, string(content), "subdir/test3.txt")
	assert.Contains(t, string(content), "This is test file 3")
	assert.NotContains(t, string(content), "ignored.exe")
	assert.NotContains(t, string(content), "image.jpg")
}

func TestPDFFlattener(t *testing.T) {
	// Setup
	tempDir, err := os.MkdirTemp("", "goflat_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test files and directories
	createTestFiles(t, tempDir)

	// Create context with logger
	logger, _ := zap.NewDevelopment()
	ctx := obs.WithLogger(context.Background(), logger)

	// Run PDFFlattener
	flattener := &PDFFlattener{}
	err = flattener.Flatten(ctx, tempDir)
	assert.NoError(t, err)

	// Check if output file exists
	outputPath := filepath.Base(tempDir) + pdfExtension
	_, err = os.Stat(outputPath)
	assert.NoError(t, err)

	// TODO: Add more specific checks for PDF content if needed
}

func TestProcessFiles(t *testing.T) {
	// Setup
	tempDir, err := os.MkdirTemp("", "goflat_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test files and directories
	createTestFiles(t, tempDir)

	// Create context with logger
	logger, _ := zap.NewDevelopment()
	ctx := obs.WithLogger(context.Background(), logger)

	// Test processFiles
	processedFiles := make(map[string]bool)
	err = processFiles(ctx, tempDir, func(path string, info os.FileInfo) error {
		processedFiles[path] = true
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, processedFiles[filepath.Join(tempDir, "test1.txt")])
	assert.True(t, processedFiles[filepath.Join(tempDir, "test2.txt")])
	assert.True(t, processedFiles[filepath.Join(tempDir, "subdir", "test3.txt")])
	assert.False(t, processedFiles[filepath.Join(tempDir, "ignored.exe")])
	assert.False(t, processedFiles[filepath.Join(tempDir, "image.jpg")])
}

func TestIgnoreExtensions(t *testing.T) {
	// Setup
	tempDir, err := os.MkdirTemp("", "goflat_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test files with ignored extensions
	for _, ext := range ignoreExtensions {
		filename := "test" + ext
		err := os.WriteFile(filepath.Join(tempDir, filename), []byte("ignored content"), 0644)
		require.NoError(t, err)
	}

	// Create context with logger
	logger, _ := zap.NewDevelopment()
	ctx := obs.WithLogger(context.Background(), logger)

	// Run TextFlattener
	flattener := &TextFlattener{}
	err = flattener.Flatten(ctx, tempDir)
	assert.NoError(t, err)

	// Check content of output file
	outputPath := filepath.Base(tempDir) + txtExtension
	content, err := os.ReadFile(outputPath)
	assert.NoError(t, err)
	for _, ext := range ignoreExtensions {
		assert.NotContains(t, string(content), "test"+ext)
	}
}

func TestErrorHandling(t *testing.T) {
	// Test with non-existent directory
	ctx := context.Background()
	flattener := &TextFlattener{}
	err := flattener.Flatten(ctx, "/non/existent/path")
	assert.Error(t, err)

	// Test with unreadable file
	tempDir, err := os.MkdirTemp("", "goflat_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	unreadableFile := filepath.Join(tempDir, "unreadable.txt")
	err = os.WriteFile(unreadableFile, []byte("test content"), 0000)
	require.NoError(t, err)

	err = flattener.Flatten(ctx, tempDir)
	assert.Error(t, err)
}

func createTestFiles(t *testing.T, dir string) {
	files := map[string]string{
		"test1.txt":        "This is test file 1",
		"test2.txt":        "This is test file 2",
		"subdir/test3.txt": "This is test file 3",
		"ignored.exe":      "This should be ignored",
		"image.jpg":        "This should be ignored",
	}

	for path, content := range files {
		fullPath := filepath.Join(dir, path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		require.NoError(t, err)
		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)
	}
}
