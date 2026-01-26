// Package watcher provides file system monitoring for Vaultwarden database changes.
package watcher

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/mijolabs/vaultage/backup"
)

// WalFileName is the SQLite write-ahead log file that indicates database changes.
const WalFileName = "db.sqlite3-wal"

// Watch monitors the Vaultwarden data directory for changes to the WAL file
// and triggers backups after the debounce period.
// It blocks until the context is cancelled.
func Watch(ctx context.Context, cfg Config) error {
	walFilePath := cfg.DataDir + "/" + WalFileName

	log.Printf("watching %s (debounce: %s)", walFilePath, cfg.Debounce)
	log.Printf("exclude attachments: %t", cfg.ExcludeAttachments)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("creating fsnotify watcher: %w", err)
	}
	defer watcher.Close()

	if err := watcher.Add(cfg.DataDir); err != nil {
		return fmt.Errorf("adding watch on data directory: %w", err)
	}

	backupFn := func() error {
		return backup.Perform(cfg.Config)
	}

	return runLoop(ctx, watcher, cfg.Debounce, backupFn)
}

// logCooldown suppresses repeated log messages within this duration.
// SQLite WAL operations often trigger multiple fsnotify events in rapid
// succession (2-3 events within milliseconds). This cooldown prevents
// noisy logs while still resetting the debounce timer for each event.
const logCooldown = time.Second

// runLoop processes file system events and triggers backups after debounce.
func runLoop(ctx context.Context, watcher *fsnotify.Watcher, debounce time.Duration, backupFn func() error) error {
	var debounceTimer *time.Timer
	var lastLogTime time.Time

	for {
		select {
		case <-ctx.Done():
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			return ctx.Err()

		case err, ok := <-watcher.Errors:
			if !ok {
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				return nil
			}
			log.Printf("watcher error: %v", err)

		case event, ok := <-watcher.Events:
			if !ok {
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				return nil
			}

			if event.Name != WalFileName {
				continue
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
				continue
			}

			if debounceTimer != nil {
				debounceTimer.Stop()
			}

			debounceTimer = time.AfterFunc(debounce, func() {
				if err := backupFn(); err != nil {
					log.Printf("backup error: %v", err)
				}
			})

			if time.Since(lastLogTime) >= logCooldown {
				log.Printf("detected change in WAL file: %s (%s) - backup scheduled in %s", event.Name, event.Op.String(), debounce)
				lastLogTime = time.Now()
			}
		}
	}
}
