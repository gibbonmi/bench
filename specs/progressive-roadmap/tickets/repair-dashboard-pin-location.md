# Repair: move the story-21 regression pin to its declared fence

Blocked by: none
Writes: internal/dashboard, internal/status/status_test.go

## What to build

Review finding Spec #1 (`reviews/progressive-roadmap.md`). Coverage row PR18
/ story 21's seam is "dashboard render test" and `internal/dashboard/` is a
declared ownership fence in `specs/progressive-roadmap/spec.md`, but
`TestDashboardRoadmapTextAndSequenceRenderFromSplitTree` landed in
`internal/status/status_test.go` and asserts `roadmap.RoadmapText` /
`roadmap.RecommendedSequence` directly rather than through
`dashboard.Snapshot` (or whatever the dashboard's own render entry point is
— read `internal/dashboard/dashboard.go` to find it). Move the test (or add
an equivalent one) into `internal/dashboard`, exercising the dashboard's own
rendering path against a split-tree fixture, so a dashboard-side reader swap
would red it. Remove the test from `internal/status/status_test.go` once its
replacement lands (don't leave the same behavior pinned twice).

## Acceptance

- [ ] A test in `internal/dashboard` renders roadmap text and recommended sequence from a split tree through the dashboard's own render entry point.
- [ ] `internal/status/status_test.go` no longer carries the PR18/story-21 pin.
- [ ] `go test ./...` and `bench gate` stay green.
