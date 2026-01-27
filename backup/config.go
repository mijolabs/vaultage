package backup

// Config holds the configuration needed for performing backups.
type Config struct {
	DataDir            string
	OutputDir          string
	ExcludeAttachments bool
	ExcludeConfigFile  bool
	WithoutEncryption  bool
	AgePassphrase      string
	AgeKeyFile         string
}
