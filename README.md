```
____   ____            .__   __
\   \ /   /____   __ __|  |_/  |______     ____   ____
 \   Y   /\__  \ |  |  \  |\   __\__  \   / ___\_/ __ \
  \     /  / __ \|  |  /  |_|  |  / __ \_/ /_/  >  ___/
   \___/  (____  /____/|____/__| (____  /\___  / \___  >
               \/                     \//_____/      \/
```
#

**Note: This project is a WIP and still incomplete. Don't use it yet.**

**Vaultage** is a command-line utility that monitors your [Vaultwarden](https://github.com/dani-garcia/vaultwarden) database for writes and applies a trailing debounce strategy to create a backup after a configurable period. The backup along with any attachments and the instance configuration will be archived with tar, then optionally encrypted using the [age](https://github.com/FiloSottile/age) encryption algorithm.


## Installation

### Docker (Recommended)

```bash
docker pull ghcr.io/mijolabs/vaultage:latest
```

### From Source

```bash
go install github.com/mijolabs/vaultage@latest
```

## Usage

### Command Line

To perform a one-time backup:

```bash
vaultage backup /path/to/vaultwarden/data --output-dir /path/to/backups
```

Or to start the watcher:
```bash
vaultage watch /path/to/vaultwarden/data --output-dir /path/to/backups
```

### Docker Compose

Add it as a second service in your Vaultwarden `docker-compose.yml` file:

```yaml
services:
  vaultwarden:
    ...

  vaultage:
    image: ghcr.io/mijolabs/vaultage:latest
    restart: unless-stopped
    environment:
      VAULTAGE_DATA_DIR: "/data"
      VAULTAGE_OUTPUT_DIR: "/backups"
      VAULTAGE_DEBOUNCE: "10s"
      VAULTAGE_EXCLUDE_ATTACHMENTS: "false"
      VAULTAGE_EXCLUDE_CONFIG_FILE: "false"
      VAULTAGE_WITHOUT_ENCRYPTION: "false"
      VAULTAGE_AGE_PASSPHRASE: "${VAULTAGE_AGE_PASSPHRASE}"
      # VAULTAGE_AGE_KEY_FILE: "/keys/age.key"                          # Or use key file instead. Path inside container
    volumes:
      - "${VAULTAGE_DATA_DIR:-./data}:/data:ro"                         # Mount Vaultwarden data directory
      - "${VAULTAGE_OUTPUT_DIR:-./data}:/backups"                       # Mount backup output directory
      # - "${VAULTAGE_AGE_KEY_FILE:-./keys/age.key}:/keys/age.key:ro"   # Mount Age key file if using key-based encryption
    command: ["watch", "/data"]
```

Then run:

```bash
docker compose up -d
```

## Configuration

All configuration options can be set via command-line flags or environment variables.

| Flag                    | Environment Variable           | Type     | Default    | Description                             |
| ----------------------- | ------------------------------ | -------- | ---------- | --------------------------------------- |
| `--data-dir`            | `VAULTAGE_DATA_DIR`            | string   | *required* | Path to Vaultwarden data directory      |
| `--output-dir`          | `VAULTAGE_OUTPUT_DIR`          | string   | `.`        | Directory for backup files              |
| `--debounce`            | `VAULTAGE_DEBOUNCE`            | duration | `10m`      | Quiet period before backup is performed |
| `--exclude-attachments` | `VAULTAGE_EXCLUDE_ATTACHMENTS` | bool     | `false`    | Exclude attachments from backup archive |
| `--exclude-config-file` | `VAULTAGE_EXCLUDE_CONFIG_FILE` | bool     | `false`    | Exclude config.json from backup archive |
| `--age-passphrase`      | `VAULTAGE_AGE_PASSPHRASE`      | string   | -          | Passphrase for Age encryption           |
| `--age-key-file`        | `VAULTAGE_AGE_KEY_FILE`        | string   | -          | Path to Age key file for encryption     |

### Duration Format

The `--debounce` flag accepts Go duration strings, e.g.:
- `10m` - 10 minutes
- `1h` - 1 hour
- `30s` - 30 seconds
- `1h30m` - 1 hour and 30 minutes

### Boolean Environment Variables

Boolean environment variables accept `true`, `1`, or `yes` (case-insensitive) as truthy values.

## How It Works

1. Vaultage monitors the Vaultwarden WAL file (`db.sqlite3-wal`) for changes
2. When a change is detected, a debounce timer starts
3. After the debounce period with no new changes, a backup is created
4. The backup uses SQLite's Online Backup API to safely copy the database
5. The backup archive includes the database, config file, and attachments (unless excluded)
6. If configured, the archive is encrypted using Age encryption
