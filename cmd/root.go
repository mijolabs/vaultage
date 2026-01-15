package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// RootCmd creates the root Cobra command for the vaultage CLI.
// It accepts a context for cancellation.
func RootCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaultage",
		Short: "Vaultwarden backups with Age encryption.",
	}

	cmd.AddCommand(Watch(ctx))

	return cmd
}
