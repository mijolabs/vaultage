package cmd

import (
	"context"
	"fmt"
	"os"

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

// Returns the value of the environment variable or the default.
func envOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func RootCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaultage",
		Short: "Vaultwarden backups with Age encryption",
	}

	cmd.AddCommand(Backup(ctx))
	cmd.AddCommand(Watch(ctx))

	originalHelp := cmd.HelpFunc()
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println(banner())
		originalHelp(cmd, args)
	})

	return cmd
}
