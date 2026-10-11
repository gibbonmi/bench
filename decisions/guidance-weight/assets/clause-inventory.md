# Clause inventory of the always-loaded guidance

Consumed by: the guidance-weight map (#9, #10, #11, #14, #15) and the FT100 always-loaded spec.
Drift: an edit of `.bench/BENCH.md`, `AGENTS.md`, or an anchor row that names either file.
Retire when: the FT100 always-loaded spec lands.

## Recommendation, scope, and evidence status

Recommendation: move B11 and B26–B33 to one build reference, merge the six light-path statements, and remove four repeated copies. Keep every other clause always-loaded until an FT231 result.

Scope: 60 clauses at main `e7ba5386`. B1–B43 are in `.bench/BENCH.md`, and A1–A17 are in `AGENTS.md`. Q1 asks where each clause acts and where else it lives. Q2 asks which failure or direction added it. Q3 asks for the headroom of each build destination.

Evidence status: five reader pairs returned on 2026-10-11, with one Sonnet/high and one Haiku/xhigh reader for each slice. The agreed set rests on both returns. The coordinator closed each disagreement with one targeted read. Byte sizes and anchor counts come from deterministic commands, not from the readers.

## Facts

### Q1 and Q2 — each clause

Bytes are the clause's lines at `e7ba5386`, counted with `awk` over the line ranges. Anchors count the require rows of `bench anchors` whose line falls in the clause. "M" marks a shared-rule marker of `internal/conformance/validity_checks_test.go:251-264`. Origin names the commit that added the clause and the record class that it cites.

| Clause | Bytes | Anchors | Acts in | Other homes | Origin |
|---|---|---|---|---|---|
| B1 | 352 | 0 | every session; Codex start and resume | `.bench/BENCH-reference.md:153-189` | 724c7a2f; failure at `specs/cli-desktop-consistency/spec.md:13` |
| B2 | 299 | 1, M | every session | `.agents/skills/bench-craft-synthesis/SKILL.md:20-27` | c21af5ea; FT187 failure at `specs/ft187-communication-surface-cut/spec.md:9` |
| B3 | 288 | 1 | every session | `.bench/BENCH-reference.md:210-211`, `:377-383`; `projects/benchkit.md:29` | 2c3743b5; d1b2b9d2, failure at `roadmap/FT89.md:44` |
| B4 | 215 | 2 | every session | `docs/adr/0025-one-owner-bounds-every-public-bench-response.md:3`; `CONTEXT.md:190-196` | 196d3353 and f759d89e; no failure at origin |
| B5 | 297 | 0 | every session | `CONTEXT.md:190-196` | 4b962e9a, failure at `roadmap/FT254.md:74`; 7b2a2858 |
| B6 | 491 | M | every session | `.agents/skills/bench-craft-delegate/SKILL.md:91`, `:114` | c21af5ea; 617d6fd9, a fabricated done-claim learning |
| B7 | 480 | M | every multi-cycle stage | `.agents/skills/bench-craft-line/SKILL.md:9`, `:13-15` | c21af5ea; f89aee09, direction |
| B8 | 252 | M | doc and ADR writing | `.agents/skills/bench-craft-adr/SKILL.md:7-13`, `:43` | c21af5ea; no failure |
| B9 | 482 | 2, M | every session | `docs/adr/0017-the-worktree-commit-runs-the-fast-lane.md`; `CONTEXT.md:8-16` | c21af5ea; 45581e89; bed3aa17, failure PG19 |
| B10 | 240 | 1 | every session; the exception only in the staleness pass | `.agents/skills/bench-implement-spec/references/staleness-pass.md:44-46` | bed3aa17, failure PG20; ec30797e, direction |
| B11 | 155 | 1 | spec build | `docs/adr/0023-each-ticket-gets-a-fresh-author.md:17` | bed3aa17; no failure |
| B12 | 122 | 1 | fix loop | `.agents/skills/bench-craft-line/SKILL.md:76-95` | bed3aa17, failure PG18 |
| B13 | 322 | 0 | every session | `.agents/skills/bench-craft-spec/references/ste-prose.md:1-8` | fedbf484; a comparison, no failure |
| B14 | 170 | 0 | every session | none | 3e323a6d; no failure |
| B15 | 211 | 0 | every session | `.agents/skills/bench-craft-grill/SKILL.md:22-24`, partial | 2567e03c, a learning about options with no recommendation |
| B16 | 218 | 0 | every session | `.bench/BENCH-reference.md:134-135`; `docs/adr/0013-the-route-is-a-projection-of-the-board.md:13` | 8ff078e3, direction |
| B17 | 160 | 1, M | every session | none | 3eb1c9a5, FT152 story 15 |
| B18 | 94 | 0 | every session | none | 120b4bc8; no failure |
| B19 | 733 | contract check | every phase | none | 9f4d0c45, failure at `specs/structured-phase-conversation.md:5-12` |
| B20 | 79 | 0 | every session | none | 5fd37891, direction |
| B21 | 295 | 1, M | every session | none | 120b4bc8 and 5fd37891, direction |
| B22 | 416 | M | every session | `.agents/commands/bench-drain.md:197-225`; `.agents/commands/bench-final-check.md:76` | cea2f42b; no failure |
| B23 | 235 | 0 | every session | `.bench/BENCH-reference.md:42-44`; `CONTEXT.md:245-249` | d263bec7, direction |
| B24 | 528 | 1 | every session | `.agents/commands/bench-drain.md:13`, `:214-220` | e0f9d5e7, aa55b83f, and 583f71e8, direction |
| B25 | 486 | 0 | every session, as orientation | `README.md:134-140` | d8cbd48e; no failure |
| B26 | 535 | 4 | spec build | `docs/adr/0023-each-ticket-gets-a-fresh-author.md:6-8` | 41f19bce, failure at `specs/fresh-ticket-authors/spec.md:11` |
| B27 | 313 | 2 | spec build | `CONTEXT.md:251` | 53fc8a34, direction |
| B28 | 343 | 3 | spec build | `docs/adr/0023-each-ticket-gets-a-fresh-author.md:15` | 53fc8a34, direction |
| B29 | 492 | 1 | spec build | `docs/adr/0023-each-ticket-gets-a-fresh-author.md:16-17` | 41f19bce; failure at `roadmap/FT293.md:104` |
| B30 | 772 | 4 | spec build | `docs/adr/0021-benchmark-workflow-orchestration.md:10` | 1212e8ec; direction in decision ticket 5 of 53fc8a34 |
| B31 | 177 | 0 | spec build | `.agents/commands/bench-implement-spec.md:21` | 1e0fb140 and d2a2e1ae; no failure |
| B32 | 390 | 2 | delegated run | `docs/adr/0023-each-ticket-gets-a-fresh-author.md:13` | fd2b5de9, direction at `specs/delegated-implementation/spec.md:5` |
| B33 | 431 | 4 | delegated run | `.agents/skills/bench-craft-delegate/references/delegation-discipline.md:11-13` | fd2b5de9, direction |
| B34 | 268 | 1, M | every session | `.agents/commands/bench-implement-spec.md:8` | cea2f42b; 75749bd9, a learning |
| B35 | 548 | 3 | every session | `docs/adr/0028-the-commitment-gates-spec-implementations.md:9-11` | 7c75cc8f, failure in FT107 |
| B36 | 761 | 6, M | every session | `docs/adr/0026-work-class-sets-the-default-delivery-order.md`; `.agents/commands/bench-drain.md:214` | 0f7fd876, direction |
| B37 | 355 | 1 | every phase close | `.agents/commands/bench-drain.md:85`, `:313` | 53964fcd, failure and direction |
| B38 | 350 | 1 | every phase close | none | 53964fcd, direction |
| B39 | 330 | 3 | every session | `docs/adr/0028-the-commitment-gates-spec-implementations.md:7-12` | 583f71e8, failure at `specs/light-path-commitment-exemption/spec.md:15` |
| B40 | 702 | 5 | every session outside a build; every drain | `.agents/commands/bench-drain.md:217-223` | 79512031 and 583f71e8, direction |
| B41 | 492 | 2, M | every phase | `docs/adr/0014-main-receives-writes-only-through-landings.md:3` | 9ffb8278, failure at `specs/parallel-landings/spec.md:18-21` |
| B42 | 410 | 1, M | every session | `.agents/commands/bench-implement-spec.md:56` | 3eb1c9a5, failure at `specs/implement-spec-full-run.md:63-65` |
| B43 | 275 | 0 | spec sign-off | `.agents/skills/bench-craft-tickets/SKILL.md:37` | 28615134, a learning about a build past sign-off |
| A1 | 304 | 1 | every session | `CLAUDE.md:3-6`; `projects/benchkit.md:165-171` | 7b6cf5f1, failure in `specs/shared-rules.md` |
| A2 | 329 | 0 | kit edit | `.bench/BENCH-reference.md:98` | 1e85fe24, failure at `specs/shared-sections.md:9-12` |
| A3 | 436 | 0 | every session | `.agents/skills/bench-craft-review/SKILL.md:36-40` | 978f68e6, direction |
| A4 | 472 | 0 | review and author charges | `docs/adr/0006-independent-test-expectations-preserve-omission-oracles.md:3` | e1a5ebea, failure in the FT78 review |
| A5 | 150 | 0 | an edit of a Markdown check | none | 9968f271, failure in FT323 |
| A6 | 382 | 0 | a dependency edit | none | 1ca682ea, direction |
| A7 | 709 | 2 | every phase close | `.bench/BENCH-reference.md:360-363`; `.agents/commands/bench-drain.md:338` | 838b0f77, direction; d0263c23, failure in FT109 |
| A8 | 255 | 0 | every session | none | 343e2440; no failure |
| A9 | 206 | 0 | every session | none | 0f574dd6, a two-hour `rg` hang |
| A10 | 265 | 0 | every session | `.bench/BENCH-reference.md:42-43` | 00405e39; no failure |
| A11 | 258 | 0 | every session | none | 343e2440; no failure |
| A12 | 165 | 1 | every session | none | bed3aa17, failure in FT107 |
| A13 | 170 | 0 | every session | `.agents/skills/bench-craft-delegate/SKILL.md:104-105`, partial | 5f1b2453; no failure |
| A14 | 117 | 1 | every session | none | bed3aa17, direction in FT107 |
| A15 | 102 | 1 | every session | none | bed3aa17, failure in FT107 |
| A16 | 262 | 1 | every session | `.bench/BENCH.md:15-17`, partial | bed3aa17, failure in FT107 |
| A17 | 128 | 1 | a hand run of the system suite | `projects/benchkit.md:676-677` | ee3fc042; no failure |

### Q3 — build destinations

Each phase file and skill below has an exact budget row equal to its line count, so it has no headroom. `.agents/commands/bench-implement-spec.md` holds 93 lines against 93, and `bench-write-spec.md` holds 87 against 87. `bench-debug.md`, craft-line, and craft-delegate are also full. `bench-review-implementation.md`, `bench-final-check.md`, and every reference file have no row. The precedent for a build reference is `.agents/skills/bench-implement-spec/references/staleness-pass.md`, which `.agents/commands/bench-implement-spec.md:21` charges.

Seventeen lines point to the build rules in `.bench/BENCH.md`. They are in `bench-implement-spec.md` (6), `bench-review-implementation.md` (4), `bench-final-check.md:194`, `bench-write-spec.md:50`, `bench-debug.md:157`, craft-spec (2), and craft-delegate (2). The forbid rows at `internal/anchors/registry_retained_workflow.go:295-301` stop the review and build phases from restating these rules.

## Inferences

Under #2, a clause moves when it acts only in some phases. That set is B11 and B26–B33, which act only in a spec build or a delegated run, and B43, which acts only at spec sign-off. B10's exception and B12 also act only in some phases, but an on-demand owner already holds each one, so each is a repeated copy.

Three clauses act in every phase, so they stay. B19 governs every phase conversation, and B37 and B38 govern every phase close. A copy in each phase file would repeat each rule many times.

The project-owned clauses of `AGENTS.md` that act in some phases stay in `AGENTS.md`. These are A2 and A4–A7. A kit phase file ships to every linked repo, so it cannot hold a rule that only this repo owns.

The light-path statements B34, B35, B39, B40, B42, and the last sentence of B24 act in any session that sizes a change. So they stay always-loaded, and #10 merges them.

| Change | Bytes out of every session |
|---|---|
| #9: B11 and B26–B33 move, and a one-line pointer stays | about 3,450 |
| #10: the light-path statements merge | about 350 |
| #14: the B10 exception, B12, and A17 go | about 330 |
| #15: B43 moves | about 275 |
| #8: the kit's `CLAUDE.md` takes the managed form | about 290 |

## Contradictions and stale facts

- `CONTEXT.md:249` says that a parked idea graduates only when `/bench-shape-idea` pulls it into a map. `.bench/BENCH.md:179-180` says that a drain moves it to the board. FT89 owns this class.
- `.agents/commands/bench-drain.md:291` sends the drain's own diff approval to "the batch approval rule", which is B43 about spec sign-off. #15 gives the drain line its own words.
- `.agents/commands/bench-implement-spec.md:93` lets `--auto-approve` pre-approve "spec and ticket expansions". `.bench/BENCH.md:124` keeps unrelated scope, a material acceptance change, and a weakened guarantee for the reviewer. No file says how the two rules meet.
- `.agents/commands/bench-write-spec.md:86` asks for "one `capture/learnings.md` entry" without the `bench learning` verb. `.bench/BENCH.md:161-165` tells the agent to use the verb.
- `.agents/skills/bench-craft-delegate/SKILL.md:114` says to run the gate against a delegate's work. `roadmap/FT89.md:7-10` already records this conflict with ADR 0023.
- The B34 sentence "Behavior defects run focused regression checks, then the gate" came from the learning promoted in 75749bd9. The same learning gave the fix-and-gate approval that 2c7db3a8 removed, so the sentence now has no named phase.

## Unknowns and unverified edges

- The FT71 build figures at `specs/fresh-ticket-authors/spec.md:11` in commit 28dfcc8d have no other record in the tree.
- The failures behind A10, A13, and A14 are unknown, and so is the failure behind the B32 and B33 concurrency rules.
- One reader claimed that 05c76cb6 removed a B19 scope sentence about CLI output and artifacts. No coordinator read confirms it.
- The origin cells are the readers' agreed set. The coordinator re-opened only the disputed ones, so a later deletion must re-open its own origin lines.

## Verification record

- [x] The output opens with the recommendation, the scope, and the evidence status.
- [x] Facts, inferences, and the move proposal sit in separate sections, and no test result applies.
- [x] One synthesized section answers each question.
- [x] A table compares the byte effect of each change, with its consequence.
- [x] Every contradiction and every residual unknown is kept.
- [x] No diagram is drawn, because the tables show every material relation.
- [x] The unverified edges are stated.
- [x] Every load-bearing claim cites a path and a line, or a commit.
- [x] A validation plan ends the output.

## Validation plan

1. Before the spec locks, run `bench anchors` on each source file and on each destination, and confirm that each moved needle has its new home.
2. For each ticket, run `bench preflight build` with `--propose-writes` and add the anchor and fixture closure files that it lists.
3. After the moves, run the docs-currency-workflow check, the shared-rule check, the structured-phase contract check, and the prose budget check.
4. Measure each file with `wc -c` before and after, and report the delta for each file.
