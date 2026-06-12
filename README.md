isabella was here
i hate github it sucks im mad

# Plans for scaling + infra 

- Api + Scheduler for sims + Webapp Frontend + Caddy -> Hetzner VPS (anywhere from 12+ per month, scales manually)
  - 12/mo = CCX13, dedicated cores so sims cant starve the api
  - webapp needs sveltekit adapter-node to run behind caddy (adapter-auto wont pick it)
  - have caddy serve the api on the same origin as the app (app.mingas.com/api/* -> go container) and we never have to deal with CORS
- Product Marketing Website (Static HTML) -> Cloudflare Pages (Free)
- Database hosting + Management -> Either Hetzner or Supabase (more likely supabase), generous free tier, then 25/mo pro (realistically 35-75 with usage fees)
  - free tier pauses the db after 7 days of inactivity, fine for dev, upgrade to pro the day we have real users
  - direct connections are ipv6 only (ipv4 is a paid addon), hetzner vps has ipv6 so ok
  - remember to enable the postgis extension on the project
- Backups + object storage (sim outputs, 3d scans, nightly pg_dump even with managed db) -> Cloudflare R2 (free tier, zero egress fees)
- Uptime monitoring -> UptimeRobot pinging /healthz (free) or write my own uptime monitoring script (yippee)
- Error tracking -> Sentry free tier (go + sveltekit)
- CI/CD -> GitHub Actions, frontends auto deploy via Pages, backend builds image -> GHCR -> ssh to vps, compose pull && up -d
### Other costs
- Domain names -> 12 bucks per year, cloudflare registrar sells at cost ~10.50 (mingas.com can be used for product website, app.mingas.com, admin.mingas.com, subdomains are free)
