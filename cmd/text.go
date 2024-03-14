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
	txtCmd = &cobra.Command{
		Use:   "txt",
		Short: "Converts to a flattened TXT file",
		Run:   runTxt,
	}
)

func init() {
	rootCmd.AddCommand(txtCmd)
	txtCmd.Flags().StringVarP(&folderPath, "folder", "f", "", "Folder path to convert (required)")
	txtCmd.MarkFlagRequired("folder")
}

func runTxt(cmd *cobra.Command, args []string) {
	// Initialize the context with a cancel function.
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()
	// Setup shutdown signal cancellation.
	ctx, cancel = signal.NotifyContext(ctx, syscall.SIGQUIT, syscall.SIGTERM)
	defer cancel()

	// Get the logger and update it.
	logger := obs.Logger(ctx).With(zap.String("command", "txt flatten"))
	ctx = obs.WithLogger(ctx, logger)

	folderPath, _ := cmd.Flags().GetString("folder")

	logger.Info("Flattening TXT", zap.String("folder", folderPath))
	err := flatten.TextFlatten(ctx, folderPath)
	if err != nil {
		logger.Error("Flattening TXT failed", zap.Error(err))
	} else {
		logger.Info("Flattening TXT complete", zap.String("folder", folderPath))
	}

} // func runTxt()
