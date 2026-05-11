# vps-deploy

**vps-deploy** is a Go CLI tool for deploying, rolling back, checking status, listing, destroying, and copying files to services on remote VPS hosts over SSH.

## Quick Start

```bash
# Install
curl -L https://github.com/DynamicKarabo/vps-deploy/releases/latest/download/vps-deploy-linux-amd64.tar.gz | tar xz
sudo mv vps-deploy /usr/local/bin/

# Or build from source
git clone https://github.com/DynamicKarabo/vps-deploy.git
cd vps-deploy
GOTOOLCHAIN=go1.23.6 go build -o vps-deploy .
```

## Commands

| Command | Description |
|---------|-------------|
| `deploy <service>` | Deploy a service — SSH → pull → up → health check |
| `rollback <service>` | Rollback a service — runs its rollback_command |
| `status <service>` | Check a service's health via its health_check_url |
| `list` | Show all configured services with docker status table |
| `init` | Interactive deploy.yaml config generator |
| `copy <svc> <local> <remote>` | SCP a file to a service's host using configured SSH key |
| `destroy <svc>` | `docker compose down` with optional `--volumes` / `--rmi` |

## Config (`deploy.yaml`)

```yaml
services:
  filebrowser:
    host: 178.105.76.236
    user: root
    key_path: ~/.ssh/vps_key
    deploy_commands:
      - mkdir -p /root/filebrowser/data /root/filebrowser/config /root/filebrowser/database
      - cd /root/filebrowser && docker compose pull
      - cd /root/filebrowser && docker compose up -d
    health_check_url: http://localhost:8081/
    rollback_command: cd /root/filebrowser && docker compose down
```

## Usage

```bash
# Deploy a service
vps-deploy deploy filebrowser

# Check health
vps-deploy status filebrowser

# List all services
vps-deploy list

# Generate config interactively
vps-deploy init

# Copy a file to a service host
vps-deploy copy filebrowser ./docker-compose.yml /root/filebrowser/docker-compose.yml

# Destroy a service
vps-deploy destroy shiori --volumes

# Rollback
vps-deploy rollback filebrowser
```

## Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `deploy.yaml` | Path to config file |

## Destroy Flags

| Flag | Description |
|------|-------------|
| `--volumes` | Remove named volumes (`docker compose down -v`) |
| `--rmi` | Remove images (`docker compose down --rmi all`) |
| `--path` | Custom compose directory (default: `/root/<service>`) |

## CI/CD

On tag push (`v*`), GitHub Actions builds binaries for 5 platforms and creates a release:

- `vps-deploy-linux-amd64.tar.gz`
- `vps-deploy-linux-arm64.tar.gz`
- `vps-deploy-darwin-amd64.tar.gz`
- `vps-deploy-darwin-arm64.tar.gz`
- `vps-deploy-windows-amd64.zip`

## License

MIT
