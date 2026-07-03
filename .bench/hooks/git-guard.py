# The destructive-git analyzer behind .bench/hooks/block-dangerous-git.sh — the
# single source that both classifies a command (argv[1] -> deny label or empty
# line) and enumerates the deny classes (--describe-classes), so the guard's
# enforcement and its advertised manifest cannot drift. Invoked by the hook via
# python3; not a hook itself (bench guards aggregates *.sh only).
import os
import posixpath
import re
import shlex
import subprocess
import sys

CONTROL_OPS = {";", "&&", "||", "|", "&", "(", ")"}
# shlex punctuation set = its default (`();<>|&`) plus newline, so a bare newline
# lexes as its own operator token instead of being folded into whitespace.
PUNCT_CHARS = "();<>|&\n"
_OP_CHARS = set(PUNCT_CHARS)
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

# The deny table: every destructive class the analyzer blocks, keyed to the verb
# label shown to the agent. classify() returns from here and --describe prints
# these values, so the enforcement and the advertisement share one source.
DENY_LABELS = {
    "push": "git push",
    "reset": "git reset --hard",
    "clean": "git clean -f",
    "branch-force": "git branch -f",
    "branch-delete": "git branch -D",
    "checkout": "git checkout path",
    "switch": "git switch --force",
    "restore": "git restore path",
    "rebase": "history rewrite",
    "stash": "git stash drop",
    "amend": "git commit --amend",
    "update-ref": "git update-ref -d",
    "tag": "git tag -d",
    "reflog": "git reflog expire",
    "worktree": "git worktree remove --force",
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
    # A bare newline is a command separator in shell, but shlex folds it into
    # whitespace — so a multi-line block (the common way an agent batches git)
    # would scan as one command and only its first verb would be classified.
    # Lex newline as punctuation and drop it from whitespace so it emits as its
    # own operator token; a newline *inside* quotes stays in the word (letters
    # keep it out of the operator-only class below), so a wrapper's inner string
    # survives intact for the recursive scan.
    try:
        lexer = shlex.shlex(s, posix=True, punctuation_chars=PUNCT_CHARS)
        lexer.whitespace_split = True
        lexer.whitespace = lexer.whitespace.replace("\n", "")
        raw = list(lexer)
    except ValueError:
        # Malformed quoting: fall back to a plain split, still honoring newlines
        # as boundaries so a multi-line block can't slip through unclassified.
        raw = []
        for idx, line in enumerate(s.split("\n")):
            if idx:
                raw.append(";")
            raw.extend(line.split())
    # Collapse any operator-only token that carries a newline (`\n`, `;\n`, `&&\n`,
    # `\n\n&&\n`) to a plain control op so scan() sees the boundary; an all-newline
    # run becomes `;`.
    out = []
    for tok in raw:
        if "\n" in tok and all(c in _OP_CHARS for c in tok):
            tok = tok.replace("\n", "") or ";"
        out.append(tok)
    return strip_redirections(out)


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
        return DENY_LABELS["push"]
    if subcommand == "reset" and "--hard" in args:
        return DENY_LABELS["reset"]
    if subcommand == "clean" and (
        "--force" in args or any(arg.startswith("-") and "f" in arg[1:] for arg in args)
    ):
        return DENY_LABELS["clean"]
    if subcommand == "branch" and branch_verdict(args):
        if any(arg in {"-f", "--force"} for arg in args) and not any(
            arg in {"-D", "-d", "--delete"} for arg in args
        ):
            return DENY_LABELS["branch-force"]
        return DENY_LABELS["branch-delete"]
    if subcommand == "checkout" and checkout_verdict(args, via_xargs):
        return DENY_LABELS["checkout"]
    if subcommand == "switch" and switch_verdict(args):
        return DENY_LABELS["switch"]
    if subcommand == "restore" and restore_verdict(args, via_xargs):
        return DENY_LABELS["restore"]
    if subcommand == "rebase":
        return DENY_LABELS["rebase"]
    if subcommand == "stash" and stash_verdict(args):
        return DENY_LABELS["stash"]
    if subcommand == "commit" and "--amend" in args:
        return DENY_LABELS["amend"]
    if subcommand == "update-ref" and "-d" in args:
        return DENY_LABELS["update-ref"]
    if subcommand == "tag" and ("-d" in args or "--delete" in args):
        return DENY_LABELS["tag"]
    if subcommand == "reflog" and reflog_verdict(args):
        return DENY_LABELS["reflog"]
    if subcommand == "worktree" and worktree_verdict(args):
        return DENY_LABELS["worktree"]
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


if len(sys.argv) > 1 and sys.argv[1] == "--describe-classes":
    seen = []
    for label in DENY_LABELS.values():
        if label not in seen:
            seen.append(label)
    print(", ".join(seen))
else:
    try:
        tokens = tokenize(sys.argv[1])
    except ValueError:
        tokens = sys.argv[1].split()
    print(scan(tokens, allow_wrapper=True) or "")
