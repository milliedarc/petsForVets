# Pets for Vets

A web app for pet owners and their veterinary clinic. Owners keep a profile for
each pet (species, breed, age, microchip, and so on) and request repeat
prescriptions online, then follow each request from *Pending* to *Ready* (or
see why it was rejected).

**Stack:** Vue 3 + TypeScript (Vite, PrimeVue, Bootstrap) on the frontend, and a
Go API on PostgreSQL with hand-written SQL on the backend.

## Architecture

```
Browser ──► Vite dev server (:5173) ──/api/*──► Go API (:8080) ──► PostgreSQL (:5432, Docker)
            Vue 3 + TypeScript                  net/http, sqlc, goose
```

- **Frontend** (`src/`): Vue 3 single-page app, written in TypeScript with
  `strict` type checking. It talks to the API through a small typed client in
  `src/api/`.
- **API** (`api/`): Go standard library HTTP server, no web framework. Queries
  are plain SQL files that [sqlc](https://sqlc.dev) turns into type-safe Go
  code; no ORM. The schema lives in versioned [goose](https://github.com/pressly/goose)
  migrations that are applied automatically when the server starts.
- **Auth**: email + password (bcrypt), with sessions in a Postgres table and an
  `HttpOnly`, `SameSite` cookie. Vite proxies `/api` so the browser sees a
  single origin, so no CORS setup is needed.
- **Data rules in the database**: foreign keys (a pet's breed must belong to its
  species), `CHECK` constraints, and a partial unique index (at most one open
  prescription request per prescription). Every query is scoped to the
  logged-in owner.

The backend has its own [README](api/README.md): request flow, endpoints, how to
add a query, and notes for PHP developers.

## Getting started

### Prerequisites

- [Node.js](https://nodejs.org/) 20.19+ or 22.12+
- [Go](https://go.dev/dl/) 1.27+
- Docker, via [Docker Desktop](https://www.docker.com/products/docker-desktop/)
  or [Colima](https://github.com/abiosoft/colima) on macOS

To change SQL queries or migrations you'll also want
[sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) and
[goose](https://github.com/pressly/goose#install). Running the app doesn't need them.

### Run it

```sh
npm install
docker-compose up -d     # PostgreSQL on localhost:5432
npm run dev              # Go API on :8080 + Vue on :5173; Ctrl+C stops both
```

Open <http://localhost:5173> and create an account on the registration page.
The first start of the API creates the database schema and the reference data
(species, breeds, medicines).

### Sample data

Prescriptions are normally written by clinic staff (that side of the app isn't
built yet). To give your account's pets a couple of sample prescriptions, add a
pet first and then run:

```sh
docker-compose exec -T db psql -U petsforvets -v email=you@example.com < api/dev/seed_dev_data.sql
```

The clinic team view is only available to staff accounts. To make yours one:

```sh
docker-compose exec db psql -U petsforvets -c "UPDATE users SET veterinary = true WHERE email = 'you@example.com';"
```

## Scripts

| Command | What it does |
|---|---|
| `npm run dev` | API and frontend together, with log lines prefixed `[api]` / `[web]` |
| `npm run api` | Just the Go API (builds `api/bin/server`, then runs it) |
| `npm run dev:web` | Just the Vite dev server |
| `npm run build` | Type-check with `vue-tsc`, then build the frontend into `dist/` |
| `npm run type-check` | Type-check only |
| `docker-compose down` | Stop Postgres (data is kept); add `-v` to wipe it |

## Project structure

```
api/                     Go backend (see api/README.md)
  cmd/server/            entry point
  internal/db/           migrations, SQL queries, sqlc-generated code
  internal/httpapi/      HTTP handlers, auth, routing
  dev/                   development-only sample data
src/                     Vue frontend
  api/                   typed API client (fetch wrapper + one file per resource)
  components/            shared components (navbar, pet cards, ...)
  views/                 pages: dashboard, pet details, prescriptions, account
  types/                 TypeScript types matching the API's JSON
docker-compose.yml       local PostgreSQL
```

## Status

| Feature | |
|---|---|
| Accounts: register, log in/out, change password | ✅ |
| Pets: add, view, edit, delete; species and breeds | ✅ |
| Repeat prescriptions: view, request, cancel | ✅ |
| Pet photos | ⏳ upload not yet ported to the new API |
| Clinic team dashboard (review requests, write prescriptions) | 🚧 planned |
| Contact details and addresses | 🚧 planned |
| Automated tests | 🚧 planned |

The app started as a Vue frontend on [PocketBase](https://pocketbase.io/). The
backend was later rewritten as the Go + PostgreSQL API above, to own the schema,
the SQL and the auth.

## Recommended IDE setup

[VS Code](https://code.visualstudio.com/) with the
[Vue - Official](https://marketplace.visualstudio.com/items?itemName=Vue.volar)
and [Go](https://marketplace.visualstudio.com/items?itemName=golang.go)
extensions, or a JetBrains IDE (WebStorm / GoLand).
