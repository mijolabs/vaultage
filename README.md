# Vaultage
### Automated Vaultwarden backups with Age encryption

**Vaultage** is designed to run as a sidecar container alongside Vaultwarden. It monitors the Vaultwarden database for changes and creates a backup after a configurable debounce period. By default, the backup archive includes attachments and configuration files. Once created, the archive is securely encrypted using the Age encryption algorithm.

Status: WIP
