# Review — cli-grammar-and-capability-evidence

Diff reviewed: `9732ebe..fde0b86` on `main` (FT87 slice 3, 78 files, +2792/−278).
Three axes, run as parallel read-only delegates.

**Round 1 fixes landed.** The four hard Standards violations, the empty-string
commit path, and their cross-axis duplicates are closed — the `help`-in-positional
regression, the three duplicated usage literals, the comment-satisfiable routing
check, the empty positional resolving to cwd, and the `t.SkipNow()` guard gap.
What remains below is unfixed.

## Standards

**7 findings remaining**, all judgment calls. Worst: `readSkipTally` silently
de-enforces strict mode on an unreadable log.

- `internal/gate/capability_skips.go:67-71` — `readSkipTally` discards
  `os.ReadFile`'s error and returns an empty tally; `TestSkipRowsStateZeroExplicitly:67`
  pins "absent log" to zero. Under `BENCH_REQUIRE_CAPABILITIES=1` that is
  enforcement, so a log the gate created but cannot read back reads identically to a
  fully capable runner — and the release workflows are exactly where that matters.
  `craft-gate`, "Choose the fail posture out loud": "Fail closed for enforcement…
  An unguarded pass-through is silent de-enforcement."
- `internal/models/models.go:70-72` — comment claims any `help` form "the grammar
  rejects with a usage line at exit 2"; `--help`, `-h`, and bare `help` all exit 0.
  `craft-comments`, "Aging".
- `internal/contract/axi/axi_roadmap_context_test.go:119` and
  `internal/contract/axi/axi_grammar_test.go:21` — story-number provenance and
  "now handles" narration; `craft-comments`, "The register".
- `internal/roadmap/roadmap.go:19` — doc comment does not open with the symbol name.
- Import grouping regressed in six files —
  `internal/canary/canary_concurrency_test.go:5`,
  `internal/conformance/gate_entry_test.go:4`, `internal/guards/guards_test.go`,
  `internal/outline/outline_test.go`, `internal/preflight/evidence_test.go`,
  `internal/preflight/integration_test.go` — the `internal/capability` import is
  interleaved alphabetically into a previously stdlib-only block, against the
  tree's separated-module-group convention (see `internal/contract/helper.go:3-12`).
- `projects/benchkit.md:136-186` — the new gate knob
  (`BENCH_REQUIRE_CAPABILITIES=1`) and three new conformance checks do not reach
  the profile, whose Gate section already documents the sibling knob
  (`BENCH_CANARY_INNER=1`) and names the conformance checks. Invariant #3.
- `internal/spec/spec.go:232`, `internal/worktree/list.go:18`,
  `cmd/bench/main.go:328` — three hand-rolled help parsers still stand, blessed
  permanently by `subcommand_routing_test.go:42`'s `whyNested` note ("each leaf owns
  its own grammar"). **Reviewer decision:** scoped follow-up or closed exemption.

## Spec

**5 findings remaining.** 27 coverage rows: 26 covered with a real red signal, 1
weak, 0 uncovered. Worst: six adopt subcommands exempted from the routing registry
with a factually wrong reason.

- `internal/conformance/subcommand_routing_test.go:78-89` — `setup`, `link`, `init`,
  `doctor`, `unlink`, `upgrade` are exempted under `whyNested` ("dispatches a
  subcommand tree rather than a flat argv"), but `internal/adopt/adopt.go:15-27`
  dispatches each to a leaf and `internal/adopt/doctor.go:181-189` is a hand-rolled
  flat `switch` whose default makes `bench doctor -h` exit 2. Fails stories 2 and 4;
  the fail-closed check the spec demanded is opened by an exemption reason nothing
  grades.
- `internal/spec/spec.go:229` (`specArg`) — no `--` case, so `bench spec implemented
  -- -x.md` is still rejected as an unknown flag. That is exactly the
  inexpressibility story 5 exists to close, on a subcommand whose sole argument is a
  path.
- `specs/cli-grammar-and-capability-evidence.md:220`, `:228` — the spec states the
  skip line is "written to stdout before the skip" and that the gate "tees that
  stream"; the build uses a `BENCH_SKIP_LOG` file side channel
  (`internal/capability/capability.go:90`, `internal/gate/capability_skips.go:3-11`).
  The reversal is correct (`go test` without `-v` discards package stdout) and is
  captured in `.bench/learnings.md:10`, but the staged spec now misdescribes what
  shipped. **Amend before the implemented flip.**
- `internal/usage/parse.go:80-82` — repeated-flag rejection is scope creep: not among
  the seven rules in "Grammar semantics, in full", it narrows previously-accepted
  invocations (`bench diff --commit a --commit b`), and is pinned by new AXI rows
  (`internal/contract/axi/axi_grammar_test.go:43-50`) no coverage row asked for.
  **Reviewer decision:** keep and amend the spec, or revert.
- Row 20 (story 10, concurrency) is weak by reinterpretation: the named wrong
  implementation ("an unguarded collector under `-race`") cannot exist under the file
  transport. `internal/gate/capability_skips_test.go:92` substitutes six bash phases
  appending cross-process — defensible, but the row's text is stale.

## Coverage

**6 findings remaining.** Worst: `gateEnv`'s skip-log isolation is unasserted, so a
canary fixture's skip can turn a real release red.

- `internal/gate/gate.go:150` — `gateEnv` strips `BENCH_SKIP_LOG` so a canary's inner
  run cannot append to the outer tally, but `gate_test.go:64-79` asserts only
  `BENCH_KIT`/`BENCH_WRAPPER`/`BENCH_GATE`. Delete the `capability.LogEnv` clause and
  the suite stays green; live consequence under `BENCH_REQUIRE_CAPABILITIES=1` (wired
  into both release workflows, `internal/conformance/workflow_checks_test.go:33,66`)
  is a canary fixture's skip turning a real release red.
- `internal/commit/commit.go:265-268` — a *deleted* named directory never enters
  `dirs` (`isDir` uses `Lstat` on the worktree), so `rm -r sub && bench commit -m x
  sub` reports every `sub/*` deletion as outside the named set — the one path the
  reviewer did name. The file-deletion case is covered; the directory case is not,
  and this slice introduced the directory rule.
- `internal/commit/commit.go:254-259` — `bench commit -m x .` can never succeed:
  `underAny` compares against `"./"` while porcelain paths are `internal/x.go`, so
  every file is listed as an offender. `.` is the obvious spelling of "everything I
  changed" and no test names it.
- `internal/usage/parse.go:83-88` — a value flag consumes a flag-shaped token with no
  look-ahead: `bench commit --spec -m "msg" file` sets `specSlug = "-m"` and fails at
  spec resolution; `bench commit -m --spec x file` commits with the message
  `--spec`. No case in `parse_test.go` gives a value flag a `-`-prefixed value.
  **Reviewer decision:** getopt-style consume (document it) or `MissingArg`.
- `internal/bounds/bounds.go:40-42` — the `inner < 0` clamp is never executed;
  `bounds_test.go:32-41` feeds only positive registry constants, so the guard can be
  deleted green. Zero is untested, and a duration above ~`MaxInt64/1.5` wraps to a
  negative deadline.
- `internal/capability/capability.go:108,110` — `Render` interpolates `Reason`
  verbatim. A reason containing `\n` produces a two-line record that
  `readSkipTally`'s `strings.Split(..., "\n")` counts as two skips
  (`internal/gate/capability_skips.go:73`); a reason pushing the line past 4096 bytes
  breaks the single-`Write` atomicity the design rests on (`capability.go:159-164`
  states the bound; nothing enforces it).

Note: the Edge inventory's FIFO dismissal ("git lists only tracked and untracked
regular paths") is wrong — `git status --untracked-files=all` does list an untracked
FIFO — but git refuses it at `add`, so the failure is loud. Flagged, not a finding.
