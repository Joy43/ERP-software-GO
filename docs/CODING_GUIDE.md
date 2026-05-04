# Coding Guide

This guide defines coding conventions for the ASSMI Super Shop ERP backend.

## 1) Core Principles

- Keep code simple, readable, and maintainable.
- Prefer explicit behavior over magic/hidden abstractions.
- Follow clear separation of concerns: handler -> service -> repository.
- Make changes small and focused.
- Optimize for correctness first, then performance.

## 2) Project Structure Rules

Use the existing structure consistently:

- `cmd/server`: app entrypoint only.
- `internal/app`: application bootstrap and dependency wiring.
- `internal/config`: environment/config loading and validation.
- `internal/database`: DB connection and shared DB setup.
- `internal/middleware`: auth, RBAC, request-level middleware.
- `internal/platform`: infrastructure adapters (cache, queue, external clients).
- `internal/services/*`: domain modules.
- `internal/router`: root router and route registration.
- `internal/shared/*`: shared internal helpers (logger, response, utils, validator).
- `migrations`: schema migrations.
- `tests/integration` and `tests/e2e`: non-unit tests.

## 3) Domain Module Layout

Each domain module should follow:

- `*.model.go`: domain and persistence models.
- `*.repository.go`: DB access logic only.
- `*.service.go`: business rules/use-cases.
- `*.handler.go`: HTTP transport logic only.
- `routes.go`: module-level route registration.

Current canonical service modules:

- `auth/{user,role,permission}`
- `product/{category,item,inventory}`
- `order/{cart,purchase,payment}`
- `pos/checkout`
- `reporting/sales`

## 4) Naming Conventions

### Folders and packages

- Use lowercase package names.
- Avoid stutter and ambiguous names.
- Keep naming consistent (prefer singular domain names already used).

### Files

- Use suffix-based files by layer:
  - `user.model.go`
  - `user.repository.go`
  - `user.service.go`
  - `user.handler.go`

### Symbols

- Exported names: `PascalCase`
- Unexported names: `camelCase`
- Constants: `PascalCase` for exported, `camelCase` for unexported.
- Avoid one-letter variable names except tiny loop counters.

## 5) Layer Responsibilities

### Handler layer

- Parse/validate request input.
- Call service methods.
- Convert result/error to API response format.
- No SQL/query logic.
- No business decision logic.

### Service layer

- Enforce business rules and workflows.
- Use repositories through interfaces where practical.
- Manage transactions when business operation spans multiple repository calls.
- Return typed/domain errors.

### Repository layer

- Execute DB operations and mapping only.
- Keep functions deterministic and focused.
- Do not include HTTP concepts.

## 6) API and Response Standards

- Use a consistent response envelope via shared response helpers.
- Keep error responses predictable and machine-readable.
- Validate input at boundaries (handler + validator helpers).
- Add API version prefix in root routes (`/api/v1`).

Recommended error response shape:

```json
{
  "success": false,
  "message": "validation failed",
  "error_code": "VALIDATION_ERROR",
  "details": {}
}
```

## 7) Error Handling

- Never swallow errors.
- Wrap errors with context before returning.
- Convert low-level errors to domain-safe messages at service/handler boundaries.
- Do not leak DB/internal details to API responses.

## 8) Logging and Observability

- Use structured logs (key/value), not ad-hoc plain text.
- Include request ID/correlation ID in request-scoped logs.
- Log at appropriate levels (debug/info/warn/error).
- Do not log secrets, passwords, or tokens.

## 9) Configuration and Secrets

- Read config from environment variables.
- Keep defaults explicit and safe.
- Add new variables to `.env.example` when introduced.
- Never commit real secrets.

## 10) Database and Migrations

- Every schema change must have a migration.
- Use forward-only migration scripts in production workflows.
- Keep migration names descriptive and timestamped.
- Repositories should be the only layer issuing queries.

## 11) Testing Standards

### Unit tests

- Place unit tests next to code using `*_test.go`.
- Focus on service business rules and edge cases.
- Use table-driven tests where useful.

### Integration/E2E tests

- Put integration tests in `tests/integration`.
- Put end-to-end tests in `tests/e2e`.
- Keep test data isolated and repeatable.

Minimum CI checks:

- `go test ./...`
- `go vet ./...`
- formatting and lint checks

## 12) Code Review Checklist

Before opening PR:

- Build and tests pass locally.
- New behavior includes tests.
- No dead code or commented-out code.
- Naming and structure follow this guide.
- API changes are documented.
- Migration included for schema updates.

## 13) Git and PR Conventions

- Keep PRs small and single-purpose.
- Commit messages should be clear and action-focused.
- Prefer conventional style (example: `feat(auth): add role assignment endpoint`).
- Include summary, risk, and test evidence in PR description.

## 14) Definition of Done

A task is done when:

- Code follows this guide.
- Relevant tests are added/updated.
- No new lint/vet issues are introduced.
- Docs/config updates are included where required.
- Reviewer can understand the change without extra explanation.
