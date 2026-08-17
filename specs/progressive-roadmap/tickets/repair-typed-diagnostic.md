# Repair: one typed Diagnostic instead of a formatted-string round-trip

Blocked by: none
Writes: internal/roadmap/tree.go, internal/roadmap/context_parse.go, internal/roadmap/tree_validation.go, internal/roadmap/tree_test.go

## What to build

Review findings Standards F2, Coverage C4 (`reviews/progressive-roadmap.md`).

`internal/roadmap/tree.go` currently formats each integrity diagnostic as a
plain `"<path>: <reason>"` string (7 `fmt.Sprintf` call sites), and
`internal/roadmap/context_parse.go:189` re-derives the path by
`strings.Cut(d, ": ")` — lossy for a legal basename containing `": "` (e.g. a
row under `roadmap/` named so that `roadmap/x: y.md` is the unrecognized-file
path; the cut yields `source = "roadmap/x"`, a path that does not exist).

Add a small `Diagnostic{Path, Reason string}` type in the roadmap package with
a `String()` returning today's exact `"<path>: <reason>"` format (so
`ValidateRoadmapTree`'s `[]string` return and every existing diagnostic string
assertion in tests keep matching byte-for-byte). `ParseDocument` builds and
returns `[]Diagnostic` internally (or returns both forms — pick whichever
keeps `ValidateRoadmapTree`'s public signature as `[]string` for the
conformance registry binding, converting once at that boundary).
`context_parse.go` consumes `Diagnostic.Path`/`Diagnostic.Reason` directly
instead of re-parsing the formatted string.

## Acceptance

- [ ] One diagnostic path+reason pair is authored once (in the `Diagnostic` construction) and consumed by both the conformance check and the context renderer — no second parse of the formatted string.
- [ ] A diagnostic whose path contains `": "` renders the correct `parse_failures.source` in `--context` output (add a test: an unrecognized file basename containing `: `).
- [ ] Every existing diagnostic-string assertion (conformance canary `EXPECT` values, `tree_test.go` diagnostic checks) still passes unchanged — the formatted string is byte-identical to today's.
- [ ] `go test ./...` and `bench gate` stay green.
