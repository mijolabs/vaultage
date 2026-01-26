package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func RootCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaultage",
		Short: "Vaultwarden backups with Age encryption.",
	}

	cmd.AddCommand(Watch(ctx))

	return cmd
}
