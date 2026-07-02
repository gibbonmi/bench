#!/usr/bin/env bash
# Canary fixture: a bench.sh whose `diff` IGNORES the recorded benchBase key and
# always answers with a merge-base resolution. The wave-2 AXI recorded-base
# contract must go red with "AXI diff recorded-base contract failed" (its exact
# `base:`/`method: recorded` assertions catch the dropped preference), proving
# the check still bites if the resolution ladder loses the recorded rung. Every
# other subcommand is a harmless no-op so setup steps reach the assertion.
set -euo pipefail
case "${1:-}" in
  diff)
    printf 'branch: whatever\n'
    printf 'base: 0000000000000000000000000000000000000000\n'
    printf 'method: merge-base\n'
    printf 'files[0]{status,path}:\n'
    ;;
  *) echo "ok" ;;
esac
