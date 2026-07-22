# Review — consumer-payload-and-phase-contract

Three-axis pass over `feebe4a..7c9d79f`. Most findings are closed; this file
carries only what is still open.

## Standards

0 findings open. The triplicated allowlist parser/matcher is collapsed onto the
exported `kitpayload` helpers, and the change-narrating comment is gone.

One item is deferred to the reviewer, repo-wide rather than per-diff: the
`FT<n> story <n>` provenance tags conflict with `craft-comments`, which forbids
provenance, but they match roughly a dozen pre-existing sites. Deciding it once
for the whole repo is cheaper than re-litigating it in each review.

## Spec

1 finding open.

**Story 2's coverage row is phrased for a behavior the spec does not want.**
The row's red signal reads "an asset added to the tree but absent from the
allowlist is not written", which implies per-file allowlisting under
`.agents/skills/`. The spec's own edge inventory
(`specs/consumer-payload-and-phase-contract.md:290`) says the opposite —
tree-typed rows walk their directory by design — and story 2's text is a
refactor to one canonical allowlist, not a granularity change. The landed test
(`internal/contract/surface/link_lifecycle_test.go:74-88`) asserts the only
version of that claim the design supports: a file under a kit-only subtree is
not written.

Reading: the row's phrasing is imprecise, not the implementation. The fix is a
row rewording, not code — but rewording a spec row is the reviewer's call, so
nothing was changed. If the reviewer instead wants per-file granularity under
allowlisted trees, that is a new story with its own design, not a defect fix.

## Coverage

0 findings open in this pass's scope.

Both blockers are closed and verified end-to-end: the packed tarball now carries
`consumer_payload.go` and `.bench/consumer-payload.json`, and `go build
./cmd/bench` succeeds inside a real extracted tarball; the kit-only guards now
fail on an emptied allowlist, held by the `kit-only-allowlist-emptied` canary
fixture. `RequiredBuildPackAssets` derives root-level Go sources and every
source's `//go:embed` targets, so a future root package or embed joins the
expectation without a second registry.

Five findings outside this pass's scope are parked in `IDEAS.md` for the next
`/bench-what-next` drain — three confirmed (`bench upgrade` prerelease→release
is a silent no-op; a symlink inside an allowlisted tree hard-fails link and
upgrade for every consumer; `.bench/skills-index.sh` parses the allowlist with a
key-order-dependent `sed` that word-splits space-bearing sources) and two
suspected (duplicate and traversal allowlist rows; `bench upgrade` argument and
manifest-state edges). The symlink one is the most consumer-visible of the three.
