package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

var boolMap = map[string]bool{
	"1":     true,
	"true":  true,
	"yes":   true,
	"0":     false,
	"false": false,
	"no":    false,
}

// Returns the value of the environment variable as a bool, or the default.
func envBoolOrDefault(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	if b, ok := boolMap[strings.ToLower(val)]; ok {
		return b
	}
	return defaultVal
}

// Returns the value of the environment variable as a string, or the default.
func envStringOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// Returns the value of the environment variable as a duration, or the default.
func envDurationOrDefault(key string, defaultVal time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	d, err := time.ParseDuration(val)
	if err != nil {
		return defaultVal
	}
	return d
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
