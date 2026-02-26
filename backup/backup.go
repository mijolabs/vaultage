// Package backup provides encrypted backup functionality for Vaultwarden data.
package backup

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	// The name of the Vaultwarden SQLite database file.
	dbFileName = "db.sqlite3"
	// The name of the attachments directory.
	attachmentsDirName = "attachments"
	// The name of the Vaultwarden config file.
	configFileName = "config.json"
)

func Perform(ctx context.Context, cfg Config) error {
	// Ensure output directory exists
	outputDir := cfg.OutputDir
	if outputDir == "" {
		outputDir = "."
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Gather in-memory db bytes and any on-disk files
	archiveEntries, err := getArchiveEntries(cfg)
	if err != nil {
		return err
	}

	// Generate output filename
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("vaultage-%s.tar", timestamp)
	outFilePath := filepath.Join(outputDir, filename)

	// Create in-memory tar archive
	archiveBuf := &bytes.Buffer{}
	if err := CreateArchive(archiveBuf, archiveEntries); err != nil {
		return fmt.Errorf("creating archive: %w", err)
	}
	archiveBytes := archiveBuf.Bytes()

	if cfg.WithoutEncryption {
		log.Printf("writing unencrypted backup: %s (%s)", outFilePath, formatSize(int64(len(archiveBytes))))
		return writeBytesToDisk(outFilePath, archiveBytes)
	}

	outEncryptedFilePath := outFilePath + ".age"
	log.Printf("creating encrypted backup: %s", outEncryptedFilePath)
	passphrase := cfg.AgePassphrase
	if passphrase == "" {
		passphrase, err = promptForPassphrase()
		if err != nil {
			return err
		}
	}

	encryptedArchiveBytes, err := encryptWithPassphrase(archiveBytes, passphrase)
	if err != nil {
		return err
	}

	log.Printf("writing encrypted backup: %s (%s)", outEncryptedFilePath, formatSize(int64(len(archiveBytes))))
	return writeBytesToDisk(outEncryptedFilePath, encryptedArchiveBytes)
}

func getArchiveEntries(cfg Config) ([]ArchiveEntry, error) {
	log.Printf("enumerating archive entries...")

	// Backup SQLite database to memory
	dbPath := filepath.Join(cfg.DataDir, dbFileName)
	dbData, err := BackupToMemory(dbPath)
	if err != nil {
		return nil, fmt.Errorf("backing up database: %w", err)
	}

	// Collect files to archive
	archiveEntries := []ArchiveEntry{
		{
			Name: dbFileName,
			Data: dbData,
			Mode: 0644,
		},
	}

	// Add attachments directory if it exists and not excluded
	if !cfg.ExcludeAttachments {
		attachmentsPath := filepath.Join(cfg.DataDir, attachmentsDirName)
		if info, err := os.Stat(attachmentsPath); err == nil && info.IsDir() {
			archiveEntries = append(
				archiveEntries,
				ArchiveEntry{
					Name: attachmentsDirName,
					Path: attachmentsPath,
				},
			)
		}
	}

	// Add config.json if it exists and not excluded
	if !cfg.ExcludeConfigFile {
		configPath := filepath.Join(cfg.DataDir, configFileName)
		if _, err := os.Stat(configPath); err == nil {
			archiveEntries = append(
				archiveEntries,
				ArchiveEntry{
					Name: configFileName,
					Path: configPath,
				},
			)
		}
	}

	return archiveEntries, nil
}

func writeBytesToDisk(outputFilePath string, data []byte) error {
	if err := os.WriteFile(outputFilePath, data, 0644); err != nil {
		// Clean up partial file on error
		os.Remove(outputFilePath)
		return fmt.Errorf("writing backup file: %w", err)
	}

	info, err := os.Stat(outputFilePath)
	if err != nil {
		return fmt.Errorf("getting file info: %w", err)
	}

	log.Printf("write successful: %s (%s)", outputFilePath, formatSize(info.Size()))

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
