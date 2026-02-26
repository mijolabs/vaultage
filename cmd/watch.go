package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

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

			cfg := resolveBackupFlags(cmd)
			cfg.DataDir = dataDir

			// Resolve watch-specific debounce flag
			debounce, _ := cmd.Flags().GetDuration("debounce")
			if !cmd.Flags().Changed("debounce") {
				debounce = envDurationOrDefault("VAULTAGE_DEBOUNCE", debounce)
			}

			// Validate age options
			if !cfg.WithoutEncryption {
				if (cfg.AgePassphrase == "") == (cfg.AgeKeyFile == "") {
					return fmt.Errorf(
						"watch mode requires exactly one of --age-passphrase or --age-key-file " +
							"to be set via cli flags or env vars",
					)
				}
			}

			watchCfg := watcher.Config{
				Config:   cfg,
				Debounce: debounce,
			}

			return watcher.Watch(ctx, watchCfg)
		},
	}

	addBackupFlags(cmd)
	cmd.Flags().Duration(
		"debounce",
		10*time.Minute,
		"trailing quiet period before backup is performed (env: VAULTAGE_DEBOUNCE)",
	)

	return cmd
}
