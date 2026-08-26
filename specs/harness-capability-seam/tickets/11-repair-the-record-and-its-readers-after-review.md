# Repair the record and its readers after review

Blocked by: 06-add-the-bench-harnesses-verb.md, 09-document-the-harness-record.md
Writes: internal/harnesses/harnesses.go, internal/harnesses/harnesses_test.go, internal/harnesses/command_test.go, internal/lines/lines.go, internal/status/route.go, .bench/BENCH-reference.md, specs/harness-capability-seam/spec.md

## What to build

The review found seven small defects in the record, its readers, and the
docs. Each repair below is one accepted predicate from
`reviews/harness-capability-seam.md`.

Target B. `TestRecordWalk` asserts that the `none` row's `headless
execution` cell holds `no`, so a flip to `yes` reds.

Target C. The `Route` doc comment in `internal/status/route.go` no longer
says the lead never changes with the harness. It says the ladder is not
re-ranked per harness, and a formless harness skips a phase action.

Target D. The `lines.CellFault` message names no harness. It says the row
binds any provider, so the id must be provider-qualified.

Target E. `internal/harnesses/harnesses.go` builds every `Headless` path
and every headless cell source from the `srcAdapters` constant. No row
spells `.bench/adapters/<name>` as a literal.

Target K. The Files bullet in `.bench/BENCH-reference.md` says the record
holds one row per harness and names each headless adapter. It does not
claim to be the one source of the adapter list.

Target N. The spec's Implementation decisions say the ladder is not
re-ranked per harness. Only the command cell differs for a harness that can
run the lead. HC28's seam text says `wrapper parity test`.

Target O. Every row's `effort selection` cell holds `unknown` with no
source and no date. The reference's effort rule describes Bench's adapter
surface, not the harness.

## Acceptance

- [ ] `TestRecordWalk` reds when the `none` row's `headless execution` cell is `yes`.
- [ ] `rg -n 'opencode' internal/lines/lines.go` finds no message text.
- [ ] `rg -n '"\.bench/adapters/' internal/harnesses/harnesses.go` finds only the `srcAdapters` constant.
- [ ] `bench harnesses codex` shows `effort selection` as `unknown` with empty source and date.
- [ ] `.bench/BENCH-reference.md` does not say `one source of that adapter list`.
- [ ] The spec no longer says the lead's state and why never depend on the harness.
- [ ] `bench coverage --check specs/harness-capability-seam/spec.md` stays green.
