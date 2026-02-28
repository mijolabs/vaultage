package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mijolabs/vaultage/backup"
)

// Creates a Cobra command that performs a backup of Vaultwarden data,
// archives it, and optionally encrypts it
func Backup(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Short: "One-shot backup and archive of Vaultwarden data",
		Use:   "backup [data dir]",
		RunE: func(cmd *cobra.Command, args []string) error {
			dataDir, err := validateDataDirFromArgs(args)
			if err != nil {
				return fmt.Errorf("validating data directory: %w", err)
			}

			cfg := resolveBackupFlags(cmd)
			cfg.DataDir = dataDir

			// Validate mutually exclusive age options
			if !cfg.WithoutEncryption {
				if cfg.AgePassphrase != "" && cfg.AgeKeyFile != "" {
					return fmt.Errorf("--age-passphrase and --age-key-file are mutually exclusive")
				}
			}

			return backup.Perform(ctx, cfg)
		},
	}

	addBackupFlags(cmd)

	return cmd
}
