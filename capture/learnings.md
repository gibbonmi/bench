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

