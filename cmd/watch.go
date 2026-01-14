package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const WalFileName = "db.sqlite3-wal"

func Watch(cfg *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch [data dir]",
		Short: "Watch for changes and perform backups",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Require data dir path
			if len(args) < 1 {
				return fmt.Errorf("missing path to vaultwarden data directory\n")
			}
			dataDir := strings.TrimSuffix(args[0], "/")
			cfg.Set("data_dir", dataDir)

			// Validate data dir path
			if info, err := os.Stat(dataDir); err != nil {
				if errors.Is(err, os.ErrNotExist) || !info.IsDir() {
					return fmt.Errorf("invalid data dir `%s`\n", dataDir)
				}
				return fmt.Errorf("%s", err.Error())
			}

			// Validate required age options
			age_passphrase := cfg.GetString("age_passphrase")
			age_key_file := cfg.GetString("age_key_file")
			if age_passphrase == "" && age_key_file == "" {
				return fmt.Errorf(
					"either `age-passphrase` or `age-key-file` must be provided\n",
				)
			}
			// Validate mutually exclusive age options
			if age_passphrase != "" && age_key_file != "" {
				return fmt.Errorf(
					"`age-passphrase` and `age-key-file` are mutually exclusive\n",
				)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := watchDataDir(cfg)
			if err != nil {
				cobra.CheckErr(err)
			}
			return nil
		},
	}

	cmd.Flags().Duration(
		"debounce",
		10*time.Minute,
		"trailing quiet period before backup is performed",
	)
	cmd.Flags().Bool(
		"exclude-attachments",
		false,
		"exclude attachments in backup archive",
	)
	cmd.Flags().Bool(
		"exclude-config-file",
		false,
		"exclude config.json in backup archive",
	)
	cmd.Flags().String(
		"age-passphrase",
		"",
		"age passphrase for backup encryption",
	)
	cmd.Flags().String(
		"age-key-file",
		"",
		"age key file for backup encryption",
	)

	// Replace dashes in flags with underscores to match env vars
	cmd.Flags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ReplaceAll(name, "-", "_"))
	})

	err := cfg.BindPFlags(cmd.Flags())
	if err != nil {
		cobra.CheckErr(err)
	}

	// fmt.Println(cfg.AllSettings())

	return cmd
}

func watchDataDir(cfg *viper.Viper) error {
	walFilePath := cfg.GetString("data_dir") + "/" + WalFileName
	log.Printf("👀 watching %s (debounce: %s)\n", walFilePath, cfg.GetString("debounce"))
	log.Printf("📦 exclude attachments: %t\n", cfg.GetBool("exclude_attachments"))

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("cannot create fsnotify watcher: %v", err)
	}
	defer watcher.Close()

	go watcherLoop(watcher, cfg)

	err = watcher.Add(cfg.GetString("data_dir"))
	if err != nil {
		log.Fatalf("cannot watch data directory: %v", err)
	}

	select {}
}

func watcherLoop(watcher *fsnotify.Watcher, cfg *viper.Viper) {
	walFilePath := cfg.GetString("data_dir") + "/" + WalFileName

	var debounceTimer *time.Timer
	for {
		select {
		case err, ok := <-watcher.Errors:
			if !ok { // Channel closed
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				return
			}
			log.Printf("Watcher error: %v", err)
		case event, ok := <-watcher.Events:
			if !ok { // Channel closed
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				return
			}
			if event.Name != walFilePath {
				continue
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
				continue
			}

			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			debounceTimer = time.AfterFunc(cfg.GetDuration("debounce"), func() {
				if err := performBackup(cfg); err != nil {
					log.Printf("Backup error: %v", err)
					return
				}
			})
			log.Printf("Detected change in WAL file: %s (%s) — backup scheduled in %s", event.Name, event.Op.String(), cfg.GetDuration("debounce"))
		}
	}
}
