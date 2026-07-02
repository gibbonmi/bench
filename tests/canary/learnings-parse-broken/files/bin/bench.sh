#!/usr/bin/env bash
# Canary fixture driver: a minimal bench.sh sourcing ONLY the query module, so the
# AXI learnings contract exercises this fixture's bench-query.sh — whose open-heading
# regex has been mangled so real `[open]` entries are never matched.
set -euo pipefail
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/bench-query.sh"
case "${1:-}" in
  learnings) shift; learnings "$@" ;;
  *) echo "usage: bench {learnings}" ;;
esac
