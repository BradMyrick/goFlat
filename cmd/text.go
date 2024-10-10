package cmd

import (
	"os/signal"
	"syscall"

	"github.com/BradMyrick/goFlat/flatten"
	"github.com/BradMyrick/goFlat/obs"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var txtCmd = &cobra.Command{
	Use:   "txt",
	Short: "Converts a folder structure to a flattened text file",
	Run:   runTxt,
}

func init() {
	rootCmd.AddCommand(txtCmd)
}

func runTxt(cmd *cobra.Command, args []string) {
	ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGQUIT, syscall.SIGTERM)
	defer cancel()

	logger := obs.Logger(ctx).With(zap.String("command", "txt flatten"))
	ctx = obs.WithLogger(ctx, logger)

	logger.Info("Flattening to text", zap.String("folder", folderPath))
	err := flatten.TextFlatten(ctx, folderPath)
	if err != nil {
		logger.Error("Text flattening failed", zap.Error(err))
	} else {
		logger.Info("Text flattening complete", zap.String("folder", folderPath))
	}
}
