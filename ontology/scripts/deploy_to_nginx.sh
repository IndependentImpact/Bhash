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

if [ ! -d "$DEPLOYMENT_DIR" ]; then
    echo "❌ Deployment directory not found: $DEPLOYMENT_DIR" >&2
    exit 1
fi

if [ "${USE_SUDO,,}" = "false" ]; then
    SUDO_CMD=()
else
    SUDO_CMD=(sudo)
fi

set -x
"${SUDO_CMD[@]}" "$RSYNC_BIN" -av --delete "$DEPLOYMENT_DIR"/ "$DEST"
