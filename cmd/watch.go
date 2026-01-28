package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/mijolabs/vaultage/backup"
	"github.com/mijolabs/vaultage/watcher"
)

// Creates a Cobra command that monitors the Vaultwarden data directory
// for database changes and triggers encrypted backups using Age encryption.
func Watch(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Watch for changes and perform backups",
		Use:   "watch [data dir]",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate data directory
			dataDir, err := validateDataDirFromArgs(args)
			if err != nil {
				return fmt.Errorf("validating data directory: %w", err)
			}

			// Get flag values
			debounce, _ := cmd.Flags().GetDuration("debounce")
			outputDir, _ := cmd.Flags().GetString("output-dir")
			excludeAttachments, _ := cmd.Flags().GetBool("exclude-attachments")
			excludeConfigFile, _ := cmd.Flags().GetBool("exclude-config-file")
			withoutEncryption, _ := cmd.Flags().GetBool("without-encryption")
			agePassphrase, _ := cmd.Flags().GetString("age-passphrase")
			ageKeyFile, _ := cmd.Flags().GetString("age-key-file")

			// Validate age options
			if !withoutEncryption {
				// Both empty OR both non-empty → invalid
				if (agePassphrase == "") == (ageKeyFile == "") {
					return fmt.Errorf(
						"watch mode requires exactly one of --age-passphrase or --age-key-file " +
							"to be set via cli flags or env vars",
					)
				}
			}

			cfg := watcher.Config{
				Config: backup.Config{
					DataDir:            dataDir,
					OutputDir:          outputDir,
					ExcludeAttachments: excludeAttachments,
					ExcludeConfigFile:  excludeConfigFile,
					WithoutEncryption:  withoutEncryption,
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
		envDurationOrDefault("VAULTAGE_DEBOUNCE", 10*time.Minute),
		"trailing quiet period before backup is performed",
	)
	cmd.Flags().String(
		"output-dir",
		envStringOrDefault("VAULTAGE_OUTPUT_DIR", "."),
		"directory for backup files",
	)
	cmd.Flags().Bool(
		"exclude-attachments",
		envBoolOrDefault("VAULTAGE_EXCLUDE_ATTACHMENTS", false),
		"exclude attachments in backup archive",
	)
	cmd.Flags().Bool(
		"exclude-config-file",
		envBoolOrDefault("VAULTAGE_EXCLUDE_CONFIG_FILE", false),
		"exclude config.json in backup archive",
	)
	cmd.Flags().Bool(
		"without-encryption",
		envBoolOrDefault("VAULTAGE_WITHOUT_ENCRYPTION", false),
		"disable encryption for backups",
	)
	cmd.Flags().String(
		"age-passphrase",
		os.Getenv("VAULTAGE_AGE_PASSPHRASE"),
		"age passphrase for backup encryption",
	)
	cmd.Flags().String(
		"age-key-file",
		os.Getenv("VAULTAGE_AGE_KEY_FILE"),
		"age key file for backup encryption",
	)

	return cmd
}
