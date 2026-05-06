# Ploggo Memory & Instructions

## Testing
- **Postgres Tests**: The test `//avalon/storage/postgres:go_default_test` requires Docker (or Podman) to be running locally to spin up test containers via `testcontainers-go`. If Docker is not running, these tests will fail with a connection error.
