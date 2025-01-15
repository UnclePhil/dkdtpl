# Caddy Ingress stack with signed local CA
## purpose
offer a global ingress stacke deployed with caddyserver and some plugins

Variant with signed certificate for the local CA

## source
- caddyserver.com
- flaviostuz caddy-docker-proxy 
- and my own version of flavio dockerfile


## prometheus query
rate(caddy_admin_http_requests_total[30s])

rate(caddy_http_requests_total[30s])
rate(caddy_http_request_errors_total[30s])