#!/usr/bin/env bash
# Install Python dependencies for ontology conversion utilities.
set -euo pipefail

PYTHON_BIN=${PYTHON_BIN:-python3.12}

if ! command -v "$PYTHON_BIN" >/dev/null 2>&1; then
    echo "❌ Python binary not found: $PYTHON_BIN" >&2
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VENV_DIR="$SCRIPT_DIR/venv"

if [ -d "$VENV_DIR" ]; then
    echo "🔁 Reusing existing virtual environment at $VENV_DIR"
else
    echo "🐍 Creating virtual environment with $PYTHON_BIN at $VENV_DIR"
    "$PYTHON_BIN" -m venv "$VENV_DIR"
fi

# shellcheck disable=SC1091
source "$VENV_DIR/bin/activate"

echo "⬆️  Upgrading pip..."
python -m pip install --upgrade pip

echo "📦 Installing conversion dependencies..."
pip install rdflib jinja2 pyshacl owlrl

echo "✅ Environment ready. Activate it later with: source $VENV_DIR/bin/activate"
