# Review Findings Persistence

Status: implemented

## Problem

`/bench-review-implementation` can surface actionable Standards, Spec, or
Coverage findings that the gate cannot see, but today those findings live only in
the chat handoff. If the session is compacted, interrupted, or resumed by a
different agent before the fix pass, the next `/bench-implement-spec` has to
reconstruct the review from memory or rerun it.

## Solution

Teach `/bench-review-implementation` to write a transient
`reviews/<spec-slug>.md` artifact only when the review finds actionable issues
that need a later fix. The artifact groups findings by axis, preserves the
citations and the worst issue per axis, and exists only until the green fix
commit deletes it. Clean reviews and reviewer-accepted residual risks stay in the
chat handoff and do not create a file.

## User stories

1. As the next implementer, I want actionable review findings persisted at
   `reviews/<spec-slug>.md`, so a resumed fix session can pick up the exact
   Standards, Spec, and Coverage work without replaying chat history.
   Line: claude-fable-5 / high. This changes the review command's guidance prose,
   a kit content surface covered by the project leverage override.
2. As the reviewer, I want the persisted file grouped by review axis with
   citations and the worst issue per axis, so the artifact keeps the same
   separation and evidence standard as the review handoff.
   Line: claude-fable-5 / high. Same command-prose surface; the quality bar is the
   semantic contract rather than code volume.
3. As the reviewer, I want clean reviews and accepted residual risks to avoid
   creating `reviews/<spec-slug>.md`, so the directory means "there is fix work to
   do" rather than "a review happened".
   Line: claude-fable-5 / high. This is a phase rule agents will follow repeatedly.
4. As the implementer closing findings, I want the green fix commit to delete the
   matching review artifact, so stale resolved findings cannot become a false
   pickup signal later.
   Line: claude-fable-5 / high. The deletion rule belongs beside the phase handoff
   and affects future agent behavior.
5. As an agent running a no-spec review, I want findings to remain chat-only unless
   the reviewer supplies an explicit slug, so reviews without a durable feature
   boundary do not invent misleading artifact names.
   Line: claude-fable-5 / high. The exception prevents guidance from creating a
   second source of feature identity.
6. As the kit maintainer, I want a conformance anchor that requires the review
   command to mention the artifact path, actionable-findings-only creation, and
   green-fix deletion, so future prose edits cannot silently drop the pickup rule.
   Line: claude-opus-4-8 / medium. This touches the oracle's conformance layer, so
   the implementation needs gate-level judgment even though the check is small.

## Implementation decisions

- `/bench-review-implementation` owns the write and retire guidance. No new CLI
  writer is introduced; the reviewing agent writes the markdown artifact as part
  of the existing phase when the conditions apply.
- The artifact path is `reviews/<spec-slug>.md`, where `<spec-slug>` is derived
  from the reviewed `specs/<spec-slug>.md` path. If the reviewer supplies a slug
  for a no-spec review, use that exact slug; otherwise no file is written.
- The artifact format is markdown with one section per axis:
  `## Standards`, `## Spec`, and `## Coverage`. Each section records the finding
  count, the worst issue for that axis, and each actionable finding with its file
  or doc citation. Empty axes may say `0 findings`; they are kept only when at
  least one other axis has actionable findings.
- Clean reviews and reviews where the human accepts all residual risks do not
  create `reviews/<spec-slug>.md`. A review artifact is a pickup file, not a review
  log.
- The fixing `/bench-implement-spec` session deletes the matching artifact in the
  same green commit that resolves the findings. The command prose should name this
  obligation so the pickup file cannot outlive the work it describes.
- The repo does not commit an empty `reviews/` directory or `.gitkeep`; the
  directory appears only when there is active review debt.
- The conformance layer pins the guidance with literal anchors rather than a new
  parser. The behavior is agent-facing process, so the oracle should catch the
  command losing the rule, not validate every possible review file shape.

## Testing decisions

- This is guidance-first kit content. The primary seam is the command prose
  consumed by agents; the gate-observable seam is the conformance workflow-anchor
  check that reads `.agents/commands/bench-review-implementation.md`.
- Add anchors in `internal/conformance/checks_test.go` for:
  `reviews/<spec-slug>.md`, actionable findings, no clean-review artifact, and
  deletion in the green fix commit. Add one canary fixture that removes the path or
  rule and proves the check bites.
- Focus manual red evidence on the current missing guidance. A direct probe before
  implementation already returns red:

  ```sh
  rg -q 'reviews/<spec-slug>\.md' .agents/commands/bench-review-implementation.md
  # exit 1 today
  ```

### Seam diagram

Seam 1 - review phase guidance (what agents follow):

    trigger: /bench-review-implementation
        |
        v
    branch diff + spec/standards/tests
        -> [ three-axis semantic review command ]
        -> chat handoff + optional reviews/<spec-slug>.md
              tests attach here by requiring the command to name the pickup rule

Seam 2 - conformance anchor (what the gate can enforce):

    trigger: bench gate / TestRootConformance
        |
        v
    kit tree
        -> [ checkWorkflowAnchors on bench-review-implementation.md ]
        -> diagnostics when the artifact guidance disappears
              tests attach here with a canary fixture that removes the anchor

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | actionable findings are persisted at `reviews/<spec-slug>.md` | command guidance + conformance anchor | already red: `rg -q 'reviews/<spec-slug>\.md' .agents/commands/bench-review-implementation.md` exits 1 | the command cannot instruct agents to create the pickup file without naming the path |
| 2 | persisted findings stay grouped by Standards, Spec, and Coverage with citations and worst issue per axis | command guidance + conformance anchor | conformance fixture removes axis/worst-issue guidance and must report the missing anchor | the check bites if the artifact stops preserving the review's evidence shape |
| 3 | clean reviews and accepted residual risks do not create a review artifact | command guidance + conformance anchor | conformance fixture removes the no-clean-artifact phrase and must report the missing anchor | the check prevents the pickup file from becoming a generic review log |
| 4 | a green fix commit deletes the matching review artifact | command guidance + conformance anchor | conformance fixture removes the delete-on-green-fix rule and must report the missing anchor | the check bites if resolved findings can remain as stale pickup work |
| 5 | no-spec reviews stay chat-only unless a slug is supplied | command guidance | not TDD-able as a runtime behavior; anchored prose names the exception | no durable spec slug means no deterministic path for a persistent artifact |
| 6 | future prose edits cannot drop the persistence contract silently | conformance anchor + canary bite-proof | new canary fixture must fail with the targeted missing-anchor diagnostic | a rotted or removed check fails the meta-test instead of becoming decorative |

### Edge inventory

- Spec path with directories or spaces: derive the slug from the basename of the
  `specs/*.md` file and keep the artifact under `reviews/`. The command prose
  should avoid shell snippets that would mishandle spaces.
- No spec path: no artifact unless the reviewer explicitly supplies a slug.
  Covered by story 5.
- Clean review: no artifact. Covered by story 3.
- Accepted residual risk: no artifact, because there is no later fix task.
  Covered by story 3.
- Mixed review with one actionable axis and two clean axes: create the file; keep
  all three axis headings so the reader sees the full review disposition.
- Existing stale artifact for the same slug: the next review overwrites or replaces
  it with the current findings, not appends a second log. **Won't handle** in code;
  this remains agent judgment because there is no writer yet.
- Findings without file-line citations: the artifact must preserve whatever
  citation the review axis supplied; low-evidence findings should be fixed during
  review aggregation, not normalized by the artifact.
- Missing trailing newline in the artifact: markdown readers tolerate it; no
  special handling.
- Invocation through different harnesses: the rule lives in the portable command
  markdown, so every harness that loads Bench commands sees the same guidance.

## Out of scope

- **Automated `bench review` writer or validator.** This spec only preserves the
  command-phase pickup rule. A CLI writer would need a real schema, parsing, and
  runtime contracts. Estimate: about 8 edits and 4 gate runs.
- **Persistent review history.** `reviews/<spec-slug>.md` is transient active debt,
  not an audit log.
