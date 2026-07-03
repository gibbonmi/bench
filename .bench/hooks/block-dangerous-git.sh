#!/usr/bin/env bash
# PreToolUse guard: the agent has no destructive git authority. This makes
# invariant #4 ("you assist, you don't decide where a decision is mine")
# enforceable for the operations that can silently destroy a shift's work or
# bypass the merge — not just aspirational.
#
# Threat model: this is an honest-mistake layer, not an evasion-resistant
# boundary. It stops a well-meaning agent from reflexively running destructive
# git; it does not try to survive deliberate evasion. Wrapper scanning goes
# exactly one level deep (a `sh -c`/`bash -c`/`zsh -c` string) by design — a
# wrapper found inside that string is not re-expanded. The backstops for a
# misaligned agent are the git pre-push hook and bench's pooled-worktree
# isolation, not this script.
#
# Note the boundary: this intercepts the AGENT's Bash tool calls. Bench's own
# controlled rollback inside `bench shift` runs in-process (not through the
# agent's shell), so the harness can still reset/clean a failed iteration while
# the agent itself cannot reach for those commands. That asymmetry is the point.
#
# Wire under hooks.PreToolUse with matcher "Bash". Exit 2 blocks and returns the
# message to the agent. `--describe` (first arg) prints the guard manifest and
# exits 0 without reading stdin, so `bench guards` can aggregate the deny surface
# without a collision; the denies clause is generated from the same DENY_LABELS
# table the analyzer denies from, so the advertisement cannot drift.
set -euo pipefail

# The analyzer program: the sibling git-guard.py, a single source used both to
# classify a command (deny) and to enumerate the deny classes (--describe).
# Resolved relative to this script so the pair travels together through the
# whole-directory hook installs — by parameter expansion only, because the
# python3-missing degradation path must still resolve it with no external
# tools on PATH. Absent or EMPTY, it must fail closed — an empty program exits
# 0 printing nothing, which the allow path below would misread as "no verdict".
hook_dir="${BASH_SOURCE[0]%/*}"
[[ "$hook_dir" == "${BASH_SOURCE[0]}" ]] && hook_dir=.
GUARD_PY="$hook_dir/git-guard.py"

if [[ "${1:-}" == "--describe" ]]; then
  printf 'name: block-dangerous-git\n'
  printf 'boundary: PreToolUse:Bash\n'
  if [[ ! -s "$GUARD_PY" ]]; then
    printf 'denies: manifest unavailable (analyzer missing)\n'
  elif command -v python3 >/dev/null 2>&1; then
    printf 'denies: destructive git — %s\n' "$(python3 "$GUARD_PY" --describe-classes)"
  else
    printf 'denies: manifest unavailable (python3 missing)\n'
  fi
  printf 'why: agents lack destructive-git authority; merge and history rewrites belong to the reviewer\n'
  exit 0
fi

input="$(cat)"

# Without python3 the hook can neither parse the command envelope nor run the
# analyzer, so it cannot classify at all. Fail closed on anything git-shaped —
# the destructive surface — while leaving non-git commands runnable so the shell
# stays usable. The raw envelope carries the command text, so a substring test
# catches the honest mistake without parsing. This mirrors the analyzer-missing
# branch below: same "cannot classify" condition, same deny verdict.
if ! command -v python3 >/dev/null 2>&1; then
  case "$input" in
    *git*) echo "BLOCKED: guard degraded (python3 missing) — can't classify commands, refusing anything git-shaped. Install python3 or hand back." >&2; exit 2 ;;
    *) exit 0 ;;
  esac
fi

cmd="$(printf '%s' "$input" | python3 -c 'import sys,json; print(json.load(sys.stdin).get("tool_input",{}).get("command",""))' 2>/dev/null || true)"
[[ -z "$cmd" ]] && exit 0

block() { echo "BLOCKED: \`$1\` — you don't have authority over this. The merge and any history rewrite are the user's; a failed shift is rolled back by bench, not by you. Stop and hand back." >&2; exit 2; }

[[ -s "$GUARD_PY" ]] || { echo "BLOCKED: guard analyzer missing — failing closed; restore .bench/hooks/git-guard.py (bench link) before retrying." >&2; exit 2; }

reason="$(python3 "$GUARD_PY" "$cmd")" || { echo "BLOCKED: guard analyzer error — failing closed; rephrase the command." >&2; exit 2; }
[[ -n "$reason" ]] && block "$reason"
exit 0
