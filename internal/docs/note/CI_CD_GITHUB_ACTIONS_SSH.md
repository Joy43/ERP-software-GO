# CI/CD Deployment Guide (GitHub Actions + SSH + Docker Compose)

This document explains a complete deployment flow for this project on your VPS host using GitHub Actions over SSH.

## 1) Deployment Flow

1. Push code to `main`.
2. GitHub Actions starts workflow.
3. Action connects to VPS using SSH key.
4. Server pulls latest code.
5. Server rebuilds and restarts containers with `docker-compose.prod.yml`.
6. SQL migrations are applied to MySQL.
7. Optional debug services (like phpMyAdmin) stay disabled unless explicitly enabled.

## 2) Current Project Paths

- Project root on server: `/home/vidatech-erp/htdocs/erp.vidatech.com.bd/assmi-super-shop-erp-backend`
- Production compose file: `docker-compose.prod.yml`
- Migrations directory: `migrations/`
- Existing workflow file: `.github/workflows/docker-image.yml`

## 3) Server Prerequisites

Install these on VPS:

- Docker Engine
- Docker Compose plugin (`docker compose`)
- Git
- OpenSSH server

Create and keep `.env` in project root with production values:

- `APP_NAME`
- `MYSQL_ROOT_PASSWORD`
- `MYSQL_DATABASE`
- `REDIS_PASSWORD`
- `JWT_ACCESS_SECRET`
- `JWT_REFRESH_SECRET`
- `JWT_ACCESS_TTL_MINUTES`
- `JWT_REFRESH_TTL_DAYS`
- `BCRYPT_COST`

Use `.env.example` as template and do not commit real secrets.

## 4) SSH Setup for GitHub Actions

### 4.1 Generate deploy key pair on your local machine

```bash
ssh-keygen -t ed25519 -C "github-actions-deploy" -f ./github_actions_deploy_key
```

This creates:

- `github_actions_deploy_key` (private key)
- `github_actions_deploy_key.pub` (public key)

### 4.2 Add public key to VPS user

```bash
mkdir -p ~/.ssh
cat github_actions_deploy_key.pub >> ~/.ssh/authorized_keys
chmod 700 ~/.ssh
chmod 600 ~/.ssh/authorized_keys
```

### 4.3 Add GitHub repository secrets

In GitHub repo -> Settings -> Secrets and variables -> Actions, add:

- `SSH_HOST` = VPS public IP or domain
- `SSH_PORT` = SSH port (usually 22)
- `SSH_USER` = VPS deploy user
- `SSH_PRIVATE_KEY` = full content of `github_actions_deploy_key`

Optional hardening secret:

- `SSH_KNOWN_HOSTS` = output of `ssh-keyscan -H your-host`

## 5) Recommended Workflow (Production Safe)

Use this as a reference workflow if you want stricter and migration-aware deployment.

```yaml
name: Deploy to VPS

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup SSH agent
        uses: webfactory/ssh-agent@v0.9.0
        with:
          ssh-private-key: ${{ secrets.SSH_PRIVATE_KEY }}

      - name: Add VPS host key
        run: |
          mkdir -p ~/.ssh
          if [ -n "${{ secrets.SSH_KNOWN_HOSTS }}" ]; then
            echo "${{ secrets.SSH_KNOWN_HOSTS }}" >> ~/.ssh/known_hosts
          else
            ssh-keyscan -p ${{ secrets.SSH_PORT }} -H ${{ secrets.SSH_HOST }} >> ~/.ssh/known_hosts
          fi

      - name: Deploy and migrate
        run: |
          ssh -p ${{ secrets.SSH_PORT }} ${{ secrets.SSH_USER }}@${{ secrets.SSH_HOST }} << 'EOF'
            set -euo pipefail

            cd /home/vidatech-erp/htdocs/erp.vidatech.com.bd/assmi-super-shop-erp-backend

            git fetch origin main
            git reset --hard origin/main

            docker compose -f docker-compose.prod.yml up -d --build

            for f in \
              migrations/000001_create_auth_tables.up.sql \
              migrations/000002_create_rbac_tables.up.sql \
              migrations/000003_create_admin_tables.up.sql \
              migrations/000004_add_employee_fields_to_users.up.sql; do
              docker compose -f docker-compose.prod.yml exec -T db \
                mysql -uroot -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE" < "$f"
            done

            docker compose -f docker-compose.prod.yml ps
          EOF
```

Notes:

- `phpmyadmin` has profile `debug` in `docker-compose.prod.yml`, so it does not start in normal production deploy.
- To start it manually: `docker compose -f docker-compose.prod.yml --profile debug up -d phpmyadmin`

## 6) About Your Existing Workflow

Current workflow in `.github/workflows/docker-image.yml` already does SSH deployment.

It is good as a starting point, but for production you should improve:

- Avoid `StrictHostKeyChecking=no`.
- Add explicit migration step after container startup.
- Use health checks and post-deploy smoke checks.

## 7) Post Deploy Checks

Run on VPS after deployment:

```bash
docker compose -f docker-compose.prod.yml ps
docker compose -f docker-compose.prod.yml logs --tail=80 app
docker compose -f docker-compose.prod.yml exec -T db mysql -uroot -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE" -e "SHOW TABLES;"
```

Expected:

- `app`, `db`, `redis` are Up
- `db` is healthy
- tables include `users`, `roles`, `permissions`, `user_roles`, `role_permissions`

## 8) Rollback Strategy (Simple)

If latest deploy fails:

1. SSH to server.
2. Go to project folder.
3. Checkout previous commit.
4. Rebuild and restart containers.

```bash
cd /home/vidatech-erp/htdocs/erp.vidatech.com.bd/assmi-super-shop-erp-backend
git log --oneline -n 5
git checkout <previous_commit_sha>
docker compose -f docker-compose.prod.yml up -d --build
```

## 9) Common Problems

### Problem: App is up but login returns 500

Cause: migrations not applied, tables missing.

Fix: apply SQL files from `migrations/` in sequence.

### Problem: `8081` not accessible

Cause: phpMyAdmin uses `debug` profile.

Fix:

```bash
docker compose -f docker-compose.prod.yml --profile debug up -d phpmyadmin
```

### Problem: GitHub Action cannot connect via SSH

Cause: wrong key, port, user, firewall, or missing authorized key.

Fix checklist:

- Verify `SSH_HOST`, `SSH_PORT`, `SSH_USER`.
- Verify `SSH_PRIVATE_KEY` matches server public key.
- Confirm VPS allows inbound SSH on configured port.
- Confirm deploy user shell access.

## 10) Security Checklist

- Keep `.env` out of Git.
- Use minimum-privilege deploy user instead of root.
- Prefer pinned Docker image tags in production.
- Store host key in `SSH_KNOWN_HOSTS` secret.
- Rotate JWT and DB credentials periodically.
