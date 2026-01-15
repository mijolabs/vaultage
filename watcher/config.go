package watcher

import (
	"time"

	"github.com/mijolabs/vaultage/backup"
)

// Config holds the configuration needed for the watcher.
type Config struct {
	backup.Config
	Debounce time.Duration
}
