# Review pickup: ft311-landing-completion

Base `265bce07aafc93667d086628927b1b893f90297c`, reviewed tip `067f5bc5a0e2297df1ed4a9717494c58935877d6`, three native axes at opus/medium on 2026-09-09. The cross-harness falsification pass is skipped under the closed FT311 decision.

## Standards

Count: 6 raw findings, 5 repair targets. Worst: the build edited the spec's fences and four tickets' `Writes:` lines.

- `specs/ft311-landing-completion/spec.md` fences and tickets 1, 2, 4, 5 — the build added six fence entries and widened four `Writes:` lines. The spec reserves that edit to spec-change authority. Each addition closes an existing reader of the landed record or folds one source, so it grants no new authority. Disposition: `ask-user`. The reviewer vetoes or keeps the amendments.
- `internal/worktree/land_effects_test.go` near `wantEffects` — the comment claims the bytes are independent of the owner's values, but the result cells compose from the production constants. Disposition: `auto-fix`. Reword the comment to name what is independent.
- `internal/worktree/land_effects_test.go` `brokerChangingLanding` and `brokerDestinationFixture` — the same five-statement Go main fixture skeleton is written twice. Disposition: `auto-fix`. One helper writes it.

## Spec

Count: 5 raw findings, 0 code repairs, 1 veto flag. Worst: row LC40 promises the final-check guidance states all three surfaces, but the spec's Guidance section replaces one sentence there.

- `specs/ft311-landing-completion/spec.md` row LC40 — the row text and the Guidance decision disagree, and the tree follows the decision. Disposition: `ask-user`. The amendment narrows the row text to the decision.
- `specs/ft311-landing-completion/spec.md` row LC34 — the seam names `TestRetroScaffoldListsTheTickets`, and the behavior is asserted by `TestRetroScaffoldListsUnknownWithoutATicketsDirectory`. Disposition: `auto-fix`. The amendment names the test that covers the row.

## Coverage

Count: 4 raw findings, 3 repair targets. Worst: a state file whose fence never closes passes the heading refusal, and the next run cannot parse the document.

- `internal/handoff/state_file.go` — a draft such as `prose`, a blank line, an open fence, then `## main` writes at exit 0. The next `bench handoff` then refuses the document on the open fence. Disposition: `auto-fix` inside story 29's reason, flagged for veto. Refuse an unclosed fence in the draft through the document parser's own fence rule, and add row LC44.
- `internal/worktree/land_effects.go` `cleanLandedSiblings` — an empty destination base reports `cleanup` `complete` with no plan and no stderr row. The edge table reads an empty narrowed set as complete, and a repository-wide sweep is a Won't handle. Disposition: `ask-user`. The tree keeps `complete`; the reviewer decides whether a stderr line or a `not-applicable` word is wanted.
