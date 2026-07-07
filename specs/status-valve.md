# Status valve (FT30)

Status: implemented

## Problem
Six signals fire on this repo today, but `bench status` renders only five: the sixth is silently truncated to a bare `+1 more` line, and there is no flag to expand the board. The one row most often hidden is a real pending signal (the roadmap-reconcile row). A reader who sees `+1 more` has no way, from that line, to learn how to see the rest.

Separately, several of the board's action strings recommend surfaces the reader cannot invoke. The merged-spec row says `promote-then-delete (spec-retire)`, the unresolved-map row says `craft-grill → /bench-write-spec`, and the orphaned-scratch-branch row says `delete scratch branch` — none of which is a command or a canonical phase the reader can run verbatim. This violates the platform's own "recommend in the form the harness can invoke" rule, on the very surface a cold session reads first.

## Solution
Add an overflow valve and make the actions invocable, without touching the ambient budget:

- `bench status --all` prints every signal row with the five-row budget off and no `+N more` line. The default board is unchanged except the truncation line becomes `+N more (bench status --all)`, so the remedy sits on the page that truncates.
- The budget and severity ordering are untouched; the SessionStart hook keeps calling `bench status` with no flag, so the ambient surface stays bounded (auto-expansion was rejected).
- Each recommended action becomes something the reader can run or invoke verbatim — a `bench` subcommand or a canonical `/bench-*` phase.
- An unrecognized argument still fails closed with exit 2; `--all` is the only newly accepted token.

## User stories

1. As a reader who hits the five-row budget, I want `bench status --all` to print every signal row with the budget off and no `+N more` line, so I can see the signals the budget hides. Line: claude-sonnet-5 / low. This is a thin render-mode toggle whose exact output the runtime contract pins, so the gate fully observes the behavior.

2. As a reader who sees the board truncate, I want the overflow line to read `+N more (bench status --all)` instead of a bare `+N more`, so the command that expands the board is on the page that truncated it. Line: claude-sonnet-5 / low. One format-string change with an exact-substring contract assertion.

3. As a reader with a merged spec awaiting retirement, I want the action to read `bench spec retire <slug>` instead of `promote-then-delete (spec-retire)`, so I can run the remedy verbatim. Line: claude-sonnet-5 / low. A single action-constant swap the retirement contract test pins exactly.

4. As a reader with an unresolved decision map, I want the action to read `/bench-shape-idea` instead of `craft-grill → /bench-write-spec`, so the recommendation is a canonical phase I can invoke and it points at shaping — which an open map still needs — rather than spec-writing. Line: claude-sonnet-5 / low. A single action-constant swap the decisions contract test pins exactly.

5. As a reader with an orphaned scratch branch, I want the action to read `bench worktree clean` instead of `delete scratch branch`, so the remedy is a real subcommand — landing with or after FT28, which makes `bench worktree clean` actually sweep those branches. Line: claude-sonnet-5 / low. The string edit is coupled to FT28's branch-sweep, so this story lands with-or-after that spec and this spec only pins the final string in the contract test.

6. As a reader who mistypes, I want `bench status` to exit 2 with a usage line for any unrecognized argument — `--all` is the only new accepted token, and `--all` plus extra junk (or `--allx`, `-a`) is unrecognized — so a typo never silently prints the default board. Line: claude-sonnet-5 / low. Argument parsing is fully gate-observable at `Command` and through the built binary.

## Implementation decisions

**Module boundaries (off the map's Handoff).** `internal/status` owns the flag, the render-mode toggle, and every action string; `bin/bench.sh` needs no edit (`status) route_porcelain "$@"` already forwards `--all` verbatim to the Go `Command`, and the by-path CLI, hooks, and adapters inherit that one routing); no other package changes.

**Flag parsing.** `status.Command` gains one accepted token, `--all`. The accepted set becomes: no args (default board), `--all` (full board), `-h`/`--help` (usage, exit 0). Anything else — including `--all` followed by any extra argument, and near-misses like `--allx` or `-a` — is a usage error, exit 2. Outside a repo stays the structured error, exit 1. The `-h`/`--help` usage line gains `[--all]` so the help stays honest (existing `TestCommandArgs` asserts the `usage: bench status` prefix, which still holds).

**Render-mode toggle.** `render` gains a parameter that controls whether the five-row budget applies. In default mode it keeps the existing behavior (lead line, up to five rows, and the overflow line when there are more than five). In all mode it prints the lead line and every row and emits no overflow line. The severity sort, the lead-line selection, and the clean/empty message are unchanged.

**Truncation line.** The overflow format becomes `+N more (bench status --all)` (default mode only).

**Action-string ledger (the exact final string per signal row — contract tests pin each).**

| sev | signal | final action string | change |
|---|---|---|---|
| 0 | gate (red) | `fix before commit` | unchanged |
| 1 | git | `commit on green / push` | unchanged |
| 2 | worktree (out-of-pool) | `clean up (bench worktree clean)` | unchanged (FT11 landed) |
| 2 | worktree (leased) | `resume leased worktree` | unchanged |
| 2 | worktree (orphaned branch) | `bench worktree clean` | **changed** from `delete scratch branch` — coupled to FT28 |
| 3 | drain | `/bench-what-next` | unchanged |
| 4 | structure | `split (craft-seams)` | unchanged (see Out of scope) |
| 5 | decisions | `/bench-shape-idea` | **changed** from `craft-grill → /bench-write-spec` |
| 6 | gate (stale, strong) | `re-run the gate` | unchanged |
| 6 | gate (stale, capture-only) | `re-run when convenient` | unchanged |
| 7 | specs | `bench spec retire <slug>` | **changed** from `promote-then-delete (spec-retire)` |
| 8 | reviews | `promote or delete by hand` | unchanged (genuinely manual) |
| 9 | roadmap | `/bench-what-next` | unchanged |

`<slug>` in the specs action is a literal placeholder in the constant, not interpolated from a filename.

**FT28 coupling for the orphan-branch string (single source per fact).** The orphan-branch action constant has one home, `status.appendWorktree`. FT28 (worktree-branch-sweep) claims the same edit in its Handoff (`internal/status only updates the action string`), and the flip to `bench worktree clean` is only honest once FT28's sweep makes that command delete branches. Resolution: this spec **builds with or after FT28**; whichever spec lands first performs the single constant edit, and this spec owns pinning `bench worktree clean` on the orphan row in the status contract test. The two specs must not both edit the constant — the later build asserts rather than re-edits.

**Ambient surface stays bounded.** `.bench/hooks/session-start.sh` is not touched; it keeps invoking `bench status` with no flag, so a cold session gets the five-row budget. Auto-expanding the hook was rejected in the decision map.

## Testing decisions

- **What a good test is here.** Drive the built `bench` binary (or `status.Command` for pure arg parsing) and assert the exact stdout board and exit code — the flag, the truncation line, the per-row action strings, and the row count — never internal render state.
- **Seam and prior art.** The one seam is `bench status [--all]` stdout, tested through the existing runtime status contract family in `internal/contract/runtime` (drives the built binary): extend `testRuntimeStatusBudget`, `testRuntimeStatusDecisions`, `testRuntimeStatusRetirementSignal`, and `testRuntimeStatusWarmPool`, which already pin these rows. Argument-validation edges also attach at `status.Command` via `TestCommandArgs` in `internal/status`. The `render`-signature callers in `internal/status/status_test.go` (`TestRenderClean`, `TestRenderDirtyLeadsGitOverDrainRow`, `TestRenderWorkingRoadmapAloneIsClean`, `TestRenderSurfacesOrphanedWorktreeBranch`) are refreshed for the new parameter in the same diff.
- **Fixture sweep.** No canary or conformance fixture embeds these action strings or the truncation text; the one status canary (`tests/canary/behavior-owned/status-regressed`) asserts only the "clean" all-clear line, which is unchanged, so it needs no refresh. The two Go test files above are the entire pinned surface.
- **Gate command.** `.bench/gate.sh`.

### Seam diagram

    trigger: reader runs `bench status` or `bench status --all`
             (SessionStart hook runs `bench status`, no flag — bounded)
        │
        ▼
    args ──▶ [ status.Command  →  render(root, all) ] ──▶ stdout board:
                                                            ▶ lead line
                                                            ≤5 rows (default) | all rows (--all)
                                                            +N more (bench status --all)  ← default, when >5
    (unknown arg | --all + extra) ──▶ usage line, exit 2
              ◀ tests attach here: internal/contract/runtime drives the built binary and
                asserts exact action strings, the truncation line, row count, and exit code;
                status.Command unit test drives arg parsing directly

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | On a >5-signal tree, `bench status --all` prints the 6th row and emits no `+N more` line | runtime contract (built binary) | New assertion beside `testRuntimeStatusBudget`: `--all` output contains the 6th signal's text and does not contain `+N more`; red on current HEAD because `--all` exits 2 today (no all-mode) | A stub that ignores `--all` still truncates: the 6th row is absent and the overflow line is present, failing both halves |
| 2 | Default board with >5 signals prints `+N more (bench status --all)` | runtime contract | Assert exact substring `+N more (bench status --all)` in `testRuntimeStatusBudget`; red on HEAD (line is bare `+N more`) | An always-green stub keeps the bare `+N more`, failing the exact-substring match |
| 3 | Retirement row action is `bench spec retire <slug>` | runtime contract (`testRuntimeStatusRetirementSignal`) | Assert the retirement row contains `bench spec retire <slug>` and not `promote-then-delete`; red on HEAD | Catches leaving the old `promote-then-delete (spec-retire)` constant |
| 4 | Decisions row action is `/bench-shape-idea` | runtime contract (`testRuntimeStatusDecisions`) | Assert `/bench-shape-idea` and not `craft-grill → /bench-write-spec`; red on HEAD | Catches leaving the old `craft-grill → /bench-write-spec` constant |
| 5 | Orphan-branch row action is `bench worktree clean` | runtime contract (`testRuntimeStatusWarmPool`) + unit (`TestRenderSurfacesOrphanedWorktreeBranch`) | Assert `bench worktree clean` on the orphan row and not `delete scratch branch`; red on HEAD. Honest only once FT28 lands (coupled) | Catches the old `delete scratch branch` constant; FT28 dependency noted so the action is not a dead end |
| 6 | `--all` alone exits 0; `--bogus`, `--allx`, `-a`, and `--all <extra>` each exit 2 with a usage line | `status.Command` unit (`TestCommandArgs`) + runtime | Extend `TestCommandArgs`: `--all` → exit 0; `--all extra` → exit 2; red on HEAD because `--all` is treated as unknown (exit 2) | An impl that accepts `--all` but only checks `args[0]` swallows trailing junk, failing the `--all extra` → 2 row |
| 3–5 | boundary: every action string outside these stories — gate, git, worktree (out-of-pool/leased), drain, structure, stale-gate, reviews, roadmap — stays exactly as the ledger lists | runtime contract | Already covered by existing assertions in the status contract family — not new TDD | Any accidental edit to an unchanged constant fails the existing exact-substring assertion for that row |

### Edge inventory

Edge classes walked per the profile's shell-CLI hostile-input checklist; each resolved as a coverage row above or a **Won't handle** line here.

- Unquoted / multi-word args, `--all` + extra, near-miss flags (`--allx`, `-a`) — coverage row 6.
- SessionStart auto-expansion — **Won't handle**: the hook keeps calling `bench status` with no flag and is not edited, so the ambient surface stays the bounded five-row budget (auto-expansion was rejected in the decision map).
- Control bytes / TOON refusal in git-sourced text — **Won't handle**: the changed actions are literal constants with no git-sourced interpolation (`<slug>` is a literal placeholder), so this change adds no new hostile-text surface, and the board is not TOON-rendered.
- Absent-vs-empty signal sources, trailing-newline files — **Won't handle**: signal detection is untouched by this spec (only action strings, the truncation line, and one flag change) and is already covered by the existing status tests.
- Symlink invocation, required tool missing from PATH, SIGINT mid-run, re-run idempotency, cwd deeper than root — **Won't handle**: `status` is a stateless read renderer with no writes, leases, or scratch state, so it is idempotent and unaffected by this change.
- Reach through every shipped surface (kit CLI, by-path CLI, hooks, adapters) — **Won't handle** as a new test: `route_porcelain "$@"` already forwards `--all` to the one Go impl, so no bin/bench.sh routing edit exists to regress and the multi-surface reach is inherited.

## Out of scope

- **Structure-row action rewording (`split (craft-seams)` → a phase name).** The decision map's #2 named the structure row for change ("keep the split recommendation but name it via the phase that loads the skill"), but the task's decided contracts scope it out, and no `/bench-*` phase owns the structural-split remedy — `craft-seams` is loaded by `/bench-write-spec` and `/bench-debug`, neither of which is a "split this module" phase — so a compliant target would have to be invented. This is a separate capability blocked on a new decision (which phase names the split remedy). **Flagged for veto:** keep `split (craft-seams)` this spec, or shape a follow-up if you want a phase wrapper. Estimate once the target is decided: 1 constant edit + fixture refresh, 1 gate run.
