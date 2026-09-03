# My Vehicles – Backup Agent

Backend service that stores and restores backups for the **My Vehicles** Android app. It exposes a small HTTP API the app calls to mirror its local database — one table at a time — to a server of your choice, and to pull it back down on a new device.

Built with [Goravel](https://www.goravel.dev), a Laravel-style web framework for Go. It's a from-scratch Go rewrite of an earlier NestJS implementation, kept behavior-compatible with the Android client while adding first-class support for SQLite, PostgreSQL and MySQL.

## Features

- **Full-table mirror sync** — each `POST /backup/{table}` call replaces the table's entire content: rows in the payload are kept, everything else is deleted, atomically.
- **Multi-database** — SQLite (zero-config default), PostgreSQL or MySQL, selected with a single environment variable. No raw, dialect-specific SQL anywhere in the codebase.
- **API key authentication** — every route but `/ping` and `/healthcheck` requires a shared secret on the `x-api-key` header.
- **Schema version gate** — requests carrying an unexpected `version` are rejected with `409 Conflict` before touching the database, so an outdated client can't silently corrupt newer data.
- **Strict payload validation** — unknown fields, missing required fields, invalid enum values and duplicate IDs are all rejected with `400 Bad Request`.
- **Dashboard stats** — a lightweight endpoint reporting row counts per table, handy for a quick sanity check after a restore.

## API reference

All endpoints are unprefixed (e.g. `http://localhost:3000/ping`).

| Method | Path                                | Auth      | Description                                   |
|--------|--------------------------------------|-----------|------------------------------------------------|
| GET    | `/healthcheck`                       | none      | Plain liveness probe, used by Docker's `HEALTHCHECK`. |
| GET    | `/ping?version=`                     | none      | Connectivity + schema version check.          |
| GET    | `/dashboard/stats?version=`          | API key   | Row counts for every table.                   |
| POST   | `/backup/{table}`                    | API key   | Replaces `{table}`'s content with the request body. |
| GET    | `/backup/{table}?version=`           | API key   | Returns `{table}`'s full content.             |

`{table}` is one of `vehicles`, `vehicle_refuels`, `vehicle_service_types`, `vehicle_services`, `vehicle_service_reminders`.

`POST` requests take `{"version": "<n>", "data": [...]}` and reply with `{"upserted_count": n, "deleted_count": n}`. `GET` requests reply with `{"version": "<n>", "data": [...]}`.

> [!NOTE]
> `version` must match the value baked into `backup.schema_version` (see [Configuration](#configuration)). This mirrors the Android app's own `DB_SCHEMA_VERSION` and exists to stop an out-of-sync client from overwriting the server's data with a payload it can't correctly represent yet.

## Requirements

- [Go](https://go.dev) 1.25+
- One of: SQLite (no extra setup), PostgreSQL, or MySQL
- [Docker](https://www.docker.com) and Docker Compose, if you'd rather not install Go or a database locally

## Getting started

```bash
git clone <repo-url>
cd my_vehicles_an_backup_agent
cp .env.example .env
go run . artisan key:generate
```

Set `API_KEY` in `.env` to a secret of your choice, create the SQLite data directory (skip this if you're using Postgres or MySQL), then run the migrations and start the server:

```bash
mkdir -p storage/database
go run . artisan migrate
go run .
```

The API is now listening on `http://127.0.0.1:3000`. For live reload during development, use [air](https://github.com/air-verse/air) (already configured via `.air.toml`):

```bash
air
```

## Configuration

All configuration is read from environment variables, normally supplied via `.env` (copy it from `.env.example`). Running the [production image](#production-image-docker-hub) instead of `go run`? None of these are baked into it — see that section for how to pass them at container start.

**Application**

| Variable    | Default                | Description                                                                                                    |
|-------------|-------------------------|-------------------------------------------------------------------------------------------------------------------|
| `APP_NAME`  | `Goravel`               | Display name, used in a few log/error messages.                                                                    |
| `APP_ENV`   | `local`                 | `local` or `production`. Affects error verbosity and default framework behavior.                                    |
| `APP_KEY`   | —                       | Required. 32-byte app encryption key. Generate with `go run . artisan key:generate` (writes to `.env`) or set the printed value directly as an env var. |
| `APP_DEBUG` | `true`                  | Set to `false` in production — `true` leaks stack traces in error responses.                                       |
| `APP_URL`   | `http://localhost`      | Base URL, used for link generation. Not security-relevant.                                                          |
| `APP_HOST`  | `0.0.0.0`               | Interface the HTTP server binds to. `0.0.0.0` accepts connections from other devices (e.g. the Android app over the network); `127.0.0.1` only accepts connections from the same host/container. |
| `APP_PORT`  | `3000`                  | HTTP port the server listens on.                                                                                    |

**Logging**

| Variable      | Default | Description                                    |
|---------------|---------|--------------------------------------------------|
| `LOG_CHANNEL` | `stack` | Goravel log channel.                              |
| `LOG_LEVEL`   | `debug` | `debug`, `info`, `warning` or `error`. Use `info` or `warning` in production. |

**Session**

| Variable            | Default | Description                     |
|---------------------|---------|-----------------------------------|
| `SESSION_DRIVER`    | `file`  | Goravel session driver.            |
| `SESSION_LIFETIME`  | `120`   | Session lifetime, in minutes.      |

**Database**

| Variable         | Default                     | Description                                                   |
|------------------|------------------------------|-----------------------------------------------------------------|
| `DB_CONNECTION`  | `sqlite`                    | `sqlite`, `postgres` or `mysql`.                                  |
| `DB_DATABASE`    | `storage/database/database.sqlite` | SQLite: file path (parent directory must exist). Postgres/MySQL: database name. |
| `DB_HOST`        | `127.0.0.1`                 | Postgres/MySQL only.                                              |
| `DB_PORT`        | `5432`                      | Postgres/MySQL only.                                              |
| `DB_USERNAME`    | `goravel`                   | Postgres/MySQL only.                                              |
| `DB_PASSWORD`    | `goravel`                   | Postgres/MySQL only.                                              |
| `DB_SCHEMA`      | `public`                    | Postgres only.                                                    |

**Backup agent**

| Variable  | Default    | Description                                                                     |
|-----------|-------------|-------------------------------------------------------------------------------------|
| `API_KEY` | `changeme` | Required, no default in production. Shared secret checked against `x-api-key` on every route but `/ping` and `/healthcheck`. |

**Authentication**

| Variable      | Default | Description                                                                                          |
|---------------|---------|-----------------------------------------------------------------------------------------------------------|
| `JWT_SECRET`  | —       | Unused by the backup agent itself — kept from the Goravel scaffold. Safe to leave empty.                    |

See `.env.example` for the full list, including mail, gRPC and AI-related settings also inherited from the Goravel scaffold and unused here.

> [!WARNING]
> `APP_KEY`, `API_KEY` and `DB_PASSWORD` are secrets. Never commit a filled-in `.env`, and never bake them into a Docker image — pass them at container start instead (see [Production image](#production-image-docker-hub)).

## Docker

There are two Dockerfiles:

- **`Dockerfile`** — development. Bakes `.env` straight into the image, so `docker compose up --build` just works with no extra setup. Not meant to leave your machine.
- **`Dockerfile.prod`** — production. No `.env`, no secrets, no source baked in — just the compiled binary, `public/` and `resources/`, running as a non-root user. All configuration is supplied at container start (see [Configuration](#configuration)). This is what `scripts/build.sh` builds and what the published Docker Hub image is made from.

A ready-to-use `compose.yml` is provided, defaulting to SQLite with just the database directory (not the whole `storage/`, which also holds logs) bind-mounted to `./data/sqlite` on the host — no separate database service to run:

```bash
cp .env.example .env
# edit .env: set API_KEY
docker compose up -d --build
docker compose exec app /www/main artisan migrate
```

> [!TIP]
> Because the container's `ENTRYPOINT` runs the HTTP server directly, run one-off Artisan commands (migrations, `db:seed`, etc.) with `docker compose exec app /www/main artisan <command>`.

### Production image (Docker Hub)

For a real deployment, use the published image, [`ddavideccia/my-vehicles-an-backup-agent`](https://hub.docker.com/r/ddavideccia/my-vehicles-an-backup-agent), instead of building locally. Point `compose.yml` at it with `image:` instead of `build: .`:

```yaml
services:
  app:
    image: ddavideccia/my-vehicles-an-backup-agent:latest
    ports:
      - "3000:3000"
    env_file: .env
    environment:
      DB_CONNECTION: sqlite
      DB_DATABASE: storage/database/database.sqlite
    volumes:
      - ./data/sqlite:/www/storage/database
    restart: unless-stopped
```

```bash
cp .env.example .env
# edit .env: set APP_KEY, API_KEY, and anything else in the Configuration table above
docker compose up -d
docker compose exec app /www/main artisan migrate
```

> [!IMPORTANT]
> The image has no `.env` inside it — every variable in [Configuration](#configuration) must be supplied from outside, via `env_file`/`environment` in `compose.yml` as above, or your orchestrator's secrets store. They're read once at process startup: editing `.env` on the host doesn't affect an already-running container. Apply changes with `docker compose up -d` (recreates the container with the new values), not `docker compose restart` (reuses the environment the container already has).

Publishing a new version of the image (maintainers only) is a `make build` away — see [`scripts/build.sh`](scripts/build.sh), which builds `Dockerfile.prod` for both `linux/amd64` and `linux/arm64`, tags it with the current commit SHA and `latest`, and pushes both tags to Docker Hub. Requires `docker login` first; skip the push with `PUSH=0 make build` (builds a single, host-platform image loaded into the local daemon instead — multi-platform output can only be pushed straight to the registry, not loaded locally).

### Choosing a different database

Each block below binds the database's storage to a host directory, so backing it up is just a matter of copying that directory (or scheduling a tool like `pg_dump`/`mysqldump` against the running container for consistent, hot backups). Pick the one matching your `DB_CONNECTION` and adjust `.env` to match. They use `build: .` for the dev image; for production swap that line for `image: ddavideccia/my-vehicles-an-backup-agent:latest` as shown [above](#production-image-docker-hub).

<details>
<summary>SQLite (default, already in <code>compose.yml</code>)</summary>

SQLite needs no separate database service: bind-mount just `storage/database/` — not all of `storage/`, which also holds logs — directly to the host, so `database.sqlite` lives there as a plain file.

```yaml
services:
  app:
    build: .
    ports:
      - "3000:3000"
    env_file: .env
    environment:
      DB_CONNECTION: sqlite
      DB_DATABASE: storage/database/database.sqlite
    volumes:
      - ./data/sqlite:/www/storage/database
    restart: unless-stopped
```

> [!TIP]
> Backing up SQLite is as simple as copying the file: `cp ./data/sqlite/database.sqlite ./backups/database-$(date +%F).sqlite`. Since the whole database is a single file, this is safe to do while the container is stopped, or while it's running as long as no write is in flight — for a byte-for-byte consistent snapshot under load, use `sqlite3 ./data/sqlite/database.sqlite ".backup ./backups/database-$(date +%F).sqlite"` instead of a plain `cp`.

</details>

<details>
<summary>PostgreSQL</summary>

```yaml
services:
  app:
    build: .
    ports:
      - "3000:3000"
    env_file: .env
    environment:
      DB_CONNECTION: postgres
      DB_HOST: postgres
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: ${DB_DATABASE:-goravel}
      POSTGRES_USER: ${DB_USERNAME:-goravel}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-goravel}
    volumes:
      - ./data/postgres:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USERNAME:-goravel}"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped
```

> [!TIP]
> `./data/postgres` holds Postgres's raw data files — stopping the stack and archiving that directory is a valid physical backup. For a logical backup you can restore selectively or across major versions, prefer `docker compose exec postgres pg_dump -U ${DB_USERNAME} ${DB_DATABASE} > backup.sql`.

</details>

<details>
<summary>MySQL</summary>

```yaml
services:
  app:
    build: .
    ports:
      - "3000:3000"
    env_file: .env
    environment:
      DB_CONNECTION: mysql
      DB_HOST: mysql
    depends_on:
      mysql:
        condition: service_healthy
    restart: unless-stopped

  mysql:
    image: mysql:8
    environment:
      MYSQL_DATABASE: ${DB_DATABASE:-goravel}
      MYSQL_USER: ${DB_USERNAME:-goravel}
      MYSQL_PASSWORD: ${DB_PASSWORD:-goravel}
      MYSQL_ROOT_PASSWORD: ${DB_PASSWORD:-goravel}
    volumes:
      - ./data/mysql:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped
```

> [!TIP]
> As with Postgres, `./data/mysql` is a valid physical backup target when the stack is stopped. For hot, logical backups use `docker compose exec mysql mysqldump -u root -p"${DB_PASSWORD}" ${DB_DATABASE} > backup.sql`.

</details>

> [!WARNING]
> The examples above pass plaintext passwords via environment variables, which is fine for local development but not for production. Use Docker/Compose secrets (or an external secrets manager) for credentials, and terminate TLS in front of the app before exposing it to the internet.

## Testing

```bash
go test ./tests/...
```

Feature tests spin up the full HTTP stack against a dedicated SQLite database (`tests/feature/.env`), exercising authentication, schema versioning, payload validation and the mirror-sync behavior end to end.

## Project structure

```
app/
  http/
    controllers/   # thin HTTP handlers: ping, healthcheck, dashboard, backup
    middleware/     # API key + schema version gates
    requests/       # envelope decoding and payload validation
  models/           # one struct per table, plus shared enums
  services/         # generic sync/fetch/count logic shared by every table
config/             # one file per concern (database.go, backup.go, ...)
database/migrations/
routes/             # route registration (api.go is the one that matters here)
tests/feature/      # end-to-end HTTP tests
```
