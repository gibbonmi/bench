#!/usr/bin/env bash
# Canary fixture CLI: implements only `gate`, while the fixture's HANDOFF.md
# also documents a `bench frobnicate` that no longer exists.
case "${1:-}" in
  gate) echo ok ;;
esac
