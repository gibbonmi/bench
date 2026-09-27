# Review evidence with stable file pages

Status: staged

Roadmap: FT337

Decision source: `roadmap/FT337.md` (named reviewed artifact).

Verification log: 2 iteration(s) to accept — Opus/high accepted after B1 and N1–N8 fixes. Final edits cover partition refusals, repeated paths, decoder ownership, and durable tests.

## Problem

A review charge stores the complete diff as one source. An early edit moves
later page boundaries, even when the later file patch stays unchanged.
An author record commit after the charge also makes the prepared tip stale.

## Solution

Prepare each file patch as a separate source through the existing pager.
Keep the complete diff, including the spec, tickets, and review record.
An unchanged patch keeps its body and page digests across rounds.

Commit the author record before the review charge. Keep narrow rounds
provisional, record their reads and findings, and compare one full control
against the same frozen diff. The reviewer then decides permanent adoption.

## User stories

Line: gpt-5.6-sol / high.
Implementation-line reason: File identity and exact reconstruction form the hardest chunk. Existing seams suffice, but new behavioral tests must cover them. Guidance also steers future sessions.
Harder chunks: RE2.

### Stable and complete evidence

1. As an axis, I want stable pages for unchanged patches, so that another file's growth does not invalidate those bytes.
2. As an axis, I want stable pages after record edits, so that bookkeeping does not invalidate unchanged code.
3. As an axis, I want file identity independent of ordinals, so that an inserted file cannot misidentify retained bytes.
4. As an axis, I want changed patches to have changed digests, so that I cannot reuse stale content.
5. As a reviewer, I want the complete frozen diff, so that reduced retrieval does not hide a changed file.
6. As an axis, I want separate sources for equal-content files, so that each file keeps its own path.
7. As an axis, I want the requested source pair, so that later chunks do not review earlier work again.
8. As an axis, I want bounded pages for Unicode patches, so that retrieval preserves every byte.
9. As an axis, I want one selected file stream, so that I can retrieve missing evidence without other file bodies.
10. As an axis, I want current membership before reuse, so that a matching digest does not grant authority.

### Current and assessable review rounds

11. As a coordinator, I want author records committed before preparation, so that dispatch receives current evidence.
12. As an axis, I want stale charges refused, so that a later record commit cannot change my review subject.
13. As a reviewer, I want each narrow round's reads and findings, so that its actual coverage stays visible.
14. As a reviewer, I want a full control on the same diff, so that additional findings expose narrow review misses.
15. As a reviewer, I want the control's disposition before adoption, so that cheap review cannot establish its own adequacy.
16. As a fresh axis, I want my own required sources, so that another session's receipt cannot replace evidence I read.

## Implementation decisions

Use the existing diff owner to expose ordered file patches from its coherent
frozen snapshot. Preflight turns that snapshot into generated sources.
The pack owner continues to derive page boundaries and digests.
Do not duplicate pagination, path decoding, or changed-path knowledge.

The diff sources contain the existing prefix, each file patch, and the existing
suffix in reconstruction order. Use these exact descriptors. RE5 checks the
fragments, and RE9 selects files by their descriptors.

| Fragment | Role | Path | Bytes |
| --- | --- | --- | --- |
| Prefix | `diff-prefix` | `diff-prefix` | Existing tables through `diff_body:` and its newline |
| File patch | `diff-file` | Repository path defined below | Complete verbatim patch, including its headers |
| Suffix | `diff-suffix` | `diff-suffix` | Existing help table after the patch body |

An empty diff has the prefix and suffix, with no file source.
Every fragment uses `kind=generated` and `required=true`.
Every diff fragment's shared row uses `kind=diff`.

Use the tip path for each surviving file, including additions, modifications, and renames.
Use the base path for a deleted file. This rule also covers patches without `---` or `+++` lines.
A rename keeps both spellings in its verbatim patch. RE9's identity cases
cover an empty addition, empty deletion, binary change, mode change, and pure rename.

Keep the producer's ambient rename behavior and complete output bytes.
The inventory uses `--no-renames`, but the patch body can contain a rename.
Match paths by identity, never by inventory position. RE3 and RE9 exercise
rename detection both enabled and disabled.

If the body cannot partition into identified patches consistent with the inventory,
refuse the charge through the existing evidence failure path. Publish no partial
artifact. Keep ambient output unchanged. RE5 covers external-diff and forced-color
outputs that the partition cannot represent.

A type change can produce two patches for one path. Keep both sources in patch
order, with the same path and distinct pack-local source IDs. A path maps to
an ordered digest list. To retrieve one file, read each matching source in order.
RE3 preserves this mapping, and RE9 covers a file-to-symlink change.

One shared decoder reads patch path syntax. It decodes C-quoted paths in
`diff --git`, `---`, `+++`, and the rename or copy path headers.
The `rename from`, `rename to`, `copy from`, and `copy to` paths are already
repository-relative. Strip display prefixes only from headers that carry them.

For headerless patches, use the `diff --git` paths and extended headers.
The `new file mode` and `deleted file mode` headers identify absent sides.
Prefer explicit rename or copy paths when present. Never split unquoted paths
at each space. Resolve them against the frozen path inventory when needed.
Parse path headers only in the patch header region, before its content.

`internal/git` owns the shared patch-path decoder. Move the existing consumers
C-quote rule there. Both diff and consumers already import this package.
Keep the decoder pure and independent of both callers. Keep the consumers
hunk interpretation unchanged, except that a trailing space stays path text
by reviewer decision on R14. RE5 exercises hostile paths, and RE9 requires
each headerless patch's exact descriptor.

All fragments remain required and retain their actual producer provenance.
Shared evidence rows bind every diff fragment in reconstruction order.
The existing schema permits repeated generated sources. Keep its version and
public command grammar. The registry owns its shared-row description and
regenerates the format reference. RE5 and RE9 cover the published result.

Compare an unchanged file by its path and exact patch bytes under the same
base. A changed base can change the patch despite unchanged tip bytes.
Source ordinals remain local to each pack. The artifact identity still changes
with the source tip and provenance. Digest equality never replaces current
membership, role, or requiredness. RE3, RE4, RE7, and RE10 grade these cases.

The implementation phase states one order: author record commit, review charge,
then axis dispatch. Keep the current stale-tip refusal when a later commit
moves the tip. RE11 and RE12 cover the lawful and stale sequences.

Narrow review stays provisional. After terminal returns, record each round in
this spec. Record the frozen pair, evidence identity, axis, reads, findings,
and available usage. Keep unavailable usage unknown. Run one independent
full-retrieval control across all axes on RE2's same frozen diff.

The control receives no narrow findings before its return. Record agreement,
additional findings, and their dispositions. Check complete patch retrieval
against raw `git diff` bytes from the frozen pair.
A miss uses the existing repair route. A missing control leaves adoption
undecided. The reviewer owns permanent adoption.

The orchestrator owns RE10 and RE13 through RE16 at the RE2 review checkpoint.
It records native returns and the comparison after the ticket's green commit.
These rows remain mandatory before final reconciliation. Their absence blocks
that checkpoint, without assigning a ticket author to review its own work.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| RE1 / 1-order-review-preparation.md | Author records precede prepared review evidence | RE11, RE12 | Chunk-chain mutations and the current-binding scenario | no |
| RE2 / 2-page-review-diffs-by-file.md | Complete evidence retains unchanged file digests | RE1, RE2, RE3, RE4, RE5, RE6, RE7, RE8, RE9, RE10, RE13, RE14, RE15, RE16 | File-source scenarios, differential reconstruction, and the narrow/control comparison | yes |

Review RE1 before RE2 starts. Both tickets name shared command registries, so
their writes remain serial. RE2 consumes RE1's preparation order during review.
Each ticket includes its behavior and tests.

### Completion plan

This version 2 plan records the real author assignments of the build.
The build runs in Claude Code. The user selected `opus` at high effort for
each ticket author. This model is the Claude mid binding of the declared line.
The user selected `fable` for each review axis and for the full control.
The review-owned control comparison remains a required acceptance checkpoint.

```bench-completion-plan
{
  "version": 2,
  "execution": {
    "mode": "delegate",
    "run_id": "ft337-review-evidence-file-pages-full-20260926",
    "orchestrator_session": "claude:session_0156tkEZcRSowaafegWfFZJP",
    "author_limit": 1,
    "assignments": {
      "1-order-review-preparation.md": [
        {
          "session": "claude:bench-writer/re-t1-author",
          "assignment": "re-t1-author",
          "model": "opus",
          "effort": "high",
          "source": "c8c444ffae2fb1578cfa54a22fa632590ffbc322",
          "native_ref": "claude:agent/re-t1-author-20260926@c8c444ffae2fb1578cfa54a22fa632590ffbc322"
        },
        {
          "session": "claude:bench-writer/re-t1-repair-c1",
          "assignment": "re-t1-repair-c1",
          "model": "opus",
          "effort": "low",
          "source": "fa1b8717f0e6426fb28898cd9b0f72ba1d3b2839",
          "native_ref": "claude:agent/re-t1-repair-c1-20260926@fa1b8717f0e6426fb28898cd9b0f72ba1d3b2839",
          "predecessor": "claude:bench-writer/re-t1-author",
          "trigger": "user-directed",
          "stopped": "claude:agent/re-t1-author returned its final report and holds no write charge",
          "preserved": "fa1b8717f0e6426fb28898cd9b0f72ba1d3b2839"
        },
        {
          "session": "claude:bench-writer/re-t1-repair-c2",
          "assignment": "re-t1-repair-c2",
          "model": "opus",
          "effort": "high",
          "source": "ca686c6fd62844452303fc72db401d4adb8e438b",
          "native_ref": "claude:agent/re-t1-repair-c2-20260926@ca686c6fd62844452303fc72db401d4adb8e438b",
          "predecessor": "claude:bench-writer/re-t1-repair-c1",
          "trigger": "user-directed",
          "stopped": "claude:agent/re-t1-repair-c1 returned its final report and holds no write charge",
          "preserved": "ca686c6fd62844452303fc72db401d4adb8e438b"
        }
      ],
      "2-page-review-diffs-by-file.md": [
        {
          "session": "claude:bench-writer/re-t2-author",
          "assignment": "re-t2-author",
          "model": "opus",
          "effort": "high",
          "source": "ed63e73259e48d881db2815b282f4c4e381e4725",
          "native_ref": "claude:agent/re-t2-author-20260926@ed63e73259e48d881db2815b282f4c4e381e4725"
        },
        {
          "session": "claude:bench-writer/re-t2-repair-c1",
          "assignment": "re-t2-repair-c1",
          "model": "opus",
          "effort": "low",
          "source": "a3979037a8ea741ab3c941f8e7b2d99d4965c528",
          "native_ref": "claude:agent/re-t2-repair-c1-20260926@a3979037a8ea741ab3c941f8e7b2d99d4965c528",
          "predecessor": "claude:bench-writer/re-t2-author",
          "trigger": "user-directed",
          "stopped": "claude:agent/re-t2-author returned its final report and holds no write charge",
          "preserved": "a3979037a8ea741ab3c941f8e7b2d99d4965c528"
        },
        {
          "session": "claude:bench-writer/re-t2-repair-c2",
          "assignment": "re-t2-repair-c2",
          "model": "opus",
          "effort": "low",
          "source": "bdd6766f03fba23b768e874733592fadc664758e",
          "native_ref": "claude:agent/re-t2-repair-c2-20260926@bdd6766f03fba23b768e874733592fadc664758e",
          "predecessor": "claude:bench-writer/re-t2-repair-c1",
          "trigger": "user-directed",
          "stopped": "claude:agent/re-t2-repair-c1 returned its final report and holds no write charge",
          "preserved": "bdd6766f03fba23b768e874733592fadc664758e"
        },
        {
          "session": "claude:bench-writer/re-t2-repair-c3",
          "assignment": "re-t2-repair-c3",
          "model": "opus",
          "effort": "low",
          "source": "b0f6171b9b3cccb2bd4eb5692b58c98c091ee65a",
          "native_ref": "claude:agent/re-t2-repair-c3-20260927@b0f6171b9b3cccb2bd4eb5692b58c98c091ee65a",
          "predecessor": "claude:bench-writer/re-t2-repair-c2",
          "trigger": "user-directed",
          "stopped": "claude:agent/re-t2-repair-c2 returned its final report and holds no write charge",
          "preserved": "b0f6171b9b3cccb2bd4eb5692b58c98c091ee65a"
        },
        {
          "session": "claude:bench-writer/re-t2-repair-c4",
          "assignment": "re-t2-repair-c4",
          "model": "opus",
          "effort": "low",
          "source": "085a0f5b4f90ba14ef7c8d1a8177b767e45d5a86",
          "native_ref": "claude:agent/re-t2-repair-c4-20260927@085a0f5b4f90ba14ef7c8d1a8177b767e45d5a86",
          "predecessor": "claude:bench-writer/re-t2-repair-c3",
          "trigger": "user-directed",
          "stopped": "claude:agent/re-t2-repair-c3 returned its final report and holds no write charge",
          "preserved": "085a0f5b4f90ba14ef7c8d1a8177b767e45d5a86"
        }
      ]
    }
  },
  "chunks": [
    {
      "id": "RE1",
      "tickets": [
        "1-order-review-preparation.md"
      ],
      "verification": [
        {
          "id": "1-workflow",
          "command": "bench test --check docs-currency-workflow",
          "ticket": "1-order-review-preparation.md"
        },
        {
          "id": "1-record-order",
          "command": "go test -count=1 -parallel=2 ./internal/preflight/evidencecmd",
          "ticket": "1-order-review-preparation.md"
        }
      ]
    },
    {
      "id": "RE2",
      "tickets": [
        "2-page-review-diffs-by-file.md"
      ],
      "verification": [
        {
          "id": "2-file-evidence",
          "command": "go test -count=1 -parallel=2 ./internal/diff ./internal/git ./internal/consumers ./internal/chargeevidence ./internal/preflight/...",
          "ticket": "2-page-review-diffs-by-file.md"
        },
        {
          "id": "2-ports",
          "command": "bench test --check injected-port-registry",
          "ticket": "2-page-review-diffs-by-file.md"
        },
        {
          "id": "2-workflow",
          "command": "bench test --check docs-currency-workflow",
          "ticket": "2-page-review-diffs-by-file.md"
        }
      ]
    }
  ],
  "final_verification": [
    {
      "id": "coverage",
      "command": "bench coverage --check specs/review-evidence-file-pages/spec.md"
    },
    {
      "id": "file-evidence",
      "command": "go test -count=1 -parallel=2 ./internal/diff ./internal/git ./internal/consumers ./internal/chargeevidence ./internal/preflight/..."
    },
    {
      "id": "workflow",
      "command": "bench test --check docs-currency-workflow"
    }
  ]
}
```

## Testing decisions

Drive `preflight.Command` through the existing `preflighttest` fixture.
Retrieve the published artifact through the normal evidence command.
Pure partition tests belong under the diff owner. They supplement the complete
producer-to-reader check. Pack format tests provide the paging precedent.

Extend `TestEvidenceReviewProvenanceRows` for every generated fragment.
It retains the existing provenance guarantee as source cardinality changes.

RE5 compares old and new outputs from deterministic Git fixtures.
Pin commit dates, identities, content, and relevant Git configuration.
Capture the old full output before the refactor, including its framing bytes.
The new producer must match that stored baseline. Raw `git diff` from the
frozen pair independently checks each patch body. Do not regenerate expected
bytes through the new partition owner.

Use test-first development for digest stability and complete reconstruction.
Observe the old monolithic source fail before the production change.
An independent expectation requires a demonstrated omission or swap mutation.
Record its restored green result. Do not implement future tests during specification.

The gate's ordinary test phase executes these package tests.
The registered `docs-currency-workflow` check executes guidance anchors.
Its existing workflow fixture owner proves each changed anchor's mutation.
This work introduces no gate phase or new executable trust route.

### Seam diagram

```text
frozen base + tip -> diff snapshot -> file patches -> preflight sources
                                                          |
                                                          v
                                                 existing pack pager
                                                          |
                                                          v
                                              evidence command reader

committed author record -> review charge -> check-current -> fresh axes
                                                          |
                                                          v
                                         narrow round + same-diff control
```

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| RE1 | 1 | Growing `a.go` leaves every `z.go` patch page digest unchanged under one base | Planned `TestReviewFilePageStability` through `preflight.Command` | The first patch crosses a page boundary, so global pagination changes later digests |
| RE2 | 2 | Separate spec, ticket, and review-record edits preserve the untouched code patch digests | Planned `TestReviewFilePageStability` document cases | Each document grows independently, so excluding only one cannot pass |
| RE3 | 3 | Adding an earlier file preserves the retained file's path-to-digest mapping | Planned `TestReviewFileIdentity` with both rename settings | The inserted source moves ordinals, so ordinal-only reuse selects the wrong body |
| RE4 | 4 | Editing a file patch changes its source digest | Planned `TestReviewFileIdentity` changed-patch case | Retaining the previous body or digest returns the old value |
| RE5 | 5 | Retrieved fragments reconstruct the complete original frozen diff byte for byte | Planned `TestReviewFileReconstruction` at the diff and preflight seams | An omitted patch or framing fragment changes the independent baseline comparison |
| RE6 | 6 | Equal-content files retain separate manifest path memberships | Planned `TestReviewFileIdentity` duplicate-content case | Equal-length paths can share later page digests, so page deduplication can lose one membership |
| RE7 | 7 | A later charge contains only its requested predecessor-to-tip diff | Planned `TestReviewFileReconstruction` predecessor-base case | A fallback to main includes the earlier chunk's sentinel patch |
| RE8 | 8 | A multi-page Unicode patch reconstructs exactly through the existing page protocol | Planned `TestReviewFileReconstruction` large-file case | A split byte, repeated page, or omitted page changes the body |
| RE9 | 9 | Selection by the declared file identity returns only that file's complete patch | Planned `TestReviewFileSelectedStream` with the identity cases | A base-path rename, empty descriptor, or inventory-position match selects no source or the wrong body |
| RE10 | 10 | Reuse requires the new manifest's membership, role, and requiredness | review-owned: inspect RE2's retrieval record | A digest-only reuse decision lacks the source descriptors |
| RE11 | 11 | Guidance orders the author record commit before the review charge | `docs-currency-workflow` with chunk-chain omission and order-swap mutations | Commit-before-dispatch alone misses the charge-before-commit mutation |
| RE12 | 12 | Current binding refuses a review charge after a later record commit | Planned `TestReviewRecordChargeOrder` through `preflight.Command` | The fixture first binds after commit, then moves the record and requires the stale-tip refusal |
| RE13 | 13 | This spec records every narrow round's actual reads and findings | review-owned: compare native returns with the round table | An omitted axis or read inventory prevents reconciliation with its native return |
| RE14 | 14 | A full-retrieval control reviews the narrow round's exact frozen diff | review-owned: compare transcripts and source pins against the frozen pair's raw Git patch bytes | A different pair or incomplete retrieval cannot establish the control |
| RE15 | 15 | Permanent adoption waits for the control comparison and reviewer disposition | review-owned: inspect the comparison and decision | Cost reduction alone supplies neither missed-finding analysis nor the decision |
| RE16 | 16 | Each fresh axis reads its own required context | review-owned: inspect each axis's native read inventory | A transferred receipt delivers no local source read |

### Edge inventory

The audience includes the kit and linked repositories that use review charges.
RE5's differential fixture covers the producer's following shapes:

- An empty diff and a diff containing only the spec or review record.
- Added, modified, deleted, and renamed files, including a pure rename.
- An empty addition and deletion with no `---` or `+++` headers.
- A mode-only patch, binary marker, and committed symlink patch.
- A file-to-symlink change with two patches for the same path.
- An external-diff or forced-color body that cannot partition, with no partial publication.
- A missing final newline and content that resembles a patch header.
- Paths with spaces, globs, quotes, backslashes, and non-ASCII bytes.
- Permitted tab, newline, and return bytes under the current renderer.
- Refused control bytes, with no partial publication.

RE1 and RE8 cover page-size edges. RE6 covers equal-content members.
RE7 covers a base other than main. RE12 covers a record-only tip change.
Existing movement and missing-source tests remain in focused verification.

Won't handle: stable digests after a changed patch or base — the current manifest identifies the new bytes.
Won't handle: delivery through another axis's receipt — RE16 keeps the fresh axis as the in-scope reader.
Won't handle: automatic permanent adoption — RE15 keeps the reviewer as the in-scope decision owner.
Won't handle: an ambient `diff.noprefix` path outside the frozen inventory — the charge refuses and publishes nothing, by reviewer decision on R13.

## Ownership fences

- `.agents/commands/bench-implement-spec.md`
- `internal/anchors`
- `tests/canary/workflow-guidance-anchors`
- `internal/preflight`
- `internal/preflight/evidencecmd`
- `internal/diff`
- `internal/git`
- `internal/consumers`
- `internal/chargeevidence`
- `.agents/skills/bench-craft-delegate/references/charge-evidence-format.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/preflight_version_test.go`
- `cmd/bench/spill_support_test.go`
- `ROADMAP.md`
- `decisions/jev-advisor/assets/jev-benchmark-trial.md`
- `docs/research/aibadger-assessment.md`
- `docs/research/parallel-implementation-wave.md`
- `docs/research/roadmap-review-2026-09-25.md`
- `roadmap/FT346.md`
- `roadmap/FT347.md`
- `roadmap/FT348.md`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/conformance/injected_ports_registry_test.go`
- `reviews/review-evidence-file-pages.md`
- `projects/benchkit.md`
- `tests/canary/guidance-prose-budgets/over-budget-skill`
- `tests/canary/line-routing/line-binding-prose-drift`
- `tests/canary/skill-description-budgets/budget-table-missing`
- `tests/canary/skill-description-budgets/description-folded`
- `tests/canary/skill-description-budgets/description-missing`
- `tests/canary/skill-description-budgets/over-budget-command`
- `tests/canary/skill-description-budgets/over-budget-description`

Reviewer disposition: approved on 2026-09-26, with the ticket graph and the fences.
The debug prerequisite is integrated. Further notes names its landing.
This phase includes that repair. The implementation tickets consume its canonical
guidance and add only the preparation order and file-page behavior.

## Out of scope

No separate capability is deferred: 0 deferred edits, 0 deferred gate runs.
This build adds no automatic persistent cache, provider-token target, new pack
version, or broader review-policy redesign. Readers reuse available bytes under
the existing membership rule.

## Further notes

### Facts and reader sweep

Evidence date: 2026-09-25.
Evidence base: `93282b885183fb976de5d8af24d94de9c861279c`.
Refresh when a cited producer, reader, registry, or debug prerequisite changes.

The factual questions are the producer, its readers, preparation order, and
mutation ownership. Inspected local owners answer all four. No external API
or upstream document determines this contract.

| Question | Source fact | Consequence |
| --- | --- | --- |
| Why do pages move? | `internal/preflight/review.go:51` collects one diff. `internal/chargeevidence/pack.go:104` pages each source independently. | Split file sources and keep the pager |
| Does review always use main? | `internal/preflight/review.go:56` passes the explicit base and tip | Preserve explicit-base behavior |
| Who owns membership? | `internal/chargeevidence/manifest.go:15` defines source descriptors. `internal/chargeevidence/metadata.go:43` writes shared rows. | Keep membership in the existing schema |
| Who decodes patch paths? | `internal/consumers/hunks.go:64` decodes only `---` and `+++` paths | Share its C-quote rule and extend headerless patch identity |
| Who grades order? | `internal/anchors/registry_chunk_chain.go:7` pins commit before dispatch | Extend that owner with charge order |
| Who executes the guidance check? | `internal/conformance/registry/registry.go` registers `docs-currency-workflow` and its fixture family | Extend the existing executed check |

A Git fixture confirmed the inventory and patch cardinalities can differ.
With rename detection enabled, six inventory paths produced five patches.
With rename detection disabled, the same pair produced six patches.
The pure rename, empty addition, binary change, and mode change lacked
`---` and `+++` headers. The fixture used a committed base and tip.
Evidence: `/tmp/ft337-git-identities-ici8r1fo` on 2026-09-25.

A second fixture produced deletion and addition patches for one file-to-symlink
path. Forced color prefixed patch headers with ANSI escapes. An external diff
replaced the patch body. Evidence: `/tmp/ft337-git-review2-z6o6az97` on 2026-09-25.
RE5 and RE9's planned durable cases supersede these temporary observations.

The roadmap's phrase "from `main`" is stale factual shorthand. Preserve the
current explicit chunk-base decision. This correction remains subject to
reviewer veto.

The reader sweep used `rg --hidden`, excluding `.git`, across the whole tree.
It included scripts and workflows. The changed fact is diff-source cardinality,
not a field or pack version. This inventory covers readers and direct helpers.

| Reader or owner | Disposition |
| --- | --- |
| `internal/preflight/review.go`: `collectReviewEvidence`, `reviewChargePack`, `reviewMetadata` | RE2 changes composition |
| `internal/diff/diff.go`, `range.go`, `snapshot.go` | RE2 preserves output and exposes its file snapshot |
| `internal/consumers/hunks.go` and tests | RE2 moves the C-quote rule to `internal/git` |
| `internal/git` | RE2 owns the shared decoder without caller imports |
| `internal/preflight/review_charge_test.go` | RE2 updates collector and provenance assumptions |
| `internal/preflight/evidencecmd/evidence_review_test.go` | RE2 replaces fixed generated-source counts with file membership assertions |
| Evidence command, consumer, mode, export, and budget tests under `internal/preflight/evidencecmd` | RE2 verifies traversal, binding, export, and bounds |
| `internal/preflight/preflighttest/fixture.go` | RE2 extends the shared fixture |
| Schema, metadata, manifest, pack, read, and export owners under `internal/chargeevidence` | RE2 preserves the format and updates necessary source-shape assumptions |
| `internal/chargeevidence/reference.go` and its generated format reference | RE2 derives the shared-row description from its registry |
| Format and refusal tests under `internal/chargeevidence` | RE2 retains independent mutation coverage |
| `internal/anchors/registry_chunk_chain.go` and its test | RE1 changes the order rule |
| `internal/anchors/registry_charge_binding.go` and its test | RE2 preserves format-reference anchors |
| `tests/canary/workflow-guidance-anchors` | RE1 owns the changed order mutation |
| `internal/tickets/registry_data.go` | Excluded: it requires the listed command registries but changes no binding |
| `internal/conformance/injected_ports_registry_test.go` | Both tickets retain existing applicable bindings |
| Scripts and workflows | Excluded: no generated diff-source cardinality reader was found |
| `internal/reviewrecord` | Excluded: its digest binds results, not page layout |

No shipped-surface claim introduces a repository-only path. The format reference
already belongs to the linked kit. Tests and fixtures remain kit-local.

### Source-sentence-to-row table

| Reviewed source clause | Rows |
| --- | --- |
| "The fix pages the diff source by file" | RE1, RE2, RE3, RE5, RE6, RE8, RE9 |
| "an unchanged file keeps its page digests" | RE1, RE2, RE3, RE4, RE7 |
| "a reader never reads bytes again that it already holds" | RE9, RE10, RE16 |
| "That narrow shape is provisional" | RE15 |
| "what each narrow round read and found" | RE13 |
| "one full-retrieval control review of the same diff" | RE14, RE15 |
| "the author record commit, then the review charge, then the axis dispatch" | RE11, RE12 |

Each clause occurs once in `roadmap/FT337.md`. Its ledger records FT71's nine
axis runs, moved-tip charge, and final narrow round. Those observations set no
new numeric cost target.

### Pre-review proof checklist

- Cited symbols: the owner inventory resolves each named production function.
- Import edges: preflight imports diff and chargeevidence. Diff and consumers both import git. The shared decoder adds no reverse import. No cross-package test helper is promised.
- Source-row clauses and occurrences: the preceding table covers the reviewed artifact and its ledger.
- Promised field labels: `sources`, `pages`, `shared_evidence`, `role`, `path`, `required`, and `sha256` retain their grammar.
- Changed-function callers: `reviewChargePack` calls the collector and metadata owner. The diff command and preflight share the snapshot producer.
- Copy survival: RE5 rejects an extra monolithic payload through exact reconstruction and source-membership assertions.
- Rendered-shape readers: `review.go` and `evidence_review_test.go` assume one diff source. The generated reference describes shared-row order.

### Red-mutation owners

| Rows | Mutation owner and probe |
| --- | --- |
| RE1, RE2, RE3 | `internal/preflight/review.go`: replace file inputs with the joined diff input |
| RE4 | `internal/preflight/review.go`: substitute the previous patch in the changed-file case |
| RE5 | The file snapshot under `internal/diff`: omit the deleted-file patch while retaining its inventory entry |
| RE6, RE9 | `internal/preflight/review.go`: bind the second file to the first file source |
| RE7 | `internal/preflight/review.go`: replace the predecessor base with main |
| RE8 | The generated file input under `internal/preflight`: omit its final page of bytes |
| RE11 | `internal/anchors/registry_chunk_chain.go`: remove charge-order enforcement while retaining commit-before-dispatch |
| RE12 | `internal/preflight/charge_pack.go:115`: accept the stale source tip |
| RE10, RE13, RE14, RE15, RE16 | Review-owned: reject missing reads, unequal pairs, incomplete retrieval, or unapproved adoption |

Resolve the exact edited symbol before the mutation. Record each behavioral red
and restored green before implementation completion. Use `bench probe` for
finished-tree mutations. A compile failure proves no behavioral predicate.

### Flagged additions and unknowns

Flagged additions: none.
File sources and differential fixtures implement the reviewed per-file option.
The debug prerequisite is integrated from landing `3119dc72d48ad44a0a382bdf13a437a877b3e538`.
The real control result and usage remain unknown at spec time.

### Narrow and control evidence

The orchestrator filled this table at RE2's review checkpoint from native
returns. Each row names its frozen pair and evidence identity. Ticket 2's
green commit precedes these records. Usage counts the session tokens that the
harness reported.

| Round / axis | Mode | Frozen pair / identity | Reads | Findings | Usage |
| --- | --- | --- | --- | --- | --- |
| RE2 r1 / Standards | narrow | `4d9e6aa8..46e287c8` / `sha256:a8bd30e5` | one `git diff`, the standards files, targeted code; no evidence page | R5, R6, R7, R8 | 127353 tokens |
| RE2 r1 / Spec | narrow | `4d9e6aa8..46e287c8` / `sha256:a8bd30e5` | one `git diff`, the spec, ticket 2, coverage rows, the author record; no evidence page | R10 | 125710 tokens |
| RE2 r1 / Coverage | narrow | `4d9e6aa8..46e287c8` / `sha256:a8bd30e5` | one `git diff`, targeted code, four consumer pages of `s30`; four probes | R11, R12, R13 | 146411 tokens |
| RE2 r1 / all axes | full control | `4d9e6aa8..46e287c8` / `sha256:a8bd30e5` | all 31 sources and all 44 pages to stream end; no `git diff` | R5, R9, R10, R13, R14, R15 | 227918 tokens |

The joined `diff-file` sources equal raw `git diff 4d9e6aa8 de294013` byte
for byte, so the control retrieval was complete.

Comparison: the narrow round found 8 raw findings with 3 sessions and 399474
tokens. The control found 6 raw findings with 1 session and 227918 tokens.
Both found R5, R10, and R13. Only the narrow round found R6, R7, R8, R11, and
R12, which include the two gaps that silent probes proved. Only the control
found R9, R14, and R15. R14 is an inherited decoding defect, and the narrow
Coverage axis gave R15 as advice only.

No axis reused bytes from another
artifact, so RE10 found no reuse decision to grade.

The reviewer extended the RE2 repair allowance by one cycle for two
checkpoint gate reds, and by one more cycle for finding R18 only.
The reviewer fenced the eight paths that the destination merge of `main`
brought in, so a normal RE2 merge round reviews them.

Reviewer disposition: the narrow shape stays provisional and is not the
permanent rule. The reviewer decided R13 as a Won't-handle refusal with a
test, and R14 as a fix in the RE2 repair.
Validation plan: run focused tests, prove mutations, review each chunk, and record the control before final reconciliation.
