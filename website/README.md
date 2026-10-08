# VERTEX/ctxbank website

The public site for CTXbank — landing page, setup guide, and the Connect page
that opens a user's **local** dashboard. Pure static files: no build step, no
backend, no tracking.

## What's here

- `index.html` — one page, three tabs (`#how`, `#setup`, `#connect`)
- `styles.css` — Soft Harbor design tokens (light default, dark toggle)
- `app.js` — tabs, copy buttons, theme, Connect-page localhost probe

## The Connect model

The site never sees project data. The Connect tab asks for the user's project
token (from `ctx token`) and probes `http://localhost:<port>/api/status` with
it. The local `ctx ui` server answers only with a valid token (CORS is enabled
on `/api/*` for exactly this). On success, "Open dashboard" launches the local
dashboard in a new tab. If `ctx ui` isn't running, the user gets a plain-English
nudge to start it.

## Deploy

Serve the directory as-is from any static host. Intended home:
`https://ctxbank.vertexdevstudio.tech` (Cloudflare Pages / Vercel / any static
host — point the subdomain at the deployed directory).

```bash
# quick local preview
cd website && python3 -m http.server 8080
# then open http://localhost:8080
```

## Design

Built to the VERTEX Soft Harbor design language
(`~/workspace/design-language/DESIGN_LANGUAGE.md`): flat colour, calm motion,
sentence case, one idea per section. Light is the default; dark is an opt-in
toggle persisted in `localStorage`.
