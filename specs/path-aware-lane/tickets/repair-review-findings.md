# Repair the accepted review findings

Blocked by: add-document-family-checks.md, record-the-selection-in-guidance.md
Writes: docs/adr/0017-the-worktree-commit-runs-the-fast-lane.md, CHANGELOG.md, internal/gate/lane_select.go, internal/gate/lane_select_test.go, internal/gate/lane_test.go, internal/gate/authorization/lane_test.go, internal/packagesurface/assets.go, internal/packagesurface/assets_test.go, specs/path-aware-lane/spec.md, reviews/path-aware-lane.md

## What to build

The review at tip `fc39350f` accepted six repair targets. This ticket closes
them on the integration source. The predicates are S1, S2, S4, and S5 from
the Standards axis, P1 from the Spec axis, and F3 and F4 from the Coverage
axis. `reviews/path-aware-lane.md` states each one.

The ADR and the changelog say what the lane does not run: no `vet` check and
no `build` check on a Markdown commit. The file comment in
`internal/gate/lane_select.go` states what the file holds and not why it is
there. One test asserts that every path class names only checks the kit lane
declares. `packagesurface.EmbedTargets` treats an absent `cmd/` or
`internal/` as empty, and the fixtures that planted a placeholder Go source
for it lose that source.

The spec's PL5 edge line drops the glob clause. `TestSelectLaneByClass` gains one mixed-class row and four prefix-boundary
rows. The resolved findings leave `reviews/path-aware-lane.md`; the
`ask-user` findings stay.

## Acceptance

- [ ] S1: neither guidance file says "no Go toolchain"; each says the commit runs no `vet` and no `build` check.
- [ ] S2: the file comment carries no placement rationale.
- [ ] S4: a test reds when a class names a check the kit lane does not declare.
- [ ] S5: `EmbedTargets` over a root with no `cmd/` returns the `internal/` targets, and no lane fixture plants a placeholder Go source.
- [ ] P1: the PL5 edge line names a space and a non-ASCII byte only.
- [ ] F3: entries for `a.go` and `bin/x.sh` return every declared check and the classes `go-source,unknown`.
- [ ] F4: `roadmapx/a.md` and `docs/ROADMAP.md` return `prose` alone; `specs/decisions/x.md` returns `prose` alone; `capture/retros` as a file returns every declared check and `unknown`.
- [ ] the gate `test` phase stays green for `internal/gate`, `internal/packagesurface`, and `internal/conformance`.
