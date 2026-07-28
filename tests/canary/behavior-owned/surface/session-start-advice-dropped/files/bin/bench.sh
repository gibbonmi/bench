#!/usr/bin/env bash
# Canary fixture support: a bench.sh the fixture's session-start.sh can resolve by
# path, so the advisory branch is reached. The planted regression is in
# dot-bench/hooks/session-start.sh (the dropped doctor pointer).
set -uo pipefail
case "${1:-}" in
  status) : ;;
  *) : ;;
esac
