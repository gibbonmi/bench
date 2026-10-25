# The Markdown lane runs the fixture bite test

Blocked by: none
Writes: internal/gate/lane.go, internal/gate/lane_select.go, internal/gate/lane_select_test.go, internal/gate/lane_run_test.go, internal/conformance/registry/registry.go, internal/conformance/fixture_bite_test.go, projects/benchkit.md, CHANGELOG.md, specs/lp-lane-fixture-bite/tickets/01-the-markdown-lane-runs-the-fixture-bite-test.md
Covers: none

## What to build

A retained canary fixture quotes the physical line bytes of the paths that it pins. A reflow of a pinned guidance file can pass the commit lane and then red the landing gate in the fixture bite test. The kit lane gets one more built-in check, `fixture-bite`, which runs that test alone through `bench test --package`. A new path class, `canary-pin`, selects the check. The class claims each path that the composed checkout's canary inventory pins, and the lane reads those pins beside the embed targets. The registry owns the name of the test, and the lane and the proof-architecture test read that one name. The profile lane table advertises the new row, and the `profile-lane-table` check grades it.

## Acceptance

- [ ] A commit dry run of a reflow that breaks a fixture anchor fails the lane on `fixture-bite`.
- [ ] A Markdown commit that changes no pinned path does not run `fixture-bite`.
- [ ] The lane derives the pinned paths from the canary inventory of the composed checkout, not from a copied list.
- [ ] The profile lane table lists `fixture-bite` with its argv and the `canary-pin` class.
