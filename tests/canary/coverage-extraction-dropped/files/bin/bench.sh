#!/usr/bin/env bash
# Canary fixture: a bench.sh whose `coverage` claims a mapped state but DROPS
# every extracted row (rows[0] on a two-row map). The wave-2 AXI coverage
# extraction contract must go red with "AXI coverage extraction contract failed"
# (its rows[2] header assertion catches the thinned table), proving the check
# still bites if extraction stops surfacing data rows. Every other subcommand is
# a harmless no-op so setup steps reach the assertion.
set -euo pipefail
case "${1:-}" in
  coverage)
    printf 'spec: %s\n' "${2:-}"
    printf 'state: mapped\n'
    printf 'rows[0]{story,seam,red_signal}:\n'
    ;;
  *) echo "ok" ;;
esac
