# after-merge-owner

Status: staged

## Problem
Post-merge duties — spec retirement, orphaned-review-pickup disposal, worktree/branch cleanup, the decision-map sweep, and the HANDOFF.md lifecycle — are scattered onto phases reached only for other reasons. `bench status` already computes and ranks these housekeeping rows, but no command consumes them, so the tail gets done only when a session happens to notice. Two of the duties are structural gaps: `decisions/` accumulates closed maps for long-shipped work with no lifecycle end (asymmetric with specs' promote-then-delete), and HANDOFF.md is an unowned, hand-written cold-session note that rots by design while `.bench/BENCH.md` already says cold pickup must not depend on it.

## Solution
Charge `/bench-final-check` with an exit duty that consumes the housekeeping status rows on a green landing on the default branch, close the `decisions/` lifecycle with promote-then-delete symmetry (including a one-time sweep of the maps whose work has already shipped), and retire HANDOFF.md by moving any still-true durable prose into CONTEXT.md and deleting the file. After this build, the post-merge tail has a named owner in both the default-branch and topic-branch flows, `decisions/` holds only maps for unshipped work, and there is no unowned handoff document to go stale.

## User stories

1. As an agent closing `/bench-final-check`, I want a post-merge exit duty that, on a green landing on the default branch, reads `bench status` and runs the housekeeping rows it flags — spec retirement (`bench spec retire` + the `spec-retire:` commit, sev 7), orphaned-review-pickup disposal (promote or delete the pickup by hand, sev 8), and worktree/branch cleanup (`bench worktree clean`, sev 2) — and, on a topic branch, states that these duties defer to the next default-branch session whose SessionStart status re-surfaces the same rows, so the tail is owned in both flows without a new phase. The duty explicitly leaves the roadmap-reconcile (sev 9) and capture-drain (sev 3) rows to `/bench-what-next`, so the two phases route rows to owners without overlapping.
Line: claude-fable-5 / high. The exit-duty wording loads into every final-check session, so under the craft-line leverage override the top tier is spent where a routing error would misdirect the post-merge tail for every future build.

2. As a cold session fixing my vocabulary, I want `decision map` defined as a Core term in CONTEXT.md — a closed `decisions/<topic>.md` whose Handoff produced a spec — with its lifecycle end stated as promote-then-delete once the work ships, so the term the ambient dashboard's signal list already uses has a canonical definition and its retirement rule is legible.
Line: claude-fable-5 / high. This defines ubiquitous language a cold session reads to fix its vocabulary, and an imprecise term compounds across sessions, so the leverage override routes it to the top tier.

3. As a reviewer, I want a one-time sweep that, for each closed map whose work has shipped, first promotes any durable decision not yet recorded in an ADR, CONTEXT.md, or the profile, then deletes the map — leaving `decisions/` holding only maps for unshipped work (the current FT25–FT36 batch) — with every dangling reference to a deleted map fixed in the same reviewed diff, so the shipped-map archive stops accumulating and stays symmetric with spec retirement. The sweep checks each map before deleting; it is not a bulk-rm.
Line: claude-fable-5 / high. Each shipped map needs a per-map durable-decision promotion judgment before it can be deleted, which the map's uncertainty flag forbids collapsing into a bulk-rm, so the judgment — not the delete — sets the top tier.

4. As a session picking the repo up cold, I want HANDOFF.md retired — its still-true durable prose moved into CONTEXT.md first, then the file deleted and its `package.json` `files[]` entry removed in the same commit — so cold pickup rests only on the owned mechanisms (SessionStart status, ROADMAP.md, CONTEXT.md, the profile) and no unowned handoff document is left to rot or to red the gate.
Line: claude-fable-5 / high. Choosing which HANDOFF prose is still durable and where it lands in CONTEXT is a promotion judgment, so the tier follows the judgment even though the file removal and the one-line package.json edit that rides with it are mechanical.

## Implementation decisions

- **`.agents/commands/bench-final-check.md`** owns the exit-duty prose. Extend the Exit handoff (do not add a phase — map #1 rejected a new after-merge phase) with a post-merge-tail duty that gates on a green landing on the default branch, consumes the sev 7 / sev 8 / sev 2 rows via their owner actions, and states the topic-branch deferral. Deep-vs-thin: the duty routes each row to its owner and must **not** restate `/bench-what-next`'s job — it hands sev 9 (roadmap reconcile) and sev 3 (capture drain) to that phase by name. It composes with, and does not duplicate, the existing same-session `reviews/<slug>.md` skip-fix-pass disposal already in the Exit handoff; the sev 8 duty covers the distinct orphaned case (a pickup whose spec is already gone).
- **`CONTEXT.md`** owns the `decision map` term. Add it to Core terms with the promote-then-delete lifecycle end; the term already appears in the `signal` line's parenthetical, so this supplies the missing definition rather than a new usage.
- **`decisions/` sweep** is a one-time content change: per-map promotion check, then delete the shipped-work maps, keeping the unshipped FT25–FT36 batch maps (including this spec's own `after-merge-owner.md`, retired later when this spec ships). The map estimates roughly 19 shipped-work maps and ~17 deletions; the exact set is whatever the per-map shipped-and-promoted check resolves at build time, not a hardcoded list.
- **`HANDOFF.md`** is deleted, not regenerated (map #3), after its durable prose moves to CONTEXT.md.
- **`package.json`** loses its `"HANDOFF.md"` `files[]` entry in the deletion commit. **Deviation from map #3, flagged for veto:** the map assigned the `files[]` removal to FT26, but `checkPackageFiles` reds the gate (`package.json files[] missing HANDOFF.md`) the moment the file is deleted with the entry left behind, so invariant #4 (repo stays green) forces the one-line entry removal into this commit. This completes the deletion decision rather than reopening it; the broader FT26 first-install package work stays out of scope.
- **No Go code changes.** The status rows (sev 2/7/8/3/9) already exist in `internal/status`; this spec consumes them, it does not add or alter them.

## Testing decisions
- **What a good test is here:** the docs-conformance layer catches stale or unknown command tokens in the edited prose, and the full gate catches a dangling `files[]` entry or reference after the deletions. Semantic correctness — whether the duty routes rows to the right owners, whether the term is defined well, whether every deleted map's durable content was promoted, whether HANDOFF's durable prose survived — is review-observed, because no conformance check pins it and adding one is a Go change this spec excludes.
- **Seams tested (both already in the gate):** the root-conformance docs checks (`checkColdPickupCLILists` + the stale-command-reference scan) over the edited command file and CONTEXT.md, and `checkPackageFiles` over the deletions. Prior art: the `stale-command-reference` canary fixture and the `tests/canary/workflow-guidance-anchors/*` fixtures.
- **Review seam (non-gate, advisory):** `/bench-review-implementation` owns the semantic rows below; it is not a test-attach point, so it carries no diagram.
- **Gate command:** `.bench/gate.sh` (the project gate).

### Seam diagram

Seam A — docs-conformance token scan over the edited prose:

    trigger: bench gate → go test ./internal/conformance -run '^TestRootConformance$'
        │
        ▼
    bench-final-check.md ──▶ [ stale-command-reference +        ] ──▶ diags == ∅  (green)
    CONTEXT.md           ──▶ [ cold-pickup CLI-list scans       ] ──▶ unknown `bench <cmd>` /
                                                                       stale `/bench-*` token → red
                                  ◀ tests attach here: edit prose, run gate, read diags

Seam B — package/conformance gate over the deletions:

    trigger: bench gate  (sweep + HANDOFF deletion staged)
        │
        ▼
    package.json files[] ──▶ [ checkPackageFiles: stat every    ] ──▶ every entry exists → green
    deleted HANDOFF.md   ──▶ [ files[] entry; doc references     ] ──▶ missing entry / dangling
    deleted maps         ──▶ [                                   ]      reference → red
                                  ◀ tests attach here: stage deletions, run gate; a left-behind
                                    files[] entry reds it

### Acceptance coverage map
| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | Exit-duty prose cites only real `bench` routes and valid `/bench-*` commands | Seam A (docs conformance) | already covered — the existing scan reds `TestRootConformance` on an unknown `bench <cmd>` or stale `/bench-*` token | The scan enumerates `bin/bench.sh` routes and valid command files; a mistyped retire subcommand or a misspelled phase name in the duty prose fails it |
| 1 | Duty routes sev 7/8/2 to owner actions, leaves sev 9/3 to `/bench-what-next`, and states the topic-branch deferral | Semantic review | review-only, not TDD-able — no conformance check pins the duty's row-routing semantics and adding one is a Go change out of scope | Only a reader can tell a mis-routed row (a duty claiming the sev 9 reconcile) or a restatement of what-next from a correct hand-off; the gate sees tokens, not intent |
| 2 | CONTEXT.md defines `decision map` as a Core term with its promote-then-delete lifecycle end | Semantic review | review-only, not TDD-able — no check enumerates CONTEXT Core terms, and a term-anchor check would be a Go change out of scope | Presence and correctness of a definition is a reading judgment; the stale-reference scan guards only command tokens in the file, not that the term exists |
| 3 | After the sweep, `decisions/` holds only unshipped-work maps, and every deleted map's durable decisions were promoted to an ADR, CONTEXT.md, or the profile first | Semantic review of the one-time diff | review-only, not TDD-able — the sweep is a one-time reviewed diff, not gated behavior; a bulk-rm is exactly what the map's uncertainty flag forbids | Whether a map held an unpromoted durable decision is a per-map judgment no gate encodes; only review catches a decision lost to deletion |
| 4 | HANDOFF.md's still-true durable prose is moved to CONTEXT.md before the file is deleted; nothing durable is lost | Semantic review | review-only, not TDD-able — no gate diffs the deleted prose against CONTEXT, and the file is gone by design | Only a reader comparing the old HANDOFF against CONTEXT can confirm durable prose survived; git keeps the old file but the gate does not compare it |
| 3–4 | The full gate stays green after HANDOFF.md and the shipped maps are deleted — no dangling `files[]` entry or reference | Seam B (`bench gate`: `checkPackageFiles` + `TestRootConformance`) | real red — deleting HANDOFF.md while its `package.json` `files[]` entry stays reds `checkPackageFiles` with `package.json files[] missing HANDOFF.md` | `checkPackageFiles` stats every `files[]` entry; a deletion that leaves the entry behind, or a doc still pointing at a removed map, reds the gate |

Degenerate-implementation check: a "delete HANDOFF.md, skip the durable-prose move, forget the `files[]` entry" shortcut reds row 6 at the gate; a "bulk-rm every map without the promotion check" shortcut passes the gate (closed maps do not move `UnresolvedCount`) but fails row 4 at review — which is why the sweep is scoped as a reviewed diff, not gated behavior.

### Edge inventory
Edge classes walked per behavior; this is a prose/content change with no new input surface, so the shell hostile-input checklist is mostly n/a.

- **Unpromoted durable decision in a shipped map** (map's uncertainty flag) — resolved as row 4: the per-map promotion check runs before any delete.
- **Durable prose in HANDOFF.md not yet in CONTEXT** — resolved as row 5: the move precedes the delete.
- **Dangling reference to a deleted file** (`files[]` entry, or a doc pointing at a removed map) — resolved as row 6: the gate reds until the reference is fixed in the same commit.
- **Topic-branch final-check** (deferral path) — resolved as row 2: the duty defers the default-branch-gated rows and says so.
- **Re-run idempotency of the exit duty** — Won't handle: the duty reads `bench status` and acts only on rows that fire, so a second run with no housekeeping row is a no-op by construction.
- **Interrupt / partial state mid-sweep** — Won't handle: the sweep is one reviewed diff committed atomically; a partial deletion is an unfinished diff the reviewer sees before commit, not a runtime hazard.
- **Shell hostile inputs (spaces/glob/control-bytes/PATH/symlink/cwd)** — Won't handle: n/a, prose-and-content change with no new argument or input surface (map Handoff item 6).
- **`package.json` `files[]` entries for files other than HANDOFF.md** — Won't handle: only the entry for the one file this spec deletes is touched; the rest of the package shape is FT26.

## Out of scope
- **FT26 first-install package fixes** — a separate capability (the npm package shape and link-time behavior, not the post-merge tail): ship the `projects/` profile templates, ignore-or-arch-check `.bench/dist` at link time, resolve the `.bench/bin` local CLI in the scaffolded gate, and add fresh-install smoke coverage. This spec removes only HANDOFF.md's `files[]` entry, because deleting the file forces it green. Estimate to build FT26: ~8 edits, ~4 gate runs.
- **A conformance pin for the exit-duty text and a CONTEXT-term anchor check** — a separate capability (hardening the gate, distinct from writing the prose): a Go check that asserts the duty prose and the `decision map` term exist, converting the review-only rows above into gate-observed ones. The map scoped this build prose-only with no Go changes, so it is deferred. Estimate: ~3 edits (check + canary fixture + wiring), ~3 gate runs.
