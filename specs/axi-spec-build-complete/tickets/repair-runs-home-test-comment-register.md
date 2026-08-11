# Repair the remaining change-history test comment

Blocked by: none
Ownership fence: `internal/specbuild/render_test.go`
Integration surfaces: none crosses
Contracts: none crosses
Closure: CM2/timeless-runs-home-test-comment

## What to build

Close the accepted Standards finding (S1) from the Terra/xhigh review of
candidate `721182b6ac70ce6e6f5dced0b70bb2f794b4182a`: the doc comment above
`TestRenderRunsHomeRendersStatusActionForOrdinaryRuns` in
`internal/specbuild/render_test.go` reads "renders its first run's status
action in the help block exactly as it always has" — "exactly as it always
has" is change-history narration (implying prior behavior differed), against
`.agents/skills/bench-craft-comments/SKILL.md`'s timeless-present register.
Fix: drop the trailing "exactly as it always has" clause so the comment
states only current behavior. Test name, body, and assertions stay
byte-identical.

## Acceptance

- [ ] [CM2] (covers local) (S1) the comment above
  `TestRenderRunsHomeRendersStatusActionForOrdinaryRuns` describes only
  current behavior with no change-history language; the test's name, body,
  and assertions are unchanged and stay green.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CM2/timeless-runs-home-test-comment | reinstate "exactly as it always has" or an equivalent change-history phrase | re-inspection against `.agents/skills/bench-craft-comments/SKILL.md`, with the existing focused test proving the edit is comment-only | read the comment against the skill's register rule and require it to name no prior-state comparison; run `go test ./internal/specbuild/... -run TestRenderRunsHomeRendersStatusActionForOrdinaryRuns` and require it green with assertions untouched |
