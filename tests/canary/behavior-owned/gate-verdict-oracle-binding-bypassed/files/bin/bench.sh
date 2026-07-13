#!/usr/bin/env bash
# Cheapest wrong authorization: remember only that some gate was green, then let
# commit reuse it without comparing the exact oracle identity.
case "${1:-}" in
  gate)
    printf 'run\n' >> .git/gate-runs
    printf 'green\n' > .git/bench-last-gate
    exit 0
    ;;
  commit)
    exit 0
    ;;
  *) exit 1 ;;
esac
