# Resolve cited test names and require the review-pickup fence member

Blocked by: none
Writes: internal/coverage, internal/preflight, internal/spec

## What to build

`bench coverage --check` grades two new facts about a mapped spec, so a
citation defect surfaces before the landing.

First, the check resolves each cited test name. A seam cell can cite a test
file and one or more backticked test names, in the shape
`internal/x/foo_test.go` (`TestName`). For each such citation, the check makes
sure the file exists and declares each cited `func TestName`. A cited name
whose leading segment is not a Go test function in that file is a violation
that names the row and the name. A row whose seam cell cites no `_test.go`
file adds no violation. The check holds no state, so a later run resolves a
renamed test again.

Second, the check requires the review pickup as a fence member. For a folder
spec `specs/<slug>/spec.md` that declares an `## Ownership fences` section,
the token list must contain `reviews/<slug>.md`. A missing member is a
violation that names the token.

The fence grammar keeps one parser. Move `fenceTokens` and its section
grammar from `internal/preflight` to `internal/spec`, and make both
`internal/preflight` and `internal/coverage` read fences through that one
function. `internal/preflight` imports `internal/coverage`, so `coverage`
must not import `preflight`. Because the shared parser ends the section at
any heading. So entries below a subsection do not count, and the pickup
check goes red with no extra code.

A historical spec stays exempt, exactly as it is for every other check.
Violation messages keep the existing `coverage map ...` grammar, because
downstream consumers match phrasings by substring.

## Acceptance

- [ ] A citation of a test name absent from the cited file is a violation that names the row and the name, at exit 1.
- [ ] A row that cites a declared test in an existing file adds no violation.
- [ ] A row whose seam cell cites no test file adds no violation.
- [ ] A folder spec with an `## Ownership fences` section that omits `reviews/<slug>.md` exits 1 and names the missing member.
- [ ] A fences section whose entries sit only below a subsection fails the pickup check, because the parser ends the section at that heading.
- [ ] `fenceTokens` has one definition, in `internal/spec`, and `internal/preflight` and `internal/coverage` both call it.
