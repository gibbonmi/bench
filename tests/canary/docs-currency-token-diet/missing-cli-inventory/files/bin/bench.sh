#!/usr/bin/env bash
# Canary fixture CLI: implements `gate`, while BENCH.md omits it.
case "${1:-}" in
  gate) echo ok ;;
esac
