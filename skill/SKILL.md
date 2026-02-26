---
name: cupi-cli
description: "CUPI CLI: query and manage Cisco Unity Connection (CUC) voicemail systems via CUPI REST API, PAWS platform admin, AST health monitoring, and DIME log collection. Supports mailbox users, distribution lists, call handlers, class of service, templates, schedules, and OS-level cluster management."
argument-hint: "[command or resource]"
allowed-tools: Bash, Read, Grep, Glob
user-invocable: true
---

## Setup

The `cupi-cli` binary must be built and available in PATH:

```bash
bash install-local.sh
# or build manually:
go build -o cupi-cli .
```

## Authentication

```bash
# Login — tests connectivity and saves credentials
cupi-cli auth login --host cuc.example.com --username admin --server prod --default

# Add additional credential types to an existing server
cupi-cli auth set-credentials --type application --username app-user --server prod
cupi-cli auth set-credentials --type platform    --username os-admin  --server prod

# Show credential status
cupi-cli auth status [--server prod]

# Switch the default server
cupi-cli auth switch prod

# List all configured servers
cupi-cli auth list

# Logout — removes credentials from secure storage
cupi-cli auth logout [--server prod] [--type cupi|application|platform|all]
```

Passwords are stored in the **OS keystore** (macOS Keychain, Windows Credential Manager, Linux Secret Service). Server config (hostnames, usernames, default server) is stored in `~/.cupi-cli/config.json` — no passwords on disk.

| Type | Used For |
|------|---------|
| `cupi` | CUPI REST provisioning (default) |
| `application` | Application-level APIs |
| `platform` | OS admin / platform-level access (PAWS) |

## Command Discovery

Run `cupi-cli --help` to list all available commands. Run `cupi-cli <command> --help` for exact flags.

## Global Flags

| Flag | Description |
|------|-------------|
| `--server <name>` | Override the default server |
| `--output json\|table\|csv\|raw` | Output format (default: table) |
| `--debug` | Print raw API request/response to stderr |
| `--max <n>` | Limit results to N items (0 = no limit) |
| `--dry-run` | Print what would be sent without making changes |

## CUPI Resources

```bash
# Users (mailboxes)
cupi-cli users list
cupi-cli users list --max 50
cupi-cli users list --query "(alias startswith j)"
cupi-cli users get jsmith
cupi-cli users get jsmith --output json
cupi-cli users add --alias jsmith --dtmf 1001 --first-name John --last-name Smith
cupi-cli users update jsmith --display-name "John Smith" --department Engineering
cupi-cli users remove jsmith

# Distribution Lists
cupi-cli distlists list
cupi-cli distlists get allvoicemail
cupi-cli distlists add --alias mylist --display-name "My List"
cupi-cli distlists update mylist --display-name "Updated List"
cupi-cli distlists remove mylist
cupi-cli distlists members list mylist
cupi-cli distlists members add mylist <member-objectId>
cupi-cli distlists members remove mylist <member-objectId>

# Call Handlers
cupi-cli handlers list
cupi-cli handlers get "Opening Greeting"
cupi-cli handlers add --display-name "My Handler" --template-id <objectId>
cupi-cli handlers update "My Handler" --dtmf 9999
cupi-cli handlers remove "My Handler"

# Class of Service
cupi-cli cos list
cupi-cli cos get UnityMailboxTemplate
cupi-cli cos update UnityMailboxTemplate --display-name "Updated COS"

# User Templates
cupi-cli templates list
cupi-cli templates get voicemailusertemplate

# Schedules
cupi-cli schedules list
cupi-cli schedules get "Weekdays"

# System Info
cupi-cli system
cupi-cli system --output json
```

## AST — System Health Monitoring

```bash
cupi-cli ast disk                     # Disk partition usage
cupi-cli ast heartbeat                # Heartbeat rates
cupi-cli ast tftp                     # TFTP server statistics
cupi-cli ast alerts                   # All system alerts
cupi-cli ast alerts --triggered       # Only currently triggered alerts
cupi-cli ast perfmon                  # Perfmon object catalog
cupi-cli ast perfmon --output json | jq '.[0]'
```

## PAWS — Platform Administration

Requires platform credentials (`cupi-cli auth set-credentials --type platform --username admin`).

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

## Multi-server

```bash
cupi-cli users list --server lab
cupi-cli auth switch prod
```

## Phase 1 Scraper

```bash
go run ./tools/scraper/ --json api_reference.json --md api_reference.md
```
