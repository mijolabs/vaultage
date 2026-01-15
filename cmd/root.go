package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RootCmd(cfg *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaultage",
		Short: "Vaultwarden backups with Age encryption.",
	}

	cmd.AddCommand(Watch(cfg))

	return cmd
}
