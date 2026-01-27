package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mijolabs/vaultage/backup"
)

// Creates a Cobra command that performs a backup of Vaultwarden data,
// archives it, and optionally encrypts it
func Backup(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup [data dir]",
		Short: "One-time backup of data directory",
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
			outputDir, _ := cmd.Flags().GetString("output-dir")
			excludeAttachments, _ := cmd.Flags().GetBool("exclude-attachments")
			excludeConfigFile, _ := cmd.Flags().GetBool("exclude-config-file")
			withoutEncryption, _ := cmd.Flags().GetBool("without-encryption")
			agePassphrase, _ := cmd.Flags().GetString("age-passphrase")
			ageKeyFile, _ := cmd.Flags().GetString("age-key-file")

			// Validate required age options (temporarily disabled until age encryption is implemented)
			// if agePassphrase == "" && ageKeyFile == "" {
			// 	return fmt.Errorf("either --age-passphrase or --age-key-file must be provided")
			// }
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
	cmd.Flags().Bool(
		"without-encryption",
		false,
		"disable encryption for backups",
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
