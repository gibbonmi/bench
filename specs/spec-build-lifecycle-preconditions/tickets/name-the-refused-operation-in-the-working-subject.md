# Name the refused operation in the working-subject refusals

Blocked by: nothing
Ownership fence: `internal/specbuild/precondition.go`, and the `s.subject(...)`
call site in `internal/specbuild/assign.go`
Assumptions: `workingSubject` today owns two literals that both begin
`spec build start requires ...` — the no-working-branch refusal and the
dirty-checkout refusal — and `subject` is its only caller. `preconditions`
already receives the `mutation` it is acting for; `Start` calls `s.subject`
directly and must supply `mutationStart`. Re-derive both call sites from the
tree at pickup; a later ticket in this build edits the same two files.

## What to build

Every precondition-gated operation borrows `start`'s wording today, so a
refused `checkpoint` reports itself as a refused `start`. That misattribution
already cost a session: it read a `checkpoint` refusal as a `start` refusal and
concluded that parking the run was safe when it was not.

Thread the acting operation from `preconditions` through `subject` into
`workingSubject`, and interpolate it into both literals. The recomposition
refusal in `preconditions` is **not** part of this change: it keeps naming
`promote`, because `promote` genuinely is the operation that recomposes, and
pointing a refused `assign` at itself would send the reader away from the fix.

**The operation token has an absence case.** `mutation` is a bare string alias
with no zero-value guard, so a call site that forgets its argument renders
`spec build  requires a clean working checkout` — a refusal naming no
operation at all, which is strictly worse than today's wrong-but-present name.
Every `mutation` value must be one of the seven declared constants; make the
empty token unrepresentable at the message site rather than trusting call
sites, and prove it with its own criterion.

**Pin the rendered word here, at the junction.** The end-to-end enumeration
lives in `Enumerate the operation-named refusals through the real CLI`, three
tickets downstream — far enough that a disagreement about the spelling would
surface as a red in someone else's fence after a delegate round is spent. Each
`mutation` constant is spelled exactly as its `bench spec build` subcommand
verb, verified against `bin/bench.sh`'s usage block: `start`, `assign`,
`checkpoint`, `integrate`, `review`, `promote`, `abandon`. The refusal names the
**verb alone** — `abandon`, never `abandon --apply` — because the flag selects a
phase of the operation, not a different operation. Pin that as a criterion so
the downstream suite cannot disagree with this file.

## Acceptance

- [ ] OP1 — the dirty-checkout refusal names the acting operation for each of the seven `mutation` constants.
- [ ] OP2 — the no-working-branch refusal names the acting operation for each of the seven `mutation` constants.
- [ ] OP3 — an empty `mutation` never reaches either message as an empty name.
- [ ] OP4 — the recomposition refusal still names `promote` when the refused operation is `assign`.
- [ ] OP5 — each of the seven tokens renders exactly its `bench spec build` subcommand verb, with the abandon token rendering `abandon` and not `abandon --apply`.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| OP1 | interpolate the operation into the no-working-branch literal only, leaving the dirty-checkout literal hardcoded to `start` | `TestDirtyCheckoutRefusalNamesOperation` | for each of the seven constants, start a run in the fixture repo, dirty the checkout, invoke that operation, assert the refusal contains that operation and not `start` (except for `start` itself) |
| OP2 | interpolate the operation into the dirty-checkout literal only | `TestNoWorkingBranchRefusalNamesOperation` | for each of the seven constants, detach HEAD in the fixture repo, invoke that operation, assert the refusal names it |
| OP3 | pass `mutation("")` from one call site | `TestEmptyMutationNeverRendersAnonymously` | call the message site with the zero token and assert it does not produce a refusal containing `spec build  ` |
| OP4 | rewrite the recomposition refusal to interpolate the caller's operation | `TestRecompositionRefusalStillNamesPromote` | compose a run, advance the working branch, call `Assign`, assert the refusal names `promote` |
| OP5 | render the abandon token as `abandon --apply`, matching the coverage row's operator phrasing rather than the subcommand verb | `TestOperationWordMatchesSubcommandVerb` | for each of the seven tokens, render a refusal and assert the operation word equals the verb parsed from `bin/bench.sh`'s `bench spec build` usage lines, read at test time rather than restated |
