// Package backup provides encrypted backup functionality for Vaultwarden data.
package backup

import (
	"log"
)

// Perform creates an encrypted backup of the Vaultwarden data directory.
// It uses Age encryption with either a passphrase or key file as configured.
func Perform(cfg Config) error {
	log.Printf("performing backup...")

	// TODO: Implement actual backup logic:
	// 1. Create tar archive of data directory
	// 2. Optionally exclude attachments based on cfg.ExcludeAttachments
	// 3. Optionally exclude config.json based on cfg.ExcludeConfigFile
	// 4. Encrypt with Age using cfg.AgePassphrase or cfg.AgeKeyFile
	// 5. Write to backup destination

	return nil
}
