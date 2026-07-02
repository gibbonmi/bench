#!/usr/bin/env bash
# Canary fixture driver: a minimal bench.sh that sources ONLY the query module, so
# the AXI CLI contracts exercise this fixture's bench-query.sh — whose TOON emitter
# has had its escaping removed. Real dispatch, gate, link, etc. are out of scope here.
set -euo pipefail
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
. "$DIR/bench-query.sh"
case "${1:-}" in
  learnings) shift; learnings "$@" ;;
  *) echo "usage: bench {learnings}" ;;
esac
