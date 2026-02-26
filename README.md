# cupi-cli

A command-line tool for querying and managing Cisco Unity Connection (CUC) voicemail servers via CUPI REST, PAWS, AST, and DIME APIs.

## Install

**From source**

```bash
git clone https://github.com/Cloverhound/cupi-cli.git
cd cupi-cli
bash install-local.sh
```

Or build manually:

```bash
go build -o cupi-cli .
```

## Quick Start

```bash
# Authenticate — saves credentials to OS keystore
cupi-cli auth login --host cuc.example.com --username admin --server prod --default

# List mailbox users
cupi-cli users list

# Get a specific user
cupi-cli users get jsmith

# Add a user
cupi-cli users add --alias jsmith --dtmf 1001 --first-name John --last-name Smith

# Use a different server for one command
cupi-cli users list --server lab
```

## API Coverage

| Resource Group | Commands |
|---|---|
| **auth** | login, set-credentials, status, logout, list, switch |
| **users** (mailboxes) | list, get, add, update, remove |
| **distlists** (distribution lists) | list, get, add, update, remove, members list/add/remove |
| **handlers** (call handlers) | list, get, add, update, remove |
| **cos** (class of service) | list, get, update |
| **templates** (user templates) | list, get |
| **schedules** | list, get |
| **system** | system info |
| **ast** | disk, heartbeat, tftp, alerts, perfmon |
| **paws** | cluster status/replication, drs backup/status |
| **dime** | get-file (log download) |

## Authentication

`cupi-cli auth login` tests connectivity and saves credentials:

1. GET `/vmrest/users?rowsPerPage=0` with Basic Auth
2. Validates credentials against the CUC server
3. Saves server config to `~/.cupi-cli/config.json`
4. Stores password in the OS keystore

Three credential types are supported per server:

| Type | Used For |
|------|---------|
| `cupi` | CUPI REST provisioning (default) |
| `application` | Application-level APIs |
| `platform` | OS admin / platform-level access (PAWS) |

```bash
# Login (CUPI credentials)
cupi-cli auth login --host cuc.example.com --username admin --server prod --default

# Add additional credential types
cupi-cli auth set-credentials --type application --username app-user --server prod
cupi-cli auth set-credentials --type platform    --username os-admin  --server prod

# Show credential status
cupi-cli auth status [--server prod]

# Output:
# server=prod (default)  host=cuc.example.com  version=15.0
#   cupi       : admin     [set]
#   application: app-user  [set]
#   platform   : os-admin  [set]

# Logout one type or all
cupi-cli auth logout --server prod --type application
cupi-cli auth logout --server prod --type all
```

## Output Formats

| Format | Flag | Notes |
|--------|------|-------|
| Table | `--output table` | Default — auto-width ASCII table |
| JSON | `--output json` | Pretty-printed; pipe to `jq` for filtering |
| CSV | `--output csv` | Pipe to file for spreadsheet import |
| Raw | `--output raw` | Raw API response |

```bash
cupi-cli users list --output json | jq '.[].Alias'
cupi-cli users list --output csv > users.csv
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--server <name>` | Override the default server for this command |
| `--output json\|table\|csv\|raw` | Output format (default: table) |
| `--debug` | Print raw API request/response to stderr |
| `--max <n>` | Limit results to N items, 0 = no limit (default) |
| `--dry-run` | Print what would be sent without making any changes |

## Dry Run

Preview write operations without applying them:

```bash
cupi-cli --dry-run users add --alias jsmith --dtmf 1001 --first-name John --last-name Smith
cupi-cli --dry-run users remove jsmith
cupi-cli --dry-run distlists add --alias mylist --display-name "My List"
```

## Configuration

Server configuration is stored in `~/.cupi-cli/config.json`:

```json
{
  "defaultServer": "prod",
  "servers": {
    "prod": {
      "host": "cuc.example.com",
      "port": 443,
      "version": "15.0",
      "credentials": {
        "cupi":        { "username": "admin" },
        "application": { "username": "app-user" },
        "platform":    { "username": "os-admin" }
      }
    },
    "lab": {
      "host": "cuc-lab.example.com",
      "port": 443,
      "version": "12.5",
      "credentials": {
        "cupi": { "username": "admin" }
      }
    }
  }
}
```

Passwords are stored in the **OS keystore** (macOS Keychain, Windows Credential Manager, Linux Secret Service) — never in the config file.

## AST — System Health Monitoring

```bash
cupi-cli ast disk                     # Disk partition usage
cupi-cli ast heartbeat                # Heartbeat rates
cupi-cli ast tftp                     # TFTP server statistics
cupi-cli ast alerts                   # All system alerts
cupi-cli ast alerts --triggered       # Only currently triggered alerts
cupi-cli ast perfmon                  # Perfmon object catalog
```

## PAWS — Platform Administration

Requires platform credentials (`cupi-cli auth set-credentials --type platform`).

```bash
cupi-cli paws cluster status          # OS-level cluster node info
cupi-cli paws cluster replication     # Replication health check
cupi-cli paws drs status              # DRS backup/restore status
cupi-cli paws drs backup --sftp-server 10.0.0.5 --sftp-user backup --sftp-password secret --sftp-dir /backups
```

## DIME — Log File Downloads

```bash
cupi-cli dime get-file /var/log/active/syslog/CiscoSyslog > syslog.txt
cupi-cli dime get-file /var/log/active/tomcat/catalina.out --output /tmp/catalina.out
cupi-cli dime get-file syslog/CiscoSyslog --node 10.0.0.5
```

## Claude Code Integration

See [`skill/SKILL.md`](skill/SKILL.md) for the Claude Code skill definition.

## Development

See [`CLAUDE.md`](CLAUDE.md) for the full development guide including project layout, build instructions, and key file reference.

```bash
# Build
go build -o cupi-cli .

# Run integration tests (requires live CUC server)
CUPI_TEST_HOST=cuc.example.com CUPI_TEST_USER=admin CUPI_TEST_PASS=secret \
  go test ./tests/ -v -timeout 120s
```

## CUC Version Support

- **12.5** — Supported
- **14.x** — Supported
- **15.x** — Supported

## License

MIT — see [LICENSE](LICENSE)
