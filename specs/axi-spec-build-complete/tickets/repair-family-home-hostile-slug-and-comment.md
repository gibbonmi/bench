# Repair the family-home hostile-slug abort and its test comments

Blocked by: none
Ownership fence: `internal/specbuild/render.go`, `internal/specbuild/render_test.go`
Integration surfaces: none crosses
Contracts: none crosses
Closure: CM1/timeless-render-test-comments, FH1/family-home-survives-a-hostile-slug

## What to build

Close the accepted findings from the Terra/xhigh review of candidate
`72cab0bd3e1947b724d6c417cb412fa3220d22f8`: Standards S1, Spec P2, Coverage C2.

**S1 — comment register.** `internal/specbuild/render_test.go`'s two doc
comments (`TestRenderRefusalWithoutConstructibleRemedy`, "is the red mutation
for EP1"; `TestRenderRefusalWithConstructibleRemedy`, "unchanged from before
the fix") carry ticket-ID provenance and change-history language, against
`.agents/skills/bench-craft-comments/SKILL.md`'s timeless-present rule (the
same class of finding already repaired once this run for
`internal/specbuild/reclaim_test.go`). Fix: reword both comments to describe
only what each test proves, in the timeless present, naming no ticket ID,
review round, or "the fix". Test names, bodies, and assertions stay
byte-identical.

**P2/C2 — family home aborts instead of degrading.** `RenderRunsHome`
(`internal/specbuild/render.go`, the `if len(runs) > 0 { ... }` block right
before its final `return`) builds the status action from `runs[0].Slug`
through `axi.Fixed`. When that slug carries a control byte no single-line
command can quote (the same hostile-value class `EP1` just repaired for
`RenderRefusal`), `axi.Shell` refuses construction, and the function currently
does `return renderError(actionErr, ...)` — discarding `out`, the already-built
`spec_build_runs`/`spec_build_run_diagnostics` tables, and returning nothing
but a bare error at exit 1. One hostile slug in one retained run therefore
hides every other run's otherwise-renderable row, which is exactly the "never
a silent skip and never an abort that hides the healthy rows" failure the
spec's SB4 row and edge inventory (`specs/axi-spec-build-complete/spec.md`
around lines 24-26, 89) name for this family. Fix: when the action fails to
construct, do not abort — leave `actions` empty (the same "no useful action"
terminal form `axi.RenderHelp` already renders as `help[0]{command}:`) and
still return `out + axi.RenderHelp(actions), 0`. This mirrors exactly the
principle `EP1` already applied to `RenderRefusal`.

New coverage in `internal/specbuild/render_test.go`: drive `RenderRunsHome`
with a `[]Status` whose first (or only) entry carries a newline-bearing
`Slug`, alongside at least one further ordinary healthy `Status` and one
`RunDiagnostic`, and require the response still renders every row (the
hostile-sluged run's own row too — its slug still renders through TOON's
existing control-byte handling, only the *action* construction is refused)
plus `help[0]{command}:`, at exit 0 — never an error exit, never a dropped row.

## Acceptance

- [ ] [CM1] (covers local) (S1) `internal/specbuild/render_test.go`'s two test
  doc comments describe only current behavior, with no ticket ID, review
  round, or "the fix" language; test names, bodies, and assertions unchanged.
- [ ] [FH1] (covers local) (P2, C2) `RenderRunsHome` given a retained run whose
  slug cannot form a single-line action still renders every run and
  diagnostic row plus an empty `help[0]{command}:` block at exit 0, instead of
  discarding the whole response at exit 1; a run set with no hostile slug
  renders exactly as it does today, byte-for-byte.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CM1/timeless-render-test-comments | reinstate a ticket-ID or "the fix"/"before"-style change-history token in either comment | re-inspection against `.agents/skills/bench-craft-comments/SKILL.md`, with the existing focused render tests proving the edit is comment-only | read both comments against the skill's register rule and require them to name no ticket, review, or change; run `go test ./internal/specbuild/... -run TestRenderRefusal` and require it green with assertions untouched |
| FH1/family-home-survives-a-hostile-slug | restore the `return renderError(actionErr, ...)` abort for the unconstructible-action branch | focused `RenderRunsHome` test | render a run set containing one newline-slug run alongside healthy rows and a diagnostic, and require every row present plus `help[0]{command}:` at exit 0 rather than a bare error at exit 1 |
