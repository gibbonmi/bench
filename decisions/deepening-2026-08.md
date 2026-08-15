# Deepening survey 2026-08 — eight refactor candidates

Status: ready

## Destination

Turn the 2026-08-15 deepening survey (report `/tmp/architecture-review-20260815T101417.html`,
scoped by `ASSESSMENT.md` L-O1/M-C1/M-W1 and roadmap FT186/FT162/FT185) into
reviewer-decided, behavior-preserving deepenings, each specced and built through the
normal workflow with an exit test that keeps existing test logic unmodified. Candidates
verified against source: (1) the gate run as one transaction, (3) verdict record classes
as one registry, (2) one worktree eligibility verdict, (4) one adopt lifecycle decision,
(5) a shift-objective owner with per-surface projections, (6) one green-marker reader,
(7) named git readers over `git.Output`, (8) one shipped skill-frontmatter reader.

## #1: Sequence against the active FT189 build

Blocked by: none
Type: Grill

### Question

FT189 (`specs/worktree-enumeration-hang`) is building in an assignment worktree; its
first ticket `resolve-git-common-dir` is the named reader candidate 7 proposes and
candidate 2 edits the same `internal/worktree` files. Shape all eight now and defer the
builds of 2 and 7 until FT189 lands, re-scoping 7 against what FT189 leaves?

### Answer

Shape all eight now. Builds of candidates 2 and 7 start only after FT189 lands on `main`; 7 is re-scoped then against the reader FT189's `resolve-git-common-dir` ticket leaves. The other six may build before FT189 lands.

## #2: ADR exposure — revise, retire, or keep 0005/0006/0012

Blocked by: none
Type: Grill

### Question

On reading, no ADR conflicts: 0005 constrains candidate 2's shape (ownership evidence
stays conjunctive, DiscardBranch derived-after), 0006 constrains the guard tests in
3 and 8 (expectation independently authored), 0012 is the registry shape 8 follows.
Keep all three unchanged and record them as constraints, or revise/retire any?

### Answer

Keep 0006 and 0012 unchanged; both are recorded as constraints (guard-test expectations in 3 and 8 stay independently authored; 8 follows the declarative-registry shape). Revise 0005 to name the single eligibility verdict as the decided state — the revision lands with candidate 2's build, after ticket #8 fixes the verdict's shape.

## #3: One map or one map per candidate

Blocked by: none
Type: Grill

### Question

Hold all eight in this map (one ticket per candidate plus shared tickets), or split
into per-module maps now?

### Answer

One map. Shared decisions settle once here; one shape ticket per candidate; specs are cut at write-spec time per ticket #5.

## #4: Refactor lane — wait for FT108 or spec each with an exit-test rule

Blocked by: none
Type: Grill

### Question

FT108's `craft-refactor` skill does not exist. Build these under the normal spec path
with the exit test written into each spec (suite passes with test logic unmodified;
mechanical renames only; a changed assertion reverts the move), or build FT108 first?

### Answer

Normal spec path. Every spec from this map carries the exit-test rule as an acceptance row: the pre-existing suite passes with test logic unmodified, mechanical renames being the only permitted test edit; a changed assertion reverts the move and reroutes to the feature path. A defect found mid-refactor is parked and fixed as its own commit, never bundled. FT108 is not a prerequisite.

## #5: Grouping and order of the builds

Blocked by: #3
Type: Grill

### Question

Group specs by module — gate (1, 3, 6), shift (5), adopt (4), skills (8),
worktree/git (2, 7 after FT189) — in that order, or another grouping/order?

### Answer

Five specs, cut and built in this order: gate (candidates 1, 3, 6) → shift (5) → skills (8) → adopt (4) → worktree/git (2, 7; blocked on FT189 landing per #1). Each spec is one `/bench-write-spec` from this map.

## #6: Candidate 1 — the gate run as one transaction

Blocked by: #4
Type: Grill

### Question

Delete the `gateEngine` seam (12 methods, one implementation, no fake) and the four
zero-caller wrappers (`executeWithEngine`, `executeWithEngineAfterAcquire`,
`executeSubjectWithEngine`, `newWorkingTreeEvaluation`) as part of the deepening, or
keep the seam for a future fake? Does candidate 3 ride in the same spec?

### Answer

Delete the `gateEngine` interface, `productionGateEngine`, and the four zero-caller wrappers; the deepened module owns lock acquisition, owner file, under-lock drift check, the pending→terminal record pair, retain/invalidate pairing, and the crash-safe replace. Callers keep the public `Execute`/`Inspect` seam unchanged; every existing consumer test (`internal/commit`, `internal/shift`, `internal/stophook`, conformance) passes with test logic unmodified. Candidate 3 rides in the same gate spec.

## #7: Candidate 3 — verdict record classes as one registry

Blocked by: #4
Type: Grill

### Question

Registry rows carry name, exact field set, validator; the guard test the
`readyFieldClasses` comment promises lands as an independently authored expectation
(ADR 0006). Are the uncalled fixtures `partialTestRecord`/`fullTestRecord` deleted?

### Answer

One registry: each record class is one row carrying its name, its exact field set, and its validator, read by the loader (`verdict.go:671-681` switch collapses into it), by whatever reports class names, and by the reuse refusal; the Pending field set joins as a row. An unregistered class cannot compile. The guard test the `readyFieldClasses` comment promises lands as an independently authored expectation (ADR 0006), and the comment's phantom references (`storeRecordClasses`, `record_classes.go`, `TestVerdictReadyFieldsAreAllRegistered`) are replaced by the names that exist. `partialTestRecord`/`fullTestRecord` are deleted.

## #8: Candidate 2 — one worktree eligibility verdict

Blocked by: #1, #4
Type: Grill

### Question

`PlanExplicitWithOptions` decides Action/Reason by statement order (24 assignments);
`PlanAutomatic` re-derives via `HasPrefix(plan.landed, ...)`. Does the deepening
preserve the current effective precedence exactly (characterized first, anomalies parked
as separate fixes), or may it correct precedence anomalies it finds?

### Answer

Preserve the current effective precedence exactly: the spec first characterizes the order `PlanExplicitWithOptions` and `PlanAutomatic` produce today (one row per reachable Action/ReasonCode outcome), then the deepened eligibility module returns one decided verdict plus its evidence with refusal precedence declared as ordered data, and `PlanAutomatic` becomes a stricter reading of that verdict — no `HasPrefix` over `plan.landed`. `ApplyExplicit`, `PlanAutomatic`, `LandCommand`, `ReleaseCommand` all read the one verdict. Any precedence anomaly the characterization exposes is parked (`bench idea`) and fixed through `/bench-debug`, never in this spec. ADR 0005's conjunctive ownership (marker + assignment + lock + landedness + recovery agreement) and derived-after `DiscardBranch` hold; ADR 0005 is revised in this build to name the eligibility verdict as the decided state.

## #9: Candidate 4 — one adopt lifecycle decision

Blocked by: #4
Type: Grill

### Question

`PlanLifecycle` grades only `upgrade`; `link`/`setup` (`transactionalLink`) and
`unlink` (`planUnlink`) hand-roll add/change/remove. Do all four verbs execute the one
decision, and is a `--dry-run`/plan projection out of scope (new behavior)?

### Answer

All four verbs — `link`, `setup`, `upgrade`, `unlink` — execute `PlanLifecycle`'s add/change/remove decision, and the transaction shrinks to applying that decided operation list atomically; `transactionalLink` and `planUnlink` stop deriving their own. A `--dry-run`/plan projection is out of scope (new observable behavior).

## #10: Candidate 5 — shift objective owner

Blocked by: #4
Type: Grill

### Question

Six surfaces, six treatments; the durable commit subject (`loop.go:288`) is raw. What
policy does the durable subject follow — same sanitizer as the banner, a bounded
projection, or iteration-only with no objective text? This closes ASSESSMENT C-08.

### Answer

One objective module hands out per-surface projections — banner line, prompt body, scratch bytes, done.sh predicate argument, durable commit subject — and no surface obtains the unprojected text. The durable subject becomes `shift: iteration <i> — <sanitize.Preview(objective)>`: escaped, 120-rune bounded with the `… (N bytes)` suffix, identical policy to the banner. This is the one behavior change in the shift spec and carries its own acceptance row; it closes ASSESSMENT C-08 / backlog rank 6.

## #11: Candidate 6 — one green-marker reader

Blocked by: #4
Type: Grill

### Question

`engine.go:133` reads the marker without peeling, `land.go:365` peels `^{commit}`,
`authorization.go` classifies a dangling symbolic ref. One reader — peel plus dangling
classification everywhere (the authorization reading becomes the reader)?

### Answer

One marker module owns read, `^{commit}` peel, dangling-symbolic-ref classification, and compare-and-advance for `refs/bench/green/<branch>`; `engine.go:133`, `land.go:365`, and `authorization.go:60-148` all call it. Callers learn one question: does the marker authorize this tip. The only behavior delta — a non-commit marker object now peels instead of mismatching — is named in an acceptance row.

## #12: Candidate 7 — named git readers over `git.Output`

Blocked by: #1, #4
Type: Grill

### Question

Promote only incantations that recur at three or more production sites into named
readers, keeping `Output`/`Raw`/`OK` as plumbing; re-scope after FT189's reader lands?

### Answer

After FT189 lands: promote only git incantations recurring at three or more production sites (common-dir/admin-dir, peeled commit, porcelain status, …) into named `internal/git` readers; `Output`/`Raw`/`OK` remain as plumbing. Re-scope against the reader FT189's `resolve-git-common-dir` ticket lands; if it already owns the common-dir reader, this candidate consumes it rather than adding a second.

## #13: Candidate 8 — one shipped skill-frontmatter reader

Blocked by: #4
Type: Grill

### Question

This is FT89's single-sourcing slice. One shipped Go module owns frontmatter and
kit-only marking; `.bench/skills-index.sh` and conformance both invoke it; the
checked-in index is the independently authored expectation. Take that slice here and
leave the FT89 row in place?

### Answer

Take FT89's single-sourcing slice here: one shipped Go module owns skill frontmatter (`index:`/`index-note:`) and kit-only marking from `consumer-payload.json`; `.bench/skills-index.sh` and the conformance check both invoke it, so neither carries its own parser. The checked-in index is the independently authored expectation (ADR 0006), following the declarative-registry shape ADR 0012 blesses. The FT89 roadmap row stays in place. Kit edit under `craft-synthesis`.

## Not yet specified

## Spec-writer discretion

- Names and file placement of the deepened modules and registry rows, provided the
  public seams named in the tickets do not move.
- Whether a registry is a slice or a map, and how compile-time registration is enforced.
- Ticket slicing within each spec (expand–contract order), per `craft-tickets`.

## Out of scope

- A `--dry-run`/plan projection for adopt verbs (new behavior; park if wanted).
- Fixing worktree precedence anomalies inside the deepening (route through `/bench-debug`).
- Deleting `git.Output`/`Raw`/`OK`.
- FT108 `craft-refactor` skill as a prerequisite.
- Revising or retiring ADR 0006 or 0012.

- Modules that passed the deletion test in the survey: `cmd/bench/main.go` command
  registry, `internal/status.GateVerdict`, `internal/gate` `durableReplaceRecordAt` and
  `runner.schedule`, `internal/worktree/resume.go` cleanup receipts.

## Sources
