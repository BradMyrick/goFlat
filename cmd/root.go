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
		Use:   "goFlat",
		PersistentPreRunE:  rootPersistentPreRunE,
		PersistentPostRunE: rootPersistentPostRunE,
		Short: "Converts a folder structure to a flattened single file",
		Long: "goFlat is a tool to convert a folder structure to a flattened PDF " +
			"or txt file by walking the folder tree and adding the folder structure " +
			"and file contents to the file",
			}
		folderPath string


)

func init() {

	// Use environment variables that match our flags.
	viper.AutomaticEnv()

	viper.SetEnvPrefix("goFlat")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	
	rootCmd.PersistentFlags().StringP(
		"folder",
		"f",
		"",
		"Folder path to convert (required)",
	)

	viper.BindPFlag("folder", rootCmd.PersistentFlags().Lookup("folder"))
	viper.BindEnv("folder")

	rootCmd.MarkPersistentFlagRequired("folder")

	rootCmd.PersistentFlags().StringP(
		"log-level",
		"l",
		"info",
		"Log verbosity level (debug, info, warn, error, fatal, panic).",
	)
	viper.BindPFlag(
		"log.level",
		rootCmd.PersistentFlags().Lookup("log-level"),
	)
	viper.BindEnv("log.level")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() (err error) {	
	// Initialize root context.
	ctx := context.Background()
	

	// Execute the command.
	_, err = rootCmd.ExecuteContextC(ctx)
	return err
	
}

func rootPersistentPreRunE(command *cobra.Command, args []string) (err error) {



	// Configure logging.
	var logger *zap.Logger
	if logger, err = obs.NewZap(); err != nil {
		fmt.Fprintf(os.Stderr, "logger initialization error: %v\n", err)
		return
	}

	// Dump the configuration
	settings := viper.AllSettings()
	// Remove sensitive information.
	// ... e.g. database.username and database.password
	config, _ := json.Marshal(settings)
	logger.Info("configuration", zap.String("settings", string(config)))

	// Update the command context.
	command.SetContext(obs.WithLogger(command.Context(), logger))
	return

} // func cmd.rootPersistentPreRunE()

func rootPersistentPostRunE(command *cobra.Command, args []string) (err error) {

	// Flush the log.
	if err = obs.Logger(command.Context()).Sync(); err != nil {
		// invalid argument on linux for stdout/stderr.
		if _, ok := err.(*fs.PathError); ok {
			err = nil
		}
	}
	return

} // func cmd.rootPersistentPostRunE()
