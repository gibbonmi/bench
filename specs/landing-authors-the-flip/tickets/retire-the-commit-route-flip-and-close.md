# Retire the commit-route flip and close

Blocked by: retire-names-the-board-remainder.md
Writes: internal/commit/commit.go, internal/commit/commit_test.go, internal/commit/close_test.go (renamed), internal/landing/landing.go, internal/landing/landing_test.go, internal/landing/state_test.go, internal/landing/close.go, internal/spec/spec.go, cmd/bench/main.go, cmd/bench/main_test.go

## What to build

`bench commit` loses `--spec`. The grammar refuses the flag as a usage error,
and the help text and the `bench help` commit row drop it. The landing request
loses its spec field. The landing owner composes only the attributed paths: the
spec transition branch, the close branch, and the close removal at reconcile
leave. The staged-check function leaves the spec package with its last caller.
The reviewed landing keeps its transition and its close.

The tickets-only predicate, the closed-folder path, and the index-tree removal
stay, and the comments that name a commit-path close are rewritten. The
blocker is a same-file edit of `internal/spec/spec.go`, not a behavior
dependency.

Delete the commit-route close tests and the landing owner's spec-request tests
with the behavior. Keep the commit fixture harness (`landingRepo`,
`runCommand`, `headPaths`) in one test file renamed for its new content,
because the exit-3 ticket drives it next. FA6 is a new test: the reviewed
close has no test today, and the deleted commit-route tests were its only
coverage. Drive `LandReviewed` with a close path and assert the published tree.

## Acceptance

- [ ] FA1: `bench commit -m m --spec x a.txt` exits 2 with the usage line and moves no ref.
- [ ] FA2: `bench commit --help` prints no `--spec`, and the `bench help` commit row reads without `--spec`.
- [ ] FA5: `LandReviewed` with a staged spec path still publishes `Status: implemented`.
- [ ] FA6: `LandReviewed` with a close path publishes a tree without that folder.
- [ ] FA7: `bench commit -m m a.txt` on a green gate publishes `a.txt` and exits 0.
- [ ] FA8: `bench commit -m m specs/<slug>` on a tickets-only folder publishes the folder's files, leaves the checkout clean, and exits 0.
