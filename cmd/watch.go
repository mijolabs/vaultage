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

// envOrDefault returns the value of the environment variable or the default.
func envOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// Watch creates a Cobra command that monitors the Vaultwarden data directory
// for database changes and triggers encrypted backups using Age encryption.
func Watch(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch [data dir]",
		Short: "Watch for changes and perform backups",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Require data dir path
			if len(args) < 1 {
				return fmt.Errorf("missing path to vaultwarden data directory")
			}
			dataDir := strings.TrimSuffix(args[0], "/")

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

			// Get flag values
			debounce, _ := cmd.Flags().GetDuration("debounce")
			outputDir, _ := cmd.Flags().GetString("output-dir")
			excludeAttachments, _ := cmd.Flags().GetBool("exclude-attachments")
			excludeConfigFile, _ := cmd.Flags().GetBool("exclude-config-file")
			agePassphrase, _ := cmd.Flags().GetString("age-passphrase")
			ageKeyFile, _ := cmd.Flags().GetString("age-key-file")

			// Validate required age options (temporarily disabled until age encryption is implemented)
			// if agePassphrase == "" && ageKeyFile == "" {
			// 	return fmt.Errorf("either --age-passphrase or --age-key-file must be provided")
			// }
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

	cmd.Flags().Duration(
		"debounce",
		10*time.Minute,
		"trailing quiet period before backup is performed",
	)
	cmd.Flags().String(
		"output-dir",
		envOrDefault("OUTPUT_DIR", "."),
		"directory for backup files",
	)
	cmd.Flags().Bool(
		"exclude-attachments",
		false,
		"exclude attachments in backup archive",
	)
	cmd.Flags().Bool(
		"exclude-config-file",
		false,
		"exclude config.json in backup archive",
	)
	cmd.Flags().String(
		"age-passphrase",
		os.Getenv("AGE_PASSPHRASE"),
		"age passphrase for backup encryption",
	)
	cmd.Flags().String(
		"age-key-file",
		os.Getenv("AGE_KEY_FILE"),
		"age key file for backup encryption",
	)

	return cmd
}
