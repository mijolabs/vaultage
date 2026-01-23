package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mijolabs/vaultage/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	rootCmd := cmd.RootCmd(ctx)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
