# Review pickup: asd-ste100-progressive-disclosure

Base `b1abde47`, reviewed tip `7641a055`. Three axes ran at `opus` on 2026-08-22.

## Standards

Raw findings: 4. Repair targets: 3. Worst issue: a second exclusion-row parser.

- S1 — `internal/conformance/prose_mechanics_test.go` re-implements the exclusion-row grammar that `internal/prose/exclusions.go` owns. Standard: `AGENTS.md`, parsers stay single-sourced. Disposition: auto-fix. Export one row-reading seam in `internal/prose` and read rows through it; the approved set stays an independent literal.
- S2 — `projects/benchkit.md` restates the two thresholds with no anchor row. Standard: `AGENTS.md`, an enforcement and its advertisement collapse into one source. Disposition: auto-fix. The paragraph points at `ste-prose.md` instead of repeating the numbers.
- S3 — `internal/prose/prose.go` package doc restates the numerals the constants own. Standard: `craft-comments`, Aging. Disposition: auto-fix. Name the bounds; drop the numerals.
- S4 — `.claude/output-styles/simplified-technical-english.md` restates the rules and is outside this diff. Disposition: no-op.

## Spec

Raw findings: 0. Repair targets: 0. Worst issue: none. Every checkable row PD1–PD44 is observed at the tip.

## Coverage

Raw findings: 10. Repair targets: 4. Worst issue: a code-span colon hides a long sentence.

- C1 — `internal/prose/parse.go` tests the label line before it folds code spans, so `` Run `foo: bar` `` with 40 words is not graded. Probe: a 40-word line with a colon inside a span returned no finding. Disposition: auto-fix, flagged for veto because it touches the parser order. Fold code spans before the label test; add the two cases (span colon; label line inside a paragraph run).
- C2 — `internal/prose/exclusions.go` accepts a directory row with no trailing `/` and the row excludes nothing. Probe: row `docs reason` left `docs/guide.md` graded with no diagnostic. Disposition: auto-fix. A row that names a directory without a trailing `/` reds as malformed; add the unit case.
- C3 — `internal/prose/parse.go` `codeSpanPattern` stops at an inner backtick, so a double-backtick span over-counts. Disposition: auto-fix. Match balanced backtick runs; add the case.
- C4 — two `Require` rows added since base have no fixture: the write-spec frozen base-and-tip row and the review-implementation universal-claim row. Disposition: auto-fix. Add one `files/`-form fixture each under `tests/canary/workflow-guidance-anchors/`.
- No-op: frontmatter `---` inside a value fails closed. Setext and list-item headings over-report. CRLF, `?!`, `…)`, `e.g.` at sentence end, deep `testdata/`, a file named `.md`, and `\r` in a row all behave per spec. `。` is not a boundary, and no won't-handle row names it.
