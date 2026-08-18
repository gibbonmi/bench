# Baseline — Fable high reconciliation run

Recorded at run start (2026-08-18), inside `/home/mgibs/workspace/bench-audit-fable`
(linked worktree; main tree `/home/mgibs/workspace/bench`).

## Model identity

Harness: Claude Code. The system prompt reports the model as **Fable 5, model id
`claude-fable-5`**. Effort: high (as assigned). No runtime surface exposes a live
model build identifier beyond that declaration; recorded at start and re-checked at
close (see `final-reconciliation.md` § "Model continuity"). No fallback or model switch
was reported by the harness at any point in the run.

## Git state at start

```
$ git status --short
 M .claude/README.md
 M .claude/settings.json
?? audit/
$ git rev-parse HEAD
58d966e2f92f7f37eba07b6215e8eef45371b72d
$ git log -1 --oneline
58d966e2 what-next: drain journal and worktree-cleanup-eligibility retro, prioritize FT100
$ git branch --show-current
audit/fable
$ git worktree list
/home/mgibs/workspace/bench              58d966e2 [main]
/home/mgibs/workspace/bench-audit-fable  58d966e2 [audit/fable]
/home/mgibs/workspace/bench-audit-opus   58d966e2 [audit/opus]
/home/mgibs/workspace/bench-audit-sol    58d966e2 [audit/sol]
```

The two tracked modifications (`.claude/README.md`, `.claude/settings.json`) pre-date
this run: they add a `## Skill overrides` section to the README and turn three more
skills (`design`, `artifact-diagramming`, `code-review`) off in `skillOverrides`. They
were not made by this run and were left untouched. `audit/` is the untracked audit
input/output tree.

## Environment observations that shaped the run

- No Go binary was present in the worktree (`dist/bench` absent) or the main tree; the
  PATH shim `~/.local/bin/bench` → `~/workspace/bench/bin/bench.sh` therefore also
  exited 127 (`no pinned binary for this platform`). The `block-dangerous-git` hook then
  refused every Bash call containing the substring `git` — including a read-only
  `cat .bench/hooks/block-dangerous-git.sh`. Recovery: `bash scripts/go-build.sh "$PWD"
  dist/bench` (allowed as a local build; `dist/` is gitignored). This is the same cold
  start Sol (E001) and Opus (E1/E2/I.8) reported and is itself evidence.
- The Claude SessionStart hook output visible at session open was the user's global
  `gl-axi` dashboard, not Bench's; Bench's `session-start.sh` exec'd `session-inspect`
  into the missing binary and printed nothing usable.

## Which commit did each prior audit examine?

| Audit | Claimed subject | Basis | Relationship to this run's HEAD |
|---|---|---|---|
| Sol (Codex, GPT-5.6 Sol xhigh) | `58d966e2f92f7f37eba07b6215e8eef45371b72d`, branch `audit/sol` | `report.md` header, `evidence.md` E005 (`git rev-parse HEAD` recorded), E023 final status | **EXACT MATCH** |
| Opus (Claude Code, Opus 5 xhigh) | `58d966e2f92f7f37eba07b6215e8eef45371b72d`, branch `audit/opus` | `report.md` § B, `01-environment-and-oracle.md` Phase 0 block | **EXACT MATCH** |

Both worktrees are still checked out at `58d966e2` (`git worktree list` above). No
repository drift separates the three audits, so no finding is attributable to version
skew; every disagreement is a difference in what was inspected, run, or inferred, and
is adjudicated on fresh observation in `reconciliation-ledger.md`.

Sol's E006 (`bench test` fixture pollution) and its wrapper-vs-binary root findings
were run in a differently configured shell environment (Sol's `BENCH_HOME` and PATH
resolution); those are re-derived here rather than trusted (see ledger L-06, L-14).

## Working conventions for this run

- Fresh executable observation outranks prior auditor evidence; every ledger entry
  names which of the two it rests on.
- Temporary experiments: a gofmt-red mutation of `internal/toon/toon.go` (cp-aside and
  restored), a throwaway `git init` repo under the session scratchpad for the
  `bench setup` chain (removed at close), local gate runs (which write only to
  gitignored `.logs/` and the git-dir gate records). No tracked file was modified;
  no commit, push, or release command was run.
