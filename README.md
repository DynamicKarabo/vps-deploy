# vps-deploy

**vps-deploy** is a Go CLI tool for deploying, rolling back, and checking the status of services on remote VPS hosts over SSH.

## Overview

Define your services in a `deploy.yaml` file. Each service specifies a remote host, SSH credentials, commands to run on deploy, an optional health check URL, and an optional rollback command.

### Commands

| Command         | Description                                      |
|-----------------|--------------------------------------------------|
| `vps-deploy deploy <service>`   | Deploy a service — runs deploy_commands via SSH, then checks health |
| `vps-deploy rollback <service>` | Rollback a service — runs its rollback_command via SSH |
| `vps-deploy status <service>`   | Check service health by polling its health_check_url |

Global flag: `--config` (default `deploy.yaml`) to specify an alternative config path.

## Example `deploy.yaml`

```yaml
services:
  filebrowser:
    host: 178.105.76.236
    user: root
    key_path: ~/.ssh/id_ed25519
    deploy_commands:
      - docker compose pull
      - docker compose up -d
    health_check_url: http://localhost:8081/health
    rollback_command: docker compose down && docker compose -f docker-compose.rollback.yml up -d
```

## Usage

```bash
# Build
go build -o vps-deploy .

# Deploy a service
./vps-deploy deploy filebrowser

# Rollback a service
./vps-deploy rollback filebrowser

# Check service status
./vps-deploy status filebrowser
```

## License

MIT
