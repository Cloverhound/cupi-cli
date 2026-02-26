---
name: cupi-cli
description: "CUPI CLI: query and manage Cisco Unity Connection (CUC) voicemail systems via CUPI REST API, PAWS platform admin, AST health monitoring, and DIME log collection. Supports mailbox users, distribution lists, call handlers, class of service, templates, schedules, and OS-level cluster management."
argument-hint: "[command or resource]"
allowed-tools: Bash, Read, Grep, Glob
user-invocable: true
---

## Setup

The `cupi` binary must be built and available in PATH:

```bash
bash install-local.sh
# or build manually:
go build -o cupi .
```

## Authentication

```bash
# Login — tests connectivity and saves credentials
cupi auth login --host cuc.example.com --username admin --server prod --default

# Add additional credential types to an existing server
cupi auth set-credentials --type application --username app-user --server prod
cupi auth set-credentials --type platform    --username os-admin  --server prod

# Show credential status
cupi auth status [--server prod]

# Switch the default server
cupi auth switch prod

# List all configured servers
cupi auth list

# Logout — removes credentials from secure storage
cupi auth logout [--server prod] [--type cupi|application|platform|all]
```

Passwords are stored in the **OS keystore** (macOS Keychain, Windows Credential Manager, Linux Secret Service). Server config (hostnames, usernames, default server) is stored in `~/.cupi-cli/config.json` — no passwords on disk.

| Type | Used For |
|------|---------|
| `cupi` | CUPI REST provisioning (default) |
| `application` | Application-level APIs |
| `platform` | OS admin / platform-level access (PAWS) |

## Command Discovery

Run `cupi --help` to list all available commands. Run `cupi <command> --help` for exact flags.

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
cupi users list
cupi users list --max 50
cupi users list --query "(alias startswith j)"
cupi users get jsmith
cupi users get jsmith --output json
cupi users add --alias jsmith --dtmf 1001 --first-name John --last-name Smith
cupi users update jsmith --display-name "John Smith" --department Engineering
cupi users remove jsmith

# Distribution Lists
cupi distlists list
cupi distlists get allvoicemail
cupi distlists add --alias mylist --display-name "My List"
cupi distlists update mylist --display-name "Updated List"
cupi distlists remove mylist
cupi distlists members list mylist
cupi distlists members add mylist <member-objectId>
cupi distlists members remove mylist <member-objectId>

# Call Handlers
cupi handlers list
cupi handlers get "Opening Greeting"
cupi handlers add --display-name "My Handler" --template-id <objectId>
cupi handlers update "My Handler" --dtmf 9999
cupi handlers remove "My Handler"

# Class of Service
cupi cos list
cupi cos get UnityMailboxTemplate
cupi cos update UnityMailboxTemplate --display-name "Updated COS"

# User Templates
cupi templates list
cupi templates get voicemailusertemplate

# Schedules
cupi schedules list
cupi schedules get "Weekdays"

# System Info
cupi system
cupi system --output json
```

## AST — System Health Monitoring

```bash
cupi ast disk                     # Disk partition usage
cupi ast heartbeat                # Heartbeat rates
cupi ast tftp                     # TFTP server statistics
cupi ast alerts                   # All system alerts
cupi ast alerts --triggered       # Only currently triggered alerts
cupi ast perfmon                  # Perfmon object catalog
cupi ast perfmon --output json | jq '.[0]'
```

## PAWS — Platform Administration

Requires platform credentials (`cupi auth set-credentials --type platform --username admin`).

```bash
cupi paws cluster status          # OS-level cluster node info
cupi paws cluster replication     # Replication health check
cupi paws drs status              # DRS backup/restore status
cupi paws drs backup --sftp-server 10.0.0.5 --sftp-user backup --sftp-password secret --sftp-dir /backups
```

## DIME — Log File Downloads

```bash
cupi dime get-file /var/log/active/syslog/CiscoSyslog > syslog.txt
cupi dime get-file /var/log/active/tomcat/catalina.out --output /tmp/catalina.out
cupi dime get-file syslog/CiscoSyslog --node 10.0.0.5
```

## Multi-server

```bash
cupi users list --server lab
cupi auth switch prod
```

## Phase 1 Scraper

```bash
go run ./tools/scraper/ --json api_reference.json --md api_reference.md
```
