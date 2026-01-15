package backup

// Config holds the configuration needed for performing backups.
type Config struct {
	DataDir            string
	ExcludeAttachments bool
	ExcludeConfigFile  bool
	AgePassphrase      string
	AgeKeyFile         string
}
