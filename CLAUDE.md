# go-ms-universe

Two microservices — **Products** (REST API + Kafka producer) and **Notifications** (Kafka consumer).

## Quick Reference

- `make help` — list all targets
- `make docker.up` — start Postgres, Kafka, Prometheus
- `make run.products` / `make run.notifications` — run services locally
- `make generate` — regenerate mocks after changing interfaces
- `swagger.yml` in `cmd/products/` is hand-maintained — update it when changing HTTP endpoints

## Architecture

Three layers, dependency flows inward:

```
httpapi (handlers) → product (domain + service) ← postgres (storage)
                                ↑
                          kafka (events)
```

- Define interfaces where consumed, not where implemented.
- Domain model uses unexported fields — construct via `New()`, reconstruct via `FromStorage()`.
- Return concrete structs from constructors; accept interfaces in consumers.

## Code Style

- Go 1.26 — use modern stdlib (`slog`, `errors.Join`, `slices`, `t.Context()`).
- Run `make lint.fix` first, then fix anything remaining.
- Run `make modernise` to auto-apply modern idioms.
- Write a failing test first for new features and bug fixes.

## Testing

- `make test.unit` — unit tests (no Docker needed)
- Use `t.Context()` instead of `context.Background()` in tests.
- Use `t.Cleanup()` instead of `defer` for resource teardown.
- Mock generation: `go.uber.org/mock` via `go generate ./...`.
