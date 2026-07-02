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
# message to the agent.
set -euo pipefail

input="$(cat)"
cmd="$(printf '%s' "$input" | python3 -c 'import sys,json; print(json.load(sys.stdin).get("tool_input",{}).get("command",""))' 2>/dev/null || true)"
[[ -z "$cmd" ]] && exit 0

block() { echo "BLOCKED: \`$1\` — you don't have authority over this. The merge and any history rewrite are mine; a failed shift is rolled back by bench, not by you. Stop and hand back." >&2; exit 2; }

reason="$(python3 - "$cmd" <<'PY'
import os
import posixpath
import re
import shlex
import subprocess
import sys

CONTROL_OPS = {";", "&&", "||", "|", "&", "(", ")"}
KEYWORDS = {"if", "then", "elif", "else", "do", "while", "until", "!", "{"}
WRAPPERS = {"sh", "bash", "zsh"}

GLOBAL_OPTS_WITH_ARG = {
    "-C",
    "-c",
    "--exec-path",
    "--git-dir",
    "--namespace",
    "--work-tree",
}


REDIRECT_RE = re.compile(r"^(?:[0-9]+)?(?:>>?|<<?<?)(?:[|&])?$|^&>>?$")


def strip_redirections(tokens):
    out = []
    i = 0
    while i < len(tokens):
        if REDIRECT_RE.match(tokens[i]):
            if out and out[-1].isdigit():
                out.pop()
            i += 2
            continue
        out.append(tokens[i])
        i += 1
    return out


def tokenize(s):
    try:
        lexer = shlex.shlex(s, posix=True, punctuation_chars=True)
        lexer.whitespace_split = True
        return strip_redirections(list(lexer))
    except ValueError:
        return strip_redirections(s.split())


def has_explicit_pathspec(args):
    if "--" not in args:
        return False
    return any(arg != "" for arg in args[args.index("--") + 1 :])


def is_pathspec_file_arg(arg):
    return arg == "--pathspec-from-file" or arg.startswith("--pathspec-from-file=")


def has_pathspec_file(args):
    return any(is_pathspec_file_arg(arg) for arg in args)


def restore_has_pathspec(args):
    opts_with_arg = {"-s", "--source"}
    i = 0
    while i < len(args):
        arg = args[i]
        if arg == "--":
            return any(rest != "" for rest in args[i + 1 :])
        if is_pathspec_file_arg(arg):
            return True
        if arg in opts_with_arg:
            i += 2
            continue
        if any(arg.startswith(f"{opt}=") for opt in opts_with_arg if opt.startswith("--")):
            i += 1
            continue
        if arg.startswith("-"):
            i += 1
            continue
        return True
    return False


def ref_resolves(arg):
    try:
        result = subprocess.run(
            ["git", "rev-parse", "--verify", "--quiet", f"{arg}^{{commit}}"],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
            timeout=2,
        )
        return result.returncode == 0
    except Exception:
        return False


def branch_exists(name):
    try:
        result = subprocess.run(
            ["git", "rev-parse", "--verify", "--quiet", f"refs/heads/{name}"],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
            timeout=2,
        )
        return result.returncode == 0
    except Exception:
        return True


def forced_creation_target(args, opt):
    for i, arg in enumerate(args):
        if arg == opt:
            return args[i + 1] if i + 1 < len(args) else None
        if arg.startswith(opt) and len(arg) > len(opt):
            return arg[len(opt):]
    return None


def checkout_free_args(args):
    opts_with_value = {"-b", "-B", "--orphan"}
    free = []
    i = 0
    while i < len(args):
        arg = args[i]
        if arg in opts_with_value:
            i += 2
            continue
        if arg.startswith("-"):
            i += 1
            continue
        free.append(arg)
        i += 1
    return free


def checkout_verdict(args, via_xargs):
    if via_xargs:
        return True
    if any(arg in {"-f", "--force"} for arg in args):
        return True
    if has_pathspec_file(args):
        return True
    target = forced_creation_target(args, "-B")
    if target is not None and branch_exists(target):
        return True
    if "--" in args:
        return has_explicit_pathspec(args)
    free = checkout_free_args(args)
    if not free:
        return False
    if len(free) > 1:
        return True
    return not ref_resolves(free[0])


def switch_verdict(args):
    if any(arg in {"-f", "--force", "--discard-changes"} for arg in args):
        return True
    target = forced_creation_target(args, "-C")
    return target is not None and branch_exists(target)


def restore_verdict(args, via_xargs):
    if via_xargs:
        return True
    staged = any(arg in {"--staged", "-S"} for arg in args)
    worktree = any(arg in {"--worktree", "-W"} for arg in args)
    if staged and not worktree:
        return False
    return "." in args or restore_has_pathspec(args)


def stash_verdict(args):
    sub = next((arg for arg in args if not arg.startswith("-")), None)
    return sub in {"drop", "clear"}


def is_delegate_worktree_path(arg):
    norm = posixpath.normpath(arg)
    return norm.startswith(".claude/worktrees/") or "/.claude/worktrees/" in norm


def branch_verdict(args):
    if not any(arg in {"-D", "-d", "--delete", "-f", "--force"} for arg in args):
        return False
    # Carve-out: deleting harness-delegate branches (worktree-*) is cleanup of
    # agent-created scratch, not reviewer history. Force-move (-f without a
    # delete flag) and any non-delegate name still block.
    if any(arg in {"-f", "--force"} for arg in args) and not any(
        arg in {"-D", "-d", "--delete"} for arg in args
    ):
        return True
    names = [arg for arg in args if not arg.startswith("-")]
    return not names or not all(name.startswith("worktree-") for name in names)


def reflog_verdict(args):
    sub = next((arg for arg in args if not arg.startswith("-")), None)
    return sub == "expire"


def worktree_verdict(args):
    sub = next((arg for arg in args if not arg.startswith("-")), None)
    if sub != "remove":
        return False
    if not any(arg in {"-f", "--force"} for arg in args):
        return False
    # Carve-out: force-removing harness-delegate worktrees under
    # .claude/worktrees/ is cleanup of agent-created scratch. Paths are
    # normalized so traversal out of that directory still blocks.
    paths = [arg for arg in args if not arg.startswith("-")][1:]
    return not paths or not all(is_delegate_worktree_path(p) for p in paths)


def find_subcommand(tokens, start, end):
    j = start
    while j < end:
        current = tokens[j]
        if current == "--":
            j += 1
            break
        if current in GLOBAL_OPTS_WITH_ARG:
            j += 2
            continue
        if any(current.startswith(f"{opt}=") for opt in GLOBAL_OPTS_WITH_ARG if opt.startswith("--")):
            j += 1
            continue
        if current.startswith("-"):
            j += 1
            continue
        break
    if j < end:
        return tokens[j], j + 1
    return None, None


def classify(subcommand, args, via_xargs):
    if subcommand == "push":
        return "git push"
    if subcommand == "reset" and "--hard" in args:
        return "git reset --hard"
    if subcommand == "clean" and (
        "--force" in args or any(arg.startswith("-") and "f" in arg[1:] for arg in args)
    ):
        return "git clean -f"
    if subcommand == "branch" and branch_verdict(args):
        if any(arg in {"-f", "--force"} for arg in args) and not any(
            arg in {"-D", "-d", "--delete"} for arg in args
        ):
            return "git branch -f"
        return "git branch -D"
    if subcommand == "checkout" and checkout_verdict(args, via_xargs):
        return "git checkout path"
    if subcommand == "switch" and switch_verdict(args):
        return "git switch --force"
    if subcommand == "restore" and restore_verdict(args, via_xargs):
        return "git restore path"
    if subcommand == "rebase":
        return "history rewrite"
    if subcommand == "stash" and stash_verdict(args):
        return "git stash drop"
    if subcommand == "commit" and "--amend" in args:
        return "git commit --amend"
    if subcommand == "update-ref" and "-d" in args:
        return "git update-ref -d"
    if subcommand == "tag" and ("-d" in args or "--delete" in args):
        return "git tag -d"
    if subcommand == "reflog" and reflog_verdict(args):
        return "git reflog expire"
    if subcommand == "worktree" and worktree_verdict(args):
        return "git worktree remove --force"
    return None


def resolve_prefixes(tokens, i, end):
    via_xargs = False
    while i < end:
        if re.match(r"^[A-Za-z_][A-Za-z0-9_]*=", tokens[i]):
            i += 1
            continue
        base = os.path.basename(tokens[i])
        if base == "env":
            i += 1
            while i < end and re.match(r"^[A-Za-z_][A-Za-z0-9_]*=", tokens[i]):
                i += 1
            continue
        if base in {"command", "nohup"}:
            i += 1
            continue
        if base == "timeout":
            i += 1
            while i < end and tokens[i].startswith("-"):
                i += 1
            if i < end:
                i += 1
            continue
        if base == "xargs":
            i += 1
            via_xargs = True
            while i < end and tokens[i].startswith("-"):
                i += 1
            continue
        break
    return i, via_xargs


def command_end(tokens, i):
    j = i
    while j < len(tokens) and tokens[j] not in CONTROL_OPS:
        j += 1
    return j


def scan(tokens, allow_wrapper):
    expect_command = True
    i = 0
    n = len(tokens)
    while i < n:
        tok = tokens[i]
        if tok in CONTROL_OPS:
            expect_command = True
            i += 1
            continue
        if expect_command and tok in KEYWORDS:
            i += 1
            continue
        if expect_command:
            end = command_end(tokens, i)
            j, via_xargs = resolve_prefixes(tokens, i, end)
            if j < end:
                base = os.path.basename(tokens[j])
                if base == "git":
                    subcommand, args_start = find_subcommand(tokens, j + 1, end)
                    if subcommand is not None:
                        reason = classify(subcommand, tokens[args_start:end], via_xargs)
                        if reason:
                            return reason
                elif allow_wrapper and base in WRAPPERS:
                    for k in range(j + 1, end):
                        if re.match(r"^-[A-Za-z]*c[A-Za-z]*$", tokens[k]):
                            if k + 1 < end:
                                inner = tokenize(tokens[k + 1])
                                reason = scan(inner, allow_wrapper=False)
                                if reason:
                                    return reason
                            break
            expect_command = False
            i = end
            continue
        i += 1
    return None


try:
    tokens = tokenize(sys.argv[1])
except ValueError:
    tokens = sys.argv[1].split()

print(scan(tokens, allow_wrapper=True) or "")
PY
)" || { echo "BLOCKED: guard analyzer error — failing closed; rephrase the command." >&2; exit 2; }
[[ -n "$reason" ]] && block "$reason"
exit 0
