# Review pickup: ft311-diagnostics

Frozen base: `f41917040c24637ca0c5d69441398361559ce319`
Reviewed tip: `9eed6663e1c2fc9ae656b1ffe9903f1393fa10bc`
Axes: Standards, Spec, and Coverage, each an opus medium read-only delegate on 2026-09-09.
Raw findings: 15. Repair targets after the collapse: 9. Ticket 7 carries the accepted repairs.

## Standards

Count: 8. Worst issue: two prose diagnostic sentences are derived a second time in `internal/gate/gate_prose.go` while the same diff single-sourced their sibling.

- `auto-fix` — AGENTS.md "one source per fact". `internal/gate/gate_prose.go` restates the unreadable-exclusion-file and unreadable-subject sentences that `internal/prose/exclusions.go` and `internal/prose/subject.go` own. Export the two diagnostics from `prose` beside `absentExclusionDiagnostic`.
- `auto-fix` — AGENTS.md test-expectation exception, condition two. `probeHelpNotes` in `internal/probe/refusal_test.go` and `rootConformanceTestName` in `internal/testreport/selection_facts_test.go` are independent expectations with no recorded red. The coordinator ran the two probes on the reviewed tip and recorded the verdict rows in ticket 7.
- `auto-fix` — Duplicated Code (judgment call). `anchorsFixtureNeedles` in `cmd/bench/anchor_help_test.go` is a fourth hand copy of the six AGENTS.md needles. Derive the list from the registry entries that name AGENTS.md, intersected with the fixture's lines.
- `no-op` — Duplicated Code. `sectionRunesMapped` re-walks headings and fences. Spec line "the evaluator does not change" fences the only other route. The comment-strip divergence this walk hides is a Spec finding below.
- `auto-fix` — `internal/gate/gate_prose.go` restates the `.md` subject predicate that `internal/prose/walk.go` owns. Export the predicate from `prose`.
- `auto-fix` — Global Data. `RegularFileModes` in `internal/git/staged.go` is an exported mutable map with one caller in its own file. Unexport it.
- `auto-fix` — Duplicated Code. `runesEqual` and `hasPrefixRunes` in `internal/anchors/locate.go` restate `slices.Equal`.
- `auto-fix` — craft-comments register. The PB5 comment in `internal/probe/probe_test.go` narrates ("now refuses at the baseline").

## Spec

Count: 3. Worst issue: the selection row's `baseline` cell is a constant, not the observed kind.

- `auto-fix` — Spec "The `baseline` cell is the observed baseline outcome kind, and it joins the row with the baseline run." `internal/probe/probe.go` `render` prints `string(testreport.OutcomePassed)`; `gradeBaseline` discards the passing outcome. Thread the observed kind from `gradeBaseline` into `render`. DG10 is partial until then.
- `auto-fix` for the divergence, `ask-user` for the single source. Spec: "HTML comments are stripped as the evaluator strips them." `StripHTMLComments` rescans its rewritten text, and `stripCommentsMapped` scans once. So a comment that a strip creates (`<!<!-- z -->-- needle -->`) hides the needle from the evaluator and not from the locate step. Repeat the mapped strip until no comment remains, and add the rejoined-comment case to the locate test. The full single-source refactor is parked as an idea, because the spec forbids an evaluator change.
- `no-op` — Spec Git-flag note. `ReadStagedIndex` adds `rev-parse --verify --quiet HEAD`, and DG37 already attaches the unborn-branch case in two green tests.

## Coverage

Count: 4. Worst issue: the help's own example, `bench gate-prose . --staged`, refuses on every relative root.

- `auto-fix` — Relative root. `IsWorkTreeTop` in `internal/git/staged.go` compares the absolute top against the operand as given, so `.` never matches. Absolutize both sides, and add a relative-root case through `GateProseCommand`. New row DG40.
- `auto-fix`, flagged for veto — Index-blob bounds. The named form refuses an invalid UTF-8 or oversized subject through `bounds.ClassifyNoFollow`; the staged form reads `git show` unbounded and parses it. Apply the same limit through `bounds.Read` and the same UTF-8 rule to the staged subject and policy blobs, with the named form's diagnostic sentences. New row DG41.
- `auto-fix` (test only) — ESC in a staged path. The staged form answers the shared unrepresentable-cell refusal at exit 1 and names no path. That is the named form's documented choice, so the code stays; the row is missing. New row DG42 pins the refusal.
- merged into the Spec comment-strip finding above. New row DG43 pins the rejoined comment.
