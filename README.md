# Vaultage

**Note: This project is a WIP and still incomplete.**

**Vaultage** is a command-line utility that monitors your [Vaultwarden](https://github.com/dani-garcia/vaultwarden) database for changes and creates a backup after a configurable debounce period.

By default, the backup archive includes attachments and configuration files. Once created, the archive is securely encrypted using the [Age](https://github.com/FiloSottile/age) encryption algorithm.

It is recommended to run as a sidecar container alongside Vaultwarden.

## Installation

### Docker (Recommended)

```bash
docker pull ghcr.io/mijolabs/vaultage:latest
```

Or from Docker Hub:

```bash
docker pull mijolabs/vaultage:latest
```

### From Source

```bash
go install github.com/mijolabs/vaultage@latest
```

## Usage

### Docker Compose

Create a `docker-compose.yml` file:

```yaml
services:
  vaultage:
    image: ghcr.io/mijolabs/vaultage:latest
    environment:
      VAULTAGE_DATA_DIR: /data
      VAULTAGE_OUTPUT_DIR: /backups
      VAULTAGE_DEBOUNCE: 10m
    volumes:
      - vaultwarden-data:/data:ro
      - ./backups:/backups
    restart: unless-stopped

volumes:
  vaultwarden-data:
    external: true
```

Then run:

```bash
docker compose up -d
```

### Command Line

```bash
vaultage watch --data-dir /path/to/vaultwarden/data --output-dir /path/to/backups
```

## Configuration

All configuration options can be set via command-line flags or environment variables.

| Flag | Environment Variable | Type | Default | Description |
|------|---------------------|------|---------|-------------|
| `--data-dir` | `VAULTAGE_DATA_DIR` | string | *required* | Path to Vaultwarden data directory |
| `--output-dir` | `VAULTAGE_OUTPUT_DIR` | string | `.` | Directory for backup files |
| `--debounce` | `VAULTAGE_DEBOUNCE` | duration | `10m` | Quiet period before backup is performed |
| `--exclude-attachments` | `VAULTAGE_EXCLUDE_ATTACHMENTS` | bool | `false` | Exclude attachments from backup archive |
| `--exclude-config-file` | `VAULTAGE_EXCLUDE_CONFIG_FILE` | bool | `false` | Exclude config.json from backup archive |
| `--age-passphrase` | `VAULTAGE_AGE_PASSPHRASE` | string | - | Passphrase for Age encryption |
| `--age-key-file` | `VAULTAGE_AGE_KEY_FILE` | string | - | Path to Age key file for encryption |

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

## License

MIT
