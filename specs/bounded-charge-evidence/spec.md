# Bounded charge evidence

Status: staged

Decision source: specs/bounded-charge-evidence/decisions/bounded-charge-evidence.md

Verification log: 2 iteration(s) to accept — Sol/high independently verified repairs for five initial findings and two producer-local checkpoint omissions. Implementation evidence remains planned only.

## Problem

A charge consumer needs complete pinned inputs before action.
The compact response names omitted evidence, then directs the consumer to an oversized full response.
The renderer makes evidence completeness depend on one response containing every source body.
The baseline reproduction returned 77,162 bytes and exceeded the observed transport ceiling.

That measurement belongs to commit `9deb0a7af31712427ff47d6fd0515e8458dbde0d`.
Spec authoring uses `b20db7d4a31f7b1d49d3a68907e377dace2a09be`.
The later debug-loop fixture change does not revise the historical measurement.

## Solution

Explicit charge preparation creates one immutable evidence artifact in repository-common storage.
Stateless preflight reads return bounded fragments of its manifest and exact source bytes.
The consumer verifies membership and complete delivery before the phase permits action.
Every operation stays within 48,000 encoded stdout bytes.

Build and review use one evidence protocol.
Review preparation freezes generated captures once per attempt.
Historical reads survive worktree release.
A separate current-action check validates the current assignment and source pins.

Approval, the task supplement, and available consumer context remain separate requirements.

## User stories

Line: gpt-5.6-sol / high.
Implementation-line reason: immutable publication and bounded verification need careful composition at known preflight and filesystem seams.
The format is precise, but production lifecycle and transport proof remain unimplemented.
Harder chunks: CE-C1B, CE-C3.

1. As an author, I want bounded access to the complete selected build charge, so that I can start with every required input.
2. As a consumer, I want separate response and evidence states, so that one page cannot imply complete delivery.
3. As a consumer, I want the full manifest before source pages, so that I can verify membership independently.
4. As a consumer, I want exact source bytes, so that escaping and missing newlines cannot change evidence.
5. As a consumer, I want bounded metadata and errors, so that large values cannot break the transport.
6. As an operator, I want a deterministic logical identity, so that identical evidence can share one artifact.
7. As a consumer, I want strict pack and manifest validation, so that invalid storage cannot supply evidence.
8. As an author, I want the existing source refusals, so that preparation cannot bless unsafe or stale inputs.
9. As an author, I want one movement retry, so that a charge cannot mix two source states.
10. As a consumer, I want historical reads from another worktree, so that assignment release cannot destroy prepared evidence.
11. As an author, I want a separate current-action check, so that historical evidence cannot authorize a changed source.
12. As a consumer, I want full artifact verification, so that corruption outside the selected page remains detectable.
13. As a consumer, I want stateless source navigation, so that I can resume without a server reading session.
14. As a reviewer, I want frozen generated review evidence, so that every axis receives the same captured bytes.
15. As a reviewer, I want declared collector provenance, so that changed captures or invocation inputs change the identity.
16. As an operator, I want a quota over temporary and published bytes, so that concurrent preparation cannot fill the disk without consent.
17. As an operator, I want explicit cleanup with an exact plan, so that only the reviewed targets can disappear.
18. As a reader, I want exclusion between cleanup and active operations, so that cleanup cannot remove bytes during a read or write.
19. As an operator, I want safe interrupted preparation, so that a partial write cannot expose a usable artifact.
20. As an operator, I want typed storage refusals, so that repair does not silently replace evidence.
21. As an author, I want separate approval and supplement checks, so that integrity cannot grant permission to act.
22. As a retained consumer, I want verified reuse of available context, so that unchanged evidence need not be delivered twice.
23. As a fresh consumer, I want complete required context, so that another consumer receipt cannot replace my evidence.
24. As a review coordinator, I want an evidence-preserving native handoff, so that the next harness can resume the same review.
25. As a user, I want removal of charge --full, so that the obsolete oversized route cannot remain the next action.
26. As a user, I want unchanged ordinary preflight and write-spec behavior, so that this change stays within the approved charge scope.
27. As a maintainer, I want one source-policy owner, so that the manifest and guidance cannot silently disagree.
28. As a consumer, I want production transport and lifecycle proof, so that a disposable prototype cannot stand in for the final implementation.

## Implementation decisions

### Owners and trust

Preflight owns selectors, required-source policy, source checks, generated collectors, and current-action checks.
A new `internal/chargeevidence/` package owns the immutable pack, bounded reads, integrity verification, quota, and cleanup.
Its interface accepts prepared values and an already resolved repository-common location.
It does not select sources, execute collectors, or infer approval.

The canonical executable format registry lives in `internal/chargeevidence/schema.go`.
It owns header layout, versions, field names, types, order, and page policy.
The pack writer, strict reader, and response encoder consume that registry.
The shipped format reference is its generated projection.

Preflight retains the distinct source-selection and current-action policies.
The format registry does not select evidence or infer requiredness.
The existing TOON adapter remains the only text encoder.
Each implementation checkpoint registers only its implemented schema families.

The highest consumer seam remains `preflight.CommandWithVersion`.
The new storage seam hides publication, lock lifetime, safe file access, and physical offsets.
Two adapters justify that seam: preparation writes packs, while separate commands read and clean them.
Local filesystem operations form a local-substitutable dependency.

Tests drive its public operations rather than edit private offsets.

The preparation response or a verified handoff supplies the trusted expected identity.
The pack does not authenticate itself.
Neither the canonical manifest nor a stored executable launches a validator.
Existing Bench executable selection remains the execution trust root.

This feature adds no executable bootstrap chain.

The following terms apply throughout this feature.
The implementation adds them to the existing glossary without copying workflow instructions.

| Term | Meaning | Avoid |
| --- | --- | --- |
| Prepared evidence set | The required inventory, descriptors, metadata, and exact captures from one successful preparation. | Complete response, prompt bundle. |
| Evidence identity | SHA-256 of the canonical manifest bytes. | Assignment identity, approval. |
| Source identity | SHA-256 of one exact source body. | Role identity, permission. |
| Delivery coverage | Actual verified byte ranges available to one consumer. | Final cursor, comprehension. |
| Current-action binding | A verified current assignment and source pair associated with an evidence identity. | Artifact lifetime, origin location. |

### Public grammar

```text
bench preflight build <slug> --charge --ticket <basename> --base <b> --source-tip <t> [--max-store-bytes <n>]
bench preflight review <slug> --charge --base <b> --source-tip <t> [--max-store-bytes <n>]
bench preflight evidence <id> [--source <source-id>] [--cursor <cursor>]
bench preflight evidence <id> --verify
bench preflight evidence <id> --check-current
bench preflight evidence-clean [--cursor <cursor>] [--apply <fingerprint>]
```

Preflight owns one executable operation registry for argument grammar and public operation descriptions.
`cmd/bench/main.go` projects that inventory through the existing command registry.
Root help and preflight help therefore advertise the same implemented operations.
The registry grows with each producer checkpoint, without future-command advertisements or runtime feature flags.

The artifact identifier is `sha256:` followed by 64 lowercase hexadecimal digits.
The source selector is a manifest source ID, never a path.
The quota operand is a positive decimal byte count within the supported unsigned 64-bit range.
Mutually exclusive read modes cannot combine.

Cleanup apply cannot combine with a cursor.
Invalid grammar uses exit 2.
Operational refusal uses exit 1.
Successful preparation, a complete page, full verification, and cleanup use exit 0.

Charge `--full` has no supported form after CE-C2.
The serial migration preserves each existing legacy route only until its corresponding phase consumer migrates.
No final alias, small-response exception, or migration command remains.
The implementation phase's separate full-run control remains unchanged.

### Canonical manifest profile 1

The profile uses flat TOON tables in the exact order below.
Every block appears once, including zero-row blocks.
Every column appears in its declared order.
The profile uses two-space indentation, comma delimiters, and one final LF per block.

There are no blank separators.
The pinned existing encoder supplies string escaping and numeric-looking string quoting.
Integers use unsigned decimal values without leading zeroes.
Booleans use typed TOON booleans.

Unknown blocks, columns, versions, duplicate identifiers, or noncanonical bytes refuse.
Decode and canonical re-encode must reproduce the exact manifest bytes.

| Block | Ordered fields | Cardinality and order |
| --- | --- | --- |
| profile | version, hash, page_bytes | One row: integer 1, string sha256, integer 8192. |
| selection | mode, spec, ticket, base, source_tip | One all-string row. Review uses an empty ticket string. |
| sources | id, role, kind, path, required, bytes, sha256 | Source order. Required is boolean and bytes is integer. Other cells are strings. |
| pages | source, index, offset, bytes, sha256 | Source order, then increasing page index. Index, offset, and bytes are integers. |
| producers | source, name, version, cwd | Source order. All cells are strings. Canonical sources have no producer row. |
| arguments | source, index, value | Source order, then argument index. Index is integer. Other cells are strings. |

Source IDs are `s` followed by their one-based decimal ordinal.
They describe position, not authority.
Metadata is the first required source, with role `metadata`, kind `derived`, and empty path.
Canonical repository sources follow in their existing declared order.

Generated review sources follow in diff, consumers, coverage order.
The source-policy owner derives this inventory and its cardinality.

The metadata body is canonical TOON from the same typed prepared structure.
It contains complete fence, writes, coverage, checks, return requirements, and axis descriptors when applicable.
Review completion facts also reside there.
Metadata excludes the assignment identifier and absolute checkout path.

Each repeated list has one row per member instead of a comma-joined semantic string.
The following metadata blocks define the exact schema.
Every block appears once in this order, including zero-row blocks.
All cells are strings, and all source cells contain manifest source IDs.

| Metadata block | Ordered fields | Row order |
| --- | --- | --- |
| charge | axis, ticket, access | One build row with empty axis, or the existing review-axis inventory order. |
| fence | path | Declared spec fence order. |
| writes | path | Selected ticket write order. Review has zero rows. |
| coverage | row | Selected ticket coverage order. Review has zero rows here. |
| checks | source | Existing build check-source order or review skill source. |
| returns | source | Existing return-source order. |
| shared_evidence | kind, source | Diff, consumers, coverage order. Build has zero rows. |
| completion_evidence | record, source_digest, plan_digest, record_state, detail | Existing completion facts. Build has zero rows. |

The build charge ticket cell contains the selected ticket source ID.
Review uses the spec source ID in that cell and access `read-only`.
Build uses access `write-within-fence`.
The manifest source inventory replaces the redundant all-evidence handle string.

The registry generates these tables in the shipped format reference; no separately maintained schema inventory survives.

Canonical sources use the exact committed path and bytes.
Generated sources use their collector name as the path and role.
Canonical source kinds are `repository`, generated source kinds are `generated`, and the metadata kind is `derived`.
The existing named source fields define canonical roles, written with hyphens between words.

The selection uses resolved repository-relative spec and ticket paths, mode build or review, and resolved Git pins.
Each producer row records the collector name, Bench version, and repository-relative working directory.
Argument rows record the ordered semantic invocation with repository-relative paths and resolved source pins.
They exclude incidental absolute checkout locations.

This provenance records declared inputs, not a hermetic environment closure.
Changed provenance or captured bytes changes the evidence identity.

A source page contains at most 8,192 raw bytes.
The writer takes the longest UTF-8 prefix within that limit.
Pages cover each source exactly once without gaps or overlaps.
Zero-byte optional sources have no pages and retain their exact empty digest.

Required canonical sources remain nonempty under the existing source policy.
The profile forbids normalization of source bytes.
Source content must be valid UTF-8 and representable by the shared TOON adapter.
Every `sha256` cell contains exactly 64 lowercase hexadecimal digits as a string.

The artifact identifier adds the `sha256:` prefix outside the manifest digest cells.

The manifest binds its selectors, pins, source roles, paths, requiredness, body digests, and ordered page hashes.
Metadata bytes and producer arguments therefore participate in the same identity.
Assignment identifiers and checkout paths belong only to the current-action binding.
Identical logical evidence from sibling worktrees consequently has one identity.

An encoder upgrade that changes canonical bytes requires explicit profile compatibility handling.
It cannot silently rewrite an existing identity.

### Physical pack version 1

One regular file contains a 24-byte header, the canonical manifest, and uncompressed raw sources in manifest order.
The name derives only from the validated evidence identity.

| Header bytes | Meaning |
| --- | --- |
| 0 through 7 | The seven ASCII bytes `BENCHEV`, then one NUL byte. |
| 8 through 11 | Unsigned little-endian container version 1. |
| 12 through 15 | Reserved unsigned little-endian value 0. |
| 16 through 23 | Unsigned little-endian manifest byte length. |

The reader derives each physical source offset from the validated preceding lengths.
There is no second physical-offset registry.
Every addition and range comparison checks overflow before allocation or seek.
The expected total length must equal the opened file length.

The reader refuses unsupported versions, invalid ranges, overlap, gaps, truncation, and trailing data.
The format uses no archive extraction, database, or compression codec.
The shipped reference permits an independent reader without Bench-private helpers.

### Responses and independent verification

All new responses use the existing TOON adapter.
One shared final encoder guard covers success, usage, and operational refusal.
The bound includes every table, error, escaped cell, and next action.
Long diagnostic operands are identified by bounded type and digest, not echoed in full.

Refusal retains a bounded recovery instruction.
It never reports an incomplete encoded record as a complete response.

Preparation returns the artifact identity and a bounded current-action binding.
Long selectors and inventories use manifest references, with their required delivery state explicit.
The preparation response contains no source bodies.
Its exact first action is `bench preflight evidence <id>`.

| Response block | Ordered fields | Meaning |
| --- | --- | --- |
| prepared | evidence, mode, base, source_tip, assignment, selection, metadata, sources, pages, manifest_bytes, response_complete, delivery, next | One row. Selection is manifest:selection and metadata is s1. Delivery is unverified. Assignment uses its existing fixed format. |
| page | evidence, stream, source, index, offset, bytes, total, sha256, content, response_complete, stream_end, next | One fragment. Manifest uses source empty and stream manifest. Source pages use stream source. |
| verified | evidence, manifest_verified, pages_verified, sources_verified, delivery | Full integrity result. Delivery remains unverified. |
| current | evidence, assignment, base, source_tip, current, delivery | Current binding result. Delivery remains unverified. |
| cleanup | fingerprint, targets, bytes, response_complete, stream_end, next | Bounded plan orientation and exact successor. |
| targets | id, kind, bytes | Exact cleanup target rows, paged by the cleanup cursor. |
| applied | fingerprint, removed, remaining, complete, next | Terminal cleanup disposition. A partial failure uses exit 1. |

Identity, state, selector, content, and next cells are strings.
Counts, offsets, lengths, and page indices are integers.
Completeness, verification, stream-end, and current fields are booleans.
`pages_verified` and `sources_verified` are counts.

There is no `delivery_complete=true` response field.

Manifest transport splits its exact byte stream into at most 8,192-byte UTF-8 fragments.
This permits a single manifest scalar to exceed the response limit.
Manifest fragments carry transport hashes, but the consumer trusts them only after the complete manifest matches the expected identity.
Source pages use the manifest descriptors and hashes.

Each read validates the complete stored manifest and the selected source page before output.
Each page's response is complete even when its evidence stream continues.

Cursors use a versioned ASCII grammar: `v1.<hex-id>.<m|s>.<source-ordinal>.<page-index>`.
All numbers use canonical unsigned decimal notation.
The manifest cursor uses source ordinal zero.
The reader rejects foreign identities, mismatched source selectors, malformed numbers, and out-of-range pages.

No cursor authorizes a filesystem path.
No cursor stores server progress.
Repeated and out-of-order reads create or change no consumer-progress state or reading log.
Operation locks and existing ordinary telemetry may change; neither may store a reading session or consumer cursor.

Default traversal reads the manifest, required metadata, canonical sources, then generated sources.
Each nonterminal result supplies its exact successor command.
The last result has `stream_end=true` and an empty next value.
Explicit source traversal ends after that source and never claims the default traversal finished.

An arbitrary final cursor does not establish earlier delivery.

The independent verification fixture starts with a trusted expected identity.
It reconstructs the manifest without producer-private code, then checks the manifest hash and strict profile.
It verifies actual returned source ranges, membership, source digests, and exact byte coverage.
It rejects missing, duplicate, overlapping, and final-only result sets.

Its independent expectations require recorded omission mutations under the repository's duplication exception.
Full artifact verification reads every page and checks every complete source digest.
Neither interface certifies consumer comprehension, approval, or the task supplement.

### Publication, capacity, and current action

Storage lives beneath the verified Git common directory in `bench-charge-evidence`.
The existing Git common-directory owner resolves that location.
Only explicit preparation creates an absent store.
Reads refuse an absent artifact without changing the store.

An empty store is valid.
No-follow access rejects unsafe path components and nonregular pack objects.
Reader validation uses the opened file, with replacement checks around read operations.
It never substitutes current checkout bytes or regenerated review output.

Preparation gathers and validates the complete candidate inside the existing movement-checked attempt.
It writes a private temporary pack and verifies its complete manifest, pages, and source hashes.
Publication occurs only after the final movement check succeeds.
One movement discards the candidate and retries once.

A second movement discards the retry and refuses.
The publish step checks that its candidate still belongs to the successful attempt.
Atomic publication never overwrites an existing artifact.
An identical existing artifact is fully verified before reuse.

A corrupt existing artifact refuses rather than being repaired implicitly.

The default quota is 1,073,741,824 bytes over published packs and temporary pack bytes.
There is no smaller per-artifact limit.
Preparation can select a different positive quota through its explicit operand.
The quota applies to the entire store during that operation.

A cross-process writer lock serializes capacity calculation, temporary growth, verification, and publication.
The validated candidate lengths determine its exact pack size before temporary growth.
The writer counts published and existing temporary bytes while holding that lock.
It refuses before the candidate can exceed the selected total quota.

A crash releases the platform lock; any orphan temporary bytes still consume capacity.
This design needs no persistent reservation ledger.

Capacity refusal preserves prior artifacts and returns one deterministic bounded recovery instruction.
Before cleanup exists, its exact recovery instruction is `retry with --max-store-bytes <required-bytes>`.
Required bytes equal the observed store bytes plus the exact candidate size.
The caller repeats its original preparation with that option; the refusal need not echo long selectors.

An unrepresentable quota refuses with the bounded capacity classification instead of a wrapping integer.

Once cleanup exists, the first recovery action is exactly `bench preflight evidence-clean`.
The explicit larger-quota retry remains the documented alternative.
No operation automatically evicts evidence.

Readers and writers hold a shared repository evidence lock for their operation.
Writers additionally hold the writer lock through capacity calculation, temporary growth, verification, and publication.
Cleanup requires the exclusive evidence lock.
It refuses immediately if an active reader or writer prevents acquisition.

Locks cover the complete open, validation, read, or publication lifetime.
The implementation uses supported platform primitives and proves cross-process exclusion.

Check-current validates the manifest selectors against a clean, active current assignment and the current source pair.
It also retains required-source comparisons and the current preflight verdict checks.
It returns that current assignment binding, never the origin assignment as authority.
A released origin remains readable but its old binding cannot authorize action.

A different active assignment needs its own successful current binding and independent phase approval.
A moved source tip requires a new preparation, even when some source bodies remain unchanged.

### Cleanup and consumer action

Cleanup plans include every published artifact and every proven orphan temporary artifact in this namespace.
They do not inspect unrelated repository storage.
The plan sorts exact target identities and kinds deterministically.
Its SHA-256 fingerprint commits to the complete target list and their byte lengths and observed file identities.

A cleanup cursor binds that fingerprint and the next target position.
If the store changes, continuation refuses and asks for a fresh plan.

Apply acquires the exclusive lock and reconstructs the complete plan.
A fingerprint mismatch refuses before any deletion.
The target list must have been completely delivered before the operator authorizes apply.
The CLI does not invent a reading receipt for that operator.

Apply rechecks each target's no-follow file identity before deletion.
It stops on the first deletion failure and reports bounded counts plus a fresh-plan recovery action.
It does not claim transaction rollback after earlier authorized deletions.
An interrupted apply leaves a subset removed and requires a fresh plan.

An absent store produces an empty plan without creating directories.

The canonical build and review phases require verified delivery and available required context before action.
They also require a current-action binding, reviewer approval, and the complete task supplement.
Retained consumers can reuse exact available source bytes only after verifying new role, requiredness, and manifest membership.
Fresh consumers retrieve their own required context.

Neither a terminal cursor nor a transferred receipt replaces that context.
Review axes share one prepared evidence identity and preserve their independent source derivation.
A capable-harness handoff carries that trusted identity and its exact retrieval command.
Thin adapters continue to load the canonical phase instructions.

They gain no private source list or retrieval sequence.

## Implementation chunks

Each chunk receives independent Standards, Spec, and Coverage review before its successor starts.
Each ticket remains one serial green checkpoint.
The retained author owns implementation, tests, probes, and repairs.

Old CE-C1 is split into CE-C1A, CE-C1B, CE-C1C, and CE-C1D.
CE-C2, CE-C3, and CE-C4 retain their original outcome boundaries; their ticket numbers change.
These stable IDs replace the previous four-ticket completion plan before implementation starts.

| stable chunk ID / ticket | delivered outcome | row count | price and read limit |
| --- | --- | --- | --- |
| CE-C1A / 1-validate-legacy-prepared-packs.md | Preserve legacy build output through a validated in-memory pack and its generated format reference. | 41 | 8 focused production files; at most 1,600 source lines. Read charge.go, preparation.go, charge_test.go, source adapters, and TOON before codec changes. |
| CE-C1B / 2-publish-bounded-build-evidence.md | Prepare persistent build evidence and traverse its manifest and bodies through bounded responses. | 54 | 10 focused files; at most 2,000 source lines. Start from the reviewed pack API, preflight command, Git common-directory owner, and filesystem precedent. |
| CE-C1C / 3-verify-current-evidence.md | Verify all bytes, bind current action, and navigate sources independently across released assignments. | 23 | 8 focused files; at most 1,500 source lines. Read the reviewed store API, assignment checks, and the production system harness. |
| CE-C1D / 4-activate-bounded-build-guidance.md | Activate the complete bounded build action path with every existing prerequisite. | 6 | 6 owner files plus their exact fixtures; at most 1,600 source lines. Read the build phase, preparation anchors, glossary, and conformance owner. |
| CE-C2 / 5-freeze-bounded-review-evidence.md | Freeze bounded review evidence and remove the remaining charge full route. | 23 | 9 owner files plus their exact fixtures; at most 1,900 source lines. Read review collectors, review phase, dispatch anchors, and the shared schema API. |
| CE-C3 / 6-clean-evidence-explicitly.md | Deliver exact cleanup plans and safe apply behavior on supported native platforms. | 21 | 8 focused files; at most 1,600 source lines. Read the store locks, cleanup parser, public inventory, and native filesystem test owner. |
| CE-C4 / 7-preserve-consumer-context.md | Preserve retained reuse and native handoff, with actual supported-harness transport proof. | 7 | 7 owner files plus their exact fixtures; at most 1,500 source lines. Read both phase handoffs, reuse anchors, and the transport record contract. |

The row counts include table-driven hostile cases, not separate implementation passes.
Every ticket is priced for one retained author context, one focused red-green cycle, and its listed mutation set.
Stop at the named read limit if an unresolved owner requires broader design; update the plan before expanding scope.
Do not let the limit replace a required source read or hide an unresolved blocker.

CE-C1A is a safe prefactor with a real legacy consumer and an exact differential oracle.
CE-C1B delivers preparation and default traversal, with safe storage from its first persistent write.
CE-C1C adds source navigation, public verification, and current binding over that working path.
CE-C1D activates build guidance only after every required producer exists.

Legacy build --full survives through CE-C1C; legacy review --full survives through CE-C1D.
CE-C1D removes the build route, and CE-C2 removes the remaining review route.
The existing phase prerequisites apply throughout this migration.
No intermediate checkpoint claims the final feature or incomplete action path is ready.

Public commands use only subsets of the approved final grammar plus the existing legacy routes awaiting their named migration.
CE-C1B exposes no source selector, verify, current, or cleanup mode.
CE-C1C adds source selection, verify, and current; CE-C3 adds cleanup.
Required response schemas and inventory checks land with each producer.

## Testing decisions

Use the public preflight command for preparation, page traversal, current checks, and cleanup responses.
Use the public evidence-store interface for pack corruption and publication failures that a valid command cannot construct.
Use deterministic filesystem failure injection for write, verification, rename, and deletion failures.
Register any new preflight port with its real-producer test.

Test substitution stays in process and cannot claim to affect a launched executable.

Existing precedents are `TestChargeProjectionAndFullRetrieval` and `TestLegacyPreflightDifferential` in `internal/preflight/charge_test.go`.
The latter defines the preserved ordinary-preflight family and remains a differential exit row.
The former supplies the CE-C1A differential, then is replaced when its consumer migrates.
`TestPreparedBuildGuidanceAnchorsBiteIndependently` and `TestPreparedReviewGuidance` establish the current guidance mutation seam.

Their owner is `docs-currency-workflow` through the registered conformance runner.
New guidance mutations must execute that owner and prove restoration.
The new synthetic guidance test functions also run through the focused conformance package selector in the completion plan.

Cross-process behavior extends `internal/systemtest/` under its existing owner and disposable repository budget.
It uses the selected exact-source binary through `bench test --check system`.
It must drive production Bench assignment release, not raw Git worktree removal.
Ordinary decision tests gain no subprocesses or repositories.

Do not copy package-private test helpers across package seams.

The new named tests below are requirements, not claims of current executable coverage.
Rows sharing a test name use separately named table cases for each predicate.
Every independent expectation records its named omission or mutation and observed red before acceptance.
The format reader fixture independently reconstructs bytes without the production decoder or private writer constants.

Schema cases independently assert the exact registered field set, order, type, cardinality, and zero-row blocks.
Every declared metadata and response field receives a recorded omission mutation under the duplication exception.
Public inventory cases compare exact help with an independent omission oracle, following the existing help inventory precedent.

The production parser and help project one operation registry; test expectations remain independent only to detect named omissions.
The format projection check compares the shipped reference with registry-generated output.
Mutating either the shipped document or one generated field must produce a recorded red.

A missing descriptor, missing last source, changed requiredness, and trusted-self-root mutation must each turn it red.

| Planned seam | Test file owner |
| --- | --- |
| preflight | New command, source-refusal, and response-budget test files under `internal/preflight/`. |
| review | New `internal/preflight/evidence_review_test.go`. |
| clean | New `internal/preflight/evidence_cleanup_test.go`. |
| format | New `internal/chargeevidence/format_test.go`. |
| store | New `internal/chargeevidence/store_test.go`. |
| system | New `internal/systemtest/charge_evidence_test.go`, with the existing system owner. |
| guidance | New `internal/conformance/charge_evidence_guidance_test.go`, through registered workflow checks. |
| inventory | Existing `cmd/bench/help_inventory_test.go`, with producer-local cases in `TestEvidenceHelpInventory`. |

These planned names do not pretend that future test files exist during spec staging.
The implementation replaces each planned seam label with the exact executed test citation before completion.

The response-budget operation matrix is preparation, manifest read, source read, verify, check-current, cleanup plan, cleanup apply, usage, and operational refusal.
Each operation has a small case and a representable expansion case.
CE-C1B closes build, default-read, and refusal rows; CE-C1C closes verification and current-response rows.
CE-C2 closes the review preparation row when that producer lands.

CE-C3 closes cleanup plan and apply rows when those producers land.
CE157 records the precleanup retry contract at CE-C1B.
CE175 updates the shared recovery oracle at CE-C3, preserving the larger-quota option as the alternative.

The earlier row does not require an obsolete runtime branch after cleanup becomes available.
The final reconciliation checks the complete operation matrix through the shared guard.
The matrix includes a 40 KB spec, a larger-than-budget ticket, 5,000 metadata entries, and a large generated review diff.
It also includes an oversized manifest scalar, numeric-looking digests, quotes, tabs, newlines, backslashes, and multibyte Unicode.

Removing any operation from the final bound guard must make its independent budget case red.

The source-refusal matrix starts after valid grammar, a clean active assignment, a staged spec, and valid explicit pins.
It then changes one required source state at a time.
Where a special file cannot enter Git, the fixture drives the no-follow source adapter below the checkout guard.
CE-C1A drives the existing legacy preparation and requires its exact refusal before any successful charge output.
CE-C1B reruns those cases and additionally proves that no artifact becomes published.

This distinction prevents an unreachable source-kind assertion from passing at an earlier guard.

Native storage evidence enumerates Linux, macOS, and Windows from the current release target policy.
Use production writer/reader tests on each supported native filesystem path.
Record unavailable evidence as pending, never as a pass from cross-compilation.
The Linux shaping timings establish no production latency threshold.

Transport records separately exercise the Codex and Claude native tool-result paths.
Each record names the harness version, source tip, expected identity, page count, maximum stdout bytes, and reconstructed source hashes.
The consumer must read actual production command results and independently reconstruct the manifest and required bodies.
Keep approval and comprehension claims out of those records.

Other AGENTS harnesses use the documented textual protocol without a new harness-specific adapter.

### Seam diagram

```text
approved task + clean assigned source
    -> preflight preparation and source policy
    -> chargeevidence store -> immutable repository-common pack

expected identity + stateless cursor
    -> preflight evidence -> validated pack read -> bounded TOON response
    -> independent consumer verification -> available required context

current assignment + frozen selectors
    -> preflight check-current -> current-action binding
    -> phase approval and supplement checks -> authorized action

cleanup plan fingerprint
    -> exclusive evidence lock -> revalidated targets -> explicit deletion
```

### Acceptance coverage map

Fixture names below are planned table cases attached to the named test.
Each why clause names the cheapest wrong implementation or its exact mutation.
The story inventory is authoritative for the row-to-story mapping.

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| CE1 | 1 | A prepared build reconstructs the selected ticket and every required canonical source exactly. | planned preflight seam, `TestEvidenceBuildRoundTrip` | Omit the delegation procedure from the source policy. |
| CE2 | 1 | A ticket larger than 48,000 bytes remains fully retrievable. | planned preflight seam, `TestEvidenceLargeTicket` | Keep the ticket inline or omit its overflow. |
| CE3 | 2 | Preparation reports delivery as unverified. | planned preflight seam, `TestEvidencePreparationState` | Set delivery_complete from successful preparation. |
| CE4 | 2 | A page reports only its response completeness and stream state. | planned preflight seam, `TestEvidencePageState` | Set delivery_complete when the last page is requested. |
| CE5 | 3 | Default traversal returns the complete manifest before metadata or bodies. | planned preflight seam, `TestEvidenceTraversalOrder` | Start source delivery before the final manifest fragment. |
| CE6 | 3 | A manifest scalar larger than the response limit reconstructs exactly. | planned preflight seam, `TestEvidenceLargeManifestScalar` | Page only between manifest rows. |
| CE7 | 3 | The independent reader rejects a manifest whose hash differs from the trusted identity. | planned format seam, `TestEvidenceIndependentReader` | Trust the artifact self-reported identity. |
| CE8 | 3 | The independent reader rejects a source page outside verified manifest membership. | planned format seam, `TestEvidenceIndependentReader` | Trust a matching page digest without membership. |
| CE9 | 4 | A source without a final newline reconstructs without an added newline. | planned preflight seam, `TestEvidenceExactText` | Frame source content by lines. |
| CE10 | 4 | Unicode and escapable control bytes reconstruct exactly. | planned preflight seam, `TestEvidenceExactText` | Split UTF-8 bytes inside a code point or normalize escapes. |
| CE11 | 4 | Numeric-looking identifiers decode as strings. | planned preflight seam, `TestEvidenceTypedCells` | Pass a digest to the encoder as a number. |
| CE12 | 4 | Representable path text remains data throughout retrieval. | planned preflight seam, `TestEvidenceHostilePaths` | Interpret spaces, globs, quotes, or Unicode as shell syntax. |
| CE13 | 5 | Build preparation contains at most 48,000 encoded stdout bytes. | planned preflight seam, `TestEvidenceResponseBudget` | Bypass the final byte guard for build preparation. |
| CE14 | 5 | Oversized required metadata reconstructs without omission. | planned preflight seam, `TestEvidenceLargeMetadata` | Bound bodies but drop fence or coverage entries. |
| CE15 | 5 | An oversized invalid operand produces a bounded usage refusal. | planned preflight seam, `TestEvidenceBoundedErrors` | Echo the complete hostile operand. |
| CE16 | 6 | Identical complete logical inputs produce identical evidence identities. | planned format seam, `TestEvidenceDeterminism` | Include a timestamp in canonical bytes. |
| CE17 | 6 | Sibling assignment and checkout locations do not change logical identity. | planned system seam, `TestEvidenceLocationIndependence` | Hash the current assignment or absolute checkout path. |
| CE18 | 6 | A changed mode changes the evidence identity. | planned format seam, `TestEvidenceIdentityInputs` | Omit mode from the manifest commitment. |
| CE19 | 6 | A changed spec selector changes the evidence identity. | planned format seam, `TestEvidenceIdentityInputs` | Omit spec selector from the manifest commitment. |
| CE20 | 6 | A changed ticket selector changes the evidence identity. | planned format seam, `TestEvidenceIdentityInputs` | Omit ticket selector from the manifest commitment. |
| CE21 | 6 | A changed base changes the evidence identity. | planned format seam, `TestEvidenceIdentityInputs` | Omit base from the manifest commitment. |
| CE22 | 6 | A changed source tip changes the evidence identity. | planned format seam, `TestEvidenceIdentityInputs` | Omit source tip from the manifest commitment. |
| CE23 | 6 | A changed role changes the evidence identity. | planned format seam, `TestEvidenceIdentityInputs` | Omit role from the manifest commitment. |
| CE24 | 6 | A changed path changes the evidence identity. | planned format seam, `TestEvidenceIdentityInputs` | Omit path from the manifest commitment. |
| CE25 | 6 | A changed requiredness changes the evidence identity. | planned format seam, `TestEvidenceIdentityInputs` | Omit requiredness from the manifest commitment. |
| CE26 | 6 | Changed source bytes change the evidence identity. | planned format seam, `TestEvidenceIdentityInputs` | Omit source bytes from the manifest commitment. |
| CE27 | 6 | Changed required metadata changes the evidence identity. | planned format seam, `TestEvidenceIdentityInputs` | Omit required metadata from the manifest commitment. |
| CE28 | 7 | The writer emits the documented canonical TOON profile. | planned format seam, `TestEvidenceCanonicalProfile` | Alter block order, column order, types, escapes, indentation, or final newline. |
| CE29 | 7 | The reader refuses an unsupported manifest profile. | planned format seam, `TestEvidenceManifestRefusals` | Accept a profile version the decoder does not implement. |
| CE30 | 7 | The reader refuses noncanonical manifest bytes. | planned format seam, `TestEvidenceManifestRefusals` | Accept duplicate blocks, fields, IDs, or alternative serialization. |
| CE31 | 7 | The reader refuses unsupported container version. | planned format seam, `TestEvidencePackRefusals` | Accept the fixture with unsupported container version. |
| CE32 | 7 | The reader refuses invalid header marker. | planned format seam, `TestEvidencePackRefusals` | Accept the fixture with invalid header marker. |
| CE33 | 7 | The reader refuses nonzero reserved header bits. | planned format seam, `TestEvidencePackRefusals` | Accept the fixture with nonzero reserved header bits. |
| CE34 | 7 | The reader refuses overflowing length arithmetic. | planned format seam, `TestEvidencePackRefusals` | Accept the fixture with overflowing length arithmetic. |
| CE35 | 7 | The reader refuses overlapping page ranges. | planned format seam, `TestEvidencePackRefusals` | Accept the fixture with overlapping page ranges. |
| CE36 | 7 | The reader refuses gapped page ranges. | planned format seam, `TestEvidencePackRefusals` | Accept the fixture with gapped page ranges. |
| CE37 | 7 | The reader refuses truncated container. | planned format seam, `TestEvidencePackRefusals` | Accept the fixture with truncated container. |
| CE38 | 7 | The reader refuses trailing container data. | planned format seam, `TestEvidencePackRefusals` | Accept the fixture with trailing container data. |
| CE39 | 7 | A selected-page read refuses a changed page body. | planned store seam, `TestEvidencePageCorruption` | Return bytes before the page digest check. |
| CE40 | 7 | A selected-page read refuses a changed manifest. | planned store seam, `TestEvidenceManifestCorruption` | Read page offsets before validating the expected manifest digest. |
| CE41 | 8 | Preparation refuses a required source in the absent state. | planned preflight seam, `TestEvidenceRequiredSourceStates` | Return a successful prepared charge from the absent fixture. |
| CE42 | 8 | Preparation refuses a required source in the empty state. | planned preflight seam, `TestEvidenceRequiredSourceStates` | Return a successful prepared charge from the empty fixture. |
| CE43 | 8 | Preparation refuses a required source in the symlink state. | planned preflight seam, `TestEvidenceRequiredSourceStates` | Return a successful prepared charge from the symlink fixture. |
| CE44 | 8 | Preparation refuses a required source in the FIFO state. | planned preflight seam, `TestEvidenceRequiredSourceStates` | Return a successful prepared charge from the FIFO fixture. |
| CE45 | 8 | Preparation refuses a required source in the socket state. | planned preflight seam, `TestEvidenceRequiredSourceStates` | Return a successful prepared charge from the socket fixture. |
| CE46 | 8 | Preparation refuses a required source in the device state. | planned preflight seam, `TestEvidenceRequiredSourceStates` | Return a successful prepared charge from the device fixture. |
| CE47 | 8 | Preparation refuses a required source in the directory state. | planned preflight seam, `TestEvidenceRequiredSourceStates` | Return a successful prepared charge from the directory fixture. |
| CE48 | 8 | Preparation refuses unsupported source bytes. | planned preflight seam, `TestEvidenceRequiredSourceStates` | Replace unsupported bytes during encoding. |
| CE49 | 8 | Preparation refuses a dirty source checkout. | planned preflight seam, `TestEvidencePreparationRefusals` | Return a successful prepared charge after a dirty-source verdict. |
| CE50 | 8 | Preparation refuses a missing selected ticket. | planned preflight seam, `TestEvidencePreparationRefusals` | Prepare a generic spec charge without a ticket. |
| CE51 | 8 | Preparation refuses a source-tip mismatch. | planned preflight seam, `TestEvidencePreparationRefusals` | Accept the current tip without comparing the requested pin. |
| CE52 | 8 | Preparation refuses an inactive assignment. | planned preflight seam, `TestEvidencePreparationRefusals` | Treat any worktree as an active assignment. |
| CE53 | 9 | One movement discards the first attempt before the retry publishes. | planned preflight seam, `TestEvidenceMovementPublication` | Publish the first attempt before its movement check finishes. |
| CE54 | 9 | A second movement publishes no artifact from either attempt. | planned preflight seam, `TestEvidenceMovementPublication` | Publish the terminal failed attempt. |
| CE55 | 10 | A sibling worktree reads the same prepared bytes. | planned system seam, `TestEvidenceSiblingRead` | Locate storage beneath the originating checkout. |
| CE56 | 10 | Production assignment release preserves historical evidence reads. | planned system seam, `TestEvidenceReleaseRead` | Delete evidence when Bench releases the originating assignment. |
| CE57 | 11 | Check-current refuses a moved source tip. | planned preflight seam, `TestEvidenceCurrentBinding` | Treat artifact integrity as current source validity. |
| CE58 | 11 | Check-current refuses a released assignment binding. | planned preflight seam, `TestEvidenceCurrentBinding` | Reuse the prior preparation assignment after release. |
| CE59 | 11 | Check-current returns the current verified assignment binding for unchanged pins. | planned preflight seam, `TestEvidenceCurrentBinding` | Require an origin path or silently reuse the old assignment. |
| CE60 | 12 | Full verification refuses corruption in an unread page. | planned store seam, `TestEvidenceVerifyAll` | Verify only the last requested page. |
| CE61 | 12 | Full verification refuses a source digest mismatch. | planned store seam, `TestEvidenceVerifyAll` | Verify page hashes without verifying the reconstructed source. |
| CE62 | 13 | An explicit source selector returns only that declared source stream. | planned preflight seam, `TestEvidenceSourceNavigation` | Treat the source identifier as a filesystem path. |
| CE63 | 13 | The reader refuses malformed or foreign cursors. | planned preflight seam, `TestEvidenceCursorRefusals` | Accept a cursor from another artifact or source. |
| CE64 | 13 | Every nonterminal read supplies the exact successor command. | planned preflight seam, `TestEvidenceNextActions` | Emit a command family or require guessed operands. |
| CE65 | 13 | Repeated and out-of-order reads create or change no consumer-progress state. | planned system seam, `TestEvidenceStatelessProcesses` | Persist a consumer cursor or reading log across otherwise successful separate processes; compare progress state before and after reads. |
| CE66 | 14 | A review preparation captures each collector once per attempt. | planned review seam, `TestEvidenceReviewCollectors` | Collect separately for each review axis. |
| CE67 | 14 | Review retrieval never reruns a collector. | planned review seam, `TestEvidenceReviewCollectors` | Regenerate diff or consumer output during reads. |
| CE68 | 14 | Each new review preparation reruns its collectors. | planned review seam, `TestEvidenceReviewCollectors` | Reuse cached collector results from a prior preparation. |
| CE69 | 14 | Every review axis resolves the same generated source identities. | planned review seam, `TestEvidenceReviewAxes` | Prepare different captured bytes for one axis. |
| CE70 | 14 | Oversized review captures reconstruct exactly. | planned review seam, `TestEvidenceLargeReview` | Keep generated output inline or truncate it. |
| CE71 | 14 | Incomplete consumer output refuses review preparation. | planned review seam, `TestEvidenceReviewRefusals` | Publish when the collector declares truncated output. |
| CE72 | 14 | Changed completion metadata changes the review evidence identity. | planned review seam, `TestEvidenceReviewMetadata` | Leave completion evidence outside the manifest commitment. |
| CE73 | 15 | Changed collector provenance changes the evidence identity. | planned review seam, `TestEvidenceReviewProvenance` | Hash captures but omit producer version or arguments. |
| CE74 | 16 | The default quota includes published artifacts and temporary writes. | planned store seam, `TestEvidenceQuota` | Count only completed packfiles. |
| CE75 | 16 | An explicit larger quota admits an artifact above the default quota. | planned store seam, `TestEvidenceQuotaOverride` | Apply a hidden per-artifact cap. |
| CE76 | 16 | Capacity refusal preserves every existing artifact. | planned store seam, `TestEvidenceQuota` | Evict old evidence to make room. |
| CE77 | 16 | Concurrent preparations cannot exceed the selected store quota. | planned system seam, `TestEvidenceQuotaConcurrency` | Release the writer lock between capacity calculation and temporary growth; force both writers to observe the same free bytes. |
| CE78 | 17 | Cleanup returns a bounded plan of exact deletion targets. | planned clean seam, `TestEvidenceCleanupPlan` | Return a summary count without retrievable target identities. |
| CE79 | 17 | Cleanup apply refuses a changed plan fingerprint before deletion. | planned clean seam, `TestEvidenceCleanupStale` | Delete the old target set after the store changes. |
| CE80 | 17 | Cleanup without apply deletes nothing. | planned clean seam, `TestEvidenceCleanupPlan` | Delete while enumerating the plan. |
| CE81 | 17 | An idle handle removed by explicit cleanup subsequently refuses retrieval. | planned clean seam, `TestEvidenceCleanupApply` | Substitute current checkout bytes for the deleted artifact. |
| CE82 | 17 | A cleanup deletion failure returns the exact unfinished disposition. | planned clean seam, `TestEvidenceCleanupFailure` | Report complete after a target deletion fails. |
| CE83 | 18 | Cleanup refuses while a reader holds the operation lock. | planned system seam, `TestEvidenceCleanupReaderExclusion` | Delete a pack while another process reads it. |
| CE84 | 18 | Cleanup refuses while a writer holds the operation lock. | planned system seam, `TestEvidenceCleanupWriterExclusion` | Delete a temporary artifact during publication. |
| CE85 | 18 | Cleanup includes an orphan temporary file only after writer exclusion succeeds. | planned clean seam, `TestEvidenceCleanupOrphans` | Infer writer death from temporary-file age. |
| CE86 | 19 | Interruption before source verification exposes no published handle. | planned system seam, `TestEvidenceInterruptedPublication` | Rename the temporary pack before final source verification. |
| CE87 | 19 | Interruption during artifact verification exposes no published handle. | planned system seam, `TestEvidenceInterruptedPublication` | Expose the pack before its complete integrity check. |
| CE88 | 19 | Publication failure at the terminal step exposes no partial artifact. | planned store seam, `TestEvidencePublicationFailure` | Return a successful handle after a failed atomic publication. |
| CE89 | 19 | Identical concurrent writers reuse one fully verified published artifact. | planned system seam, `TestEvidenceConcurrentPublication` | Replace an existing artifact without verifying it. |
| CE90 | 19 | Distinct concurrent writers retain their distinct artifacts. | planned system seam, `TestEvidenceConcurrentPublication` | Use one temporary publication name for different sets. |
| CE91 | 20 | The artifact reader refuses a symlink object without consuming it. | planned store seam, `TestEvidenceStoreKinds` | Open the symlink fixture as a regular pack. |
| CE92 | 20 | The artifact reader refuses a FIFO object without consuming it. | planned store seam, `TestEvidenceStoreKinds` | Open the FIFO fixture as a regular pack. |
| CE93 | 20 | The artifact reader refuses a socket object without consuming it. | planned store seam, `TestEvidenceStoreKinds` | Open the socket fixture as a regular pack. |
| CE94 | 20 | The artifact reader refuses a device object without consuming it. | planned store seam, `TestEvidenceStoreKinds` | Open the device fixture as a regular pack. |
| CE95 | 20 | The artifact reader refuses a directory object without consuming it. | planned store seam, `TestEvidenceStoreKinds` | Open the directory fixture as a regular pack. |
| CE96 | 20 | An invalid artifact identifier refuses before path resolution. | planned store seam, `TestEvidenceStorePaths` | Allow traversal or an absolute artifact path. |
| CE97 | 20 | A corrupt existing artifact refuses replacement during preparation. | planned store seam, `TestEvidenceExistingCorruption` | Overwrite corruption with the new candidate. |
| CE98 | 20 | An unsafe store directory refuses preparation and retrieval. | planned store seam, `TestEvidenceStoreDirectoryStates` | Follow a replaced store directory. |
| CE99 | 20 | An absent store creates only during explicit preparation. | planned store seam, `TestEvidenceStoreDirectoryStates` | Create storage during an ordinary read. |
| CE100 | 20 | An empty existing store reports an empty cleanup plan. | planned clean seam, `TestEvidenceCleanupEmpty` | Treat an empty store as a corrupt artifact. |
| CE101 | 21 | Build guidance requires approval independently of verified evidence. | planned guidance seam, `TestEvidenceBuildGuidance` | Remove the approval prerequisite. |
| CE102 | 21 | Build guidance requires the complete task supplement independently of evidence. | planned guidance seam, `TestEvidenceBuildGuidance` | Let verification supply the supplement. |
| CE103 | 21 | Build guidance requires available required context before action. | planned guidance seam, `TestEvidenceBuildGuidance` | Replace available evidence with a completion claim. |
| CE104 | 22 | Reuse requires verified role and requiredness in the new manifest. | planned guidance seam, `TestEvidenceConsumerGuidance` | Reuse by body digest alone. |
| CE105 | 23 | A fresh consumer retrieves its required evidence despite another consumer receipt. | planned guidance seam, `TestEvidenceConsumerGuidance` | Transfer only the terminal cursor. |
| CE106 | 23 | The independent delivery check rejects missing, duplicate, and final-only page sets. | planned format seam, `TestEvidenceDeliveryCoverage` | Accept final stream state without exact page coverage. |
| CE107 | 24 | A native handoff preserves the expected identity and exact retrieval action. | planned guidance seam, `TestEvidenceHandoffGuidance` | Carry only the originating checkout path. |
| CE108 | 25 | Charge --full receives exit 2 through the bounded usage path. | planned preflight seam, `TestEvidenceRemovedFull` | Retain a small-response alias. |
| CE109 | 25 | Executable next actions contain no charge --full route. | planned preflight seam, `TestEvidenceNextActions` | Keep an obsolete retrieval action. |
| CE110 | 25 | Canonical phase guidance contains no charge --full route. | planned guidance seam, `TestEvidenceRetiredRoute` | Restore the retired flag in either phase. |
| CE111 | 26 | Ordinary preflight preserves its enumerated baseline behavior. | existing `internal/preflight/charge_test.go` (`TestLegacyPreflightDifferential`) | Change a baseline exit, check, or output field. |
| CE112 | 26 | The implementation phase full-run control remains present. | planned guidance seam, `TestEvidenceUnchangedRoutes` | Remove the phase-level --full control. |
| CE113 | 26 | Write-spec retains its decision-source and author-fork contract. | planned guidance seam, `TestEvidenceUnchangedRoutes` | Add charge retrieval to the write-spec phase. |
| CE114 | 27 | Required-source inventories derive from one executable policy. | planned preflight seam, `TestEvidenceSourcePolicy` | Keep a second required-source inventory in the legacy renderer; omit a policy descriptor and observe a differing prepared inventory. |
| CE115 | 28 | Codex transport reconstructs the required production evidence without truncation. | review-owned: `reviews/bounded-charge-evidence.md`, CE-C4 Codex transport record | Accept only a local encoder round-trip. |
| CE116 | 28 | Claude transport reconstructs the required production evidence without truncation. | review-owned: `reviews/bounded-charge-evidence.md`, CE-C4 Claude transport record | Infer Claude delivery from the Codex result. |
| CE117 | 28 | Production storage checks pass on each supported native platform. | review-owned: `reviews/bounded-charge-evidence.md`, CE-C3 native storage record | Treat Linux prototype results as cross-platform proof. |
| CE118 | 20 | An absent store read refuses without creating the store. | planned store seam, `TestEvidenceStoreDirectoryStates` | Create a lock directory during a missing-artifact read. |
| CE119 | 17 | An absent store returns an empty cleanup plan without creation. | planned clean seam, `TestEvidenceCleanupEmpty` | Create the namespace while planning cleanup. |
| CE120 | 13 | A nested working directory retrieves the same artifact and source bytes. | planned system seam, `TestEvidenceNestedDirectory` | Resolve evidence relative to the current directory. |
| CE121 | 17 | Cleanup refuses an unsafe target before deleting any target. | planned clean seam, `TestEvidenceCleanupUnsafe` | Follow a symlink encountered in the target inventory. |
| CE122 | 17 | Interrupted cleanup requires a fresh plan for the remaining targets. | planned system seam, `TestEvidenceCleanupInterrupt` | Reuse an old fingerprint after a partial deletion. |
| CE123 | 11 | Check-current refuses a dirty assignment checkout. | planned preflight seam, `TestEvidenceCurrentBinding` | Skip the clean-source predicate for an existing artifact. |
| CE124 | 11 | Check-current refuses changed required-source bytes. | planned preflight seam, `TestEvidenceCurrentBinding` | Compare only the Git pair without the current required inputs. |
| CE125 | 7 | The pack reader refuses a page range beyond its source length. | planned format seam, `TestEvidencePackRefusals` | Accept an in-file range that crosses into the next source. |
| CE126 | 7 | The manifest reader refuses a changed scalar type. | planned format seam, `TestEvidenceManifestRefusals` | Coerce a numeric digest or quoted requiredness. |
| CE127 | 14 | A failed collector refuses preparation without publishing a handle. | planned review seam, `TestEvidenceReviewRefusals` | Publish an earlier collector result after a later collector failure. |
| CE128 | 8 | A missing required tool preserves a bounded preparation refusal. | planned preflight seam, `TestEvidencePreparationRefusals` | Publish a handle despite source inspection failure. |
| CE129 | 20 | A replaced open artifact refuses a changed selected page. | planned store seam, `TestEvidenceReplacementRace` | Validate a path and then read unchecked replacement bytes. |
| CE130 | 16 | Invalid quota operands refuse before temporary storage changes. | planned preflight seam, `TestEvidenceQuotaOperands` | Accept zero, negative, malformed, or overflowing quota input. |
| CE131 | 5 | The manifest read response contains at most 48,000 encoded stdout bytes. | planned preflight seam, `TestEvidenceResponseBudget` | Bypass the final byte guard for manifest read. |
| CE132 | 5 | The source read response contains at most 48,000 encoded stdout bytes. | planned preflight seam, `TestEvidenceResponseBudget` | Bypass the final byte guard for source read. |
| CE133 | 5 | The verify response contains at most 48,000 encoded stdout bytes. | planned preflight seam, `TestEvidenceResponseBudget` | Bypass the final byte guard for verify. |
| CE134 | 5 | The check-current response contains at most 48,000 encoded stdout bytes. | planned preflight seam, `TestEvidenceResponseBudget` | Bypass the final byte guard for check-current. |
| CE135 | 5 | The review preparation response contains at most 48,000 encoded stdout bytes. | planned review seam, `TestEvidenceResponseBudget` | Bypass the final byte guard for review preparation. |
| CE136 | 5 | The cleanup plan response contains at most 48,000 encoded stdout bytes. | planned clean seam, `TestEvidenceResponseBudget` | Bypass the final byte guard for cleanup plan. |
| CE137 | 5 | The cleanup apply response contains at most 48,000 encoded stdout bytes. | planned clean seam, `TestEvidenceResponseBudget` | Bypass the final byte guard for cleanup apply. |
| CE138 | 5 | The usage refusal response contains at most 48,000 encoded stdout bytes. | planned preflight seam, `TestEvidenceResponseBudget` | Bypass the final byte guard for usage refusal. |
| CE139 | 5 | The operational refusal response contains at most 48,000 encoded stdout bytes. | planned preflight seam, `TestEvidenceResponseBudget` | Bypass the final byte guard for operational refusal. |
| CE140 | 21 | Build guidance requires verified delivery before action. | planned guidance seam, `TestEvidenceBuildGuidance` | Treat orientation as sufficient evidence. |
| CE141 | 21 | Build guidance requires a current-action binding before action. | planned guidance seam, `TestEvidenceBuildGuidance` | Permit action after the source moves. |
| CE142 | 21 | Review guidance requires verified delivery before action. | planned guidance seam, `TestEvidenceReviewGuidance` | Omit verified delivery from the migrated review instructions. |
| CE143 | 21 | Review guidance requires available required context before action. | planned guidance seam, `TestEvidenceReviewGuidance` | Omit available required context from the migrated review instructions. |
| CE144 | 21 | Review guidance requires a current-action binding before action. | planned guidance seam, `TestEvidenceReviewGuidance` | Omit a current-action binding from the migrated review instructions. |
| CE145 | 21 | Review guidance requires reviewer approval before action. | planned guidance seam, `TestEvidenceReviewGuidance` | Omit reviewer approval from the migrated review instructions. |
| CE146 | 21 | Review guidance requires the complete task supplement before action. | planned guidance seam, `TestEvidenceReviewGuidance` | Omit the complete task supplement from the migrated review instructions. |
| CE147 | 26 | The legacy build charge preserves its enumerated output through the validated in-memory pack. | planned preflight seam, `TestLegacyPreparedPackDifferential` | Compare pre-refactor exact bytes and exits; omit one prepared source from the legacy projection and require differential red. |
| CE148 | 7 | The shipped format reference equals the canonical format registry projection. | planned format seam, `TestEvidenceFormatProjection` | Edit a shipped field label or omit one generated field; the independent projection comparison must turn red. |
| CE149 | 5 | Build metadata conforms to the exact registered metadata schema. | planned preflight seam, `TestEvidenceBuildMetadataSchema` | Omit each declared block or field, change its type, reorder columns, or comma-join a list in separate recorded mutation cases. |
| CE150 | 1 | Build preparation accepts exactly its declared argument forms. | planned preflight seam, `TestEvidenceBuildGrammar` | Remove a required pin or ticket guard; permit a read-only flag, duplicate flag, extra operand, or invalid quota in separate cases. |
| CE151 | 13 | Default evidence reads accept exactly their declared argument forms. | planned preflight seam, `TestEvidenceReadGrammar` | Accept an unknown flag, extra operand, missing cursor value, or duplicate cursor; independently enumerate permitted default-read forms. |
| CE152 | 27 | The public help inventory projects exactly the implemented preparation and read grammar. | planned inventory seam, `TestEvidenceHelpInventory` | Keep main.go limited to review and build, omit evidence, or advertise a later operation; independent expected inventory makes each mutation red. |
| CE153 | 2 | Build preparation returns the complete registered prepared response schema. | planned preflight seam, `TestEvidencePreparedSchema` | Omit manifest_bytes or next, change a field type, or reorder fields; mutate every declared field independently. |
| CE154 | 3 | Manifest fragments return the complete registered page response schema. | planned preflight seam, `TestEvidenceManifestPageSchema` | Omit total or next, change a field type, or reorder fields; mutate every declared field independently. |
| CE155 | 4 | Source fragments return the complete registered page response schema. | planned preflight seam, `TestEvidenceSourcePageSchema` | Omit source or next, change a field type, or reorder fields; mutate every declared field independently. |
| CE156 | 1 | Build preparation returns the exact manifest-first retrieval command. | planned preflight seam, `TestEvidencePreparationNext` | Return a source page, verify command, generic hint, or missing next instead of bench preflight evidence plus the exact prepared identity. |
| CE157 | 16 | Before cleanup lands, capacity refusal returns the exact larger-quota retry instruction. | planned preflight seam, `TestEvidenceCapacityRecovery` | Replace the CE-C1B retry instruction with a generic capacity error; CE175 owns the later shared-oracle migration. |
| CE158 | 12 | Full verification accepts only its exclusive declared argument form. | planned preflight seam, `TestEvidenceVerifyGrammar` | Accept verify with cursor, source, check-current, duplicate flags, or unrelated flags; each separate grammar case must refuse before store access. |
| CE159 | 11 | Current checking accepts only its exclusive declared argument form. | planned preflight seam, `TestEvidenceCurrentGrammar` | Accept check-current with cursor, source, verify, duplicate flags, or unrelated flags; each separate grammar case must refuse before store access. |
| CE160 | 12 | Full verification returns the complete registered verified response schema. | planned preflight seam, `TestEvidenceVerifiedSchema` | Omit a verification count or delivery, change a field type, or reorder fields; mutate every declared field independently. |
| CE161 | 11 | Current checking returns the complete registered current response schema. | planned preflight seam, `TestEvidenceCurrentSchema` | Omit assignment or delivery, change a field type, or reorder fields; mutate every declared field independently. |
| CE162 | 27 | The public help inventory projects the implemented source, verify, and current grammar. | planned inventory seam, `TestEvidenceHelpInventory` | Omit a new read mode or advertise cleanup early; compare the exact inventory with an independently maintained omission oracle. |
| CE163 | 14 | Review preparation accepts exactly its declared argument forms. | planned review seam, `TestEvidenceReviewGrammar` | Accept a ticket selector, absent pin, read-only flag, duplicate flag, or extra operand; each independent grammar case refuses before collectors. |
| CE164 | 14 | Review metadata conforms to the exact registered metadata schema. | planned review seam, `TestEvidenceReviewMetadataSchema` | Omit each populated review block or field, change its type, reorder columns, or comma-join a repeated list in separate mutation cases. |
| CE165 | 14 | Review preparation returns the complete registered prepared response schema. | planned review seam, `TestEvidenceReviewPreparedSchema` | Omit manifest_bytes or next, change a field type, or reorder fields; mutate every declared field independently. |
| CE166 | 14 | Review preparation returns the exact manifest-first retrieval command. | planned review seam, `TestEvidenceReviewPreparationNext` | Return an axis source, verify command, generic hint, or missing next instead of bench preflight evidence plus the exact prepared identity. |
| CE167 | 25 | The public help inventory projects review preparation without charge full. | planned inventory seam, `TestEvidenceHelpInventory` | Retain charge --full or omit the review artifact route; the independent exact inventory expectation must turn red. |
| CE168 | 17 | Cleanup accepts exactly its declared argument forms. | planned clean seam, `TestEvidenceCleanupGrammar` | Accept apply with cursor, missing fingerprint, duplicate flag, extra operand, or evidence-only flag; each case refuses before deletion. |
| CE169 | 17 | Cleanup plans return the complete registered cleanup response schema. | planned clean seam, `TestEvidenceCleanupSchema` | Omit fingerprint or next, change a field type, or reorder fields; mutate every declared field independently. |
| CE170 | 17 | Cleanup target pages return the complete registered targets response schema. | planned clean seam, `TestEvidenceCleanupTargetsSchema` | Omit id, kind, or bytes, change a field type, or reorder fields; mutate every declared field independently. |
| CE171 | 17 | Cleanup apply returns the complete registered applied response schema. | planned clean seam, `TestEvidenceCleanupAppliedSchema` | Omit remaining or next, change a field type, or reorder fields; mutate every declared field independently. |
| CE172 | 27 | The public help inventory projects exactly the implemented cleanup grammar. | planned inventory seam, `TestEvidenceHelpInventory` | Omit evidence-clean, hide apply, or invent an alias; compare the exact inventory with the independent omission oracle. |
| CE173 | 25 | Build charge full refuses after the build guidance migration. | planned preflight seam, `TestEvidenceRemovedBuildFull` | Keep the old build --full dispatch branch after CE-C1D; the exact build grammar must return bounded exit 2. |
| CE174 | 13 | Source-selected reads accept exactly their declared argument forms. | planned preflight seam, `TestEvidenceSourceGrammar` | Accept duplicate or missing-value source flags, extra operands, mismatched cursor pairing, or verify/current conflicts; separately prove permitted source and cursor forms. |
| CE175 | 16 | After cleanup lands, capacity refusal returns `bench preflight evidence-clean` first, with the explicit larger-quota alternative. | planned preflight seam, `TestEvidenceCapacityRecovery` | Keep the old quota-first expectation, emit a generic failure, or omit the alternative; the migrated shared oracle requires the complete recovery contract. |

### Edge inventory

The audience is every repository that links Bench, including this kit repository.
The following checklist attaches the profile's shell and Markdown input classes to this feature.

| Edge class | Concrete disposition |
| --- | --- |
| Spaces, globs, quotes, and Unicode | Preserve representable source paths as data through the shared encoder. |
| Missing final newline | Preserve exact bytes and their digest. |
| Absent required source | Refuse preparation before publication. |
| Empty required source | Refuse independently of absence. |
| Absent tickets directory | Preserve ordinary preflight behavior, but refuse a task-specific charge without a selected ticket. |
| Empty tickets directory | Preserve its distinct ordinary preflight result and refuse charge selection. |
| Absent evidence store | Reads refuse without creation, preparation creates, and cleanup returns an empty plan. |
| Empty evidence store | Reads refuse a missing artifact and cleanup returns an empty plan. |
| Unsafe files or store ancestors | Refuse without following links or consuming special files. |
| Missing tool | Preserve the existing source or collector refusal and bounded recovery output. |
| Nested working directory | Resolve the same repository-common store and semantic collector arguments. |
| Interrupted preparation | No handle becomes published before source and artifact verification finish. |
| Interrupted cleanup | Report incomplete deletion and require a newly fingerprinted plan. |
| Repeated preparation | Verify and reuse an identical artifact without overwriting it. |
| Repeated cleanup apply | Refuse the now-stale fingerprint if the target inventory changed. |
| Concurrent quota growth | Hold the writer lock from capacity calculation through publication. |
| Read or write during cleanup | Exclusive-lock refusal precedes any deletion. |
| Malformed cursor or identifier | Refuse before filesystem path resolution. |
| Unknown profile or pack version | Refuse without a permissive fallback. |
| Very large descriptor or error | Fragment required content and bound diagnostic operands. |
| Consumer stops early | Delivery remains incomplete despite a valid artifact or final cursor. |
| Body unchanged, role changed | Verify new membership before reuse. |

No acceptance fixture swaps a package variable across a subprocess.
The existing review observer remains an in-process collector-count probe.
The system owner supplies explicit process coordination for lifecycle cases.

Won't handle: a hermetic collector environment — the review caller receives exact captures and declared provenance instead.
Won't handle: arbitrary editor blocking — the build caller follows the honest-agent phase prerequisites.
Won't handle: write-spec evidence transport — the spec author retains the approved author-fork route.
Won't handle: general CLI pagination — existing non-charge callers retain their current output contracts.

Won't handle: compression or cross-artifact body deduplication — the evidence reader consumes exact uncompressed pack bytes.
Won't handle: kernel or privileged hostile modification that bypasses filesystem locks — digest checks still detect changed returned evidence.

## Ownership fences

Reviewer disposition: the user pre-approved the shaped scope and authorized this spec-and-ticket landing.
The following implementation fence is proposed within that scope and receives the independent Sol/high review.
It grants no implementation authority until the implementation phase starts.

- `internal/preflight`
- `internal/chargeevidence`
- `.agents/skills/bench-craft-delegate/references/charge-evidence-format.md`
- `reviews/bounded-charge-evidence.md`
- `internal/toon/toon_test.go`
- `internal/conformance/data_handling_test.go`
- `internal/systemtest`
- `internal/conformance/injected_ports_registry_test.go`
- `internal/git/worktree_admin.go`
- `cmd/bench/main.go`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/conformance/charge_evidence_guidance_test.go`
- `internal/anchors/registry_ft311_preparation.go`
- `internal/conformance/ft311_preparation_test.go`
- `.agents/commands/bench-implement-spec.md`
- `tests/canary/workflow-guidance-anchors/prepared-build-approval`
- `CONTEXT.md`
- `tests/canary/workflow-guidance-anchors/delegated-dispatch-declaration`
- `tests/canary/workflow-guidance-anchors/delegated-entry-refusals`
- `tests/canary/workflow-guidance-anchors/delegated-resumption-contents`
- `tests/canary/workflow-guidance-anchors/implement-spec-coverage-task-seeding`
- `tests/canary/workflow-guidance-anchors/implement-spec-cross-harness-pointer`
- `tests/canary/workflow-guidance-anchors/implement-spec-entry-validation`
- `tests/canary/workflow-guidance-anchors/implement-spec-incapable-harness`
- `tests/canary/workflow-guidance-anchors/implement-spec-inline-exception`
- `tests/canary/workflow-guidance-anchors/implement-spec-mandatory-delegation-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-read-only-helper`
- `tests/canary/workflow-guidance-anchors/implement-spec-status-flip-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-worktree-before-preflight`
- `tests/canary/workflow-guidance-anchors/implement-spec-write-delegation`
- `tests/canary/workflow-guidance-anchors/line-anchor-missing`
- `tests/canary/workflow-guidance-anchors/prepared-build-freshness`
- `tests/canary/workflow-guidance-anchors/dg-25`
- `tests/canary/workflow-guidance-anchors/dg-26`
- `tests/canary/workflow-guidance-anchors/dg-29`
- `tests/canary/workflow-guidance-anchors/dg-29-verification-target`
- `tests/canary/workflow-guidance-anchors/dg-30`
- `tests/canary/workflow-guidance-anchors/dg-31`
- `tests/canary/workflow-guidance-anchors/dg-31-contradiction-trigger`
- `tests/canary/workflow-guidance-anchors/implement-spec-adoption-freshness`
- `tests/canary/workflow-guidance-anchors/implement-spec-prose-owner-transfer`
- `tests/canary/docs-currency-token-diet/signal-vocabulary-drift`
- `tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-coverage-map-term`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-parts`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-decision-map-term`
- `tests/canary/workflow-guidance-anchors/context-reader-sweep-term`
- `tests/canary/workflow-guidance-anchors/context-ticket-vocabulary`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `.agents/commands/bench-review-implementation.md`
- `tests/canary/workflow-guidance-anchors/prepared-review-shared-evidence`
- `tests/canary/workflow-guidance-anchors/coverage-axis-anchor`
- `tests/canary/workflow-guidance-anchors/delegated-axis-exclusions`
- `tests/canary/workflow-guidance-anchors/delegated-chunk-tip-review`
- `tests/canary/workflow-guidance-anchors/prepared-review-axis-returns`
- `tests/canary/workflow-guidance-anchors/prepared-review-blast-evidence`
- `tests/canary/workflow-guidance-anchors/prepared-review-capable-handoff`
- `tests/canary/workflow-guidance-anchors/prepared-review-inline-axis-route`
- `tests/canary/workflow-guidance-anchors/prepared-review-legacy-entry-points`
- `tests/canary/workflow-guidance-anchors/prepared-review-native-dispatch`
- `tests/canary/workflow-guidance-anchors/prepared-review-runtime-capability`
- `tests/canary/workflow-guidance-anchors/review-base-merged-main-tip`
- `tests/canary/workflow-guidance-anchors/review-clean-terminal-result`
- `tests/canary/workflow-guidance-anchors/review-cross-harness-opt-in`
- `tests/canary/workflow-guidance-anchors/review-falsification-accept-routing`
- `tests/canary/workflow-guidance-anchors/review-falsification-dispositions`
- `tests/canary/workflow-guidance-anchors/review-persistence-anchor`
- `tests/canary/workflow-guidance-anchors/review-preflight-explicit-base`
- `tests/canary/workflow-guidance-anchors/review-repair-ticket-covers`
- `tests/canary/workflow-guidance-anchors/review-repair-ticket-owner`
- `tests/canary/workflow-guidance-anchors/review-standing-falsification`
- `tests/canary/workflow-guidance-anchors/review-universal-claim-bar`

Build-time rewrites exclude every `specs/*/spec.md` and every existing implementation ticket.
An in-scope plan expansion follows the operating guide before its writes occur.
The last ticket touching a package reconciles that package's full acceptance inventory.
No writer edits the separate debug-loop-guidance assignments.

## Out of scope

The following separate capabilities retain their own future decision and spec work.

| Capability | Derived estimate | Surviving caller |
| --- | --- | --- |
| Write-spec delivery protocol | 6 edits, 2 gate runs: phase, fork adapter, evidence owner, tests, anchors, fixtures. | Existing approved-source author fork. |
| Hermetic review collectors | 5 edits, 2 gate runs: loader, environment contract, capture producer, tests, guidance. | Exact prepared review capture. |
| Universal edit blocker | 6 edits, 3 gate runs: two harness adapters, shared authorization owner, hooks, tests, guidance. | Canonical phase action checks. |
| Compression and shared body storage | 4 edits, 2 gate runs: container writer, reader, format profile, lifecycle tests. | Uncompressed pack reader. |
| Global CLI pagination | 5 edits, 3 gate runs: query registry, shared projection, caller migration, tests, help. | Existing non-charge queries. |

## Further notes

### Inherited implementation-entry risk

The current implementation phase requires charge --full before this protocol exists.
A large approved spec can therefore reproduce the same bounded-output failure at implementation entry.
Staging this spec does not fix that executable defect.
This spec grants no manual reconstruction or raw-output workaround.

A future implementation start must resolve any observed entry contradiction through the existing reviewer authority.
This document task authorizes no production implementation.

### Reverified decision sources

Every structured source entry was reread before seam selection.
The approved report and all 15 decision tickets were read in place.
The complete topic and index moved into this spec without a top-level copy.
All moved repository-relative references changed in the same authoring pass.

| Structured source | Result |
| --- | --- |
| Compiled research report | Decisions and prototype limits retained. Historical measurements remain pinned to 9deb0a7. |
| internal/preflight/charge.go | Shared full/compact renderer and required-source loader remain the defect owner. |
| .agents/commands/bench-implement-spec.md | Full-charge requirement remains present with separate approval and supplement prerequisites. |
| .agents/commands/bench-review-implementation.md | Prepared shared evidence and native handoff remain current. |
| .agents/commands/bench-write-spec.md | Author-fork and approved-source rules remain unchanged and excluded from implementation. |

The seam exploration spent 18 targeted source and enforcement reads before selecting the existing preflight seam.
Subsequent reads confirmed system ownership, format encoding, and ticket-plan enforcement.
No source entry was unavailable.
Main advanced to `f613e42792d1b0a936834fde3a4449d5da31901d` during authoring.

The author reread both changed phase clauses and the current bounded repair policy from that commit.
Their new evidence-only correction exception changes repair review cadence, not preparation or delivery.
The implementation preserves that exception and its owning policy.
The assignment remains based on b20db7d until the coordinator composes the destination delta.

### Reader sweep and enforcement closure

The sweep used `rg --hidden` and excluded `.git/`.
It included `.mjs`, workflow files, hidden adapters, prose, tests, anchors, and canary mutations.
The shipped-surface words are charge, evidence, and preflight.
The new format reference resides in the shipped `.agents/` payload.

Repo-only specs and review records remain inputs owned by each linked repository.

| Reader or direct helper | Disposition |
| --- | --- |
| preflight command, preparedCommand, renderCharge, renderReviewCharge, renderReviewPacket | Migrate the build pack in CE-C1A, activate build guidance in CE-C1D, and migrate review in CE-C2. |
| renderChargePacket, chargeInvocation | Replace the full/omitted route and remove surviving charge --full generation. |
| loadChargeSources, readChargeSource, preparationCheckoutRefusal, preparationTicket | Preserve source and assignment policy through the new output path. |
| collectReviewEvidence, completeConsumerEvidence, completionEvidenceTable | Freeze outputs and bind their complete metadata in CE-C2. |
| buildChargeSourceSet.list, reviewChargeSourceSet.list, reviewEvidence.sources | Remain the source inventory derivation, or migrate together to one replacement owner. |
| internal/preflight charge, review_charge, proposal, delegated_evidence tests | Update charge expectations without changing proposal or ordinary preflight behavior. |
| build and review phase instructions | Migrate their retrieval and action instructions in their corresponding chunk. |
| registry_ft311_preparation and registry_ft311_review_dispatch | Replace exact retired-route anchors while preserving unrelated instructions. |
| prepared-build-approval and prepared-review-shared-evidence fixtures | Update exact mutation operands alongside each phase migration. |
| prepared-review-capable-handoff fixture | Add the trusted identity and retrieval action to the existing handoff contract. |
| internal/conformance/ft311_preparation_test.go | Derive the changed anchor inventory without adding a second source list. |
| internal/conformance/injected_ports_registry_test.go | Register any changed observer or new injected port with its actual producer test. |
| cmd/bench/main.go and help_inventory_test.go | Project the live preflight operation registry and independently detect omitted public inventory rows. |
| internal/chargeevidence/schema.go and shipped format reference | One executable schema owner with a generated reference and biting projection check. |
| Thin Codex adapters and Claude links | No protocol copy or independent source list. Canonical phase sources remain their owner. |
| .mjs scripts and workflows | No charge-field reader found. Release-preflight evidence is a different domain. |
| Historical specs, roadmap records, and the compiled research report | Historical mentions remain evidence, not executable retrieval instructions. |

The author read the public grammar, both current charge producers, the TOON adapter, and the two FT311 anchor registries.
The author also read the injected-port registry, ordinary architecture census, system owner, and completion-plan parser.
The current tests named as precedents were read in this session.
No reused test receives a new behavioral claim without a new fixture case.

Ticket closure derives from preflight missingClosures and the ticket binding registry.
Each ticket names its exact required fixture directories and bound registry files.
Directory entries omit a trailing slash because pathCovered adds the descendant separator.
These companions preserve existing checks and grant no unrelated behavioral scope.

### Fixed pre-review proof checklist

- Cited symbols: the reader table resolves the production symbols to `internal/preflight/` definitions read in this session.
- Import edges: preflight imports the shared TOON adapter, Git owner, existing collectors, and proposed evidence store. No reverse import is proposed.
- Source-row clauses and occurrences: the table below maps each approved decision clause to its executable or review-owned rows.
- Promised field labels: the tables specify the contract; one executable format registry drives production encoding, decoding, and generated reference tables.
- Changed-function callers: command calls chargeCommand, which calls preparedCommand. The prepared owner calls the build or review renderer.
- Changed-function callers: renderReviewCharge calls renderReviewPacket. Both packet paths call renderChargePacket, which owns the retired projection.
- Changed-function callers: chargeInvocation supplies the old build and review next actions. The new next-action owner replaces both consumers.
- Copy survival: the source-policy row fails when an independent required-source list survives. The retired-route rows fail when charge --full remains executable.

### Flagged additions

No product capability was added beyond the approved map.
The following engineering choices make the approved behavior implementable and reviewable.

| Choice | Derivation |
| --- | --- |
| 24-byte little-endian header and explicit marker | Defines the custom container approved by decision 14. |
| 8,192-byte UTF-8 pages and typed manifest schema | Makes decisions 8 and 12 concrete within the encoded response bound. |
| Versioned identity-bound cursor | Implements decision 13 without server progress. |
| Metadata as a required source | Binds large inventories and completion facts under decisions 8 and 12. |
| Shared operation lock plus serialized writer lifetime | Implements decisions 7 and 10 without a persistent reservation ledger. |
| Executable format and operation registries | Keep serialization, generated documentation, parsing, and public inventory single-sourced under the repository standard. |
| Four serial checkpoints replacing CE-C1 | Delivers the approved behavior within retained author contexts while preserving phase prerequisites. |
| Independent consumer fixture and actual transport records | Supplies the production validation required by decision 5. |
| Explicit partial-cleanup disposition | Preserves honest outcomes when decision 10's apply operation encounters a filesystem failure. |

### Source-sentence-to-row table

The source is the compiled decision ticket's Answer unless the table names the report.
Repeated occurrences in the report are historical evidence for the same predicate, not another policy owner.

| Source clause and occurrences | Acceptance rows |
| --- | --- |
| Decision 2: "All required pinned evidence must be verified and available to the consumer before action." Report action barrier. | CE1, CE2, CE3, CE4, CE5, CE6, CE7, CE8, CE9, CE10, CE11, CE12, CE101, CE102, CE103, CE105, CE106, CE140, CE141, CE142, CE143, CE144, CE145, CE146 |
| Decision 3: "The artifact remains readable across worktrees and after its originating assignment releases." Report lifecycle. | CE55, CE56, CE57, CE58, CE59, CE62, CE63, CE64, CE65, CE86, CE87, CE88, CE89, CE90, CE120, CE123, CE124 |
| Decision 4: "Write-spec keeps its current authoring, approval, fork, and resume behavior." Report scope. | CE111, CE112, CE113 |
| Decision 5: "Production validation must cover bounded manifest delivery, encoding expansion, the verification interface, and supported harness transports." | CE5, CE6, CE7, CE8, CE9, CE10, CE11, CE12, CE13, CE14, CE15, CE105, CE106, CE115, CE116, CE117, CE131, CE132, CE133, CE134, CE135, CE136, CE137, CE138, CE139 |
| Decision 6: "Remove charge --full from the supported grammar." Report migration. | CE108, CE109, CE110, CE111, CE112, CE113 |
| Decision 7: "Implementation must validate supported platforms, production Bench assignment release, concurrent storage safety, and hostile filesystem behavior." | CE28, CE29, CE30, CE31, CE32, CE33, CE34, CE35, CE36, CE37, CE38, CE39, CE40, CE55, CE56, CE83, CE84, CE85, CE86, CE87, CE88, CE89, CE90, CE91, CE92, CE93, CE94, CE95, CE96, CE97, CE98, CE99, CE100, CE115, CE116, CE117, CE118, CE125, CE126, CE129 |
| Decision 8: "Each preparation, retrieval, verification, cleanup, and error response contains at most 48,000 encoded stdout bytes." | CE13, CE14, CE15, CE78, CE79, CE80, CE81, CE82, CE119, CE121, CE122, CE131, CE132, CE133, CE134, CE135, CE136, CE137, CE138, CE139 |
| Decision 9: "Before action, the phase requires complete verified delivery, current source checks, reviewer approval, and the task supplement." | CE57, CE58, CE59, CE101, CE102, CE103, CE104, CE105, CE106, CE107, CE123, CE124, CE140, CE141, CE142, CE143, CE144, CE145, CE146 |
| Decision 10: "The quota includes published artifacts and temporary writes." "The owner performs no automatic eviction." | CE74, CE75, CE76, CE77, CE78, CE79, CE80, CE81, CE82, CE83, CE84, CE85, CE86, CE87, CE88, CE89, CE90, CE119, CE121, CE122, CE130 |
| Decision 11: "Changed generated bytes or declared provenance produce a changed evidence identity." | CE66, CE67, CE68, CE69, CE70, CE71, CE72, CE73, CE127 |
| Decision 12: "The identity hashes the complete canonical manifest bytes." "Full artifact verification checks every stored page and source digest." | CE5, CE6, CE7, CE8, CE9, CE10, CE11, CE12, CE16, CE17, CE18, CE19, CE20, CE21, CE22, CE23, CE24, CE25, CE26, CE27, CE28, CE29, CE30, CE31, CE32, CE33, CE34, CE35, CE36, CE37, CE38, CE39, CE40, CE60, CE61, CE125, CE126 |
| Decision 13: "Each read returns an exact next action or an explicit stream end." "Reading an existing artifact never reruns its collectors." | CE3, CE4, CE62, CE63, CE64, CE65, CE66, CE67, CE68, CE69, CE70, CE71, CE72, CE104, CE105, CE106, CE107, CE120, CE127 |
| Decision 14: "The reader derives physical offsets from validated lengths instead of maintaining a second offset registry." | CE28, CE29, CE30, CE31, CE32, CE33, CE34, CE35, CE36, CE37, CE38, CE39, CE40, CE125, CE126 |
| Decision 15: "This approval does not authorize production implementation." | Phase scope only. This change stages documents and implements no production behavior. |
| Report: "The source-set descriptors should remain the one policy owner for renderer, artifact, verifier, and required-next-action derivation." | CE114 |

The review repair adds the following exact protocol predicates from the same closed sources.

| Source clause and occurrences | Acceptance rows |
| --- | --- |
| Decision 4: preserve ordinary behavior; safe prefactoring of the existing charge consumer. | CE147 |
| Decision 6: remove charge full as each canonical consumer migrates. | CE173 |
| Decisions 12 and 14: the canonical manifest and documented single pack format; repository one-source rule. | CE148 |
| Decisions 8 and 12: complete metadata and exact bounded response fields. | CE149, CE153, CE154, CE155, CE160, CE161, CE164, CE165, CE169, CE170, CE171 |
| Decision 13: exact grammar, mutually exclusive read modes, and manifest-first preparation action. | CE150, CE151, CE156, CE158, CE159, CE163, CE166, CE168, CE174 |
| Decision 10: deterministic explicit capacity recovery without eviction. | CE157, CE175 |
| Project public-help inventory rule applied to the decision 13 grammar. | CE152, CE162, CE167, CE172 |

### Completion plan

The review pickup stays absent until a real review or proof record exists.
Its fence is present from the first implementation ticket.

```bench-completion-plan
{"version":1,"chunks":[{"id":"CE-C1A","tickets":["1-validate-legacy-prepared-packs.md"],"verification":[{"id":"preflight","command":"bench test --package ./internal/preflight"},{"id":"format","command":"bench test --package ./internal/chargeevidence"},{"id":"mutation","command":"bench test --package ./internal/chargeevidence","probe":"omit a declared field from the shipped format projection"}]},{"id":"CE-C1B","tickets":["2-publish-bounded-build-evidence.md"],"verification":[{"id":"preflight","command":"bench test --package ./internal/preflight"},{"id":"store","command":"bench test --package ./internal/chargeevidence"},{"id":"inventory","command":"bench test --package ./cmd/bench"},{"id":"system","command":"bench test --check system"},{"id":"mutation","command":"bench test --package ./internal/preflight","probe":"omit manifest_bytes from the prepared response"}]},{"id":"CE-C1C","tickets":["3-verify-current-evidence.md"],"verification":[{"id":"preflight","command":"bench test --package ./internal/preflight"},{"id":"store","command":"bench test --package ./internal/chargeevidence"},{"id":"inventory","command":"bench test --package ./cmd/bench"},{"id":"system","command":"bench test --check system"},{"id":"mutation","command":"bench test --package ./internal/preflight","probe":"accept --verify with --cursor"}]},{"id":"CE-C1D","tickets":["4-activate-bounded-build-guidance.md"],"verification":[{"id":"preflight","command":"bench test --package ./internal/preflight"},{"id":"inventory","command":"bench test --package ./cmd/bench"},{"id":"guidance","command":"bench test --check docs-currency-workflow"},{"id":"guidance-cases","command":"bench test --package ./internal/conformance --run TestEvidence"},{"id":"mutation","command":"bench test --check docs-currency-workflow","probe":"permit build action without a current binding"}]},{"id":"CE-C2","tickets":["5-freeze-bounded-review-evidence.md"],"verification":[{"id":"preflight","command":"bench test --package ./internal/preflight"},{"id":"store","command":"bench test --package ./internal/chargeevidence"},{"id":"inventory","command":"bench test --package ./cmd/bench"},{"id":"guidance","command":"bench test --check docs-currency-workflow"},{"id":"guidance-cases","command":"bench test --package ./internal/conformance --run TestEvidence"},{"id":"mutation","command":"bench test --package ./internal/preflight","probe":"rerun a collector during retrieval"}]},{"id":"CE-C3","tickets":["6-clean-evidence-explicitly.md"],"verification":[{"id":"preflight","command":"bench test --package ./internal/preflight"},{"id":"store","command":"bench test --package ./internal/chargeevidence"},{"id":"inventory","command":"bench test --package ./cmd/bench"},{"id":"system","command":"bench test --check system"},{"id":"mutation","command":"bench test --package ./internal/preflight","probe":"skip cleanup fingerprint revalidation before deletion"}]},{"id":"CE-C4","tickets":["7-preserve-consumer-context.md"],"verification":[{"id":"preflight","command":"bench test --package ./internal/preflight"},{"id":"guidance","command":"bench test --check docs-currency-workflow"},{"id":"guidance-cases","command":"bench test --package ./internal/conformance --run TestEvidence"},{"id":"mutation","command":"bench test --check docs-currency-workflow","probe":"permit a final cursor to replace available required context"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/bounded-charge-evidence/spec.md"},{"id":"preflight","command":"bench test --package ./internal/preflight"},{"id":"store","command":"bench test --package ./internal/chargeevidence"},{"id":"inventory","command":"bench test --package ./cmd/bench"},{"id":"system","command":"bench test --check system"},{"id":"guidance","command":"bench test --check docs-currency-workflow"},{"id":"guidance-cases","command":"bench test --package ./internal/conformance --run TestEvidence"}]}
```
