package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mijolabs/vaultage/backup"
)

// Creates a Cobra command that performs a backup of Vaultwarden data,
// archives it, and optionally encrypts it
func Backup(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Short: "One-time backup of data directory",
		Use:   "backup [data dir]",
		RunE: func(cmd *cobra.Command, args []string) error {
			dataDir, err := validateDataDirFromArgs(args)
			if err != nil {
				return fmt.Errorf("validating data directory: %w", err)
			}

			// Get flag values
			outputDir, _ := cmd.Flags().GetString("output-dir")
			excludeAttachments, _ := cmd.Flags().GetBool("exclude-attachments")
			excludeConfigFile, _ := cmd.Flags().GetBool("exclude-config-file")
			withoutEncryption, _ := cmd.Flags().GetBool("without-encryption")
			agePassphrase, _ := cmd.Flags().GetString("age-passphrase")
			ageKeyFile, _ := cmd.Flags().GetString("age-key-file")

			// Validate mutually exclusive age options
			if !withoutEncryption {
				if agePassphrase != "" && ageKeyFile != "" {
					return fmt.Errorf("--age-passphrase and --age-key-file are mutually exclusive")
				}
			}

			cfg := backup.Config{
				DataDir:            dataDir,
				OutputDir:          outputDir,
				ExcludeAttachments: excludeAttachments,
				ExcludeConfigFile:  excludeConfigFile,
				WithoutEncryption:  withoutEncryption,
				AgePassphrase:      agePassphrase,
				AgeKeyFile:         ageKeyFile,
			}

			return backup.Perform(ctx, cfg)
		},
	}

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
