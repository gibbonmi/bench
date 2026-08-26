# Review pickup: harness-capability-seam

Frozen base `c35f1516`, reviewed tip `233baf80`. Raw findings: Standards 5,
Spec 5, Coverage 7. De-duplicated repair targets: 11 (6 auto-fix, 5
ask-user). Findings that collapse to one target carry the same letter.

## Standards

Count: 5. Worst: S1, the reference claims a one-source guarantee that three
live consumers refute.

- S1 `ask-user` (target K). `.bench/BENCH-reference.md:58-59` says
  `internal/harnesses` is the one source of the adapter list. AGENTS.md: two
  derivations of one fact must collapse into one source. The tree refutes the
  claim: `internal/gate/phases.go:281-284`, `internal/packagesurface/assets.go:23-25`,
  and `internal/conformance/entry_point_parity_test.go:81-89` each hold the
  three adapter paths and do not import `harnesses`.
- S2 `auto-fix` (target D). `internal/lines/lines.go:397` returns
  `is not provider-qualified (opencode model ids are provider/model)`. The
  rule now keys on `harnesses.AnyProvider`, so the message names a harness
  the rule no longer names. Tests match only `is not provider-qualified`.
- S3 `auto-fix` (target E). `internal/harnesses/harnesses.go` declares
  `srcAdapters = ".bench/adapters/"` and still spells
  `.bench/adapters/<name>` six times in `Headless` and the headless cell
  source. Build both from the constant (Duplicated Code).
- S4 `ask-user` (target L). `internal/conformance/harness_record_test.go:1236`
  re-derives the wiring substring rule that `internal/guards/guards.go:212`
  owns, and the comment at `:1204-1206` asserts sameness instead of
  composing it.
- S5 `ask-user` (target M). `internal/status/route.go` moves `HarnessClaude`
  from a `const` to a `var` computed by `claudeName()` over the action table.
  The answer depends on action-table order, and `route_test.go` pins the
  literal `claude` anyway (Global Data, Speculative Generality).

## Spec

Count: 5. Worst: P1, the `harness-record` oracle accepts a live symlink at an
adapter path that the spec's edge inventory refuses.

- P1 `auto-fix` (target A). Spec edge inventory: "A dangling symlink at an
  adapter path is absent, and a live symlink is refused and named."
  `harnessRecordHeadlessDiags` (`harness_record_test.go:257`) calls `exists()`,
  which is `os.Stat` (`checks_test.go:348`) and follows the link. No test and
  no code path covers the live half. Classify with the no-follow classifier,
  name the refusal, add the HC33-sibling test.
- P2 `auto-fix` (target B). HC09: "the `none` row has ... `headless
  execution: no`." `TestRecordWalk` never asserts
  `Rows[3].Mechanics[MechanicHeadless].Value == No`; a flip to `yes` stays
  green.
- P3 `ask-user` (target N) plus `auto-fix` (target C). Implementation
  decisions: "The lead's state and why never depend on the harness." HC17
  requires the lead to move for a formless harness, and `Route` does that.
  The doc comment at `route.go:76-77` restates the false universal (target
  C). The spec sentence needs an amendment to "the ladder is not re-ranked
  per harness" (target N).
- P4 `ask-user` (target O). Spec: "A `yes` or `no` cell carries a source."
  Every row holds `effort selection: no` sourced to the reference's effort
  rule, which states Bench's adapter surface, not the harness's own
  facility. The cells overclaim.
- P5 `no-op` folded into target N. HC28's seam label says
  `entry-point-parity row`; the closing test is
  `cmd/bench/main_test.go:282-296`. Amend the row's seam text.

## Coverage

Count: 7. Worst: C1, same as P1.

- C1 `auto-fix` (target A). Same as P1.
- C2 `ask-user` (target Q). `parityStaticDiags`
  (`entry_point_parity_test.go:181`) skips an unreadable static entry, so a
  deleted `scripts/release-preflight.sh`, `.agents/commands/bench.md`, or
  the workflow file leaves HC41 and HC42 ungraded with no diagnostic.
  `harness-record` reds an absent adapter; parity stays silent.
- C3 `auto-fix` (target F). `harness_record_test.go:147-148` emits
  `is not valid JSON`; no test plants `{"hooks":{`.
- C4 `auto-fix` (target F). HC35 names "a special file or symlink"; only the
  FIFO half runs (`:449`). A live symlink at `.codex/hooks.json` takes the
  same refusal branch and is untested.
- C5 `no-op`. A non-regular shim under `.bench/hooks/` is skipped by
  `parityEnumerationDiags:161` with no diagnostic. Undecided edge.
- C6 `no-op`. `harnessConfigNames` closes the config scan to
  `settings.json` and `hooks.json` inside a real directory. Undecided edge.
- C7 `no-op`. `TestRecordWalk` grades the `checked` date against
  `time.Now()`; backward clock skew reds the gate.
