package cmd

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/BradMyrick/goFlat/flatten"
	"github.com/BradMyrick/goFlat/obs"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// menu to flatten folder into a text file
var (
	pdfCmd = &cobra.Command{
		Use:   "pdf",
		Short: "Converts to a flattened PDF file",
		Run:  runPdf,
	}
)

func init() {
	rootCmd.AddCommand(pdfCmd)
}

func runPdf(cmd *cobra.Command, args []string) {
	// Initialize the context with a cancel function.
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()
	// Setup shutdown signal cancellation.
	ctx, cancel = signal.NotifyContext(ctx, syscall.SIGQUIT, syscall.SIGTERM)
	defer cancel()

	// Get the logger and update it.
	logger := obs.Logger(ctx).With(zap.String("command", "pdf flatten"))
	ctx = obs.WithLogger(ctx, logger)

	folderPath, _ := cmd.Flags().GetString("folder")

	logger.Info("Flattening PDF", zap.String("folder", folderPath))
	err := flatten.PdfFlatten(ctx, folderPath)
	if err != nil {
		logger.Error("Flattening PDF failed", zap.Error(err))
	} else{
		logger.Info("Flattening PDF complete", zap.String("folder", folderPath))
	}

} // func runTxt()


