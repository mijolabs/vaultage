// Package backup provides encrypted backup functionality for Vaultwarden data.
package backup

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	// DBFileName is the name of the Vaultwarden SQLite database file.
	DBFileName = "db.sqlite3"
	// ConfigFileName is the name of the Vaultwarden config file.
	ConfigFileName = "config.json"
	// AttachmentsDirName is the name of the attachments directory.
	AttachmentsDirName = "attachments"
)

// Perform creates a backup of the Vaultwarden data directory.
// It safely backs up the SQLite database to memory, then creates a tar archive
// containing the database and optionally the config file and attachments.
func Perform(cfg Config) error {
	log.Printf("performing backup...")

	// Ensure output directory exists
	outputDir := cfg.OutputDir
	if outputDir == "" {
		outputDir = "."
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Backup SQLite database to memory
	dbPath := filepath.Join(cfg.DataDir, DBFileName)
	dbData, err := BackupToMemory(dbPath)
	if err != nil {
		return fmt.Errorf("backing up database: %w", err)
	}

	// Collect files to archive
	entries := []ArchiveEntry{
		{
			Name: DBFileName,
			Data: dbData,
			Mode: 0644,
		},
	}

	// Add config.json if it exists and not excluded
	if !cfg.ExcludeConfigFile {
		configPath := filepath.Join(cfg.DataDir, ConfigFileName)
		if _, err := os.Stat(configPath); err == nil {
			entries = append(entries, ArchiveEntry{
				Name: ConfigFileName,
				Path: configPath,
			})
		}
	}

	// Add attachments directory if it exists and not excluded
	if !cfg.ExcludeAttachments {
		attachmentsPath := filepath.Join(cfg.DataDir, AttachmentsDirName)
		if info, err := os.Stat(attachmentsPath); err == nil && info.IsDir() {
			entries = append(entries, ArchiveEntry{
				Name: AttachmentsDirName,
				Path: attachmentsPath,
			})
		}
	}

	// Generate output filename
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("vaultage-%s.tar", timestamp)
	outputPath := filepath.Join(outputDir, filename)

	// Create the tar archive
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer file.Close()

	if err := CreateArchive(file, entries); err != nil {
		// Clean up partial file on error
		os.Remove(outputPath)
		return fmt.Errorf("creating archive: %w", err)
	}

	// Get file size for logging
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("getting file info: %w", err)
	}

	log.Printf("backup complete: %s (%s)", outputPath, formatSize(info.Size()))

	return nil
}

// formatSize returns a human-readable file size string.
func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
