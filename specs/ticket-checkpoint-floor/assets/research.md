# FT349 research: the ticket check floor and the leftover `planned` citation

Consumed by: `specs/ticket-checkpoint-floor/spec.md` (the FT349 spec).
Drift: a change to the commit lane, the preflight rows, the plan model, the citation rules, or the `Writes:` grammar invalidates the cited lines.
Retire when: the FT349 spec retires.

## Recommendation, scope, and evidence status

Recommendation: put both refusals in the build preflight and in the chunk checkpoint, which already read the completion plan and the ticket `Covers:` lines. Do not put them in the commit lane.
The floor refusal compares the planned verification commands of each ticket with a floor that the preflight derives from the ticket's `Writes:` line.
The `planned` refusal compares the seam cells of the rows that a committed ticket covers with the `planned` marker.

Scope: the commit lane, the build preflight, the completion plan, the chunk checkpoint, the citation rules, and the staged work on the same paths.
Subject tree: commit `8bb4b894`, read only.

Evidence status: ten paired readers (Sonnet and Haiku at high effort, two for each of five slices) returned.
The coordinator re-opened the load-bearing sources listed in the verification record.
The Haiku readers in R1, R2, and R3 could not run `rg`, so their negative claims cover only the files that they read. The Sonnet reader of each pair ran a full search, and the agreed set below rests on both returns.

## Facts

### Q1. What the ticket commit runs

- `bench commit` runs the lane, publishes, and stops. With `--preflight-build <slug>`, it then builds the worktree and runs `bench preflight build <slug> --source-tip <published>` (`cmd/bench/commit_chain.go:32-55`). The preflight runs after publication, so its red cannot stop the commit.
- The commit does not know its ticket. The slug is the only planning input (`internal/commit/commit.go:299-301`).
- This repo has no `.bench/phases.json`, so the commit uses the built-in kit lane, which is selective by path class (`internal/gate/manifest.go:171-188`, `internal/gate/kit_source.go:86-88`).
- The kit lane runs gofmt, prose, vet, build, structure, and document-family checks (`internal/gate/lane.go:75-84`). Its comment says "The whole-project test phase is deliberately absent" (`internal/gate/lane.go:71-74`). No path class selects a Go test (`internal/gate/lane_select.go:106-156`).
- ADR 0017 binds the lane to a declared check list that composed changes select by path class, and says the lane "takes seconds" (`docs/adr/0017-the-worktree-commit-runs-the-fast-lane.md:7-15`).
- The only commit-time reader of `Writes:` is the light-path authorization for an unbound assignment (`internal/commitment/repository/light_path.go:27-99`, `candidate.go:99-120`). It grades paths, not checks.

### Q2. What the build preflight runs, and how it reads `Writes:`

- Every `bench preflight build` form runs one pure `Decide` (`internal/preflight/decision.go:179-224`). `Decide` does no I/O and runs no test.
- The guidance runs the default form after each ticket commit (`.agents/commands/bench-implement-spec.md:48`). `--plan-only` runs the same rows without the commitment row (`internal/preflight/command.go:95`, `:157-161`).
- The rows include `completion-plan` (`internal/preflight/plan.go:36-42`), `writes-resolve`, three closure rows, `fence-writes`, `kit-pin`, `rows-owned`, and `rows-membership`.
- `Writes:` splits on `, ` (`internal/maps/fields.go:107-116`). Each entry is an exact path or a directory prefix, with an optional `(new)` suffix and no glob (`internal/tickets/writes.go:13-25`).
- No function maps a `Writes:` entry to a Go package. The nearest seam is `bench test --changed`, which maps changed regular files to packages through `go list -json -test ./...`, embed owners, and the reverse-import closure (`internal/testreport/selection.go:60-263`). It refuses a directory entry (`selection.go:105`) and selects nothing for a non-Go path outside the embed set (`selection.go:201-206`).
- `packagesurface.EmbedTargets` lists every embed target (`internal/packagesurface/assets.go:87-120`). No `//go:embed` directive lives under `cmd/bench`.
- No predicate decides "a public response changed". The public command table is `cmd/bench/main.go:70-150`. The bound registry binds `cmd/bench` paths to five registry files (`internal/tickets/registry_data.go:21-44`).
- Root conformance is `TestRootConformance` with `BENCH_CONFORMANCE_ROOT` (`internal/conformance/gate_entry_test.go:16-31`). Only the gate runs it whole (`internal/gate/phases.go:101-149`). `bench test --check <name>` runs one scope of it (`internal/testreport/named_check.go:91-99`). Delegation discipline says `bench test --package ./internal/conformance` is not the root conformance pass (`.agents/skills/bench-craft-delegate/references/delegation-discipline.md:125-128`).

### Q3. The `planned` marker and row-to-ticket traceability

- `plannedRe` is the whole word `planned` (`internal/coverage/citations.go:33-35`). Its only use exempts a cell from the stray-test-name rule (`citations.go:147-150`).
- A `planned` cell that cites a backticked path and a backticked name list is still graded. A missing file or an undeclared name is red (`citations.go:198-244`). Real planned cells avoid this with an unbackticked list, for example `specs/process-lifetime/spec.md:861`.
- `internal/coverage` reads no spec status. Build preflight refuses a non-staged spec (`internal/preflight/gather.go:48-53`). The gate's coverage sweep runs `coverage --check` on every spec (`internal/conformance/docs_workflow_checks_test.go:367-394`).
- The author rule is `.agents/skills/bench-craft-spec/references/map-discipline.md:122-125`. No guidance tells an author to rewrite a `planned` cell after the test lands.
- A ticket's `Covers:` tokens match `RowIDPattern` `([A-Z]+)[0-9]+` (`internal/tickets/tickets.go:55-65`). Readers are `rows-owned` and `rows-membership` (`internal/preflight/decision.go:367-408`) and the plan chunk rows (`internal/reviewrecord/plan.go:105-110`).
- `coverage.ParseSpec` exports the opt-in flag, the row IDs, and the violations. It does not export the seam cell (`internal/coverage/coverage.go:519-526`). A preflight refusal of a `planned` cell needs a new projection.
- No check at any stage refuses a leftover `planned` cell.
- 21 of 23 staged specs use `planned` cells. The roadmap row's count of four is stale (`roadmap/FT349.md:22`).
- `citation_form_test.go:73-97` pins that a planned bare-name cell passes and stays uncited.

### Q4. Evidence of executed checks

- `bench test` persists no run record (`internal/testreport/outcome.go:29-50`). The census skips Bench calls and records no operand (`internal/census/census.go:43-46`, `:171-175`).
- The gate's reusable green is keyed by tree and oracle and holds no package list (`internal/gate/engine.go:103-151`). The lane record holds a tree and an outcome only (`internal/gate/lane.go:319-331`).
- The spec's `bench-completion-plan` block holds the planned verification commands. Each requirement is `{id, command, probe?, ticket?}` (`internal/reviewrecord/plan.go:15-33`). A version 2 plan gives every ticket at least one requirement (`plan.go:124-150`).
- `bench record verification` copies the planned command and takes a caller-stated exit code (`internal/reviewrecord/write.go:246`, `recordcmd/command.go:363`). The chunk checkpoint refuses a missing, stale, failed, or wrong-command record (`internal/reviewrecord/check.go:27-56`, called from `internal/gate/checkpoint.go:216-221`).
- Nothing compares the planned commands with the ticket's `Writes:` line.
- Plans spell root conformance in several ways. Three specs plan `bench test --check conformance-meta` (`specs/native-record-operations/spec.md:412-414` is one).
- Some tickets carry an unparsed `## Checkpoint verification` heading with `bench test --package` bullets (`specs/shared-test-fixtures/tickets/01-share-git-fixture-owner.md:44-49`). The ticket parser does not read it (`internal/tickets/tickets.go:42-53`).

### Q5. Overlap with staged work and closed decisions

- The approved order puts FT293, FT375, and FT369 before FT349 in the priority list (`ROADMAP.md:21-27`). The commitment names FT349 as the eligible outcome (`bench commitment show`). The recommended sequence names only FT349 (`ROADMAP.md:319-321`).
- `.bench/commitment.json:288-296` fixes two criteria, `FT349.floor` and `FT349.planned`.
  - The floor criterion names three omissions: root conformance, a package that the `Writes:` line derives, and `cmd/bench` for a public response or embed change.
  - The planned criterion refuses a planned citation for a row that the committed ticket covers.
- FT293 (`specs/preflight-ownership-closure`) fences most of `internal/preflight`, `internal/coverage`, `internal/tickets`, `map-discipline.md`, and `slicing-checks.md` (`spec.md:739-840`). It moves the coverage parser into `rows.go` and `projection.go` (`spec.md:445`).
- FT375 (`specs/preflight-spec-staleness`) adds `--audit-spec` with a cited-tests row. It keeps the ordinary after-ticket preflight unchanged, and a planned test stays valid after a predecessor creates it (`spec.md:161`, `:208`, `:264`).
- FT293 and FT375 pin source censuses that an FT349 edit to the same packages makes stale (`specs/preflight-ownership-closure/spec.md:938-945`).
- `roadmap/FT380.md:6-11` keeps lane enforcement closed until a `bench commit` repro exists. `roadmap/FT215.md` leaves open whether the lane gains checks that it did not run.
- ADR 0028 is the precedent for a commit-time check against `Writes:` (`docs/adr/0028-the-commitment-gates-spec-implementations.md:9-19`).
- An open learning holds the FT349 evidence (`capture/learnings.md:207-210`). Six plan commits from one build each added a package check that the planned chunk commands omitted. A second open learning reports 55 leftover `planned` cells at FT393 final reconciliation, each with an existing test (`capture/learnings.md:192-195`).

## Inferences

- The commit lane is the wrong owner. It runs no Go test, it takes no ticket input, ADR 0017 bounds it to seconds, and FT380 and FT215 keep lane test selection open.
- The completion plan is the existing statement of "the checks the author runs". The held evidence describes plan commands that omitted packages. A plan check therefore catches the observed failure at `--plan-only` time, before any author starts, and again after each ticket commit.
- The chunk checkpoint already refuses a planned requirement with no passing record. With a complete plan, the checkpoint then proves that each floor command ran with exit 0, as far as a self-stated exit code proves it.
- The `planned` refusal needs the set of committed tickets. The chunk checkpoint knows its chunk's tickets and rows (`plan.go:105-110`). The after-ticket preflight does not know which tickets are committed.
- FT375's after-ticket rule stays intact if FT349's refusal grades only the rows of tickets whose chunk the checkpoint closes.

## Proposals for the spec author

These are options, not decisions. The spec author selects and the reviewer confirms.

| Option | Floor owner | `planned` owner | Consequence |
|---|---|---|---|
| A | A build preflight row that compares each ticket's planned commands with its derived floor | A chunk checkpoint refusal for the chunk's covered rows | Composes the plan and checkpoint seams. Catches the observed failure before dispatch. Exit codes stay self-stated. |
| B | The commit lane runs the derived floor | The commit lane refuses the cell | Conflicts with ADR 0017's seconds bound and with the open FT215 and FT380 decisions. The commit does not know its ticket. |
| C | The after-ticket preflight runs `go test` for the floor | The after-ticket preflight refuses the cell | Breaks `Decide`'s no-I/O contract. The preflight runs after publication, so it cannot stop the commit. |

## Contradictions

- The roadmap row says four staged specs use `planned`. The tree has 21.
- `roadmap/FT349.md:7-8` says a planned cell "may name a test that does not exist yet". The code exempts only a bare name outside the citation form.
- `.agents/commands/bench-implement-spec.md:44` says the ticket commits "on a green lane". A lane pass is never called green (`internal/gate/lane.go:32-37`).
- The guidance names "root conformance" with no command. Plans spell it in several ways.
- The guidance names "a changed public response or embed pattern" with no path predicate. No embed directive lives under `cmd/bench`.
- FT375 keeps the after-ticket preflight unchanged for planned tests. FT349's refusal must not grade a row whose ticket has not committed.
- FT375's cited-tests row reports "each planned test that exists" (`roadmap/FT375.md:14-15`). If FT349 and FT375 both derive that fact, the repo holds two derivations. The later spec must compose the earlier projection.
- `roadmap/FT89.md:9-10` keeps open whether the ticket check floor replaces the per-return gate run. FT349 enforces the floor but must not decide that question.
- `internal/coverage/citation_form_test.go:81` pins that a planned cell passes. A change to that expectation needs a demonstrated red under ADR 0006 (`docs/adr/0006-independent-test-expectations-preserve-omission-oracles.md:3`).

## Unknowns

- The exact six held plan commits. They sit in the history of `362eb494^2`. Candidates include `3ec313b0`, `1a9d1c4d`, `6db48cb0`, `5a341734`, `47b779cd`, `9ed6c6f8`, and `d7ad195d`. The spec author confirms them.
- Whether `bench test --changed` reaches `cmd/bench` through the reverse-import closure for a given embed owner. Nobody ran it.
- Whether the floor for a non-Go `Writes:` path (Markdown, fixture, JSON) is empty or maps to a named check.

## Validation plan

- Run `bench preflight build <slug> --plan-only` against a fixture plan that omits one package of a ticket's `Writes:` set, and expect a red floor row.
- Run the chunk checkpoint against a fixture whose covered row keeps a `planned` cell after its test exists, and expect a refusal.
- Replay the six held plan commits' parent plans through the new floor row, and expect each to red.
- Run `bench coverage --check` on every staged spec after the change, and confirm that no spec turns red before its own build.

## Verification record

- [x] The output opens with the recommendation, the scope, and the evidence status.
- [x] Facts, inferences, and proposals are in separate sections. No tested result exists.
- [x] One section holds each question.
- [x] An option table records each consequence.
- [x] Every contradiction and residual unknown is kept.
- [ ] A diagram: not needed. The seam relation is linear and the table states it.
- [x] Unverified claims are listed.
- [x] The output ends its body with a validation plan.
- [x] The coordinator re-opened `.bench/commitment.json:288-296`, `capture/learnings.md:192-210`, `internal/reviewrecord/plan.go:15-150`, and `specs/native-record-operations/spec.md:380-425`.
