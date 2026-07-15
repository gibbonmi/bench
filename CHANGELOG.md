# Bench changelog

The synthesis record. `/bench-update-kit` reads this as its baseline — what Bench
already adopted from upstream, and what it deliberately rejected — so closed
decisions stay closed and each re-synthesis is diffed against a known state. Append
one entry per `/bench-update-kit` run or learnings-sourced promotion (queued by
`/bench-what-next`, built under the synthesis discipline); don't rewrite history.

## Unreleased

- **Learnings run (2026-07-15, scope: post-worktree-lifecycle drain).** Reworded
  FT58 after the lifecycle build shipped its lock-protocol half — a live owner
  is never aged out, competing reclaimers serialize through the
  rename-and-identity-check takeover, and a successor's lease survives release,
  each with a red-capable test — leaving the pool-root hardening residue
  (propagating tighten failures, symlink/non-owned-root rejection, post-create
  mode revalidation); the row demotes MEDIUM→LOW and now records that the tree
  asserts best-effort tighten, a fork the build must put to the reviewer.
  Reworded FT94 after the resume golden grew a fourth site
  (`lifecycle_policy_test.go`). IDEAS.md was empty, so zero ideas drained.
  Promoted the one open journal entry to FT96, a rule-shaped
  `/bench-write-spec` edit naming the fast path for reviewer-closed in-session
  forks (recorded straight into a decision map, flagged in-spec for veto). The
  previous run promoted no rule edits, so no promotion-build entry was due. The
  refreshed sequence keeps FT80 and FT81 first and slots FT96 third as an XS
  synthesis promotion.

- **Roadmap reconcile (2026-07-15, scope: post-FT93-ship drain).** Removed FT93
  after both halves landed — the retained-verdict surfacing at `bf31967` and the
  (b)/(c) orphan reconcile at `4d6c290` — and retired its merged spec
  (`worktree-orphan-reconcile`). Reworded FT91 to the measured post-perf
  numbers: the parallelized gate now runs ~50s wall-clock (was ~4 min), so the
  row demotes MEDIUM→LOW and graduates only if 50s demonstrably still drags
  shift iteration; the remaining rows verify current and the FT6, FT24, FT8,
  and FT38 triggers remain unmet. Drained one idea into new row FT94: the
  `bench resume` summary line is asserted as a verbatim golden in three test
  files — test-vs-test duplication the one-source standard does not protect —
  so a shared expected-format helper single-sources the literal. Dismissed the
  one open journal entry (FT93(b)/(c) reconcile-vs-preserve) as resolved: the
  shipped spec encodes exactly the conditional compact-residue /
  report-preserved fork it proposed, and it proposed no rule change. The
  previous run promoted no rule edits, and the FT65 and FT66 promotion-build
  records remain present. The refreshed sequence leads with FT80 and FT81
  through `/bench-shape-idea`, then FT94 as a small lighter-path cleanup.

- **Learnings run (2026-07-14, scope: post-FT90-fix reconcile).** Removed FT90
  after its fix landed at `19fa3bc` — `dist/bench` is out of the declared gate
  inputs and the gitignored-declared-input conformance check plus bite fixture
  are in the tree; the other 17 rows verify current and the FT6, FT24, FT8, and
  FT38 triggers remain unmet. Drained two ideas: FT91 (gate wall-clock
  proportional to the diff; shape-first because the speed-versus-oracle cut is
  a reviewer decision) and FT92 (attributed subject-drift message plus shipping
  the gitignored-input check as consumer scaffolding). Verdicted the one open
  journal entry by reproducing its red through the accused command pair itself:
  one gitignored file makes the automatic plan retain, the retain path
  completes without a terminal release receipt, and `bench worktree release
  --request` masks the retained verdict as "terminal receipt missing" —
  promoted to FT93 (with the request-less-clean orphaned-assignment follow-on).
  The entry's alternative hypothesis — that the `--request` pair is misplaced
  plumbing — is dismissed: the pair is the documented `bench worktree --help`
  usage and absent from the plumbing enumeration, so no inventory change. The
  previous run promoted no rows, and the FT65 and FT66 promotion-build records
  remain present. The refreshed sequence leads with FT93 through `/bench-debug`
  (repro in hand), then FT80 and FT81 through `/bench-shape-idea`.

- **Roadmap reconcile (2026-07-14, scope: post-FT79-ship drain).** FT79's row
  left with its spec retirement (`eafb3fb`); the remaining 17 rows verify
  current and the FT6, FT24, FT8, and FT38 triggers remain unmet. Drained one
  idea into new row FT90: the gate subject fingerprints `dist/bench` while the
  gate's own serialized build phase rewrites that binary, so the first gate
  after any source change discards its verdict as subject drift. Code-confirmed
  in `internal/gate/subject.go` plus `.bench/gate-inputs.json`; the row's
  wording corrects the parked line's gitignore framing, since the exposure is
  the declared manifest path, not ignore handling (`TreeHash` already respects
  `.gitignore`). The journal had zero open entries, so there were no verdicts.
  The previous run promoted no rows, and the FT65 and FT66 promotion-build
  records remain present. The refreshed sequence leads with FT90 through
  `/bench-debug` (small fix, removes per-change gate friction), then FT80 and
  FT81 through `/bench-shape-idea`.

- **Learnings run (2026-07-14, scope: FT79 decision-map reconcile).** No roadmap
  row shipped; the only tree change since the prior reconcile is FT79's
  completed decision map, so all 18 rows remain current and the FT6, FT24, FT8,
  and FT38 triggers remain unmet. IDEAS.md was already empty, so zero ideas
  drained, and the journal had zero open entries, so there were no verdicts. The
  previous run promoted no rows, and the existing FT65 and FT66 promotion-build
  records remain present. The refreshed sequence is FT79 through
  `/bench-write-spec` (its decision map is complete), then FT80 and FT81 through
  `/bench-shape-idea`.

- **Roadmap reconcile (2026-07-13, scope: post-FT78 retirement and cold-session
  continuity).** Removed FT78 after `bench spec history
  oracle-bound-gate-verdicts` returned its retirement record at `3fa8dec`; the
  other 18 rows remain current. Drained one idea by folding a copy-paste,
  harness-native cold-session continuation prompt contract into FT89 rather
  than adding a duplicate guidance row. The journal had zero open entries, and
  the FT6, FT24, FT8, and FT38 triggers remain unmet. Previous promotion builds
  remain recorded. The refreshed sequence is FT79, FT80, then FT81 through
  `/bench-shape-idea`.

- **Learnings run (2026-07-13, scope: FT78 proof-closure reconcile).** Reworded
  FT78 from implementation to proof closure against its revised staged spec; no
  roadmap row shipped, and the other 18 rows remain current. The FT6, FT24, FT8,
  and FT38 triggers remain unmet. IDEAS.md was already empty, so zero ideas
  drained. Dismissed the wrong-harness-handoff learning because the canonical
  communication rule already requires translating Bench phase recommendations
  to this harness's `$bench-*` adapter; this was an application miss, not a kit
  gap. The previous run promoted no rows, and the existing FT65 and FT66
  promotion-build records remain present. The refreshed sequence is FT78 through
  `/bench-implement-spec`, then FT79 and FT80 through `/bench-shape-idea`.

- **Learnings run (2026-07-13, scope: post-FT78 staging reconcile).** Reworded
  FT78 from shaping to implementation and linked its staged
  `oracle-bound-gate-verdicts` spec; no roadmap row shipped, and the other 18
  rows remain current because the only tree changes since the prior reconcile
  are FT78's decision maps and spec. The FT6, FT24, FT8, and FT38 triggers
  remain unmet. IDEAS.md was already empty, so zero ideas drained, and the
  journal had zero open entries, so there were no verdicts. The previous run
  promoted no rows, and the existing FT65 and FT66 promotion-build records
  remain present. The refreshed sequence is FT78 through
  `/bench-implement-spec`, then FT79 and FT80 through `/bench-shape-idea`.

- **Learnings run (2026-07-12, scope: post-FT77 retirement reconcile).**
  Removed FT77 after `bench spec history ownership-safe-worktree-cleanup`
  returned its retirement record at `2f38f1c`; reworded the release-readiness
  lead and reassessment gate to describe only findings that remain open. The
  other 19 rows remain current, and the FT6, FT24, FT8, and FT38 triggers remain
  unmet. IDEAS.md was already empty, so zero ideas drained, and the journal had
  zero open entries, so there were no verdicts. The previous run promoted no
  rows, and the existing FT65 and FT66 promotion-build records remain present.
  The refreshed sequence is FT78, FT79, then FT80 through `/bench-shape-idea`.

- **Learnings run (2026-07-12, scope: post-FT77 implementation reconcile).**
  No roadmap row shipped: `bench spec history
  ownership-safe-worktree-cleanup` returned no retirement record. Reworded
  FT77 from implementation to semantic review because its spec is implemented
  on the default branch; whether a clean semantic review already happened is
  the batch's veto surface because a finding-free review leaves no artifact.
  The other 19 rows remain current, and the FT6, FT24, FT8, and FT38 triggers
  remain unmet. IDEAS.md was already empty, so zero ideas drained, and the
  journal had zero open entries, so there were no verdicts. The previous run
  promoted no rows, and the existing FT65 and FT66 promotion-build records
  remain present. The refreshed sequence is FT77 semantic review, then FT78
  and FT79 through `/bench-shape-idea`.

- **Learnings run (2026-07-12, scope: release-readiness roadmap reconcile).**
  Reworded the release-status lead so its NO-GO decision no longer masquerades
  as a feature row, restoring schema parsing without changing the decision.
  Verified all 20 roadmap rows against the authoritative snapshot; none has a
  shipped spec or history record, and the FT6, FT24, FT8, and FT38 triggers
  remain unmet. IDEAS.md was already empty, so zero ideas drained, and the
  journal had zero open entries, so there were no verdicts. The previous run
  promoted no rows, and the existing FT65 and FT66 promotion-build records
  remain present. The refreshed sequence remains FT77, FT78, then FT79 through
  `/bench-shape-idea`.

- **Learnings run (2026-07-11, scope: clean-state reconcile).** Verified all six
  roadmap rows against the authoritative snapshot; none shipped or needed
  rewording. FT71, FT6, and FT58 still lack their evidence triggers, FT24
  remains upstream-blocked, and the dated FT38 and FT8 revisits are not yet
  actionable. IDEAS.md was already empty, so zero ideas drained, and the
  journal had zero open entries, so there were no verdicts. Previous
  promotion-build records remain complete: FT65 and FT66 are present. The
  refreshed recommended sequence remains the dated FT38 and FT8 revisits.

- **Learnings run (2026-07-11, scope: post-FT72–FT75 reconcile).** Removed
  FT72, FT73, FT74, and FT75 after their responsibility splits shipped; the
  authoritative structure snapshot no longer reports any of their recorded
  file violations, and the cached gate is current and green. No parked,
  scheduled, tabled, or upstream trigger has fired. IDEAS.md was already empty,
  so zero ideas drained, and the journal had zero open entries, so there were
  no verdicts. Previous promotion-build records remain complete: FT65 and FT66
  are present. The recommended sequence is now the dated FT38 and FT8 revisits.

- **Learnings run (2026-07-11, scope: post-structure-wave reconcile).** Removed
  FT62, FT64, FT69, and FT70 after their responsibility splits shipped at
  `3e92085`; the authoritative structure snapshot now reports only
  FT72/FT73/FT74/FT75, each still reproducing at its recorded limit. No parked,
  scheduled, tabled, or upstream trigger has fired. IDEAS.md was already empty,
  so zero ideas drained, and the journal had zero open entries, so there were
  no verdicts. Previous promotion-build records remain complete: FT65 and FT66
  are present. The recommended sequence is the remaining structure wave,
  followed by the dated FT38 and FT8 revisits.

- **Learnings run (2026-07-11, scope: post-FT68 retirement reconcile).** Removed
  FT68 after `bench spec history structured-phase-conversation` confirmed its
  retirement at `8071257`. All eight structure rows still reproduce at their
  recorded limits, and no parked, scheduled, or upstream trigger has fired.
  IDEAS.md was already empty, so zero ideas drained, and the journal had zero
  open entries, so there were no verdicts. Previous promotion-build records
  remain complete: FT65 and FT66 are present. The recommended sequence now runs
  the structure passes in two priority-preserving waves: FT62/FT64/FT69/FT70,
  then FT72/FT73/FT74/FT75.

- **Learnings run (2026-07-11, scope: post-FT68 implementation + structure
  reconcile).** No roadmap row shipped: `bench spec history
  structured-phase-conversation` returned no retirement record. Reworded FT68
  from implementation to final check because its implementation and semantic
  review are complete, and added FT75 for the untracked
  `internal/conformance/docs_workflow_helpers_test.go` structure violation at
  449/400. The seven existing structure rows still reproduce at their recorded
  limits; no parked, scheduled, or upstream trigger has fired. IDEAS.md was
  already empty, so zero ideas drained, and the journal had zero open entries,
  so there were no verdicts. Previous promotion-build records remain complete:
  FT65 and FT66 are present. The recommended sequence is now FT68 final check,
  then the eight structure passes.

- **Learnings run (2026-07-11, scope: post-retirement + structure reconcile).**
  Removed FT67 after `bench spec history mandatory-implementation-delegation`
  confirmed its retirement at `3da2e7d`; FT68 remains staged. Added FT72
  (`internal/git/git.go` 403/400), FT73
  (`internal/worktree/worktree_test.go` 489/400), and FT74
  (`internal/intent/intent.go` 406/400); FT62/FT64/FT69/FT70 still reproduce at
  their recorded limits. Rechecked FT24 against the current Codex hooks surface:
  delegation still has no deny-capable tool hook and `SubagentStart` cannot stop
  the spawn, while its common model field replaces the stale no-model wording.
  The other parked and scheduled rows have not met their triggers. IDEAS.md was
  already empty, so zero ideas drained. Dismissed the grill-fog learning because
  `craft-grill` already owns the behavior and the entry identifies an application
  miss, not a kit gap. Rendered the journal's sample heading inline so it no
  longer parses as a malformed entry. Previous promotion-build records remain
  complete: FT65 and FT66 are present. The recommended sequence is now FT68,
  then the seven structure passes.

- **Learnings run (2026-07-11, scope: staged-spec + structure reconcile).**
  Reconcile found three staged specs with no roadmap rows: roadmap-context
  shipped in the big CLI update without a `--spec` mark, so this pass retires
  it as the backstop; the two unbuilt specs became FT67 (mandatory
  implementation delegation) and FT68 (structured phase conversation), ranked
  above structure debt because FT67 changes how every later build runs. Two
  new structure violations became FT69 (`binary_repair_test.go` 416/400) and
  FT70 (`lifecycle_test.go` 405/400); FT62's "one of two issues" wording was
  corrected. Drained one idea to FT71, parked pending evidence: kit-sized
  shift-session log, graduating on an observed need for durable run evidence.
  Dismissed the pre-authorized-grill learning — reviewer instruction governs
  and the entry proposes no rule change; codify a named mode only if it
  recurs. The learnings malformed-source flag in `bench roadmap --context` is
  designed behavior (the template placeholder is preserved as evidence), not a
  defect. Previous run's promotion builds checked: FT65 and FT66 build entries
  are present, so none is missing. The recommended sequence is now FT67, FT68,
  then the four structure rows.

- **Learnings run (2026-07-10, scope: post-FT66 reconcile).** Verified all seven
  remaining roadmap rows against the tree; none shipped or needed rewording, and
  FT62/FT64 still reproduce exactly at 441/400 and 407/400. The parked and
  scheduled rows have not met their graduation triggers. IDEAS.md was already
  empty and the learnings journal had zero open entries, so nothing drained or
  required a verdict. The recommended sequence remains FT62, then FT64. The
  preceding FT66 promotion-build entry records the previous run's build, so no
  promotion record is missing.

- **FT66 promotion build (2026-07-10).** Learnings-origin review guidance folded
  into the existing semantic-review process: each round now ends by integrating
  accepted findings, running focused checks, and running one final gate; another
  semantic review round opens only on a red gate or reviewer request. Focused
  root conformance and the full gate closed the promotion. Semantic review found
  the rule lacked a bite proof; a collapsed-space conformance anchor and targeted
  canary now keep its full sequence and reopening condition live.

- **Learnings run (2026-07-10, scope: post-FT65 reconcile).** Removed FT65 as
  shipped; the remaining eight roadmap rows still match the tree, including the
  two exactly reproduced structure issues, and the current Codex hook surface
  does not graduate FT24. IDEAS.md was already empty, so zero ideas drained;
  the learnings journal had zero open entries. The recommended sequence is now
  FT66, FT62, then FT64. Previous run's promotion builds checked: the FT65 entry
  below records its build; FT66 remains unbuilt, so no other entry is due.

- **FT65 promotion build (2026-07-10).** Learnings-origin defect folded into
  the existing AXI diff seam: branch-relative files and body now include
  committed, index, and tracked worktree changes from the resolved base while
  the log remains commit-only and `--commit` remains exact. The accused
  `bin/bench.sh diff --full` repro and end-to-end worktree contract went
  red→green; staged-state and exact-commit guards passed, followed by the full
  gate.

- **Learnings run (2026-07-10, scope: immediate no-op reconcile).** Verified all
  nine roadmap rows against the unchanged post-reconcile tree; none shipped or
  needed rewording, and the FT62/FT64 structure debt still reproduces exactly.
  IDEAS.md was already empty, so zero ideas drained; the learnings journal had
  zero open entries. The recommended sequence remains FT65, FT66, then FT62.
  Previous run's promotion builds checked: FT65 and FT66 remain unbuilt, so no
  build entry is due; the preceding XS promotion build entry is present.

- **Learnings run (2026-07-10, scope: post-XS-promotion reconcile).** Reworded
  FT62 after the tree grew from one to two structure issues and added FT64 for
  `internal/canary/canary_test.go`; FT62 remains open. IDEAS.md was already
  empty, so zero ideas drained. Dismissed the emitted-CLI-output entry because
  `craft-synthesis` already classifies every CLI touch as behavior and requires
  dogfood. Reproduced the dirty-default-branch diff defect through exact
  `bin/bench.sh diff --full` (`files[0]` beside a non-empty pinned Git diff) and
  promoted it to FT65; fixing the canonical CLI instead of adding a guidance
  fallback is the veto surface. Promoted the terminal repair-pass bound to FT66
  as a rule-shaped review-guidance edit. Previous run's promotion builds
  checked: the XS synthesis entry below records FT55/56/57/59/63; no build entry
  is missing.

- **XS synthesis promotion (2026-07-10, FT55/56/57/59/63).** Five rule edits
  establish that structure split-vs-grant checks both file and directory
  budgets; sanctioned `bench commit`, `--spec`, and `bench spec history`
  contracts sit at their points of use; shared-worktree delegates pin file
  tools and diagnose path/CWD mismatches; fix-pass delegates verify a
  commit-specific sentinel; and specs delivering new phase commands land with
  the command while the stale-command sweep stays fail-closed. The Codex
  top/mid/cheap binding is Sol (`gpt-5.6-sol`), Terra (`gpt-5.6-terra`), and Luna
  (`gpt-5.6-luna`); Claude projects those tiers through the `fable`, `opus`, and
  `sonnet` aliases from the same binding. Codex shifts select `workspace-write`,
  a dogfood-discovered gap now fixed. Focused binding and adapter contracts went
  red→green, then a fresh linked consumer Luna shift took its gate red→green and
  committed exactly one file only after green. A push-triggered hook audit was
  blocked by the outer guard.

- **Learnings run (2026-07-10, scope: post-FT51-54 reconcile).** FT51-54
  removed as shipped and their four specs retired (cli-hygiene, docs-batch,
  test-hardening, assess-owner). Drained one parked idea to FT62
  (helper.go structure debt, split-or-grant under the FT55 rule). One journal
  entry promoted to FT63: staged-spec posture for new phase commands —
  recommended as a `/bench-write-spec` rule (spec lands with its command)
  rather than a sweep exemption; the alternative stays flagged for veto.
  A second entry (interactive inspection piped rg through `head` against the
  AGENTS.md read-tool rule) dismissed: the existing rule is explicit, no kit
  change — the correction is behavioral. Previous run's promotion builds
  checked: FT60 backfilled last run; FT55/56/57/59 remain unbuilt, no entries
  due.

- **Learnings run (2026-07-09, scope: post-FT60 reconcile).** Reconcile-only
  pass: IDEAS.md empty, no open journal entries. FT60 removed as shipped (the
  serialized gate build phase plus two review pickups). FT61 added from the
  tree: three files over the structure budget, two grown by the FT60 build.
  Previous run's promotion builds checked: FT60's build entry was missing and
  is backfilled below; FT59 is not yet built, no entry due. Same-day
  amendment: FT38 tabled by reviewer direction, revisit on or after
  2026-08-09.

- **FT60 promotion build (2026-07-08, backfilled 2026-07-09).** Gate change
  under `craft-gate` from the 2026-07-08 drain: a serialized build phase now
  owns the `dist/bench` write and runs before the parallel phases, so contract
  fixtures never exec a mid-rewrite binary; the sequential runner derives its
  serial-first ordering from `splitSerialPhases` (one source). Review pickups
  added three Coverage tests (SIGINT during the serial build, serial-phase
  reorder, half-present build-surface table).

- **Learnings run (2026-07-08, scope: post-FT43-49 drain).** Drained one parked
  idea and two open entries. Idea parked-pending-evidence: reclaim-lock
  protocol for lease Claim (FT58) — graduates on an observed double-win.
  Promoted to roadmap: sentinel precondition for fix-pass delegate charges
  (FT59, `craft-delegate` XS edit); gate phase race — conformance rebuilds
  `dist/bench` under the running contract phase, reproduced through
  `bench gate` itself (FT60, gate change under `craft-gate`). Previous run's
  promotion builds verified present in this record (the two 2026-07-07
  entries, backfilled by FT48).

- **Learnings run (2026-07-07, scope: verification probe).** Drained the open
  verification-probe entry and built its rule in the same reviewed pass under
  `craft-synthesis` (reviewer present and directing): `craft-delegate`'s
  done-claim verification list gains the independent-probe bullet — at least
  one accepted behavior probed outside the delegate's own tests, kept constant
  across a batch, because gate-green alone cannot tell a correct build from a
  self-consistent wrong one. Journal emptied; prose-only, gate green.

- **Learnings promotion (2026-07-07, FT35 + FT36: rule edits).** Two add-only
  rule edits from the drained 2026-07-06 learnings, built under
  `craft-synthesis`: worktree-isolated delegate charges open with the ff-only
  stale-base check, with the orchestrator fast-forwarding blocked worktrees
  (`craft-delegate`); and `/bench-write-spec`'s entry contract records the
  reviewer-directed batch-drain override without weakening the default
  decision-map gate. Prose-only, gate green.

- **Learnings promotion (2026-07-06, L1: shared-tree worktree rule).** From the
  2026-07-05 shared-tree-contention entry (drained to the roadmap in a prior
  reconcile). Folded into invariant 1 in `.bench/BENCH.md` — one sentence
  generalizing the existing delegate worktree-isolation clause to your own
  side-work: when `git status` shows another writer's in-flight edits, take
  side-work to a `bench worktree` or wait, so every gate verdict answers for
  exactly one diff. Fold, not a new piece; prose-only, gate green.

- **Learnings run (2026-07-04, scope: journal close-out).** Drained three
  entries. Promoted (already shipped): review-findings persistence — the
  `/bench-review-implementation` pickup artifact at `reviews/<spec-slug>.md`
  landed with conformance anchors and a canary bite proof, resolving the
  chat-only-findings entry. Recommended-and-parked: stale-gate benign/real
  status split — a `bench status` semantics change routed to the roadmap, to
  be shaped in one session with the parked `gate_tree_hash` capture-scratch
  carve-out. Dismissed: spec-without-decision-map deviation — one-off context;
  the entry contract and ask-first rule held, no rule change.

- **Learnings run (2026-07-04, scope: review/spec guidance).** Drained six open
  entries. Promoted: `craft-synthesis` requires a fresh-session dogfood run when
  skill or command triggers changed; `craft-delegate` records that read-only
  review findings are verified and fixed by the invoking session, not a separate
  worktree; `/bench-review-implementation` has an inline-axis fallback when a
  harness forbids unsolicited sub-agents; `/bench-write-spec` now checks external
  format/library divergence and runnable byte/wire compatibility probes; Research
  assets that claim byte or wire compatibility carry their own probe; and the
  Bench profile's hostile-input checklist now includes real CLI, linked by-path
  CLI, hook, and adapter invocation surfaces. Left open: stale-gate
  classification and durable review-findings storage, both still need product
  decisions. Dogfood loop pending the current gate diagnosis; verified with
  consistency greps and targeted tests in this run.

- **Learnings run (2026-07-04, scope: line governance).** Drained six open
  entries and one roadmap item. Promoted: the line declaration now uses an
  iteration cap as the stop condition; `craft-line` keeps venue routing
  right-sized instead of requiring every story to delegate, with inline work
  allowed for tiny slices and atomic diffs when deviations are reported; ordinary
  spec authoring stays mid-tier, with a top-tier exception only for Handoff
  uncertainty plus reviewer approval; delegate charges should use the Handoff
  digest and line-ranged excerpts over whole-file read lists; and `/bench-debug`
  allows direct fix-and-gate for small single-seam fixes. Dismissed: mandatory
  facilitator delegation as too rigid for delegation overhead and atomic diffs.
  Dogfood loop pending the current gate diagnosis; verified with consistency
  greps and targeted tests in this run.

- **Learnings run (2026-07-04, scope: learnings).** Drained two open entries
  (seam-level test batching; skipped conformance test counted as a pass).
  Promoted into `craft-tdd`: the red step now says to stub minimal declarations
  in compiled languages so the one test compiles before confirming a behavioral
  red — batching the file is never the fix — and the oracle section warns that
  a skip-only run still prints `ok`, so read test output, not the summary line.
  Pruned the batching entry's own proposed "same seam" clause as a duplicate of
  prose already in the skill. Dogfood loop per proportionality: prose-only —
  consistency grep of conformance pins plus a green `bench gate`.

- **Learnings run (2026-07-02, scope: learnings).** Drained one open entry.
  Promoted: cached routing for review-axis delegates in `projects/benchkit.md`
  Lines — ~60k tokens each on mid (three axes ≈ 180k), sourced from actuals
  running 2x a ~30k declaration. Project-specific calibration, not a kit rule.
  Dogfood loop waived: prose-only routing note; verified via consistency grep
  and a green gate.

- **Learnings run (2026-07-02, scope: learnings).** Drained one open entry.
  Promoted: quantified oracle promises must name their granularity — when a
  spec's behavior or red-signal promise ranges over a set ("each check"), the
  spec enumerates the set or states per-item vs per-class explicitly
  (`/bench-write-spec` step 4). Sourced from the state-surface build reading
  "each check bites" as per-class and review catching ~9 missing canaries.
  Dogfood loop waived: prose-only spec-authoring guidance; the next real spec
  is the dogfood.

- **Learnings run (2026-07-02, scope: learnings).** Drained two open entries.
  Promoted: a batch approval covers per-spec sign-offs when the reviewer is
  unreachable — build-and-flag, with specs left as post-hoc veto surface and a
  hard stop absent a batch approval (`.bench/BENCH.md` Workflow); venue routing
  in `craft-line` — delegate a spec'd build to `bench shift` on the cheap
  binding when every story's line is cheap and the gate fully observes the
  coverage map, with `/bench-write-spec`'s handoff clauses now pointing at that
  test instead of "mechanical enough". Dogfood loop waived: prose-only guidance
  a shift can't observe; verified via consistency audit and a green gate.

- **Learnings run (2026-07-02, scope: learnings).** Dismissed: unbound-model
  delegation at the reviewer's word — should-have-applied, `craft-line`'s
  resolve-first rule and the line hook already cover it; no kit change. Per
  reviewer direction, the journal convention changed from mark-resolved to
  prune-on-resolve: `.bench/learnings.md` now holds open entries only, verdicts
  live in the CHANGELOG and integration commits (updated the journal header, the
  `bench init` scaffold, `/bench-integrate-learnings`, and `craft-synthesis`).

- **Learnings run (2026-07-02, scope: learnings).** Drained eight open entries from
  `.bench/learnings.md` (seven tagged, one untagged straggler). Promoted: derived
  out-of-scope estimates (`<n> edits, <n> gate runs`) in `/bench-write-spec` step 3
  and its template; a Won't-handle line must leave at least one in-scope caller able
  to exercise the interface (`/bench-write-spec` step 5); delegate done-claims are
  verified against the gate and `git status`, with write-delegations in isolated
  worktrees (`.bench/BENCH.md` invariant 1); explicit staging instead of blind
  `git add -A`, with unexplained working-tree files blocking the commit
  (`/bench-implement-spec`). Dismissed: cached-routing miss and Codex skill-menu leak
  (fixes already shipped), seam-widening (one instance, not a pattern), invisible
  guards (already parked on the roadmap). Parked: the shift loop's own `add -A`
  staging as a roadmap idea.
- **Learnings run (2026-07-01, scope: learnings).** Drained seven open entries from
  `.bench/learnings.md`. Promoted: value-ranked out-of-scope capture in
  `/bench-write-spec`; a standing direct fix-and-gate shortcut for concrete review
  findings in `.bench/BENCH.md`; rename/refactor hygiene in `/bench-implement-spec`;
  and command-currency hardening in the gate for `.agents/**`, non-historical
  `decisions/`, Codex `$bench-*` adapters, and exact historical markers, with
  targeted canaries. Recorded as already promoted: Bench harness invocation guidance
  belongs in `.bench/BENCH.md`, not `AGENTS.md`.
- **Safe link dogfood slice (2026-06-28).** Made `bench link` copy by default,
  preserve project-owned `AGENTS.md` content through a managed Bench block, fail on
  same-name project-owned skills/commands/hooks, record installed assets in a link
  manifest, install portable `.agents/` content with Claude/Codex adapters, and gate
  the contract through throwaway-repo link checks plus npm package inspection. Added
  `.claude/README.md` so users can see that Claude paths are adapters to `.agents/`
  and shared `.bench/hooks/`.
- **Learnings run (2026-06-28, scope: learnings).** Drained two open entries from
  `.bench/learnings.md`. Promoted: a maintainer rule into HANDOFF — any `.bench/*`
  file the kit's prose references must be scaffolded by `bench init` and locked by a
  behavioral gate check (the `learnings.md` bug; executable fix already in `724bf8c`).
  Dismissed: "gate green claimed without a run" — already governed by invariant 1 and
  `/bench-final-check`; a pre-commit gate run would be a fourth check surface. Skipped:
  generalizing gate check 1d to every `.bench/*` file — speculative, only two exist.
- **Communication made first-class + portable; four learnings drained (2026-06-30,
  scope: learnings).** Sharpened AGENTS.md "How to talk to me" — clarity over density,
  tables/lists encouraged, the finding-vs-context tension scoped (cut the derivation,
  keep the one-clause why and enough context to resume cold), and "recommend, don't
  offer a blind menu." Shipped it to consumers via a new Communication section in
  `.bench/BENCH.md`. Promoted from `.bench/learnings.md`: recommend-at-every-question-
  and-hand-off; approval-table-before-build (`/bench-write-spec` exit); scan-for-unwritten-answers-
  before-closing-a-map (`craft-grill` + `/bench-shape-idea` exits). Dismissed: the persistent
  task-list progress tracker — `TaskCreate/TaskUpdate` are Claude-Code-only, so
  mandating them in harness-shared rules leaks one harness into the core.

## 0.2.0 — app-agnostic, npx-distributable, self-maintaining

- Made the kit app-agnostic: removed the hockey/coach/puck framing and all
  Regroup-specific content from core files (AGENTS.md, every skill and command, the
  CLI). Project-specific rules now live only in `projects/<name>.md`.
- Packaged for `npx` (`package.json`, dual `bench`/`benchkit` bins). `bench link`
  detects an ephemeral npx cache and copies instead of symlinking.
- Generalized the `design-system` skill and made the Claude-Design ↔ other-harness
  transition an explicit property: consumption reads committed artifacts, so the
  authoring tool is never a workflow dependency.
- Added `/setup` (interactive per-repo configuration, mirrors
  `setup-matt-pocock-skills`) and a Pocock-structure migration path.
- Added `/resynthesize` (this maintenance command) with three quality loops.

## 0.1.0 — initial synthesis

Baseline of what Bench incorporates from upstream. Adopted, with provenance in
`README.md`:

- **From Pocock:** `/start-ideation` (decision-mapping), `/spec` (to-prd), `/build` (implement),
  `/prep-shift` (review), `/fix-bug` (diagnosing-bugs), and the `seams` (codebase-design
  + design-it-twice), `tdd-at-seams` (tdd), `adr` (domain-modeling ADR format),
  `grill` (grilling), `writing-great-skills` skills, plus the `block-dangerous-git`
  hook (git-guardrails).
- **From kunchenguid:** the `axi` skill (AXI spec), the `.bench/gate.sh` + Stop-hook
  external-oracle pattern (no-mistakes), the `bench shift` gated loop with
  notes-between-iterations (gnhf), and `bench worktree` (treehouse).
- **Deliberately rejected:** firstmate's fleet orchestrator (overkill for solo work),
  the strict iron-law unattended TDD run (the self-grading failure it produces),
  lavish-axi and the tool-specific AXI binaries (build to the spec instead). These
  are closed unless a future `/resynthesize` finds a material upstream change.
- **Bench-native (neither repo):** the declared line (model + effort), stateless-reader
  docs, the design system as visual oracle, and the gate-as-oracle / never-self-grade
  invariant that binds it together.
