package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func banner() string {
	return `
____   ____            .__   __
\   \ /   /____   __ __|  |_/  |______     ____   ____
 \   Y   /\__  \ |  |  \  |\   __\__  \   / ___\_/ __ \
  \     /  / __ \|  |  /  |_|  |  / __ \_/ /_/  >  ___/
   \___/  (____  /____/|____/__| (____  /\___  / \___  >
               \/                     \//_____/      \/
`
}

func RootCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaultage",
		Short: "Vaultwarden backups with Age encryption",
	}

	originalHelp := cmd.HelpFunc()
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println(banner())
		originalHelp(cmd, args)
	})

	cmd.AddCommand(Backup(ctx))
	cmd.AddCommand(Watch(ctx))

	return cmd
}
