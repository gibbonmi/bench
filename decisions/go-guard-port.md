# go-guard-port — slice 4 of the Go rewrite

Child of `decisions/go-rewrite.md` (#6, slice order). Slice 4 as mapped is
"`git-guard.py` absorbed as a Go package"; the grill widened it to the whole
destructive-git guard: the analyzer's only caller, `block-dangerous-git.sh`,
needs python3 twice (JSON envelope parse + analyzer), and porting the analyzer
alone would leave the hook straddling two runtimes. Absorbing the hook applies
the parent map's slice-5 shim pattern one hook early and deletes the kit's
last python3 runtime dependency. Slice 5's remaining scope after this lands:
`session-start.sh`, `stop.sh`, `check-agent-line.sh`, and the adapters'
`_line-guard.sh`. The regression net is strong: a 79-assertion allow/block
matrix in the runtime-git gate fragment drives the hook end-to-end through
the JSON envelope and survives the port unchanged.

## #1: Where does the port stop — analyzer only, or the whole hook?

Type: Grill

### Question
Port only `git-guard.py` (hook keeps python3 for the envelope parse, calls
the binary for classification), or absorb the hook: binary reads the
PreToolUse envelope on stdin, classifies, and emits the verdict, with
`block-dangerous-git.sh` reduced to a thin shim?

### Answer
**Whole hook.** Analyzer-only leaves the hook needing two runtimes and
doubles the degradation branches; absorption removes python3 from the kit
entirely (nothing else runs python once this file dies) and is the
already-decided slice-5 shim shape applied early. The fail-closed shell rim
stays in the shim: binary missing → refuse anything git-shaped, allow the
rest. Rejected: analyzer-only.

## #2: What is the port-parity net?

Type: Grill

### Question
The analyzer is 418 lines of subtle tokenizing (newline collapsing,
malformed-quote fallback, redirection stripping, one-level wrapper scan,
xargs/env prefixes); the 79-case matrix asserts verdicts, not tokenizer
internals. Matrix + Go table tests alone, or also a differential check
before the Python dies?

### Answer
**Matrix + Go table tests, plus a throwaway build-time differential pass.**
A generated command corpus runs through both the Python and Go analyzers and
the verdict diff must be empty before `git-guard.py` is deleted. The corpus
harness is never committed — it is insurance during the port, not a
permanent artifact. Rejected: matrix-only (thinnest exactly where the
tokenizer is subtlest).

## Handoff

1. **Module boundaries.** New `internal/` package for the analyzer
   (tokenizer, prefix/wrapper resolution, scan, per-verb verdicts, the deny
   table); ref/branch existence checks go through `internal/git`. `cmd/bench`
   gains a plumbing subcommand (name is the spec's call) that reads the hook
   JSON envelope on stdin and yields the verdict, with a describe-classes
   mode feeding the manifest. `block-dangerous-git.sh` becomes a thin shim;
   `git-guard.py` and the hooks `__pycache__` are deleted; the gate's
   python-parse check (gate.sh 1c) retires with the last `.py`.
2. **Contracts.** The hook's observable interface is unchanged: JSON envelope
   on stdin; exit 0 allow; exit 2 with a `BLOCKED:` message naming the deny
   label; missing/empty `tool_input.command` allows; `--describe` prints the
   name/boundary/denies/why manifest and exits 0 without reading stdin, with
   the denies clause generated from the same deny table classification uses
   (one source per fact, as today). All 79 matrix assertions and the
   carve-outs (delete `worktree-*` branches, force-remove
   `.claude/worktrees/` paths) carry verbatim. New plumbing contract shape —
   whether the binary exits 2 itself or the shim maps a stdout verdict — is
   the spec's call; either way the analyzer-error path stays fail-closed.
3. **Deep vs thin.** The binary is the deep unit: envelope parse, tokenize,
   scan, classify, and the block message live behind the subcommand. The shim
   is a one-glance pass-through with exactly two rims of its own: binary
   missing → refuse git-shaped input loudly, binary errored → refuse.
4. **Black-box assertables.** The 79-case allow/block matrix runs unchanged
   against the ported hook. New `go test` table tests: tokenizer edges
   (newline-as-separator, operator-token collapsing, malformed-quote
   fallback, redirection stripping, quoted newlines surviving into wrapper
   strings), prefix resolution (env/command/nohup/timeout/xargs), the
   one-level wrapper scan, each verdict function's option/carve-out matrix,
   and describe-classes equalling the deny table. The differential pass (#2)
   gates deletion of the Python.
5. **Gate attachment.** The unchanged shell gate is the oracle; the
   runtime-git fragment exercises the hook end-to-end and the guards
   fragments cover `--describe` and fail-closed behavior. The
   `analyzer-parse-broken` canary retires with its guarded surface; the
   fail-closed contracts that today blank/copy `git-guard.py` re-point to
   the binary-missing shape (the spec adapts them — the property guarded,
   "cannot classify → refuse", survives; its trigger changes). The
   python3-missing degradation branch dies rather than porting. New
   gate-visible case: shim with an absent binary path must still fail
   closed — reachable by the fragment, so a contract, not manual verify.
6. **Hostile-input owners.** Multi-line command blocks → tokenizer
   newline handling, table-tested. Malformed quoting → the fallback-split
   path, preserved verbatim. Worktree carve-out paths with spaces or `..`
   traversal → path normalization tests. Binary missing in hook context →
   the shim's fail-closed rim. Non-JSON stdin or an envelope without
   `tool_input.command` → allow (current behavior), asserted. Hostile
   command text (wrappers, redirects, xargs) is the analyzer's own job —
   the threat model stays honest-mistake, one wrapper level deep, by design.
7. **Uncertainty flags.** Both harnesses wire this hook
   (`.claude/settings.json` and `.codex/hooks.json`); the spec verifies the
   Codex envelope carries `tool_input.command` in the same shape the current
   python parse assumes — if it differs, the binary must accept both, and if
   that can't be settled by reading the adapters, escalate per `craft-line`
   rather than guess.
8. **Rejected alternatives.** Analyzer-only port keeping python3 for the
   envelope (#1); matrix-only parity net (#2); splitting `git-guard.py` to
   meet the 400-line structure budget (superseded — the file is deleted;
   the parked roadmap line for that split was removed with this map).
9. **Domain watch-outs.** After this slice the destructive-git guard depends
   on the platform binary: a missing binary fails closed on everything
   git-shaped, so a broken install makes git unusable for the agent until
   repair — loud and safe, but a new failure mode the old pure-text hook
   did not have. `bench guards` aggregates `*.sh` manifests only; the
   analyzer's Go home must not grow a hooks-dir `.sh` that answers
   `--describe` a second time. The guard's authority asymmetry (bench's own
   in-process rollback stays outside the agent's shell) is behavior, not
   commentary — the shim must keep intercepting only the agent's tool calls.

Dependency order: n/a — single spec.
