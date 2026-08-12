# Review pickup — axi-query-disclosure

Candidate: commit 173bffdcb3e454869801933ea4f66df01a47c0d0 (branch main, clean tree).
Snapshot: /tmp/axi-query-review-173bffdc.toon (captured once, candidate-bound).
Sources: three fresh independent axis passes (REJECT, four unique actionable
findings) plus the rerun bounded read-only Fable cross-harness pass (ACCEPT,
nine advisory findings). Every citation below re-verified against the tree at
the candidate commit. Raw findings 13; actionable 11; de-duplicated repair
targets 11 (R1 and R2 name different fixes on a shared surface and are expected
to land as one restructuring ticket).

## Standards

Count: 6 raw, 4 actionable. Worst: R1.

- **R1 (P1, ask-user)** — `cmd/bench/command_registry_test.go` spawns 21
  disposable Git repositories in ordinary untagged tests: `newAXIEnvelopeRepo`
  (lines 228–241, ~7 `gitpkg.Raw` process calls each) runs in 3 of 6 subtests ×
  7 members, plus `git worktree add` in `setupAXIWorktree`. The architecture
  census (`projects/benchkit.md:198-209`) budgets repositories to the tagged
  system package and does not count process-backed fixtures reached through
  `internal/git`. Constraint: the census allowlist must not be widened to make
  this pass; whether the census rule's scope covers this package-main test is
  the reviewer call.
- **R2 (ask-user)** — fixture-harness duplication: the minimal TempDir +
  `git init` scaffold is derived independently in `learningsRepo`
  (`internal/learnings/learnings_test.go:135`), `mapsRepo`
  (`internal/maps/maps_test.go:132`), twelve inline `git init` blocks in
  `internal/guards/guards_test.go`, and `newAXIEnvelopeRepo`
  (`cmd/bench/command_registry_test.go:228`), beside the pre-existing
  `newWorktreeRepo`. AGENTS.md code standard: "a fixture harness pasted N times
  must collapse to one source." Collapsing needs a cross-package test-support
  seam the spec's fences did not authorize.
- **R3 (auto-fix)** — dead branch: `internal/worktree/list.go:29`
  `if parsed.EndedFlags` is unreachable because `worktreeListGrammar` sets
  `HelpOnlyWhenSole` and `Result.EndedFlags` is assigned only inside the block
  that flag skips (`internal/usage/parse.go:96-104`); observed `--` behavior
  comes from the positional path.
- **R4 (auto-fix)** — `HelpOnlyWhenSole` under-documents its effect: it
  disables all flag recognition and the `--` terminator, not only mid-args
  help; the `Grammar` doc comment (`internal/usage/parse.go:23-40`) does not
  say so.
- no-op: identical split case arms `internal/usage/parse.go:68-71`; the
  `"none"` wired-cell literal derived at both `internal/guards/guards.go:225`
  and `:309` (same file, low drift risk).

## Spec

Count: 4 raw, 4 actionable. Worst: R5 (High, QD3).

- **R5 (High, QD3, auto-fix)** — envelope assertions are marker substrings,
  not structured validation: `cmd/bench/command_registry_test.go:121,133,136`
  use `strings.Contains` and `hasAXIHelpEnvelope` (line 373) is a line-prefix
  scan. QD3 (spec.md:59) requires "structured TOON stdout … and the appended or
  honest-empty `help[N]{cmd,why}` block"; malformed or trailing non-TOON stdout
  stays green. Repair: decode complete stdout with the official `toon-go`
  decoder and require a terminal schema-correct help table per member.
- **R6 (Medium, QD6, auto-fix)** — the learnings "unreadable file" QD6 state is
  exercised only through the oversized read-error route
  (`internal/bounds/classify.go:81`; fixture `ControlRecordLimit+1`); the
  open/stat class — e.g. a dangling symlink, correctly classified unreadable by
  `classify.go:75,125-129` — has no public-surface old-to-new pair, and QD6
  (spec.md:39,62) forbids one state borrowing another's oracle without a named
  alias.
- **R7 (ask-user)** — guards stale/unwired QD6 evidence never touches a real
  stale manifest: `internal/guards/guards_test.go:84-105` injects rows through
  the `enumerateGuards`/`inspectGuard` package variables with test-authored
  cells, against spec.md:38 "real fixtures (a stale guard manifest, …)".
  Mitigation: the consumer predicate (`guards.go:309`) shares
  `adopt.PrePushStale` (`internal/adopt/link_hook.go:55`) with the production
  emitter, so constant drift is impossible; only the row shape is stub-borne.
- **R8 (ask-user)** — coverage disclosure `why` mislabels the story cell as
  "row": `internal/coverage/coverage.go:444` builds
  `"check coverage row "+row[0]` where `row[0]` is the story column; traced
  execution emits `check coverage row 1,2` for map row QD6 (story "1,2"). No
  spec line pins the wording; fixing re-cuts candidate-side fixture bytes.

## Coverage

Count: 3 raw, 3 actionable. Worst: R11.

- **R9 (Medium, auto-fix)** — the public worktree terminal pair
  (`internal/worktree/list_actions_test.go:99-121`) covers only a present
  foreign registration; an owned completed assignment appears solely in the
  unit-level `actionsForRows` table (line 126), so a `listAssignmentRow`
  regression for `StateComplete` stays green at the public surface.
- **R10 (auto-fix)** — the maps bounds-classified invalid branch
  (`internal/maps/maps.go:58`) has no test pinning its carried path: mutate it
  to pass `mapName` instead of `candidate.Path` and all tests stay green while
  the disclosed repair path is wrong. Missing: a public-command case for a
  bounds-invalid (empty/oversized/invalid-UTF-8) map file.
- **R11 (ask-user)** — active assignment with a deleted tree is an undecided
  state class: `listAssignmentRow` (`internal/worktree/list.go:107`) renders
  state `active` / tree `missing`, and `actionsForRows` (`list.go:84-90`)
  advertises `path`/`exec` into the nonexistent tree and never `clean`
  (`orphanPath` is set only for foreign registrations). QD6 names
  active/orphaned/empty/non-actionable without deciding this hybrid; defining
  the class changes locked coverage, so it requires a reviewer-approved spec
  amendment, not a build-time widening.
