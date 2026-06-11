# Deploying a New Instance of Quiniela Mundial

## Context

This is a deployment guide for spinning up a second instance of the **Quiniela Mundial**
app (World Cup 2026 prediction game, Spanish UI) on a new server — same `mundial.*`
subdomain pattern, different domain, routed through Cloudflare Tunnel (cloudflared).

Source: **https://github.com/cryptze/wm-pickems**

The app is a single Docker container: a Go binary (PocketBase) serving the API and
a compiled SvelteKit SPA on one port (default `8090`). Data lives in a Docker volume.

---

## Prerequisites on the target server

- Docker + Docker Compose v2
- `cloudflared` installed and authenticated to the Cloudflare account that owns the target domain
- An existing Cloudflare Tunnel (or create a new one) connected to this server
- Git

---

## Step 1 — Clone the repo

```sh
git clone https://github.com/cryptze/wm-pickems.git /opt/mundial
cd /opt/mundial
```

---

## Step 2 — Configure environment

```sh
cp .env.example .env
```

Edit `.env`:

```env
HTTP_PORT=8090

# Optional: paid API-Football key for real-time results (free tier has no WC2026)
API_FOOTBALL_KEY=

# PocketBase admin superuser (created automatically on first boot if set)
PB_ADMIN_EMAIL=your@email.com
PB_ADMIN_PASSWORD=a-strong-password

# Google OAuth (optional — leave blank to disable)
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
```

> If Google OAuth is used, add `https://mundial.<YOUR_DOMAIN>/api/oauth2-redirect`
> as an Authorized Redirect URI in the Google Cloud Console.

---

## Step 3 — Build and start

```sh
docker compose up --build -d
```

Verify it's healthy:

```sh
docker ps --filter "name=fhun_tips"
curl http://localhost:8090/api/health
```

First boot auto-runs migrations and seeds 48 teams / 12 groups / 104 fixtures.

---

## Step 4 — Configure Cloudflare Tunnel

Edit the cloudflared config file (typically `/etc/cloudflared/config.yml` or
`~/.cloudflared/config.yml`):

```yaml
tunnel: <YOUR_TUNNEL_ID>
credentials-file: /etc/cloudflared/<YOUR_TUNNEL_ID>.json

ingress:
  # ... your existing ingress rules ...

  - hostname: mundial.<YOUR_DOMAIN>        # e.g. mundial.miempresa.com
    service: http://localhost:8090

  - service: http_status:404
```

Then restart cloudflared:

```sh
# systemd (most common)
sudo systemctl restart cloudflared

# or if running as a user service
systemctl --user restart cloudflared
```

Make sure the DNS record `mundial.<YOUR_DOMAIN>` → your tunnel CNAME exists in the
Cloudflare dashboard (Cloudflare Tunnel → Public Hostnames, or DNS tab).

---

## Step 5 — Create an admin account

```sh
docker compose exec app wm-pickems superuser create you@example.com 'a-strong-pass' --dir=/pb_data
```

The PocketBase admin panel is at `https://mundial.<YOUR_DOMAIN>/_/`.

---

## Operating

| Task | Command |
|------|---------|
| View logs | `docker compose logs -f` |
| Force result sync | `POST https://mundial.<YOUR_DOMAIN>/api/sync/refresh` (superuser auth) |
| Manual match result | `POST /api/admin/matches/{id}/result` with `{"FTHome":2,"FTAway":1,"Status":"finished"}` |
| Recompute all scores | `POST /api/admin/recompute` (superuser auth) |
| Restart | `docker compose restart` |
| Update | `git pull && docker compose up --build -d` |

---

## Backup

```sh
docker run --rm \
  -v mundial_pb_data:/d \
  -v "$PWD":/b \
  alpine tar czf /b/pb_data-backup.tgz -C /d .
```

Restore: extract the archive back into the volume before `docker compose up`.

---

## Notes

- The Docker volume is named `mundial_pb_data` (set by docker-compose).
- The container is named `fhun_tips` (legacy name kept for consistency).
- Health endpoint: `GET /api/health` returns 200 when running.
- Results sync automatically every 30 min from openfootball (free, community-updated)
  or API-Football if a paid key is configured.
