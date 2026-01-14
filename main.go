package main

import (
	"os"

	"github.com/mijolabs/vaultage/cmd"
	"github.com/spf13/viper"
)

func main() {
	cfg := viper.New()
	cfg.AutomaticEnv()

	rootCmd := cmd.RootCmd(cfg)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
