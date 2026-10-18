## Outcome

The anchor-admission-prune build landed at 3854648849e2193c95c9933b9639c6f1d858a75f from the reviewed pair d568cb05..03ebd20f, and the landing gate was green. The registry now admits five anchor classes, and the step-scoped kind is gone. An admission test refuses an unclassified Require row. ADR 0012 and six guidance sites state the rule.

The build ran ten tickets in seven chunks as a `--full --auto-approve` run with an author limit of 1. Two tickets and one chunk came from in-build fixes by reviewer direction. Ticket 4a made `writes-resolve` resolve a deletion of an earlier chunk. Ticket 6a made a parameter grammar resolve through its callers. The chunks used 0, 1, 0, 1, 1, 2, and 2 repair cycles, and no chunk went past the allowance.

## Gate-stage timings

- landing: commit 3854648849e2193c95c9933b9639c6f1d858a75f, trace 52d154c6f216e071782b6a872a3c2339
- gofmt: 173 ms
- vet: 1451 ms
- test: 187063 ms
- race: 2505 ms
- system: 59471 ms
- shellcheck: 570 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ticket 1 to 8 authors | each ticket met its rows inside its fence | held, with the repairs below | unknown | abstain | claude-opus-5-5 / medium, high for ticket 8 / implementer |
| tickets 4a and 6a authors | each preflight fix resolved the refused charge | held after one repair each | unknown | abstain | claude-opus-5-5 / medium / implementer |
| orchestrator | the AP-C1 and AP-C2 repair counts match the declared expectation | held for those two chunks | unknown | abstain | claude-opus-5-5 / high / orchestrator |

No author charge required a claims table with integer confidences, so the build holds no labeled pair. Brier mean: unknown over 0 pairs; 3 abstentions.

## Coordinator catches

- A silent coordinator probe showed that the AP-C5 repair 1 built the ROUTE test cases from the new bound constant, so the AP13 bound was unpinned. Repair 2 restored a literal seven-word case.
- The ticket 3 author was resumed twice, because `bench record ticket` binds a non-last ticket to the frozen tip by commit. The ticket 4 author then had to refreeze at the record tip with a new id prefix.
- The AP-C7 checkpoint refused the review chain, because the AP-C3 comment-only correction changed source after its freeze. AP-C7 took a refreeze at the frozen AP-C3 tip and a third review round.
- The ticket 5, ticket 6a, and ticket 7 charges each refused on a predictive preflight check: `writes-resolve`, `registry-closure`, and `charge-premises` grammar effects. Each refusal took a plan expansion or an in-build fix.
- The ticket 6a author stopped at the structure budget and split the fix into two new files. The coordinator held the citation back from the plan commit, because the lane grades the committed tree.
- One coordinator probe ran while another author ran checks in the same tree. The author's seven checks passed, so no run was affected.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-fold-retained-workflow.md | 0 | none |
| 2-prune-core-registry.md | 1 | delegate-error |
| 3-prune-build-registries.md | 0 | none |
| 4-prune-phase-registries.md | 0 | none |
| 4a-resolve-earlier-deletions.md | 1 | delegate-error |
| 5-prune-hand-pins.md | 1 | delegate-error |
| 6-retire-step-kind.md | 0 | none |
| 6a-root-parameter-grammars.md | 1 | delegate-error |
| 7-admit-by-class.md | 2 | delegate-error, delegate-error |
| 8-state-admission-rule.md | 2 | delegate-error, delegate-error |

## Agent-experience improvements

### Bench CLI

- The landing census counted 34 raw calls: rg 20, cat 5, sed 4, cd 2, and one each of awk, ls, and wc. Their Bench forms are the Grep tool or `bench consumers` for rg, the Read tool or `bench preflight evidence` for cat and sed, `bench worktree exec` for cd, none for awk, `bench worktree path` for ls, and `bench structure` for wc; add a `bench anchors --cursor` page so a large anchor table needs no raw read.
  Feeds: new
- Make `bench record ticket` compare a staged run with the chunk source by source digest, and let the last ticket's freeze call owe every chunk requirement once.
  Feeds: new
- Make `bench record ticket` skip an id prefix that the requirement id already starts with, and let a rerun supersede its earlier entry.
  Feeds: new
- Make `bench record verification --final` accept a staged excerpt, so a final run needs no hand-written excerpt file.
  Feeds: new

### Skills

- State in the implement-spec phase that a source correction after a chunk freeze moves the chunk tip, and that the next chunk base is the moved tip.
  Feeds: new
- State in the implement-spec phase that an earlier ticket of a multi-ticket chunk records at the last author's freeze before the record commit.
  Feeds: new

### Process

- Remove the predictive closure and premise checks under FT293, because three of this build's charge refusals came from them.
  Feeds: FT293
- Require a claims table with integer confidences in every author and repair return, so a build holds labeled pairs.
  Feeds: none
- Start no coordinator probe while an author runs checks in the same tree.
  Feeds: none