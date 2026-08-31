# Grade a citation against the executed tag census

Blocked by: derive-executed-tag-census.md
Writes: internal/coverage/citations.go, internal/coverage/citations_test.go
Covers: CE1, CE2, CE4, CE5, CE6, CE7, CE20, CE21, CE22

## What to build

A seam-cell citation must point at a test file the gate executes on this host.
The coverage package imports the gate package for the executed tag census. It
then evaluates the cited file's build constraints against that census.

Evaluation uses the standard `go/build/constraint` package over the file's
`//go:build` line. It also applies the GOOS and GOARCH filename suffix rule. A
term is satisfied when a census set holds it, or when it equals the host's GOOS,
the host's GOARCH, or a satisfied release tag. The citation passes when any
executed set satisfies the file. The citation is a violation when no executed set
satisfies it.

The violation names the row number, the cited file, and the constraint that
failed. A malformed `//go:build` expression is a violation that names the file.
An empty census leaves the check inapplicable, so a non-Go root stays green.

Classify the cited path before any open. Use `bounds.ClassifyNoFollow`. A path
that is not a regular file is a violation, and the check never opens it.

The delivered resolution check keeps its current messages. This ticket adds the
execution arm beside it.

## Acceptance

- [ ] CE1 — a citation to a file whose build constraints no executed tag set
      satisfies is a violation.
- [ ] CE20 — a cited file with a malformed `//go:build` expression is a violation
      that names the file.
- [ ] CE21 — a graded root with no test phase leaves the execution check
      inapplicable.
- [ ] CE22 — a cited file behind a custom tag passes when the root's phase
      manifest declares that tag.
- [ ] CE2 — the unexecuted-constraint violation names the row number, the cited
      file, and the constraint.
- [ ] CE4 — a citation into the system suite passes on the kit root.
- [ ] CE4 — a cited file with no build constraint passes.
- [ ] CE5 — a citation to a stress-tagged file is a violation.
- [ ] CE6 — a cited file with a foreign GOOS filename suffix is a violation.
- [ ] CE7 — a cited path that is not a regular file is a violation, and the check
      performs no open on it.
- [ ] `bench gate` stays green.

## Delegate charge

You work in the Bench repo on the `citation-execution-proof` spec. Read
`specs/citation-execution-proof/spec.md` first. Then read
`internal/coverage/citations.go` and `internal/coverage/citations_test.go` in
full. Read the exported census function in `internal/gate/tag_census.go`. That
census is your only source of executed tags.

Extend `checkCitation` so a resolved citation must also prove execution.
Classify the cited path with `bounds.ClassifyNoFollow` before any open. Report
a violation for a path that is not a regular file.

Parse the file's `//go:build` line with `go/build/constraint`. Apply the GOOS and
GOARCH filename suffix rule. Satisfy a term from a census set, from the host GOOS or GOARCH, from a
`go/build.Default.ReleaseTags` entry, or from a toolchain-implied tag (`unix`,
`gc`, `cgo`). Strip `_test` from the filename before you read the GOOS suffix.

Report a violation
when no executed set satisfies the file. Name the row number, the cited file, and
the constraint in that message. Report a violation for a malformed constraint
expression.

Derive the census for the base the citations resolve against. Resolve the kit
through the gate's kit rule, so the cited files and the census come from one
tree. Return no violation when the census is empty.

Add a fixture case whose root carries a phase manifest with a custom `-tags`
set. Assert a citation behind that tag passes. This case reds a hardcoded
census copy. Keep the delivered resolution
messages unchanged.

Add `TestCitationUnexecutedConstraint` in `internal/coverage`. Write fixture
files for the system tag, the stress tag, and a windows filename suffix. Write
more fixtures for an untagged file, a malformed constraint, and a planted FIFO.
Assert the exact violation text.

Run only `bench worktree exec <label> -- go test ./internal/coverage/`. Do not
commit. Do not edit the spec.
