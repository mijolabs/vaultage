package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"
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

const bannerWidth = 56

func RootCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaultage",
		Short: "Vaultwarden backups with Age encryption",
	}

	originalHelp := cmd.HelpFunc()
	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		if !c.HasParent() {
			if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w >= bannerWidth {
				fmt.Println(banner())
			}
		}
		originalHelp(c, args)
	})

	cmd.AddCommand(Backup(ctx))
	cmd.AddCommand(Watch(ctx))
	cmd.AddCommand(Generate(ctx))

	return cmd
}
