# Project one row shape across schemas

Blocked by: carry-a-schema-descriptor-through-the-parser.md, accept-the-reduced-coverage-header.md
Writes: internal/coverage/coverage.go, internal/coverage/coverage_test.go, internal/coverage/testdata

## What to build

An agent seeding implementation tasks from `bench coverage <spec>` gets one row
shape whichever schema the spec uses. The projection becomes
`rows[N]{story,behavior,seam}` — behavior is the cell that names what to build
and it exists under every accepted header, so a caller never branches on the
spec's schema. The `spec:` and `state:` lines and the AXI action list are
unchanged: a map with violations still renders one `retry after repairing
coverage map` action, and a clean map still renders one action per row.

The red already exists and must be honoured rather than regenerated.
`internal/coverage/testdata/` holds five checked-in `.stdout` goldens that
`TestCommandPreservesCheckedInPreDisclosureResponses` compares against, two of
which carry `rows[1]{story,seam,red_signal}` on line 3. Hand-edit those goldens to
the new header — never regenerate them by running the code, which would make them
agree with whatever the implementation does. Separately, `coverage_test.go`'s
`wantTable` builds its expectation by calling `toon.Table` with the
implementation's own column list, so it follows a projection change silently;
move `TestCommand`'s projection assertions to a literal expected block so a second
tautological assertion cannot mask a later change.

## Acceptance

- [ ] `(covers SR9)` `bench coverage <spec>` renders `rows[N]{story,behavior,seam}` for a
      six-column spec, with the behavior cell's text in the middle column.
- [ ] `(covers SR10)` The same command renders the identical TOON header for a reduced spec.
- [ ] `(covers SR9)` Both assertions compare against a literal expected block, not a value
      recomputed with the implementation's column list.
- [ ] `(covers SR9)` Reverting the projection to `{story,seam,red_signal}` turns both the
      literal assertions and the checked-in goldens red — demonstrated and recorded.
- [ ] `(covers SR9)` The five `testdata/*.stdout` goldens were edited by hand, not regenerated.
- [ ] `(covers SR11)` A violating spec still renders the single retry action; a clean one still
      renders one action per row.
