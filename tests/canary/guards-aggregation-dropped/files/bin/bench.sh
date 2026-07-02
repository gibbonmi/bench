#!/usr/bin/env bash
# Canary fixture: a bench.sh whose `guards` aggregation DROPS the pre-push row — it
# lists the four hook/adapter guards but never the git-layer pre-push guard. The AXI
# guards aggregation contract must go red with "AXI guards aggregation contract failed"
# (its guards[5] header assertion catches the missing fifth row), proving that check
# still bites if the aggregator stops reporting a guard. `link` and every other
# subcommand are harmless no-ops so the contract's setup step reaches the assertion.
set -euo pipefail
case "${1:-}" in
  guards)
    cat <<'TOON'
guards[4]{guard,boundary,denies}:
  block-dangerous-git,PreToolUse:Bash,destructive git
  check-agent-line,PreToolUse:Agent,unbound model
  stop,Stop,done on red
  _line-guard,sourced,unbound model
TOON
    ;;
  *) echo "ok" ;;
esac
