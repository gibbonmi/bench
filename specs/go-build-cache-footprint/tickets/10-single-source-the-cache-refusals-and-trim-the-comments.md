# Single-source the cache refusals and trim the comments

Blocked by: none

Writes: internal/gocache/command.go, internal/gocache/clean.go, internal/gocache/clean_test.go, internal/preprelease/preprelease_test.go, internal/gate/lane_test.go, internal/gate/report_test.go, internal/gate/cache_footprint_run_test.go, specs/go-build-cache-footprint/tickets/07-measure-the-footprint-after-trimpath.md

## What to repair

Review finding (Standards 1): `command.go` and `clean.go` each derive the
directory and each spell the two refusals. Collapse the derivation and its
refusals into one package function that both commands call.

Review finding (Standards 2): the T07 comment in `preprelease_test.go`
narrates the deleted producer. Keep the fact; drop the history.

Review finding (Standards 3), reviewer decision 2026-08-27: the comments on
the C12 test, the T03 test, and the R08 test argue the test's design. Trim
each to the fact the test pins.

Review finding (Standards 5): `clean_test.go` defines `equal` over string
slices. Use `slices.Equal`.

Review finding (Spec 2): record in ticket 07's checklist that the reviewer
moved C12 to the gate runner integration seam on 2026-08-27.

## Acceptance

- [ ] R14 — A directory path with an ESC byte makes `bench cache` refuse with a named reason.
- [ ] L10 — `bench cache clean` with no `go` on `PATH` refuses and names `go`.

Delivered outcome: one derivation and one refusal set serve both cache
commands, and the diff's comments state facts only.
