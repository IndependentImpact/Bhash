# NGINX Content Negotiation for Hashgraph Ontology

This document explains how to configure NGINX to serve the Hashgraph Ontology
artefacts with HTTP content negotiation so that clients receive the serialisation
that best matches their `Accept` header.

## How Content Negotiation Works Here

Each ontology module is published in four representations:

| File | MIME type | Description |
|------|-----------|-------------|
| `<module>.html` | `text/html` | Human-readable catalogue (default) |
| `<module>.ttl`  | `text/turtle` | Turtle / TriG serialisation |
| `<module>.owl`  | `application/rdf+xml` | RDF/XML serialisation |
| `<module>.jsonld` | `application/ld+json` | JSON-LD serialisation |

When a client requests `/core` (no extension), NGINX maps the incoming `Accept`
header to the correct file extension and serves the matching file.  All other
clients (browsers, web crawlers) receive the HTML catalogue.

A `Vary: Accept` response header is included so HTTP caches store and serve
separate responses per format.

## Prerequisites

- NGINX ≥ 1.20 on Ubuntu/Debian
- Ontology artefacts deployed to `/var/www/hashgraphontology.xyz/`  
  (run `ontology/scripts/deploy_to_nginx.sh` to populate this directory)
- (Optional) A TLS certificate – see [Let's Encrypt / Certbot](#tlshttps-with-certbot)

## Quick-Start Installation

### 1 – Install NGINX

```bash
sudo apt update && sudo apt install -y nginx
```

### 2 – Deploy the ontology artefacts

```bash
# Install Python dependencies (rdflib, jinja2, etc.)
bash ontology/install_requirements.sh

# Convert source TTL files to all deployment formats
ontology/venv/bin/python ontology/scripts/convert_ontologies.py \
    --source-dir ontology/src \
    --deployment-dir ontology/deployment

# Generate the landing page
ontology/venv/bin/python ontology/scripts/convert_ontologies.py \
    --deployment-dir ontology/deployment \
    --generate-index

# Sync to the NGINX web root
PYTHON_BIN=ontology/venv/bin/python bash ontology/scripts/deploy_to_nginx.sh
```

### 3 – Install the NGINX configuration

```bash
sudo cp ontology/nginx/hashgraphontology.xyz.conf \
        /etc/nginx/sites-available/hashgraphontology.xyz.conf

sudo ln -s /etc/nginx/sites-available/hashgraphontology.xyz.conf \
           /etc/nginx/sites-enabled/hashgraphontology.xyz.conf

# Remove the default site if still present
sudo rm -f /etc/nginx/sites-enabled/default

# Verify configuration and reload
sudo nginx -t && sudo systemctl reload nginx
```

### 4 – Test content negotiation

```bash
# HTML (default – what a browser receives)
curl -s -I https://hashgraphontology.xyz/core

# Turtle
curl -s -H "Accept: text/turtle" https://hashgraphontology.xyz/core | head

# JSON-LD
curl -s -H "Accept: application/ld+json" https://hashgraphontology.xyz/core | head

# RDF/XML
curl -s -H "Accept: application/rdf+xml" https://hashgraphontology.xyz/core | head
```

## How the `map` Directive Works

The configuration file defines a `map` block **outside** the `server` block
(at the `http` context level):

```nginx
map $http_accept $negotiated_ext {
    default                     ".html";
    "~*text/turtle"             ".ttl";
    "~*application/ld\+json"    ".jsonld";
    "~*application/rdf\+xml"    ".owl";
    "~*application/owl\+xml"    ".owl";
}
```

Each ontology `location` block then uses `try_files` with the variable:

```nginx
location ~ ^/(core|consensus|token|smart-contracts|file-schedule|mirror-analytics|hiero)$ {
    add_header Vary Accept always;
    try_files $uri$negotiated_ext $uri.html =404;
}
```

NGINX evaluates `$negotiated_ext` per request, appends it to the URI, checks
whether that file exists in the `root`, and falls back to the `.html` file.

> **Note**: The `map` block must appear in the `http` context (i.e. inside
> `/etc/nginx/nginx.conf` or an `include`d file), **not** inside a `server`
> block.  If you place the configuration file in `sites-available/` it is
> already loaded at the `http` context level.

## Adding New Modules

When a new ontology module is added (e.g. `governance`):

1. Add the module name to the regex in `location`:

   ```nginx
   location ~ ^/(core|consensus|...|governance)$ {
   ```

2. Add the corresponding entry in `_MODULE_META` inside
   `ontology/scripts/convert_ontologies.py` so the landing page includes the
   new card.

3. Re-run the conversion and deployment pipeline (step 2 above).

## TLS/HTTPS with Certbot

```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d hashgraphontology.xyz
```

Certbot automatically edits the server block to handle HTTPS and issues a
cron job for certificate renewal.

## Troubleshooting

| Symptom | Likely cause | Fix |
|---------|-------------|-----|
| `curl -H "Accept: text/turtle" /core` returns HTML | `map` block is inside `server {}` not `http {}` | Move the `map` block to the outer `http` context |
| 404 for `/core.ttl` | Artefacts not deployed | Run `deploy_to_nginx.sh` |
| `nginx -t` reports unknown directive `map` | NGINX version too old | Upgrade to NGINX ≥ 1.7 |
| Vary header missing | `add_header` inside an `if` block | Ensure `add_header` is in the `location` block, not inside `if` |
