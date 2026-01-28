package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Returns the value of the environment variable or the default.
func envOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// Validates the data directory path from command-line arguments.
func validateDataDirFromArgs(args []string) (string, error) {
	// Require data dir path provided in args
	if len(args) < 1 {
		return "", fmt.Errorf("missing path to vaultwarden data directory")
	}
	dataDir := strings.TrimSuffix(args[0], "/")

	// Validate data dir path
	info, err := os.Stat(dataDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("data directory does not exist: %s", dataDir)
		}
		return "", fmt.Errorf("checking data directory: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("path is not a directory: %s", dataDir)
	}

	return dataDir, nil
}
