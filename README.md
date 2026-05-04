
# assmi-super-shop-erp-backend

## Service Overview

SuperShop ERP System
- Auth Service (Authentication & Authorization)
- Product Service
- Order Service
- POS Service
- Reporting Service
- API Gateway (Routing, Rate Limiting, Auth Middleware)

## Folder Structure

supershop-erp/
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── app/
│   │   └── bootstrap.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── database/
│   │   └── mysql.go
│   │
│   ├── middleware/
│   │   ├── auth.go
│   │   └── rbac.go
│   │
│   ├── platform/
│   │
│   ├── services/
│   │   ├── auth/
│   │   │   ├── user/
│   │   │   │   ├── user.model.go
│   │   │   │   ├── user.repository.go
│   │   │   │   ├── user.service.go
│   │   │   │   └── user.handler.go
│   │   │   ├── role/
│   │   │   │   ├── role.model.go
│   │   │   │   ├── role.repository.go
│   │   │   │   ├── role.service.go
│   │   │   │   └── role.handler.go
│   │   │   ├── permission/
│   │   │   │   ├── permission.model.go
│   │   │   │   ├── permission.repository.go
│   │   │   │   ├── permission.service.go
│   │   │   │   └── permission.handler.go
│   │   │   └── routes.go
│   │   │
│   │   ├── product/
│   │   │   ├── category/
│   │   │   │   ├── category.model.go
│   │   │   │   ├── category.repository.go
│   │   │   │   ├── category.service.go
│   │   │   │   └── category.handler.go
│   │   │   ├── item/
│   │   │   │   ├── item.model.go
│   │   │   │   ├── item.repository.go
│   │   │   │   ├── item.service.go
│   │   │   │   └── item.handler.go
│   │   │   ├── inventory/
│   │   │   │   ├── inventory.model.go
│   │   │   │   ├── inventory.repository.go
│   │   │   │   ├── inventory.service.go
│   │   │   │   └── inventory.handler.go
│   │   │   └── routes.go
│   │   │
│   │   ├── order/
│   │   │   ├── cart/
│   │   │   │   ├── cart.model.go
│   │   │   │   ├── cart.repository.go
│   │   │   │   ├── cart.service.go
│   │   │   │   └── cart.handler.go
│   │   │   ├── purchase/
│   │   │   │   ├── purchase.model.go
│   │   │   │   ├── purchase.repository.go
│   │   │   │   ├── purchase.service.go
│   │   │   │   └── purchase.handler.go
│   │   │   ├── payment/
│   │   │   │   ├── payment.model.go
│   │   │   │   ├── payment.service.go
│   │   │   │   └── payment.handler.go
│   │   │   └── routes.go
│   │   │
│   │   ├── pos/
│   │   │   ├── checkout/
│   │   │   │   ├── checkout.service.go
│   │   │   │   └── checkout.handler.go
│   │   │   └── routes.go
│   │   │
│   │   └── reporting/
│   │       ├── sales/
│   │       │   ├── sales_report.service.go
│   │       │   └── sales_report.handler.go
│   │       └── routes.go
│   │
│   ├── router/
│   │   └── router.go
│   │
│   └── shared/
│       ├── logger/
│       ├── response/
│       ├── utils/
│       └── validator/
│
├── migrations/
├── tests/
│   ├── integration/
│   └── e2e/
├── deploy/
├── scripts/
├── docs/
├── go.mod
└── README.md


# run project.

docker-compose up --build


# deployment docs

- CI/CD (GitHub Actions + SSH): docs/CI_CD_GITHUB_ACTIONS_SSH.md

## HTTPS over IP (Production)

If you access the API by server IP and need HTTPS, this project now includes an Nginx reverse proxy for TLS termination.

1. Generate a self-signed certificate for your server IP:

```bash
chmod +x scripts/generate-ip-cert.sh
./scripts/generate-ip-cert.sh <SERVER_IP>
```

Example:

```bash
./scripts/generate-ip-cert.sh 203.76.120.10
```

For your domain, use:

```bash
./scripts/generate-ip-cert.sh erp.vidatech.com.bd www.erp.vidatech.com.bd 144.79.133.252
```

2. Start production stack:

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

3. Open:

```text
https://<SERVER_IP>:10443
```

If your DNS points `erp.vidatech.com.bd` to this server and your host network is forwarding 443 to the container, you can use:

```text
https://erp.vidatech.com.bd
```

Notes:
- Browser will show a warning for self-signed certificates. This is expected.
- Use a domain + trusted CA certificate (Let's Encrypt, etc.) for public production use.
- If your host has free 80/443 ports, you can run with standard ports:

```bash
HTTPS_HTTP_PORT=80 HTTPS_PORT=443 docker compose -f docker-compose.prod.yml up -d --build
```


# API Testing 

- API Testing : api_tests.http

- Run API Testing : 

- install extension : REST Client

- open api_tests.http file

- run api : Ctrl + Alt + R



# Navigate to project directory
cd "/Users/ssjoy/veda tech/assmi-super-shop-erp-backend"

# Run all migrations (create tables)
./scripts/migrate.sh up

# Check migration status
./scripts/migrate.sh status

# Rollback (delete tables)
./scripts/migrate.sh down

# Reset database (down then up)
./scripts/migrate.sh reset

# radis connect url localhost
redis://default:change-me-redis-password@127.0.0.1:6379
# Ieam setup require
https://drive.google.com/file/d/10n97y-2BZiBcYI6FnIjY5jXfzXALSePi/view?usp=drivesdk

HOW TO cached  USE THE FIX:

If permission error occurs in future:
  ./scripts/quick-clear-cache.sh

For specific user:
  ./scripts/clear-redis-cache.sh user 1

For all users:
  ./scripts/clear-redis-cache.sh users

To check cache status:
  ./scripts/clear-redis-cache.sh stats
