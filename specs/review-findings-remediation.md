# Review-findings remediation batch

Source: the 2026-07-01 whole-repo defect review. Eight build-ready findings,
each with an obvious fix direction — batched into one spec so the gate locks
all of them. The two design-heavy findings (guard rework — shipped; Codex hook
layer — own decision map) are not in this batch.

## Problem

The review verified defects the gate never sees: the shift loop destroys its
own inter-iteration memory on rollback; the Stop hook can trap an agent in a
permanent block and record a gate verdict that never ran; `bench link` breaks
or silently loses its git-hook backstop in common repo layouts and corrupts
docs that mention its own markers; the refactor prompt points agents at debt
the shift never touched; three canary fixtures pass vacuously; and five
command/skill docs make claims that contradict the machinery.

## Solution

Fix each at its existing seam, with a contract-matrix or anchor red-to-green
per fix, so every one of these regressions is a gate failure from now on.

## User stories

Shift loop:
1. As the shift agent, I want `.bench-notes.md` and `.bench-objective` to
   survive a red-gate rollback (both the iteration loop and the refactor
   phase), so a retry learns from the failed attempt instead of repeating it.
2. As the shift agent, I want the refactor-phase prompt scoped to the files
   this shift touched — the flagged files inlined into the prompt — so I never
   split unrelated legacy files on a feature branch.

Stop hook:
3. As the harness, I want the Stop hook to honor `stop_hook_active` from its
   stdin JSON (exit 0 when true), so a permanently red gate cannot loop an
   agent forever; the dead `BENCH_STOP_CHECKED` env guard is removed.
4. As the reviewer, I want a missing bench CLI to fail open — allow the stop
   with a loud stderr warning and write no gate cache — so an unfindable
   binary neither traps the agent nor forges a `red` verdict the dashboard
   then reports as real.

Link safety:
5. As a user linking from a kit path containing glob metacharacters
   (`[`, `?`, `*`), I want every file installed at its correct destination.
6. As a user linking inside a git worktree or submodule (`.git` is a file), I
   want link to succeed, installing the pre-push hook where git actually reads
   it (`git rev-parse --git-path hooks`), not to die mid-install.
7. As a user with `core.hooksPath` configured (husky et al.), I want the
   pre-push hook installed into the effective hooks directory under the same
   Bench-managed conflict rules, so the backstop is never silently absent.
8. As a user whose remote lacks `origin/HEAD`, I want link to resolve the real
   default branch (best-effort `git remote set-head origin --auto`, falling
   back to `main`), so the pre-push hook guards the actual default branch.
9. As a user whose AGENTS.md *documents* the Bench markers inside a fenced
   code block, I want marker detection and the managed-block rewrite to ignore
   fenced lines, so a relink never deletes project-owned prose.

Canary attribution:
10. As the gate maintainer, I want a meta-check that runs the inner gate once
    against an empty fixture and asserts no fixture's EXPECT matches that
    baseline output, so a vacuous canary is itself a red gate.
11. As the gate maintainer, I want the three vacuous fixtures
    (`bad-frontmatter`, `dangling-index`, `acceptance-coverage-anchor`) and
    the four weakly-attributed ones re-planted so each EXPECT can only be
    produced by its targeted check firing on the planted content.

Doc accuracy (each edit is one story; all change agent behavior):
12. `bench-final-check` names the real gate resolution chain
    (`.bench/gate.sh` → `$BENCH_GATE` → auto-detect), never
    `projects/<name>.md`; a gate anchor enforces the chain stays named.
13. `bench-debug` says "add the repro **to** the gate" — pointing the gate at
    the repro is gate-weakening and needs explicit reviewer approval.
14. `bench-write-spec` reads settled ADRs from `docs/adr/` and treats
    `decisions/` as working maps, matching `craft-adr`'s split.
15. `bench-review-implementation` lists `AGENTS.md` and `.bench/BENCH.md` as
    Standards sources (CLAUDE.md is import pointers only).
16. `craft-skills` includes itself in the model-invoked inventory
    (`craft-synthesis` invokes it mid-work).
17. `craft-cli` carves out CLIs whose conformance target is project-declared
    (bench's own stderr/exit-3 contract must not be "fixed" into TOON).
18. `craft-seams` replaces the phantom per-file budget escape hatch with the
    real options: raise the global cap deliberately, or record the exception
    for the reviewer.
19. `craft-tdd` warns that inside a shift, red-to-green must complete within
    one iteration — a stopped-at-red iteration is rolled back and deleted.

## Implementation decisions

- Rollback preservation: the two rollback sites exclude the scratch files via
  `git clean` exclude patterns; the loop's cleanup-on-exit behavior is
  unchanged.
- Stop hook: parse stdin once with the same python-one-liner idiom the guard
  uses; `stop_hook_active` → exit 0 before any gate run; missing CLI → stderr
  warning, exit 0, no cache write. The cache is written only by a real gate
  run's verdict.
- Link: quote the prefix-strip pattern; resolve the hooks directory through
  `git rev-parse --git-path hooks` everywhere the hook is written or
  preflight-checked; default-branch resolution stays in `default_branch()`
  with the set-head attempt at link time only (no network on hot paths).
- Fence-aware markers: one awk state machine (toggling on ``` lines) shared by
  the marker count, marker line, and block rewrite paths — a single source of
  truth for "is this line inside a fence".
- Refactor prompt: `refactor_prompt` takes the flagged-files output as an
  argument and inlines it; no new CLI flag, no repo-wide instruction.
- Canary meta-check: baseline output computed once per gate run from one empty
  temp repo; each EXPECT is `grep -vF`-checked against it before the fixture
  runs. Re-planted fixtures put content under the paths the targeted checks
  actually scan (e.g. `.agents/skills/`, not `.claude/skills/`), and EXPECT
  strings name the planted artifact (`skill 'orphan'`, `bench-write-spec.md`)
  so they cannot come from glob literals or missing-file noise.
- Docs: text edits at the cited lines; one new `require_anchor`-style gate
  check for story 12 only (the others are guarded by review, not anchors —
  anchor-per-edit is sediment).

## Testing decisions

- Good tests here run the real CLI/hook in a throwaway repo and assert
  observable behavior: exit codes, file presence/content, prompt content
  captured via a recording `BENCH_AGENT`, cache-file content. Prior art: every
  block in `.bench/gate-runtime-contracts.sh` and `.bench/gate-link-contracts.sh`.
- Seams: the `bench shift`/`bench link` CLI contract seam, the Stop-hook
  process seam (stdin JSON → exit code + cache file), the gate's canary
  section, and `require_anchor` for story 12. No new seams.
- Gate: `bench gate`. Done = green with all new contract cases present.
- Line (declared): shell/contract stories at **mid tier, medium effort** (binding
  in projects/benchkit.md: Sonnet 4.6; in-session delegates run Opus 4.8 per the
  alias caveat there); doc stories 12–19 at **top tier, high effort** (profile
  routes skill/command/doc authoring to top). Combined cap ~250k tokens.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | notes+objective survive red rollback into iteration 2 | shift contract (recording agent) | observed 2026-07-01: iteration 2 reports `NOTES FILE GONE` / `OBJECTIVE FILE GONE` | asserts the exact files the rollback currently deletes |
| 2 | refactor prompt contains the flagged touched files, not repo-wide `bench structure` | shift contract (prompt-recording agent) | observed: captured prompt says ``Run `bench structure` ``, zero touched-scope mentions | prompt content is the behavior; the row greps it |
| 3 | `stop_hook_active: true` on stdin → exit 0 under armed shift + red gate | stop-hook process seam | code-verified stdin is never read; red run lands in build (fixture: armed, red gate, flag true, currently exits 2) | the flag is the documented loop-breaker; the row exercises it directly |
| 4 | missing bench → exit 0 + warning, no cache file | stop-hook process seam | observed: exit 2 and cache `red <sha>` written with bench absent | asserts both the trap and the forged verdict |
| 5 | link from `dir[1]` kit path installs correctly | link contract | observed (verified repro): plan rows carried absolute paths, `.bench/bin` empty | re-runs link from a metachar path and checks installed files |
| 6 | link inside a worktree succeeds, hook lands in effective hooks dir | link contract | observed: `chmod: .../.git/hooks/pre-push: Not a directory`, exit 1 mid-install | the failing layout becomes a fixture |
| 7 | with `core.hooksPath`, hook installs there under conflict rules | link contract | observed: effective dir `.husky` vs written `.git/hooks/pre-push` | asserts hook exists at `git rev-parse --git-path hooks` |
| 8 | pre-push hook names the remote default branch when resolvable | link contract | not TDD-able red-first without a remote fixture; asserted post-build with a local bare remote whose HEAD is `master` | hook content grep proves which branch is guarded |
| 9 | fenced markers survive relink; real markers still rewritten; unclosed fences fail closed | link contract | observed (verified repro): relink deleted content between fenced marker and end marker; unclosed-fence relink duplicated the managed block | fixture AGENTS.md with fenced + real markers, content-grepped with a single-managed-block count after relink; unclosed-fence fixture must conflict |
| 10 | vacuous EXPECT = red gate | gate canary meta-check | observed: `bad-frontmatter`, `dangling-index`, `acceptance-coverage-anchor` EXPECTs match an empty fixture | the meta-check is the regression test for canary rot itself |
| 11 | re-planted fixtures bite on planted content only | gate canary | same observation; post-fix EXPECTs name planted artifacts | EXPECT strings that can only come from the targeted check |
| 12 | final-check names the real resolution chain | require_anchor | observed: `grep -qF '.bench/gate.sh' .agents/commands/bench-final-check.md` fails today | anchor red until the doc edit lands |
| 13–19 | doc edits at cited lines | review | not TDD-able (prose); verified at review phase against the cited findings | — |

## Out of scope

- **Codex hook layer** — separate capability, own decision map queued
  (`/bench-shape-idea`); research spike ~1h before it can be specced.
- **Per-path structure budgets** — separate capability (new gate surface);
  parked on the roadmap; ~2h if built.
- **Worktree lease staleness/TOCTOU hardening** — separate capability
  (concurrency model for the pool); parked; ~2–3h.
- **`BENCH_AGENT` multi-word/adapter contract** — owned by the open
  `decisions/dogfood-improvements.md` ticket #4 (reviewer decision pending
  there); not built here to avoid pre-empting that map.
- **Worktree-count status nag** — separate capability (status signal
  semantics for the warm pool); parked; ~1h.
