
# 📦 Setup Docker + Traefik + Cloudflare Tunnel (Tanpa IP Publik)

## 🧩 Tujuan
Menyediakan banyak subdomain (misal: `app.domain.com`, `api.domain.com`, dll) untuk aplikasi yang berjalan di dalam Docker, **tanpa memerlukan IP publik**, menggunakan **Cloudflare Tunnel** dan **Traefik**.

---

## 🔧 Komponen yang Digunakan

- **Docker**: Menjalankan semua service dan reverse proxy.
- **Traefik**: Reverse proxy yang menangani routing berdasarkan Host (subdomain).
- **Cloudflare Tunnel (cloudflared)**: Untuk expose service Docker ke internet melalui jaringan Cloudflare.
- **Cloudflare DNS**: Untuk mapping subdomain ke tunnel endpoint.

---

## 🏗️ Struktur Arsitektur

```text
Internet
   ↓
Cloudflare DNS (CNAME subdomain ke tunnel)
   ↓
Cloudflare Tunnel (cloudflared)
   ↓
Docker Host (localhost)
   ↓
Traefik Reverse Proxy (Docker)
   ↓
Service-service dalam Docker
```

---

## 📄 File `docker-compose.yml`

```yaml
version: '3.9'
services:
  traefik:
    image: traefik:v2.11
    command:
      - --entrypoints.web.address=:80
      - --providers.docker=true
      - --api.dashboard=true
    ports:
      - "8080:80"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock

  webapp:
    image: your-webapp-image
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.webapp.rule=Host(`app.domain.com`)"
      - "traefik.http.services.webapp.loadbalancer.server.port=80"

  api:
    image: your-api-image
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.api.rule=Host(`api.domain.com`)"
      - "traefik.http.services.api.loadbalancer.server.port=5000"
```

---

## 🌐 File `cloudflared/config.yml`

```yaml
tunnel: mytunnel-id
credentials-file: /root/.cloudflared/mytunnel-id.json

ingress:
  - service: http://localhost:8080
```

> Semua request diserahkan ke Traefik yang berjalan di port 8080 (dalam Docker)

---

## 🌍 DNS Setting di Cloudflare

| Subdomain            | Type  | Value                          |
|----------------------|-------|--------------------------------|
| `app.domain.com`     | CNAME | `mytunnel.cfargotunnel.com`    |
| `api.domain.com`     | CNAME | `mytunnel.cfargotunnel.com`    |
| `*.domain.com`       | CNAME | `mytunnel.cfargotunnel.com`    |

> Kamu bisa pakai `*` untuk wildcard agar otomatis.

---

## ✅ Kesimpulan

- Kamu **tidak butuh IP publik** untuk expose service Docker kamu ke internet.
- Cukup **satu tunnel (`mytunnel-id`)** untuk semua subdomain.
- Gunakan **Traefik** untuk mengatur routing subdomain ke service berdasarkan label Docker.
- Semua subdomain cukup diarahkan ke satu endpoint `mytunnel.cfargotunnel.com` melalui DNS Cloudflare.
- Tambah service baru cukup dengan nambah container & label Traefik saja — **tanpa perlu edit cloudflared atau DNS lagi.**

---

## 🗂️ Referensi Tambahan

- [Traefik Docs](https://doc.traefik.io/traefik/)
- [Cloudflare Tunnel Docs](https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/)
