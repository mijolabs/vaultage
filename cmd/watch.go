package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/mijolabs/vaultage/backup"
	"github.com/mijolabs/vaultage/watcher"
)

// Environment variable names
const (
	envDataDir            = "VAULTAGE_DATA_DIR"
	envOutputDir          = "VAULTAGE_OUTPUT_DIR"
	envDebounce           = "VAULTAGE_DEBOUNCE"
	envExcludeAttachments = "VAULTAGE_EXCLUDE_ATTACHMENTS"
	envExcludeConfigFile  = "VAULTAGE_EXCLUDE_CONFIG_FILE"
	envAgePassphrase      = "VAULTAGE_AGE_PASSPHRASE"
	envAgeKeyFile         = "VAULTAGE_AGE_KEY_FILE"
)

// Default values
const (
	defaultOutputDir = "."
	defaultDebounce  = 10 * time.Minute
)

// getEnv returns the value of the environment variable or the default.
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// getEnvBool returns the boolean value of the environment variable or the default.
func getEnvBool(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	// Treat "true", "1", "yes" as true (case-insensitive)
	val = strings.ToLower(val)
	return val == "true" || val == "1" || val == "yes"
}

// getEnvDuration returns the duration value of the environment variable or the default.
func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return defaultVal
	}
	return d
}

// Watch creates a Cobra command that monitors the Vaultwarden data directory
// for database changes and triggers encrypted backups using Age encryption.
func Watch(ctx context.Context) *cobra.Command {
	var (
		dataDir            string
		outputDir          string
		debounce           time.Duration
		excludeAttachments bool
		excludeConfigFile  bool
		agePassphrase      string
		ageKeyFile         string
	)

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch for changes and perform backups",
		Long: `Watch monitors the Vaultwarden data directory for database changes
and creates encrypted backups after a configurable debounce period.

All flags can also be set via environment variables with the VAULTAGE_ prefix.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate required data dir
			if dataDir == "" {
				return fmt.Errorf("data directory is required (use --data-dir or %s)", envDataDir)
			}
			dataDir = strings.TrimSuffix(dataDir, "/")

			// Validate data dir path
			info, err := os.Stat(dataDir)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("data directory does not exist: %s", dataDir)
				}
				return fmt.Errorf("checking data directory: %w", err)
			}
			if !info.IsDir() {
				return fmt.Errorf("path is not a directory: %s", dataDir)
			}

			// Validate mutually exclusive age options
			if agePassphrase != "" && ageKeyFile != "" {
				return fmt.Errorf("--age-passphrase and --age-key-file are mutually exclusive")
			}

			cfg := watcher.Config{
				Config: backup.Config{
					DataDir:            dataDir,
					OutputDir:          outputDir,
					ExcludeAttachments: excludeAttachments,
					ExcludeConfigFile:  excludeConfigFile,
					AgePassphrase:      agePassphrase,
					AgeKeyFile:         ageKeyFile,
				},
				Debounce: debounce,
			}

			return watcher.Watch(ctx, cfg)
		},
	}

	cmd.Flags().StringVar(
		&dataDir,
		"data-dir",
		os.Getenv(envDataDir),
		"path to Vaultwarden data directory",
	)
	cmd.Flags().StringVar(
		&outputDir,
		"output-dir",
		getEnv(envOutputDir, defaultOutputDir),
		"directory for backup files",
	)
	cmd.Flags().DurationVar(
		&debounce,
		"debounce",
		getEnvDuration(envDebounce, defaultDebounce),
		"quiet period before backup is performed",
	)
	cmd.Flags().BoolVar(
		&excludeAttachments,
		"exclude-attachments",
		getEnvBool(envExcludeAttachments, false),
		"exclude attachments from backup archive",
	)
	cmd.Flags().BoolVar(
		&excludeConfigFile,
		"exclude-config-file",
		getEnvBool(envExcludeConfigFile, false),
		"exclude config.json from backup archive",
	)
	cmd.Flags().StringVar(
		&agePassphrase,
		"age-passphrase",
		os.Getenv(envAgePassphrase),
		"passphrase for Age encryption",
	)
	cmd.Flags().StringVar(
		&ageKeyFile,
		"age-key-file",
		os.Getenv(envAgeKeyFile),
		"path to Age key file for encryption",
	)

	return cmd
}
