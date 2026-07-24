# Review — cli-grammar-and-capability-evidence

Diff reviewed: `9732ebe..fde0b86` on `main` (FT87 slice 3, 78 files, +2792/−278).
Three axes, run as parallel read-only delegates at `gpt-5.6-terra` / medium.

## Standards

**11 findings.** Worst: the `help` positional rule at `internal/usage/parse.go:62`
silently drops reviewer capture, which is the one thing the capture seam promises
never to do.

Hard violations:

- `internal/usage/parse.go:62` — a bare `help` anywhere before `--` is treated as a
  help request, breaking the `bench idea` seam contract in `projects/benchkit.md:63`
  ("`idea` appends one dated line"). Verified against `dist/bench`: `bench idea help
  me remember the parser` prints usage, exits 0, and creates no `IDEAS.md`. Pre-diff
  it parked the line. `MaxArgs: -1` grammars (`idea`, `commit` paths) need `help`
  recognized only at position 0, or not at all.
- `cmd/bench/main.go:90` and `:100` — `usage: bench commands --brief` is written
  twice, once as `commandsGrammar.Help` and once as a literal in the `--brief`-absent
  branch. AGENTS.md, "Code standard — one source per fact". Same defect at
  `internal/roadmap/roadmap.go:25` and `:52`.
- `internal/commit/commit.go:45` vs `:156` — same duplication, and the copies have
  already drifted: the grammar `Help` reads `-m <msg> [--spec <slug>] [--]
  <path>...`, the usage-error line omits `[--]`, so the error path hides the `--`
  escape this diff added.
- `internal/conformance/subcommand_routing_test.go:218-232` — `packageReaches` is a
  substring search for `usage.Parse` over non-test `.go` files, and every routed
  package's new grammar doc comment contains that literal (e.g.
  `internal/status/status.go:40`). Deleting the real call and keeping the comment
  leaves the check green. The bite proof at `:291` writes `var _ = usage.Parse` —
  real code — so it never exercises the vacuous state. `craft-gate`, "Prove it
  bites".

Judgment calls:

- `internal/gate/capability_skips.go:67-71` — `readSkipTally` discards
  `os.ReadFile`'s error and returns an empty tally; `TestSkipRowsStateZeroExplicitly:67`
  pins "absent log" to zero. Under `BENCH_REQUIRE_CAPABILITIES=1` that is
  enforcement, so an unreadable log reads identically to a fully capable runner.
  `craft-gate`, "Choose the fail posture out loud".
- `internal/models/models.go:70-72` — comment claims any `help` form "the grammar
  rejects with a usage line at exit 2"; `--help`, `-h`, and `help` all exit 0.
- `internal/contract/axi/axi_roadmap_context_test.go:119` and
  `internal/contract/axi/axi_grammar_test.go:21` — story-number provenance and
  "now handles" narration in comments; `craft-comments`, "The register".
- `internal/roadmap/roadmap.go:19` — doc comment does not open with the symbol name.
- Import grouping regressed in six files —
  `internal/canary/canary_concurrency_test.go:5`,
  `internal/conformance/gate_entry_test.go:4`, `internal/guards/guards_test.go`,
  `internal/outline/outline_test.go`, `internal/preflight/evidence_test.go`,
  `internal/preflight/integration_test.go` — the `internal/capability` import is
  interleaved alphabetically into a previously stdlib-only block, against the
  tree's separated-module-group convention (see `internal/contract/helper.go:3-12`
  in this same diff).
- `projects/benchkit.md:136-186` — the new gate knob
  (`BENCH_REQUIRE_CAPABILITIES=1`) and three new conformance checks do not reach
  the profile, whose Gate section already documents the sibling knob
  (`BENCH_CANARY_INNER=1`) and names the conformance checks. Invariant #3.
- `internal/spec/spec.go:232`, `internal/worktree/list.go:18`,
  `cmd/bench/main.go:328` — three hand-rolled help parsers still stand, blessed
  permanently by `subcommand_routing_test.go:42`'s `whyNested` note ("each leaf owns
  its own grammar"). Reviewer decision: scoped follow-up or closed exemption.

## Spec

**6 findings.** 27 coverage rows: **26 covered with a real red signal**, **1 weak**,
**0 uncovered**. Worst: story 4's help rule applied to variadic free-text
positionals loses reviewer input at exit 0.

- `internal/roadmap/roadmap.go:46` + `internal/usage/parse.go:62` — `bench idea help
  <text>` swallows the idea (see Standards, same root). Story 4's rule was written
  for flat query commands; a variadic free-text positional needs an exception.
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
  shipped. Amend before the implemented flip.
- `internal/usage/parse.go:80-82` — repeated-flag rejection is scope creep: not among
  the seven rules in "Grammar semantics, in full", it narrows previously-accepted
  invocations (`bench diff --commit a --commit b`), and is pinned by new AXI rows
  (`internal/contract/axi/axi_grammar_test.go:43-50`) no coverage row asked for.
- Row 20 (story 10, concurrency) is weak by reinterpretation: the named wrong
  implementation ("an unguarded collector under `-race`") cannot exist under the file
  transport. `internal/gate/capability_skips_test.go:92` substitutes six bash phases
  appending cross-process — defensible, but the row's text is stale.

Note (not a finding): row 26's second clause ("run with the constant temporarily
changed") is not automated anywhere, though
`internal/worktree/lifecycle_test.go:32-37` does derive from `bounds.LeaseStale`.
Nothing from **Out of scope** was built.

## Coverage

**10 findings.** Worst: an empty-string path argument silently widens a commit to the
whole cwd.

- `internal/usage/parse.go:98-101` + `internal/commit/commit.go:272-282` — `bench
  commit -m msg ""` (an unset shell variable) parses `""` as a positional, `rootRel`
  maps it through `filepath.Abs("")` → cwd; from a subdirectory `isDir` is true,
  `underAny` explains every dirty file beneath, and `git add -A -- :(literal)sub`
  stages all of them (`commit.go:228`, `:242`, `:131`). No test passes an empty argv
  token. Milder sibling: `-m ""` reaches `git commit -m ""` and fails with a raw git
  error instead of a usage line (`commit.go:142`).
- `internal/gate/gate.go:150` — `gateEnv` strips `BENCH_SKIP_LOG` so a canary's inner
  run cannot append to the outer tally, but `gate_test.go:64-79` asserts only
  `BENCH_KIT`/`BENCH_WRAPPER`/`BENCH_GATE`. Delete the `capability.LogEnv` clause and
  the suite stays green; live consequence under `BENCH_REQUIRE_CAPABILITIES=1` (now
  wired into both release workflows, `internal/conformance/workflow_checks_test.go:33,66`)
  is a canary fixture's skip turning a real release red.
- `internal/conformance/skip_ownership_test.go:23` — `skipMethods` is
  `{"Skip","Skipf"}`; `t.SkipNow()` is the third `testing.TB` method that ends a test
  without running it and walks straight through the guard. The bite proof (`:99-137`)
  never covers it.
- `internal/commit/commit.go:265-268` — a *deleted* named directory never enters
  `dirs` (`isDir` uses `Lstat` on the worktree), so `rm -r sub && bench commit -m x
  sub` reports every `sub/*` deletion as outside the named set — the one path the
  reviewer did name. The file-deletion case is covered; the directory case is not,
  and this slice introduced the directory rule.
- `internal/commit/commit.go:254-259` — `bench commit -m x .` can never succeed:
  `underAny` compares against `"./"` while porcelain paths are `internal/x.go`, so
  every file is listed as an offender. `.` is the obvious spelling of "everything I
  changed" and no test names it.
- `internal/usage/parse.go:62` — a help spelling in *positional* position
  short-circuits (`bench commit -m msg help` prints help instead of committing a file
  named `help`). Every case in `parse_test.go:29-37, 190-202` puts the spelling in
  flag position. Same root as the Standards/Spec headline.
- `internal/usage/parse.go:83-88` — a value flag consumes a flag-shaped token with no
  look-ahead: `bench commit --spec -m "msg" file` sets `specSlug = "-m"` and fails at
  spec resolution; `bench commit -m --spec x file` commits with the message
  `--spec`. No case in `parse_test.go` gives a value flag a `-`-prefixed value.
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
- `internal/conformance/subcommand_routing_test.go:218-233` — the routed-claim check
  is satisfied by a comment (same finding as Standards, reached independently).

Note: the Edge inventory's FIFO dismissal ("git lists only tracked and untracked
regular paths") is wrong — `git status --untracked-files=all` does list an untracked
FIFO — but git refuses it at `add`, so the failure is loud. Flagged, not a finding.

Probes that came back genuinely covered: `--` and second `--`, bare `-`, repeated
flags across `--`, glob/space filenames on commit and outline, sibling-directory
prefix leakage (`subdir` vs `sub`), zero-vs-absent skip log, cross-writer skip-log
interleaving, the exactly-at-the-window lease case.
