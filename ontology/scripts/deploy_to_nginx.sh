#!/usr/bin/env bash
# Synchronise generated ontology artefacts to an nginx web root.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOYMENT_DIR=${DEPLOYMENT_DIR:-"$SCRIPT_DIR/../deployment"}
SITE_NAME=${1:-hashgraphontology.xyz}
TARGET_ROOT=${TARGET_ROOT:-/var/www}
DEST="$TARGET_ROOT/$SITE_NAME/"
USE_SUDO=${USE_SUDO:-true}
RSYNC_BIN=${RSYNC_BIN:-rsync}
WEB_USER=${WEB_USER:-www-data}
WEB_GROUP=${WEB_GROUP:-www-data}
INDEX_FILE="$DEPLOYMENT_DIR/index.html"

if [ ! -d "$DEPLOYMENT_DIR" ]; then
    echo "❌ Deployment directory not found: $DEPLOYMENT_DIR" >&2
    exit 1
fi

if [ "${USE_SUDO,,}" = "false" ]; then
    SUDO_CMD=()
else
    SUDO_CMD=(sudo)
fi

mapfile -t ARTIFACTS < <(cd "$DEPLOYMENT_DIR" && find . -maxdepth 1 -type f ! -name "index.html" -printf '%f\n' | sort)
mapfile -t DIRECTORIES < <(cd "$DEPLOYMENT_DIR" && find . -maxdepth 1 -mindepth 1 -type d -printf '%f/\n' | sort)

cat >"$INDEX_FILE" <<'HTML'
<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Hashgraph Ontology Artefacts</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 2rem; line-height: 1.6; background: #f9fafb; color: #1f2933; }
        h1 { color: #0f172a; }
        .section { margin-bottom: 1.5rem; }
        ul { list-style: disc inside; padding-left: 0; }
        a { color: #0ea5e9; text-decoration: none; }
        a:hover { text-decoration: underline; }
        .note { color: #52616b; }
    </style>
</head>
<body>
    <h1>Hashgraph Ontology</h1>
    <p class="note">Browse the published ontology artefacts and documentation below.</p>
    <div class="section">
        <h2>Artefacts</h2>
        <ul>
HTML

    for entry in "${ARTIFACTS[@]}"; do
        printf '            <li><a href="%s">%s</a></li>\n' "$entry" "$entry" >>"$INDEX_FILE"
    done

cat >>"$INDEX_FILE" <<'HTML'
        </ul>
    </div>
    <div class="section">
        <h2>Additional Directories</h2>
        <ul>
HTML

if [ "${#DIRECTORIES[@]}" -eq 0 ]; then
    echo "            <li class=\"note\">No additional directories</li>" >>"$INDEX_FILE"
else
    for dir in "${DIRECTORIES[@]}"; do
        printf '            <li><a href="%s">%s</a></li>\n' "$dir" "$dir" >>"$INDEX_FILE"
    done
fi

cat >>"$INDEX_FILE" <<'HTML'
        </ul>
    </div>
    <p class="note">For details on each artefact, open the HTML documentation or download the RDF representations (TTL, OWL, JSON-LD).</p>
</body>
</html>
HTML

set -x
"${SUDO_CMD[@]}" "$RSYNC_BIN" -av --delete "$DEPLOYMENT_DIR"/ "$DEST"
"${SUDO_CMD[@]}" chown -R "$WEB_USER":"$WEB_GROUP" "$DEST"
