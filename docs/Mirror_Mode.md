# Mirror Mode - User Guide

## What is Mirror Mode?

Mirror mode runs a Beo Echo instance whose **whole domain is one project**. Every path on that domain, including `/`, is served by the project unchanged.

Use it when an app only works from a domain root:

- single-page apps and sites that link absolute paths (`/assets/app.js`, `/favicon.ico`)
- frontends that read their route from `window.location.pathname` (e.g. `/f/<token>`)
- an iframe or webhook that must look like a real host, not `https://beo.example.com/<alias>/...`

Without mirror mode such an app, served under `/<alias>/`, requests `/assets/...` from the domain root and gets Beo Echo's default response instead of the file.

---

## Architecture

```
               https://beo.example.com             https://app-mirror.example.com
                         │                                       │
                         ▼                                       ▼
┌───────────────────────────────────────┐   ┌───────────────────────────────────────┐
│  Main instance                        │   │  Mirror instance                      │
│  MIRROR_PROJECT_ALIAS = (empty)       │   │  MIRROR_PROJECT_ALIAS = web-app       │
│                                       │   │                                       │
│  /            → UI                    │   │  /              → project web-app     │
│  /login, /api → UI + management API   │   │  /assets/app.js → project web-app     │
│  /<alias>/... → project <alias>       │   │  /api/...       → project web-app     │
└───────────────────┬───────────────────┘   └───────────────────┬───────────────────┘
                    │                                           │
                    └─────────────► shared database ◄───────────┘
                         (projects, actions, request logs)
```

The project is still created and edited in the main instance's UI. The mirror instance only serves it.

---

## Setup

1. Create the project in the main instance as usual (mock, proxy or forwarder) and note its alias, e.g. `web-app`.
2. Run a second instance on the **same PostgreSQL database**, with its **own config volume**, and set `MIRROR_PROJECT_ALIAS`:

```bash
docker run -d \
  --platform linux/amd64 \
  -p 8081:80 \
  -e DATABASE_URL=postgresql://username:password@host:5432/beo_echo \
  -e MIRROR_PROJECT_ALIAS=web-app \
  -v $(pwd)/beo-echo-mirror-config/:/app/configs/ \
  ghcr.io/yogasw/beo-echo:latest
```

3. Point the new domain at that container.

> **Use a separate config volume.** The instance rewrites `/app/configs/caddy/Caddyfile` on start. A volume shared with the main instance would give the main instance the mirror Caddyfile and hide its UI.

> **Use PostgreSQL.** With the default SQLite, each instance has its own database file, so the mirror would not see the project created in the main UI.

---

## Behavior

| Request on the mirror domain | Served as |
|---|---|
| `GET /` | project `web-app`, path `/` |
| `GET /assets/app.js` | project `web-app`, path `/assets/app.js` |
| `GET /f/abc123` | project `web-app`, path `/f/abc123` |
| `POST /api/v1/items` | project `web-app`, path `/api/v1/items` |
| `OPTIONS /anything` (CORS preflight) | answered with the project's CORS config |

- **No Beo Echo routes.** The UI, login, `/api` (management API, health check) and `/mcp` are not registered on a mirror instance.
- **Caddy.** The generated Caddyfile sends every path to the backend; no frontend files are served. Removing `MIRROR_PROJECT_ALIAS` restores the default Caddyfile on the next start.
- **Everything project-level still applies.** Mode (mock / proxy / forwarder), endpoints, responses, rules, delays, CORS mode and [Actions](Actions_Documentation.md) (custom headers, body rewrites) work exactly as on the main instance.
- **Logs.** Request logs are written to the shared database and show up in the project's Logs tab on the main instance.
- **Custom subdomains** (system config) are skipped on a mirror instance.

---

## Environment Variable

| Variable | Default | Description |
|---|---|---|
| `MIRROR_PROJECT_ALIAS` | empty | Alias of the project to serve from the domain root. Empty = normal instance. |
