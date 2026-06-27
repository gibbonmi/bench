#!/usr/bin/env bash
# PreToolUse guard: the agent has no destructive git authority. This makes
# invariant #4 ("you assist, you don't decide where a decision is mine")
# enforceable for the operations that can silently destroy a shift's work or
# bypass the merge — not just aspirational.
#
# Note the boundary: this intercepts the AGENT's Bash tool calls. Bench's own
# controlled rollback inside `bench shift` runs in-process (not through the
# agent's shell), so the harness can still reset/clean a failed iteration while
# the agent itself cannot reach for those commands. That asymmetry is the point.
#
# Wire under hooks.PreToolUse with matcher "Bash". Exit 2 blocks and returns the
# message to the agent.
set -euo pipefail

input="$(cat)"
cmd="$(printf '%s' "$input" | python3 -c 'import sys,json; print(json.load(sys.stdin).get("tool_input",{}).get("command",""))' 2>/dev/null || true)"
[[ -z "$cmd" ]] && exit 0

block() { echo "BLOCKED: \`$1\` — you don't have authority over this. The merge and any history rewrite are mine; a failed shift is rolled back by bench, not by you. Stop and hand back." >&2; exit 2; }

case "$cmd" in
  *"git push"*)                       block "git push" ;;
  *"git reset --hard"*)               block "git reset --hard" ;;
  *"git clean -f"*|*"git clean -df"*|*"git clean -xf"*) block "git clean -f" ;;
  *"git branch -D"*|*"git branch -d "*) block "git branch -D" ;;
  *"git checkout ."*|*"git restore ."*) block "git checkout ." ;;
  *"git rebase"*|*"git push --force"*|*"git push -f"*) block "history rewrite" ;;
esac
exit 0
