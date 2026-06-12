# Quiniela Mundial — Guía de Mantenimiento

## Información general

| Item | Valor |
|------|-------|
| URL | https://mundial.maryteruya.com |
| Admin panel | https://mundial.maryteruya.com/_/ |
| Admin email | andresr@maryvisa.com |
| Carpeta del proyecto | `/root/projects/mundial` |
| Contenedor | `fhun_tips` |
| Puerto interno | `127.0.0.1:8090` |
| Red Docker | `n8n-produccion_n8n-network` + `wmpickems` |

---

## Arquitectura

```
Internet → Cloudflare Tunnel (tunnel: 138d5cd5-103c-4dab-a2df-bde02cc2e856)
         → fhun_tips:8090
         → PocketBase (Go binary) sirve API + SvelteKit SPA
         → Volumen Docker: mundial_pb_data
```

---

## Comandos cotidianos

```bash
cd /root/projects/mundial

# Ver estado
docker compose ps

# Ver logs en vivo
docker compose logs -f

# Reiniciar (no recarga .env)
docker compose restart

# Recrear con nuevas variables de entorno (sí recarga .env)
docker compose up -d

# Crear/actualizar superuser si se pierde acceso al panel
docker exec fhun_tips wm-pickems superuser upsert andresr@maryvisa.com 'NUEVO_PASS' --dir=/pb_data
```

---

## Variables de entorno (.env)

El archivo `/root/projects/mundial/.env` contiene:

```
HTTP_PORT=8090
PB_ADMIN_EMAIL=andresr@maryvisa.com
PB_ADMIN_PASSWORD=***
GOOGLE_CLIENT_ID=163710230882-s4cn10mmpgg1a6hps3iov27sd3vras6b.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=***
API_FOOTBALL_KEY=          # vacío = usa worldcup26.ir (gratis, sin clave)
RESULTS_SOURCE=auto        # auto | worldcup26ir | openfootball | apifootball
```

> **Importante:** después de editar `.env` hay que usar `docker compose up -d` (no `restart`) para que el contenedor tome los nuevos valores.

---

## Actualizar la app

```bash
cd /root/projects/mundial
git pull
docker compose up --build -d
```

---

## Cloudflare Tunnel

La ruta `mundial.maryteruya.com` está en la configuración **remota** del tunnel (no solo en el archivo local). Para modificarla hay que actualizar via API o desde el dashboard de Cloudflare Zero Trust.

El archivo local `/opt/n8n-produccion/cloudflared/config.yml` se mantiene sincronizado como referencia pero **no es el que usa cloudflared** — usa la config remota (version 6).

---

## Backup de datos

```bash
# Crear backup del volumen
docker run --rm \
  -v mundial_pb_data:/d \
  -v /root/projects/mundial:/b \
  alpine tar czf /b/pb_data-backup-$(date +%Y%m%d).tgz -C /d .

# Restaurar (con el contenedor detenido)
docker compose down
docker run --rm \
  -v mundial_pb_data:/d \
  -v /root/projects/mundial:/b \
  alpine sh -c "cd /d && tar xzf /b/pb_data-backup-FECHA.tgz"
docker compose up -d
```

---

## Google OAuth

Configurado automáticamente al arrancar desde las variables `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET`.

- Redirect URI registrado en Google Cloud Console: `https://mundial.maryteruya.com/api/oauth2-redirect`
- Si se agregan nuevas URIs en Google Console, no requiere cambios en el servidor.

---

## Sincronización de resultados

- Se ejecuta automáticamente cada 30 min desde **worldcup26.ir** (fuente gratuita, sin clave de API).
- Fallback manual disponible via `RESULTS_SOURCE=openfootball` o `apifootball` en el `.env`.
- Para forzar una sincronización manual:

```bash
curl -X POST https://mundial.maryteruya.com/api/sync/refresh \
  -H "Authorization: Bearer TOKEN_SUPERUSER"
```

- Para ingresar resultado manualmente desde el panel admin: **Collections → matches** → editar el partido.
