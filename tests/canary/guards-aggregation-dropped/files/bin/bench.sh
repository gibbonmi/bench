#!/usr/bin/env bash
# Canary fixture: a bench.sh whose `guards` aggregation DROPS the pre-push row — it
# lists the three hook guards but never the git-layer pre-push guard, and miscounts
# its own header (guards[3] where a linked fixture must aggregate four:
# block-dangerous-git, check-agent-line, stop, pre-push). The AXI guards aggregation
# contract must go red with "AXI guards aggregation contract failed" (its guards[4]
# header assertion catches the count, and the row loop catches the missing pre-push),
# proving that check still bites if the aggregator stops reporting a guard. `link` and
# every other subcommand are harmless no-ops so the contract's setup step reaches the
# assertion.
set -euo pipefail
case "${1:-}" in
  guards)
    cat <<'TOON'
guards[3]{guard,boundary,denies}:
  block-dangerous-git,PreToolUse:Bash,destructive git
  check-agent-line,PreToolUse:Agent,unbound model
  stop,Stop,done on red
TOON
    ;;
  *) echo "ok" ;;
esac
