/*
Copyright © 2026 Michael Johansson <mijolabs@remotenode.io>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RootCmd(cfg *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaultage",
		Short: "Vaultwarden backups with Age encryption.",
		// 		Long: `
		// A longer description that spans multiple lines and likely contains
		// examples and usage of using your application. For example:

		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}

	cmd.AddCommand(Watch(cfg))

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// func Execute() {
// 	if err := RootCmd.Execute(); err != nil {
// 		os.Exit(1)
// 	}
// }

// func init() {
// 	cobra.OnInitialize(initEnvVars)

// 	// rootCmd.PersistentFlags().BoolVar(&includeAttachments, "include-attachments", true, "include attachments in backup archive")
// 	// rootCmd.PersistentFlags().BoolVar(&includeConfig, "include-config", true, "include config.json in backup archive")
// }

// func initEnvVars() {
// 	viper.AutomaticEnv()
// 	// err := viper.Bind.BindPFlags(rootCmd.Flags())
// 	// if err != nil {
// 	// 	return err
// 	// }
// }
