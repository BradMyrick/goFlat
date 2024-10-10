package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/BradMyrick/goFlat/obs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	rootCmd = &cobra.Command{
		Use:                "goFlat",
		PersistentPreRunE:  rootPersistentPreRunE,
		PersistentPostRunE: rootPersistentPostRunE,
		Short:              "Converts a folder structure to a flattened single file",
		Long: `goFlat is a tool to convert a folder structure to a flattened PDF or text file.
It walks the folder tree, adding the folder structure and file contents to the output file.
Binary and image files are automatically excluded to focus on textual content.`,
	}
	folderPath string
	verbose    bool
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("goFlat")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	rootCmd.PersistentFlags().StringVarP(&folderPath, "folder", "f", "", "Folder path to convert (required)")
	viper.BindPFlag("folder", rootCmd.PersistentFlags().Lookup("folder"))
	viper.BindEnv("folder")
	rootCmd.MarkPersistentFlagRequired("folder")

	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "Log verbosity level (debug, info, warn, error, fatal, panic)")
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindEnv("log.level")

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}

func Execute() error {
	ctx := context.Background()
	return rootCmd.ExecuteContext(ctx)
}

func rootPersistentPreRunE(cmd *cobra.Command, args []string) error {
	logger, err := obs.NewZap()
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger initialization error: %v\n", err)
		return err
	}

	settings := viper.AllSettings()
	config, _ := json.Marshal(settings)
	logger.Info("configuration", zap.String("settings", string(config)))

	if verbose {
		logger.Info("verbose logging enabled")
	}

	cmd.SetContext(obs.WithLogger(cmd.Context(), logger))
	return nil
}

func rootPersistentPostRunE(cmd *cobra.Command, args []string) error {
	err := obs.Logger(cmd.Context()).Sync()
	if _, ok := err.(*fs.PathError); ok {
		err = nil
	}
	return err
}
