# Review pickup — pocock-guidance-doctrine

Composed diff `019078b6..76bf332c` plus ticket 10's uncommitted record files.
Raw axis findings: Standards 7, Spec 2, Coverage 3, cross-harness (gpt-5.6-terra
xhigh, falsification posture) 4 — raw total 16. De-duplicated repair targets: 6
auto-fix, 3 ask-user, remainder no-op. Terra's REFUTED verdict rests on reading
the spec-time who-writes-where fence as retired; PG11's "spec-authorized write
set" and `bench preflight`'s paths-authorized check keep it live, so the verdict
is treated as WEAKENED with its two concrete defects accepted.

## Auto-fix (deterministic rule or exact spec predicate, approved scope)

- AF1 `.bench/BENCH.md` — re-prune to ≤150 lines at the repo's ~77-column wrap;
  the current file meets the budget by wrap width (median line 127 cols; 1869
  words ≈ 170 lines at house wrap). Preserve every anchor needle, marker
  phrase, and the structured-phase clause block. [Spec must-fix]
- AF2 wide-line leaks — rewrap `bench-craft-grill/SKILL.md:11` (161 cols),
  `bench-craft-delegate/SKILL.md:71` (112), `bench-implement-spec.md:60` (99)
  to house width without exceeding their budgets. [Standards 4]
- AF3 `bench-review-implementation.md:22,112` — "ownership-fenced repair
  tickets" contradicts the demoted ticket fence (decision map #3: advisory
  `Writes:`, never refusal or review narrowing); reword to the slim ticket
  shape. [Terra 1, accepted part]
- AF4 `bench-debug.md:112` — "A spec-build write delegate" names the removed
  lifecycle; reword to "A write delegate". [Terra 2]
- AF5 `internal/conformance` — (a) Lstat-classify the `.agents/skills` root
  before ReadDir, fail closed on symlink/missing/non-directory, with a test
  [Terra 4]; (b) focused tests for the untested `subject missing` and `subject
  unreadable` diagnostics [Coverage 1]; (c) doc comment on `proseBudgetLines`
  starts with the symbol [Standards 6]; (d) reconcile `InputBenchkitProfile`'s
  doc comment with what the check actually reads [Standards 2]; (e) rename the
  mirror classifier's "user-invoked skill" wording — its only member is
  model-invocable [Standards 7].
- AF6 `CHANGELOG.md` (ticket 10 worktree) — drop the false "moved shell
  conventions" claim, name the four rules actually added, stop restating the
  budget numbers the profile owns. [Standards 1+3, Spec nit]
- AF7 `bench-craft-delegate/SKILL.md:120-122` — the repair-routing sentence
  points at its own section; state the route directly. [Standards 5]

## Ask-user (reviewer decision; not applied)

- AU1 spec edge inventory says overlapping exact/default classifications
  "refuse rather than pick a winner"; PG13 and the shipped checker make the
  exact row win. Code follows the reviewed decision; the inventory sentence is
  stale. Amend the spec sentence? [Coverage 2]
- AU2 `prototype` lacks `disable-model-invocation: true`; craft-skills defines
  user-invoked by that flag, but shape-idea's Prototype tickets need the model
  to charge it. Add the flag, or accept model-invocable? [Standards 7 corollary]
- AU3 `specs/pocock-guidance-doctrine/tickets/{04,11}` changed outside the
  spec's fence list — planning records of the build itself. Accept as workflow
  bookkeeping, or require a fence amendment? [Terra 3, Coverage nit]

## No-op (refuted, with derivation)

- Terra 1's wider claim: craft-spec's spec-time fence and bench-write-spec's
  fence section are required by PG11 (Coverage "enumerates ... the
  spec-authorized write set") and consumed by `bench preflight`
  paths-authorized; ticket-04/one-10 ticket files carry the schema of the
  doctrine that governed their build — historical records, not live guidance.
- Coverage write-set check, Spec axis clean list, terra's no-vanished-lessons
  sweep: recorded as clean.

## Open semantic rows

PG2, PG11, PG24 close only by fresh-session dogfood; artifact evidence
(anchors, index, symlinks, payload) is in the tree.
