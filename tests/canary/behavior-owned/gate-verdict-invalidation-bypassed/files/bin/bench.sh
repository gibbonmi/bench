#!/usr/bin/env bash
# Cheapest wrong persistence order: run the oracle directly before installing
# pending, allowing an older green to remain authoritative during execution.
case "${1:-}" in
  gate)
    exec .bench/gate.sh
    ;;
  *) exit 1 ;;
esac
