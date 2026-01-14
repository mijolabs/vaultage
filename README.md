# Vaultage
**Vaultage** is a command-line utility that will monitor your [Vaultwarden](https://github.com/dani-garcia/vaultwarden) database for changes and creates a backup after a configurable debounce period.

By default, the backup archive includes attachments and configuration files. Once created, the archive is securely encrypted using the [Age](https://github.com/FiloSottile/age) encryption algorithm.

It is recommended to run as a sidecar container alongside Vaultwarden.

Status: WIP
