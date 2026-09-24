# Fresh ticket authors review record

## FA ticket 1 author evidence

Author session: `claude:bench-writer/fta-t1-author`. Line: opus / high / cap 3 attempts. Source tip at start: `42fcfb131425e074a2bf8b536ad00767bcfe87bc`. Ticket commit: `41f19bce59bd91357c7c1a94672ed1577ba538c1`.

The author added the FA rows to the anchor registries before the guidance edit. Then the author ran the rows against the unchanged guidance, and each row went red. After the guidance edit, each row went green.

### Done claims

The coordinator fills the label cell after its own probe.

| Row | Status | Confidence | Label |
|---|---|---|---|
| FA1 | verified | 9 | |
| FA2 | verified | 9 | |
| FA3 | verified | 9 | |
| FA4 | verified | 9 | |
| FA5 | verified | 9 | |
| FA6 | verified | 9 | |
| FA7 | verified | 9 | |
| FA8 | verified | 9 | |
| FA9 | verified | 9 | |
| FA10 | verified | 9 | |
| FA11 | verified | 9 | |
| FA12 | verified | 9 | |
| FA13 | verified | 9 | |
| FA14 | verified | 9 | |

### Red then green

The red column quotes the diagnostic from `bench test --full --package ./internal/conformance --run TestFreshTicketAuthors` on the unchanged guidance. The FA6 red comes from `TestEvidenceBoundedActionSentencesArePinned/build` on the unchanged guidance. The green column is the same test after the guidance edit.

| Row | Seam | Red before the edit | Green after the edit |
|---|---|---|---|
| FA1 | Require, `registry_ft311_review_dispatch.go` | operating guide dropped the fresh author session for each ticket on the declared line | pass |
| FA2 | Require, `registry_ft311_review_dispatch.go` | operating guide dropped the limit-1 version 2 delegate plan | pass |
| FA3 | Require, `registry_retained_workflow.go` | operating guide dropped the concurrent authors and the full tier range of `--delegate` | pass |
| FA4 | Forbid, `registry_ft311_review_dispatch.go` | operating guide restored the retained implementation session | pass |
| FA5 | Require, `registry_ft311_preparation.go` | `bench-implement-spec.md` dropped the narrow author charge | pass |
| FA6 | `TestEvidenceBoundedActionGuidance`, build delivery `a narrow author read` | the verified-delivery sentence is pinned by 0 Require rows | pass |
| FA7 | Forbid, `registry_ft311_preparation.go` | `bench-implement-spec.md` restored the full evidence retrieval | pass |
| FA8 | Forbid, `registry_ft311_preparation.go` | `bench-implement-spec.md` restored the verified-delivery build rule | pass |
| FA9 | Require, `registry_ft311_preparation.go` | `bench-implement-spec.md` dropped the orchestrator's read bound | pass |
| FA10 | Require, `registry_ft311_preparation.go` | `bench-implement-spec.md` dropped the handoff refresh at each chunk checkpoint | pass |
| FA11 | Require, `registry_ft311_review_dispatch.go` | operating guide dropped the orchestrator's final reconciliation | pass |
| FA12 | Require, `registry_retained_workflow.go` | operating guide dropped the fresh repair session for each affected ticket | pass |
| FA13 | Forbid, `registry_retained_workflow.go` | operating guide restored the ticket-author repair ownership | pass |
| FA14 | Require, `registry_retained_workflow.go` | reference dropped the version 2 plan amendment or its verification for each ticket | pass |

### Probe verdicts

Each probe ran through `bench probe` after `bench worktree build`. Each restore reads `yes`.

| Row | File | Mutation | Check | Verdict | Diagnostic |
|---|---|---|---|---|---|
| FA4 | `.bench/BENCH.md` | swap: the retired authorship sentence returns | `docs-currency-workflow` | bit | operating guide restored the retained implementation session |
| FA7 | `.agents/commands/bench-implement-spec.md` | swap: the retired full-retrieval sentence returns | `docs-currency-workflow` | bit | restored the full evidence retrieval |
| FA8 | `.agents/commands/bench-implement-spec.md` | swap: the retired verified-delivery sentence returns | `docs-currency-workflow` | bit | restored the verified-delivery build rule |
| FA13 | `.bench/BENCH.md` | swap: the retired repair sentence returns | `docs-currency-workflow` | bit | operating guide restored the ticket-author repair ownership |
| FA6 | `internal/conformance/charge_evidence_guidance_test.go` | swap: the build delivery returns to `verified delivery` | `TestEvidenceBoundedActionGuidance` | bit | no registry row states the verified-delivery diagnostic |
| FA1 | `.bench/BENCH.md` | omission: the fresh author sentence leaves | `docs-currency-workflow` | bit | operating guide dropped the fresh author session for each ticket on the declared line |

Each FA4, FA7, FA8, and FA13 probe adds the retired sentence and keeps every other needle. So only the Forbid row can make the check red.

### Verification

| Check | Verdict | Elapsed |
|---|---|---|
| `bench test --check docs-currency-workflow` | pass | 781 ms |
| `bench test --package ./internal/conformance --run TestRootConformance` | pass | 4954 ms |
| `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner` | pass | 13119 ms |
| `bench test --package ./internal/...` | pass, 11 capability skips | 36726 ms for the conformance package |
| `bench gate-prose . -- <three edited Markdown files>` | pass | not reported |
| `go vet ./...` | pass | not reported |
| `bench commit` lane | pass | not reported |

### Flags for review

- The ticket text puts the verification split in the build phase too. The build phase points to `.bench/BENCH-reference.md` for the amendment contents, so that fact has one source.
- The plan-expansion paragraph of `.bench/BENCH.md` now names the orchestrator, not the retained author.
- The `.bench/BENCH.md` line budget is 185 lines. To stay in budget, the author joined the wrapped landing-shape sentence onto one line.
- The new rows use the diagnostic prefix `fresh ticket author: `. `TestFreshTicketAuthors` makes each row bite and reads the live guidance.
