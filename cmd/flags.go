package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/mijolabs/vaultage/backup"
)

// Registers shared backup settings flags on a command.
func addBackupFlags(cmd *cobra.Command) {
	cmd.Flags().String("output-dir", ".", "directory for backup files (env: VAULTAGE_OUTPUT_DIR)")
	cmd.Flags().Bool("exclude-attachments", false, "exclude attachments in backup archive (env: VAULTAGE_EXCLUDE_ATTACHMENTS)")
	cmd.Flags().Bool("exclude-config-file", false, "exclude config.json in backup archive (env: VAULTAGE_EXCLUDE_CONFIG_FILE)")
	cmd.Flags().Bool("without-encryption", false, "disable encryption for backups (env: VAULTAGE_WITHOUT_ENCRYPTION)")
	cmd.Flags().String("age-passphrase", "", "age passphrase for backup encryption (env: VAULTAGE_AGE_PASSPHRASE)")
	cmd.Flags().String("age-key-file", "", "age key file for backup encryption (env: VAULTAGE_AGE_KEY_FILE)")
}

// Reads the shared backup flags, applying env var fallbacks when a flag
// was not explicitly set on the command line.
func resolveBackupFlags(cmd *cobra.Command) backup.Config {
	outputDir, _ := cmd.Flags().GetString("output-dir")
	if !cmd.Flags().Changed("output-dir") {
		outputDir = envStringOrDefault("VAULTAGE_OUTPUT_DIR", outputDir)
	}

	excludeAttachments, _ := cmd.Flags().GetBool("exclude-attachments")
	if !cmd.Flags().Changed("exclude-attachments") {
		excludeAttachments = envBoolOrDefault("VAULTAGE_EXCLUDE_ATTACHMENTS", excludeAttachments)
	}

	excludeConfigFile, _ := cmd.Flags().GetBool("exclude-config-file")
	if !cmd.Flags().Changed("exclude-config-file") {
		excludeConfigFile = envBoolOrDefault("VAULTAGE_EXCLUDE_CONFIG_FILE", excludeConfigFile)
	}

	withoutEncryption, _ := cmd.Flags().GetBool("without-encryption")
	if !cmd.Flags().Changed("without-encryption") {
		withoutEncryption = envBoolOrDefault("VAULTAGE_WITHOUT_ENCRYPTION", withoutEncryption)
	}

	agePassphrase, _ := cmd.Flags().GetString("age-passphrase")
	if agePassphrase == "" {
		agePassphrase = os.Getenv("VAULTAGE_AGE_PASSPHRASE")
	}

	ageKeyFile, _ := cmd.Flags().GetString("age-key-file")
	if ageKeyFile == "" {
		ageKeyFile = os.Getenv("VAULTAGE_AGE_KEY_FILE")
	}

	return backup.Config{
		OutputDir:          outputDir,
		ExcludeAttachments: excludeAttachments,
		ExcludeConfigFile:  excludeConfigFile,
		WithoutEncryption:  withoutEncryption,
		AgePassphrase:      agePassphrase,
		AgeKeyFile:         ageKeyFile,
	}
}
