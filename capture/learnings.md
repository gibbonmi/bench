# Learnings — usage journal

- 2026-08-20 — FT230 spec review took 2 iterations. Stage that missed: spec
  authoring. What review caught: two step-name byte contracts on the exact
  bytes the spec deletes (`native_workflow_test.go`'s platform-first check and
  the `preflight-publish-order-bypassed` canary anchor), plus a coverage row
  citing "existing state-machine tests" that do not exist in the live build
  (`canary_shared_test.go` compiles only under tags with no fixture). Why it
  was missed: the author swept `internal/conformance` for anchors on
  `release.yml`'s publish steps but only inside `workflow_checks_test.go`, and
  cited a test seam without reading the package's test files this session —
  the read-before-cite rule was applied to code but not to the tests a
  coverage row names. Proposed rule change: `craft-spec`'s map discipline
  gains one line — a row whose seam is "existing tests" must name the test
  function, verified read this session; and a spec that deletes literal bytes
  runs one repo-wide `rg` for those bytes (canary fixtures included) before
  the map is written.

- 2026-08-20 — A mid-build spec amendment (two fence paths) cost far more than
  the two lines it changed. What happened: `bench worktree land` refuses when
  source and destination staged spec bytes differ, so the FT230 fence
  amendment forced a hand mirror commit on the destination — and because every
  `bench commit` runs the full gate (~3 min), the amendment alone cost two
  extra full gate runs (the destination handoff commit made a third) before
  land would proceed. The right behavior: an amendment that only the source
  carries should reach the destination through one cheap, sanctioned step, not
  through hand `git show > file` mirroring plus full-gate commits. Proposed
  rule change (reviewer decides the shape): either (a) `bench worktree land`
  adopts the source's spec bytes as the published truth — which is what the
  review skill already implies ("the landing publishes the source's spec
  bytes, so an amendment never routes through a hand commit on the
  destination") — and drops the identical-bytes refusal, or (b) a
  `bench spec sync <slug>` plumbing verb copies the source's spec directory to
  the destination and commits it with a scoped gate (spec files cannot break
  Go phases, so a docs-scoped fast path would do). Also worth pricing: a
  capture-only fast lane in the gate for commits touching only
  `capture/`/`specs/*.md`, since three of this build's six destination gate
  runs graded pure-markdown diffs.

