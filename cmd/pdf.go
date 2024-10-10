package cmd

import (
	"os/signal"
	"syscall"

	"github.com/BradMyrick/goFlat/flatten"
	"github.com/BradMyrick/goFlat/obs"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var pdfCmd = &cobra.Command{
	Use:   "pdf",
	Short: "Converts a folder structure to a flattened PDF file",
	Run:   runPdf,
}

func init() {
	rootCmd.AddCommand(pdfCmd)
}

func runPdf(cmd *cobra.Command, args []string) {
	ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGQUIT, syscall.SIGTERM)
	defer cancel()

	logger := obs.Logger(ctx).With(zap.String("command", "pdf flatten"))
	ctx = obs.WithLogger(ctx, logger)

	logger.Info("Flattening to PDF", zap.String("folder", folderPath))
	err := flatten.PdfFlatten(ctx, folderPath)
	if err != nil {
		logger.Error("PDF flattening failed", zap.Error(err))
	} else {
		logger.Info("PDF flattening complete", zap.String("folder", folderPath))
	}
}
