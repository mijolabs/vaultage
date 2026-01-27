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

	"github.com/mijolabs/vaultage/crypto"
)

const (
	// DBFileName is the name of the Vaultwarden SQLite database file.
	DBFileName = "db.sqlite3"
	// AttachmentsDirName is the name of the attachments directory.
	AttachmentsDirName = "attachments"
	// ConfigFileName is the name of the Vaultwarden config file.
	ConfigFileName = "config.json"
)

func GetArchiveEntries(cfg Config) ([]ArchiveEntry, error) {
	log.Printf("enumerating archive entries...")

	// Backup SQLite database to memory
	dbPath := filepath.Join(cfg.DataDir, DBFileName)
	dbData, err := BackupToMemory(dbPath)
	if err != nil {
		return nil, fmt.Errorf("backing up database: %w", err)
	}

	// Collect files to archive
	archiveEntries := []ArchiveEntry{
		{
			Name: DBFileName,
			Data: dbData,
			Mode: 0644,
		},
	}

	// Add attachments directory if it exists and not excluded
	if !cfg.ExcludeAttachments {
		attachmentsPath := filepath.Join(cfg.DataDir, AttachmentsDirName)
		if info, err := os.Stat(attachmentsPath); err == nil && info.IsDir() {
			archiveEntries = append(
				archiveEntries,
				ArchiveEntry{
					Name: AttachmentsDirName,
					Path: attachmentsPath,
				},
			)
		}
	}

	// Add config.json if it exists and not excluded
	if !cfg.ExcludeConfigFile {
		configPath := filepath.Join(cfg.DataDir, ConfigFileName)
		if _, err := os.Stat(configPath); err == nil {
			archiveEntries = append(
				archiveEntries,
				ArchiveEntry{
					Name: ConfigFileName,
					Path: configPath,
				},
			)
		}
	}

	return archiveEntries, nil
}

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
	archiveEntries, err := GetArchiveEntries(cfg)
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
		log.Printf("writing unencrypted backup: %s", outFilePath)
		return writeArchiveToDisk(outFilePath, archiveBytes)
	}

	outEncryptedFilePath := outFilePath + ".age"
	log.Printf("creating encrypted backup: %s", outEncryptedFilePath)
	passphrase := cfg.AgePassphrase
	if passphrase == "" {
		passphrase, err = crypto.PromptForPassphrase()
		if err != nil {
			return err
		}
	}

	encryptedArchiveBytes, err := crypto.EncryptWithPassphrase(archiveBytes, passphrase)
	if err != nil {
		return err
	}

	log.Printf("writing encrypted backup: %s", outEncryptedFilePath)
	writeArchiveToDisk(outEncryptedFilePath, encryptedArchiveBytes)

	return nil
}

func writeArchiveToDisk(outputFilePath string, archiveBytes []byte) error {
	if err := os.WriteFile(outputFilePath, archiveBytes, 0644); err != nil {
		// Clean up partial file on error
		os.Remove(outputFilePath)
		return fmt.Errorf("writing backup file: %w", err)
	}
	return nil
}

// // Create the tar archive in memory
// archiveData, err := CreateArchiveInMemory(entries)
// if err != nil {
// 	return nil, fmt.Errorf("creating archive: %w", err)
// }

// log.Printf("backup created successfully, size: %s", formatSize(int64(len(archiveData))))
// return archiveData, nil

// // Generate output filename
// timestamp := time.Now().Format("20060102_150405")
// filename := fmt.Sprintf("vaultage-%s.tar", timestamp)
// outputPath := filepath.Join(outputDir, filename)

// // Create the tar archive
// file, err := os.Create(outputPath)
// if err != nil {
// 	return fmt.Errorf("creating output file: %w", err)
// }
// defer file.Close()

// if err := CreateArchive(file, entries); err != nil {
// 	// Clean up partial file on error
// 	os.Remove(outputPath)
// 	return fmt.Errorf("creating archive: %w", err)
// }

// // Get file size for logging
// info, err := file.Stat()
// if err != nil {
// 	return fmt.Errorf("getting file info: %w", err)
// }

// log.Printf("backup complete: %s (%s)", outputPath, formatSize(info.Size()))

// return nil
// }

// formatSize returns a human-readable file size string.
// func formatSize(bytes int64) string {
// 	const (
// 		KB = 1024
// 		MB = KB * 1024
// 		GB = MB * 1024
// 	)

// 	switch {
// 	case bytes >= GB:
// 		return fmt.Sprintf("%.1f GB", float64(bytes)/GB)
// 	case bytes >= MB:
// 		return fmt.Sprintf("%.1f MB", float64(bytes)/MB)
// 	case bytes >= KB:
// 		return fmt.Sprintf("%.1f KB", float64(bytes)/KB)
// 	default:
// 		return fmt.Sprintf("%d B", bytes)
// 	}
// }
