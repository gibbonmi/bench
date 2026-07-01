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

reason="$(python3 - "$cmd" <<'PY'
import os
import shlex
import sys

try:
    tokens = shlex.split(sys.argv[1], posix=True)
except ValueError:
    tokens = sys.argv[1].split()

opts_with_arg = {
    "-C",
    "-c",
    "--exec-path",
    "--git-dir",
    "--namespace",
    "--work-tree",
}


def git_invocations():
    for i, token in enumerate(tokens):
        if os.path.basename(token) != "git":
            continue
        j = i + 1
        while j < len(tokens):
            current = tokens[j]
            if current == "--":
                j += 1
                break
            if current in opts_with_arg:
                j += 2
                continue
            if any(current.startswith(f"{opt}=") for opt in opts_with_arg if opt.startswith("--")):
                j += 1
                continue
            if current.startswith("-"):
                j += 1
                continue
            break
        if j < len(tokens):
            yield tokens[j], tokens[j + 1 :]


for subcommand, args in git_invocations():
    if subcommand == "push":
        print("git push")
        break
    if subcommand == "reset" and "--hard" in args:
        print("git reset --hard")
        break
    if subcommand == "clean" and (
        "--force" in args or any(arg.startswith("-") and "f" in arg[1:] for arg in args)
    ):
        print("git clean -f")
        break
    if subcommand == "branch" and any(arg in {"-D", "-d", "--delete"} for arg in args):
        print("git branch -D")
        break
    if subcommand in {"checkout", "restore"} and "." in args:
        print(f"git {subcommand} .")
        break
    if subcommand == "rebase":
        print("history rewrite")
        break
PY
)"
[[ -n "$reason" ]] && block "$reason"
exit 0
