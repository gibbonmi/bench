#!/usr/bin/env bash
# Pre-fix init: scaffolds the gate but forgets the learnings journal.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
mkdir -p "$root/.bench"
: > "$root/.bench/gate.sh"
