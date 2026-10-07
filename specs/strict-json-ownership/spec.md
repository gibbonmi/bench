# Strict JSON reads under one owner

Status: staged

Decision source: `decisions/architecture-primitives/tickets/1.md`, reviewed in `ft362-process-lifetime-map`, and its ready parent `decisions/architecture-primitives.md`

Verification log: 1 iteration to accept the spec — independent Sol6.1/xhigh review accepted SHA256 fab072a3665188f7a545335b3bb1a237821383c6dd44e9ea11472dbd1df92e17. The independent ticket review also accepted the graph at static confidence 9

Implementation approval: pending

## Problem

The vulnerability exception reader can accept malformed trailing bytes and incompatible second values. Its second decode refuses only a nil error.
The source is `internal/releasepreflight/vulnerability.go`, in `ValidateVulnerabilityPolicy`.

Strict readers also repeat recursive duplicate-key checks. The copies occur in gate verdicts, review records, and release evidence.
Their sources are `internal/gate/verdict.go`, `internal/reviewrecord/parse.go`, and `internal/releaseevidence/json_validation.go`.
The existing owner is `internal/jsonfile/decode.go`.

These readers serve the Bench executable in this repository and every repository that links the kit. Their domain and filesystem policies differ.
A common parser must preserve those policies.

## Solution

Repair the vulnerability exception reader first. It accepts one complete exception document and refuses duplicate keys and unknown typed fields.
Then move the repeated strict grammar into `jsonfile`. Caller adapters retain their domain validation and diagnostic categories.

The build preserves each caller's accepted newline policy. It does not require a final newline merely because a caller reads a file.
FT354 remains open for the separate durable-file and shell-argument outcomes. This spec has no Roadmap retirement field.

## User stories

Line: gpt-5.6-sol / high
Implementation-line reason: Caller compatibility is the hardest chunk. Existing parser seams constrain the API, but independent fixtures must prove diagnostic and schema preservation.
Harder chunks: SJ-C2, SJ-C4, SJ-C5

### Vulnerability policy

1. As a release operator, I want malformed suffixes refused, so that an exception document cannot conceal invalid bytes.
2. As a release operator, I want second values refused, so that every policy byte belongs to one document.
3. As a release operator, I want duplicate exception keys refused, so that an exception has one declared value per key.
4. As a release operator, I want unknown exception fields refused, so that a misspelled field cannot change the policy silently.
5. As a release operator, I want ordinary whitespace accepted, so that a user-authored policy needs no record framing.
6. As a release operator, I want missing and empty policy states preserved, so that current release decisions remain valid.
7. As a release operator, I want exception domain rules preserved, so that the parser change cannot extend an exception's authority.
8. As a release operator, I want scanner streams preserved, so that multiple scanner events remain valid input.

### Shared grammar

9. As a maintainer, I want duplicate keys refused recursively, so that nested records cannot replace an earlier value.
10. As a maintainer, I want typed unknown fields refused, so that strict schemas remain complete.
11. As a maintainer, I want complete consumption, so that strict readers cannot ignore a suffix.
12. As a record owner, I want explicit newline framing, so that persisted records retain their existing contract.
13. As a schema owner, I want exact key spelling where required, so that case aliases cannot bypass an exact schema.
14. As a caller author, I want classified structural errors, so that caller diagnostics do not depend on message matching.

### Caller compatibility

15. As a gate operator, I want cache and manifest decisions preserved, so that the common parser does not change authorization.
16. As a gate operator, I want sentinel categories preserved, so that a trailing document retains its current diagnostic.
17. As a reviewer, I want review documents accepted without record framing, so that fenced payloads remain readable.
18. As a reviewer, I want review limits and domain checks preserved, so that malformed review evidence remains unusable.
19. As a release operator, I want release evidence schemas preserved, so that shared grammar cannot weaken evidence validation.
20. As a release operator, I want dynamic manifest keys preserved, so that the requirement registry remains the schema authority.
21. As a release operator, I want package identity projections preserved, so that unrelated package.json fields remain valid.
22. As a maintainer, I want duplicate-tolerant readers explicitly preserved, so that consolidation cannot silently tighten their contracts.

### Ownership and evidence

23. As a maintainer, I want duplicate walkers removed, so that a future grammar correction has one implementation owner.
24. As a maintainer, I want an ownership check that detects a surviving copy, so that a partial migration cannot claim completion.
25. As a reviewer, I want differential caller evidence, so that preserved behavior has an independent oracle.
26. As a maintainer, I want each migration chunk green, so that the first defect repair can land before the refactor.

## Implementation decisions

### Terms and API

A **complete document** contains one JSON value followed only by JSON whitespace. Avoid: JSON stream, record framing.
A **persisted record** is a complete document whose caller contract requires a final newline. Avoid: every JSON file.
A **partial projection** reads selected fields while allowing unrelated fields. Avoid: permissive duplicate keys.

Retain `jsonfile.Decode`, `DecodeDocument`, and `DecodeExactDocument`. Their definitions in `internal/jsonfile/decode.go` already separate the framing contracts.
Use `DecodeDocument` for the exception array. Keep `Decode` for existing persisted-record callers.
Keep `DecodeExactDocument` for existing exact-schema callers. These choices correspond to rows SJ5, SJ22, SJ23, and SJ24.

Expose duplicate-key and trailing-data error categories through `errors.Is`. Gate adapters map those categories to their existing local sentinels.
Do not classify unknown fields by a new duplicate walker. Preserve the existing encoding/json field and type semantics.
Rows SJ25, SJ28, and SJ29 grade this contract.

Use the existing raw-document form for structural validation. `DecodeDocument` into `json.RawMessage` checks grammar before the package identity's existing partial decode.
For component objects, decode into `map[string]json.RawMessage`. The caller then applies the registry's exact key set.
Rows SJ40 through SJ45 grade these distinctions. A new structural-validation API is unnecessary.

The common owner contains no filesystem, size, null, domain, expiry, or permission policy. Caller adapters retain these decisions.
The source definitions below identify their current owners. Rows SJ6 through SJ15 and the caller differential rows grade their preservation.

### Consolidation table

| site | old rule | resulting rule | disposition |
|---|---|---|---|
| ValidateVulnerabilityPolicy | Typed unknown refusal with an incomplete suffix check | DecodeDocument with recursive duplicate refusal | SJ1–SJ15 |
| gate strictJSON | Typed decode, suffix refusal, recursive duplicate walk | Shared document grammar with local sentinel mapping | SJ26–SJ31 |
| reviewrecord decode | Bounds check, duplicate walk, typed decode, suffix refusal | Bounds check and shared document grammar | SJ32–SJ35 |
| releaseevidence decodeStrict | json.Valid, duplicate walk, typed decode | Shared document grammar with caller error context | SJ36–SJ39 |
| component decodeExactObject | Duplicate walk and exact dynamic map keys | Shared map decode and the same exact dynamic keys | SJ40–SJ42 |
| package identity | Duplicate walk and partial json.Unmarshal | Shared raw-document decode and the same partial json.Unmarshal | SJ43–SJ45 |
| freshness parseSeal | Typed unknown and suffix refusal with duplicate acceptance | Preserve this reader in place | Won't handle below |
| prospectiveartifact readRecordAt | Typed unknown and suffix refusal with duplicate acceptance | Preserve this reader in place | Won't handle below |
| existing jsonfile consumers | Explicit document, record, or exact-document API | Preserve each API selection | SJ22–SJ24 |

A migration must preserve observable diagnostic classes for competing invalid inputs. The caller baseline includes unknown-plus-suffix and duplicate-plus-suffix cases.
A thin adapter may preserve classification precedence. It must not retain a recursive grammar implementation.

## Implementation chunks

The independent spec review accepted the artifact before ticket authoring. Every chunk follows its predecessor's green checkpoint and independent semantic review.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
|---|---|---|---|---|
| SJ-C1 / 1-policy-document.md | The exception validator refuses concealed document bytes | SJ1–SJ15 | TestVulnerabilityPolicyDocumentContract, TestVulnerabilityPolicyDomainCompatibility | no |
| SJ-C2 / 2-jsonfile-gate.md | Shared error categories serve gate callers | SJ16–SJ25, SJ26–SJ31 | TestDocumentStructuralContract, TestGateJSONCompatibility | yes |
| SJ-C3 / 3-review-documents.md | Review readers use the common grammar | SJ32–SJ35 | TestReviewJSONCompatibility | no |
| SJ-C4 / 4-release-documents.md | Release readers retain strict and partial schemas | SJ36–SJ45 | TestReleaseJSONCompatibility, TestComponentJSONCompatibility, TestPackageIdentityJSONCompatibility | yes |
| SJ-C5 / 5-ownership-closure.md | The gate detects a surviving duplicate walker | SJ46–SJ50 | TestStrictJSONOwnerBites, TestStrictJSONOwnerNegativeControls | yes |

SJ-C1 uses the existing owner without waiting for new error categories. SJ-C2 supplies the shared error contract before later adapters consume it.
SJ-C5 runs only after every migration has landed. Its check does not make intermediate copies a premature red.

## Testing decisions

Test the exception policy through `ValidateVulnerabilityPolicy` with a fixed date. Test grammar through the same exported jsonfile APIs that production callers use.
These are in-process seams with no external dependency. The existing prior art is `internal/jsonfile/decode_test.go`.

Test gate manifests through their existing loader and cache verdicts through their existing load path. Test review records through Parse and ReadPlan.
The precedents are `internal/gate/manifest_test.go`, `internal/gate/admin_readers_test.go`, and `internal/reviewrecord/record_test.go`.
Tests that need a repository reuse the existing fixture owners. New gate tests run the branch-native architecture census.

Test release schemas through their real validators. Test package identity through the artifact validator with a valid archive fixture.
The precedent is `internal/releaseevidence/package_artifact_test.go`. Extend its fixture rather than paste a second archive builder.
No test crosses a package seam through another package's private test helper.

Before a caller migration, capture baseline results for the enumerated fixture family. Permanent tests compare the current result with that committed baseline.
An ordinary test does not rewrite repository production sources. Normalize only incidental path roots, never a refusal category or declared policy value.
Independent expectations require a recorded omission or swap that makes the named test red. Each acceptance row needs its own red-to-green evidence.

The final owner check joins root conformance as `strict-json-owner`. It runs in the dev tier and therefore also the ship tier.
Its subject is the graded root. Its declared input class is go-source.
The check uses the existing conformance registry and executable binding. No CLI verb, injected port, or canary family is added.

The check grades duplicate-token walkers in the migrated caller packages. It also grades the required caller-to-owner composition.
It permits local diagnostic and domain adapters. It does not forbid encoding/json for partial projections or stream readers.

An absent caller package is an absent optional surface. An unreadable or unparsable present source produces a named refusal.
Synthetic mutation fixtures prove both the owner rule and its exclusions. Root conformance proves the registered check executes.

### Seam diagram

    exception bytes + findings + fixed date
        │
        ▼
    ValidateVulnerabilityPolicy ──▶ jsonfile.DecodeDocument ──▶ policy decision
        ▲                              ▲
        │                              │
    validator tests                document tests

    cache / manifest / review / release bytes
        │
        ▼
    caller policy adapter ──▶ jsonfile document grammar ──▶ caller result
        ▲
        │
    permanent baseline and hostile-document fixtures

    graded production sources
        │
        ▼
    root conformance ──▶ strict-json-owner ──▶ named ownership diagnostic

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| SJ1 | 1 | The validator refuses [] followed by nope | planned TestVulnerabilityPolicyDocumentContract in internal/releasepreflight/vulnerability_test.go | The existing second decode accepts this malformed suffix |
| SJ2 | 2 | The validator refuses [] followed by {} | planned TestVulnerabilityPolicyDocumentContract in internal/releasepreflight/vulnerability_test.go | A valid second object cannot disappear |
| SJ3 | 2 | The validator refuses [] followed by 7, true, null, a string, or an array | planned TestVulnerabilityPolicyDocumentContract in internal/releasepreflight/vulnerability_test.go | An incompatible second value exposes the incomplete suffix check |
| SJ4 | 3 | The validator refuses repeated id, reason, or expires keys in one exception | planned TestVulnerabilityPolicyDocumentContract in internal/releasepreflight/vulnerability_test.go | The last decoded key cannot replace the first |
| SJ5 | 5 | The validator accepts [] with leading and trailing JSON whitespace without a final newline | planned TestVulnerabilityPolicyDocumentContract in internal/releasepreflight/vulnerability_test.go | Persisted-record framing would reject a valid policy |
| SJ6 | 4 | The validator refuses an exception with an extra typed field | planned TestVulnerabilityPolicyDocumentContract in internal/releasepreflight/vulnerability_test.go | A permissive decoder would discard the extra field |
| SJ7 | 6 | An absent policy with no findings returns nil | planned TestVulnerabilityPolicyDomainCompatibility in internal/releasepreflight/vulnerability_test.go | Absence cannot become an empty-file refusal |
| SJ8 | 6 | A present zero-byte policy returns the empty-file refusal | planned TestVulnerabilityPolicyDomainCompatibility in internal/releasepreflight/vulnerability_test.go | An empty file cannot become an absent policy |
| SJ9 | 6 | Present [] and null policies with no findings return nil | planned TestVulnerabilityPolicyDomainCompatibility in internal/releasepreflight/vulnerability_test.go | Shared null rejection would change current domain behavior |
| SJ10 | 7 | An uncovered finding remains an uncovered-vulnerabilities refusal | planned TestVulnerabilityPolicyDomainCompatibility in internal/releasepreflight/vulnerability_test.go | A parser repair cannot grant a missing exception |
| SJ11 | 7 | A used exception with expiry equal to today returns nil | planned TestVulnerabilityPolicyDomainCompatibility in internal/releasepreflight/vulnerability_test.go | The expiry comparison must remain strictly less than today |
| SJ12 | 7 | An exception dated before today remains expired | planned TestVulnerabilityPolicyDomainCompatibility in internal/releasepreflight/vulnerability_test.go | Shared decode success cannot bypass expiry |
| SJ13 | 7 | Duplicate exception entries retain the duplicate-exception refusal | planned TestVulnerabilityPolicyDomainCompatibility in internal/releasepreflight/vulnerability_test.go | Object-key grammar cannot replace the domain uniqueness check |
| SJ14 | 7 | An unused exception retains the unused-exception refusal | planned TestVulnerabilityPolicyDomainCompatibility in internal/releasepreflight/vulnerability_test.go | A document repair cannot widen exception coverage |
| SJ15 | 7 | Blank required values and invalid full dates retain their current domain refusals | planned TestVulnerabilityPolicyDomainCompatibility in internal/releasepreflight/vulnerability_test.go | Typed decode success alone does not establish a valid exception |
| SJ16 | 9 | Document decode refuses duplicate keys inside an object in an array | planned TestDocumentStructuralContract in internal/jsonfile/decode_test.go | A root-only duplicate check would pass |
| SJ17 | 9 | Document decode refuses repeated keys with equivalent escaped spellings | planned TestDocumentStructuralContract in internal/jsonfile/decode_test.go | Raw token spelling must not create distinct key identities |
| SJ18 | 10 | Typed document decode refuses an unknown nested field | planned TestDocumentStructuralContract in internal/jsonfile/decode_test.go | Root-only schema refusal would pass |
| SJ19 | 11 | Document decode refuses every second JSON value kind | planned TestDocumentStructuralContract in internal/jsonfile/decode_test.go | Decoding a fixed second target misses incompatible values |
| SJ20 | 11 | Document decode refuses malformed trailing bytes | planned TestDocumentStructuralContract in internal/jsonfile/decode_test.go | Success requires EOF rather than any decode error |
| SJ21 | 11 | Document decode accepts JSON whitespace after its value | planned TestDocumentStructuralContract in internal/jsonfile/decode_test.go | A byte-equality suffix check would reject valid whitespace |
| SJ22 | 12 | Existing record callers retain final-newline refusal | planned TestDocumentFramingCompatibility in internal/jsonfile/decode_test.go | Replacing Decode with DecodeDocument would remove framing |
| SJ23 | 12 | Existing document callers retain acceptance without a final newline | planned TestDocumentFramingCompatibility in internal/jsonfile/decode_test.go | Replacing document decode with Decode would add framing |
| SJ24 | 13 | Exact document decode refuses a case-alias key that ordinary typed decode accepts | planned TestDocumentFramingCompatibility in internal/jsonfile/decode_test.go | A shared non-exact reader would erase the exact schema |
| SJ25 | 14 | errors.Is identifies duplicate and trailing errors from the owner | planned TestDocumentErrorCategories in internal/jsonfile/decode_test.go | Message matching cannot satisfy the classification contract |
| SJ26 | 15, 25 | Gate manifest results match the committed baseline family | planned TestGateJSONCompatibility in internal/gate/json_compatibility_test.go | Domain and competing-refusal drift changes the stored result |
| SJ27 | 15, 25 | Cache loader results match the committed baseline family | planned TestGateJSONCompatibility in internal/gate/json_compatibility_test.go | A parser change cannot alter ready, invalid, absent, or unavailable results |
| SJ28 | 16 | Gate duplicate errors retain errDuplicateJSONName through errors.Is | planned TestGateJSONCompatibility in internal/gate/json_compatibility_test.go | A generic wrapped error would lose the local sentinel |
| SJ29 | 16 | Gate trailing errors retain errTrailingJSON through errors.Is | planned TestGateJSONCompatibility in internal/gate/json_compatibility_test.go | The manifest's trailing diagnostic depends on this sentinel |
| SJ30 | 15 | Gate documents valid without a final newline remain valid | planned TestGateJSONCompatibility in internal/gate/json_compatibility_test.go | These current readers permit unframed documents |
| SJ31 | 15 | Unknown manifest fields retain the unknown-field key diagnostic | planned TestGateJSONCompatibility in internal/gate/json_compatibility_test.go | A new owner prefix must not hide the rejected key |
| SJ32 | 17 | Parse accepts a valid review document without a final newline | planned TestReviewJSONCompatibility in internal/reviewrecord/json_compatibility_test.go | Fenced payloads trim their final newline |
| SJ33 | 17 | ReadPlan accepts the real fenced completion-plan payload | planned TestReviewJSONCompatibility in internal/reviewrecord/json_compatibility_test.go | A unit-only decode test would miss the fence producer |
| SJ34 | 18, 25 | Review results match the committed baseline family | planned TestReviewJSONCompatibility in internal/reviewrecord/json_compatibility_test.go | Null, version, schema, and evidence refusals remain caller-owned |
| SJ35 | 18 | Review byte limits retain their current refusal | planned TestReviewJSONCompatibility in internal/reviewrecord/json_compatibility_test.go | Moving grammar must not bypass bounds.ClassifyBytes |
| SJ36 | 19, 25 | Strict release results match the committed baseline family | planned TestReleaseJSONCompatibility in internal/releaseevidence/json_compatibility_test.go | Embedded registries and real validators retain their domain results |
| SJ37 | 19 | Strict release validators refuse recursive duplicate keys | planned TestReleaseJSONCompatibility in internal/releaseevidence/json_compatibility_test.go | A migrated adapter must reach the shared duplicate check |
| SJ38 | 19 | Strict release validators refuse unknown typed fields | planned TestReleaseJSONCompatibility in internal/releaseevidence/json_compatibility_test.go | A structural-only migration would weaken typed schemas |
| SJ39 | 19 | Requirement evidence retains its final-newline refusal | planned TestReleaseJSONCompatibility in internal/releaseevidence/json_compatibility_test.go | The caller still requires persisted evidence framing |
| SJ40 | 20 | A component object accepts exactly the registry's declared key set | planned TestComponentJSONCompatibility in internal/releaseevidence/json_compatibility_test.go | Hardcoded struct tags cannot replace dynamic schema authority |
| SJ41 | 20 | A component object refuses an extra or missing declared key | planned TestComponentJSONCompatibility in internal/releaseevidence/json_compatibility_test.go | A generic map decoder would admit an incomplete schema |
| SJ42 | 20 | A component object refuses duplicate keys recursively | planned TestComponentJSONCompatibility in internal/releaseevidence/json_compatibility_test.go | Partial object extraction must not hide nested duplicate values |
| SJ43 | 21 | Artifact validation accepts package.json with unrelated fields | planned TestPackageIdentityJSONCompatibility in internal/releaseevidence/package_artifact_test.go | Typed DisallowUnknownFields on the identity projection would reject it |
| SJ44 | 21 | Artifact validation refuses a duplicate package identity key | planned TestPackageIdentityJSONCompatibility in internal/releaseevidence/package_artifact_test.go | A partial json.Unmarshal alone would accept it |
| SJ45 | 21 | Artifact validation refuses duplicate keys inside an unrelated package field | planned TestPackageIdentityJSONCompatibility in internal/releaseevidence/package_artifact_test.go | Schema tolerance must not become structural tolerance |
| SJ46 | 23, 24 | The registered owner check refuses a planted recursive duplicate walker | planned TestStrictJSONOwnerBites in internal/conformance/strict_json_owner_test.go | A surviving implementation copy makes the check red |
| SJ47 | 24 | The registered owner check refuses a caller that bypasses jsonfile | planned TestStrictJSONOwnerBites in internal/conformance/strict_json_owner_test.go | Deleting a walker alone cannot establish shared ownership |
| SJ48 | 22, 24 | The owner check permits each named preservation exclusion | planned TestStrictJSONOwnerNegativeControls in internal/conformance/strict_json_owner_test.go | A broad decoder ban would silently tighten excluded contracts |
| SJ49 | 24 | Root conformance executes strict-json-owner against a planted source copy | planned TestStrictJSONOwnerRegisteredPath in internal/conformance/strict_json_owner_test.go | A correct but unregistered helper is insufficient |
| SJ50 | 26 | Every chunk closes with its complete rows and a green lane result | review-owned chunk checkpoint records | A later migration cannot hide an incomplete earlier checkpoint |

Not covered: story 8 — findingIDs retains its current multi-value scanner stream and existing test coverage

### Edge inventory

The input family includes absent, zero-byte, whitespace-only, null, empty array, and valid nonempty documents.
It includes invalid UTF-8 behavior inherited from encoding/json, escaped duplicate names, nested objects, arrays, wrong types, malformed values, and suffixes.
The compatibility oracle retains encoding/json behavior outside the required duplicate and suffix checks. It adds no Unicode normalization or numeric coercion rule.

The vulnerability family includes every second value kind. It includes malformed suffixes after both empty and nonempty exception arrays.
An invalid value such as [{}] remains a domain refusal. An empty array with findings remains an uncovered-vulnerabilities refusal.

The gate baseline includes valid framed and unframed documents, whitespace, null, unknown fields, duplicate keys, syntax errors, and suffixes.
It includes unknown-plus-suffix and duplicate-plus-suffix cases. Its cache cases cover size, mode, regular-file, absent, and unreadable states.
The review baseline uses the same grammar family with valid record and plan seeds. It also covers byte limits and invalid domain fields.

The release baseline covers embedded registries, governance policies, producer envelopes, SPDX documents, reproducibility records, and native proofs.
Its component cases use the requirement registry's field labels. Its artifact cases use the existing archive fixture with a valid identity.
The release family retains null decisions and every existing requirement newline check.

The source-check family includes an absent package, an empty package, a present unreadable source, and a present unparsable source.
It includes a renamed recursive walker and a typed decoder that bypasses jsonfile. Negative controls include strings, comments, tests, and excluded streams.
The check must not infer ownership from a historical count or a retired helper name alone.

Paths with spaces and glob characters remain fixture cases for filesystem callers. No shell argv or printed command changes occur here.
JSON strings with escaped control bytes retain caller validation. No terminal sanitizer, TOON emitter, or sink rule changes occur here.

No package variable replacement reaches a subprocess. The package artifact fixture uses the existing requirements override with its restore function.
There is no new cleanup operation, remote service, process owner, clock source, or publication step.

**Won't handle:** duplicate acceptance in freshness.parseSeal — this reader survives with its current unknown-field, suffix, and digest checks.
**Won't handle:** duplicate acceptance in prospectiveartifact.readRecordAt — this reader survives with its current file and record checks.
**Won't handle:** multi-value govulncheck events — findingIDs survives as the scanner's stream reader.
**Won't handle:** go list JSON streams — freshness.buildInputs survives as a multi-package stream reader.
**Won't handle:** other permissive JSON projections and event readers — their named families remain outside this strict-reader migration.
**Won't handle:** new path or filesystem guarantees — existing regular-file and metadata readers retain those contracts.

## Ownership fences

These are prospective build fences. The author of this planning pass writes only this spec.
The independent reviewer checks the fences before ticket authoring. A required expansion returns through the approved plan-expansion procedure.

SJ-C1

- `internal/releasepreflight/vulnerability.go`
- `internal/releasepreflight/vulnerability_test.go`

SJ-C2

- `internal/jsonfile/decode.go`
- `internal/jsonfile/decode_test.go`
- `internal/jsonfile/fields.go`
- `internal/gate/verdict.go`
- `internal/gate/json_compatibility_test.go`
- `internal/gate/testdata/strict-json-baseline.json`

SJ-C3

- `internal/reviewrecord/parse.go`
- `internal/reviewrecord/json_compatibility_test.go`
- `internal/reviewrecord/testdata/strict-json-baseline.json`
- `internal/reviewrecord/recordtest/fixture.go`

SJ-C4

- `internal/releaseevidence/json_validation.go`
- `internal/releaseevidence/component_manifest.go`
- `internal/releaseevidence/package_artifact.go`
- `internal/releaseevidence/json_compatibility_test.go`
- `internal/releaseevidence/package_artifact_test.go`
- `internal/releaseevidence/testdata/strict-json-baseline.json`

SJ-C5

- `internal/conformance/strict_json_owner_test.go`
- `internal/conformance/registry/checks.go`
- `internal/conformance/checks_test.go`
- `internal/conformance/tier_live_tree_test.go`
- `projects/benchkit.md`
- `CHANGELOG.md`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_retained_workflow.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `tests/canary/guidance-prose-budgets/over-budget-skill`
- `tests/canary/line-routing/line-binding-prose-drift`
- `tests/canary/skill-description-budgets/budget-table-missing`
- `tests/canary/skill-description-budgets/description-folded`
- `tests/canary/skill-description-budgets/description-missing`
- `tests/canary/skill-description-budgets/over-budget-command`
- `tests/canary/skill-description-budgets/over-budget-description`
- `tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-owner`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-routing`
- `tests/canary/workflow-guidance-anchors/benchkit-spec-ownership`
- `tests/canary/workflow-guidance-anchors/benchkit-system-suite-route`
- `tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns`
- `tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary`

Plan and review pickup:

- `specs/strict-json-ownership/spec.md`
- `specs/strict-json-ownership/tickets/`
- `reviews/strict-json-ownership.md`

No writer changes another spec or its tickets. No implementation writer changes the reviewed architecture map or its resolved decision tickets.
The final migration author owns the invariant for its whole affected package. Later tickets preserve the previously accepted caller contracts.

## Out of scope

Durable file replacement remains a separate outcome under architecture decision ticket 2. Its own spec derives its edit and gate-run estimate.
Shell argument spelling remains a separate outcome under architecture decision ticket 3. Its own spec derives its edit and gate-run estimate.
These are separate capabilities, rather than deferred parts of strict JSON.

A stricter freshness or prospective-owner schema needs a separate behavior decision. Estimate: two reader edits, two fixture edits, two gate runs.
Migration of unrelated permissive JSON readers needs a separate caller-contract inventory. Estimate: one inventory edit, one gate run before a build estimate.
A manual performance benchmark is outside this planning task. This spec makes no throughput or allocation claim.

## Further notes

### Source and caller inventory

The whole hidden-tree sweep searched Go, JavaScript, shell, and workflow files. It excluded only .git.
The needles were json.NewDecoder, json.Unmarshal, jsonfile, DisallowUnknownFields, duplicate helper names, JSON.parse, vuln-exceptions, and govulncheck.
The source inventory is current at base `37ac80b8adeff4ecbd57f3d813ff8937f2bcbc68`. Historical survey counts are not acceptance predicates.

The changed-function sweep used bench consumers with the following symbols:

- gate.strictJSON and gate.rejectDuplicateNames
- reviewrecord.decode and reviewrecord.uniqueJSON
- releaseevidence.decodeStrict and releaseevidence.rejectDuplicateJSONKeys
- releaseevidence.decodeExactObject and releasepreflight.ValidateVulnerabilityPolicy
- jsonfile.Decode, jsonfile.DecodeDocument, and jsonfile.DecodeExactDocument

Gate strictJSON serves loadVerdict, decodeManifest, and loadManifest. Their sources are verdict.go, manifest.go, and subject.go in internal/gate.
Review decode serves Parse and ReadPlan. Their sources are parse.go and plan.go in internal/reviewrecord.
Release decodeStrict serves the following source-defined readers:

- readReproducibility and inspectNativeProofs in artifact_proofs.go
- readReleasePlan, readReleaseArtifacts, and archiveInventory in release_plan.go
- readRollbackTarget, validateRequirementBytes, validateGovernancePolicy, and validateSPDXDocument in release_requirements.go
- inspectRequirements in requirement_inspection.go
- loadRegistry and loadRequirements in types.go

The release duplicate helper also serves decodeExactObject and validatePackageEvidence. Their sources are component_manifest.go and package_artifact.go.
The component object helper serves root, component, target, and file-entry objects. Their field lists derive from requirements.ComponentManifest.
ValidateVulnerabilityPolicy serves runner.runVulnerability. findingIDs produces the findings from the govulncheck event stream before that call.

Existing jsonfile record consumers are assessment.Store.Read, intent.readPath, repairpilot.load, and worktree.decodeMarker.
Existing document consumers are assessment.readJSON, assessment.collectHarness, commitment.ParseProposal, and commitment repository receipt readers.
The other document consumers are repairpilot.readRecordInput and worktree.loadBuildOutputs. Exact consumers are commitment.Parse, ParseProposal, and ParseEvidence.
These consumers retain their current API and require no production edits. Their package tests remain integration checks for owner API changes.

The named permissive Go families are harness, benchguard, writeguard, gitguard, stophook, lines, harnesstranscript, models, publication, adopt, and capturetx.
Other exclusions are otelrecord, testreport, test fixture decodes, canary inventory, and intent recovery projections.
JavaScript exclusions are the offline registry, binary repair, release-plan, artifact verifier, and release-evidence producer readers.
Their runtime and schema contracts remain outside this migration. No workflow calls a changed public CLI command.

The build refreshes the full inventory before dispatch. A newly found strict duplicate walker joins an approved fence before its author writes.
An unrelated permissive reader takes a named exclusion. A stronger caller policy requires a separate reviewer decision.

### Pre-review proof checklist

- Cited symbols: Definitions and consumers resolve in the source files named above
- Import edges: go list resolves jsonfile, gate, reviewrecord, releaseevidence, and releasepreflight without a cycle
- Source-row clauses and occurrences: The source-to-row table below covers the reviewed ticket and the parent destination
- Promised field labels: id, reason, expires remain exception labels, while component labels derive from requirements.ComponentManifest
- Changed-function callers: The current bench consumers inventory is recorded above, including private helpers
- Copy survival: SJ46 and SJ49 fail when a recursive walker survives outside the common owner
- Rendered-shape readers: No command output shape moves, while SJ31 preserves unknownField's existing diagnostic extraction
- Pin operators: SJ25, SJ28, and SJ29 use errors.Is, while compatibility rows compare complete normalized baseline results
- Entry reads: Existing loader filesystem reads remain caller-owned, and no new ambient read enters jsonfile
- Derived expectations: Baselines record pre-change caller results, while independent hostile-input assertions require observed mutation reds
- Consolidated rules: The consolidation table maps every changed rule to its acceptance rows or preservation exclusion
- Quantified obligations: The source inventory governs all migrated strict adapters and all recursive duplicate positions
- Workflow-step writes: none

The enforcement reads are .bench/gate.sh, the conformance registry, its executable binding, tier_test.go, and ordinary_build_census_test.go.
The fixture reads include the gate manifest and admin readers, the review record fixture, and the package artifact fixture.
DATA_HANDLING.md supplies the existing persistence and local-evidence posture. This spec changes no data path, retained value, or access mode.

### Source-to-row table

| reviewed source clause | occurrence | rows or disposition |
|---|---|---|
| Use the existing jsonfile owner | ticket 1 answer | SJ16–SJ25, SJ46–SJ49 |
| Use DecodeDocument for a complete user-authored JSON document | ticket 1 answer | SJ1–SJ6, SJ23 |
| Use Decode for a persisted record whose contract requires a final newline | ticket 1 answer | SJ22, SJ39 |
| Use DecodeExactDocument only when exact key spelling belongs to the schema contract | ticket 1 answer | SJ24 |
| First repair the vulnerability exception parser | ticket 1 answer | SJ-C1 precedes all migration chunks |
| Require complete consumption, duplicate-key rejection, and unknown-field rejection for exception documents | ticket 1 answer | SJ1–SJ6 |
| Retain valid trailing whitespace and the existing exception semantics | ticket 1 answer | SJ5, SJ7–SJ15 |
| Test malformed suffixes, second values, duplicate keys, missing files, and valid empty policies through the real validator | ticket 1 answer | SJ1–SJ9 |
| Do not impose persisted-record newline framing on that file | ticket 1 answer | SJ5 |
| Use one owner for each primitive while keeping their contracts distinct | parent Destination | SJ16–SJ25, SJ46–SJ49 |
| Separate the strict-JSON repair from file durability and printed-command changes | parent Destination | Separate out-of-scope outcomes |
| Migration batches with complete caller inventories and unchanged sentinel categories | parent Spec-writer discretion | SJ28–SJ29, caller inventory, SJ-C1–SJ-C5 |

### Flagged additions and residuals

The owner check is an explicit implementation proposal for the decided single-owner outcome. SJ46–SJ49 make its closure observable.
Baseline fixtures are evidence for the approved compatibility constraint. They add no caller behavior.
Shared error categories support the approved unchanged sentinel contract. They add no new refusal policy.

The two duplicate-tolerant strict readers remain explicit residuals. The spec does not label them migrated or claim universal strictness across all JSON readers.
The durable-file and shell-argument decisions remain closed and separate. No parent map or resolved decision ticket needs a copy into this spec folder.

This planning pass runs prose and coverage checks only. It supplies no implementation, runtime-test, mutation, gate, or performance evidence.

### Completion plan

The plan uses version 1 for the serial author checkpoints. The build records one fresh implementation author per ticket on the accepted implementation line.
Every ticket owns its complete vertical result. Review closes each frozen chunk before the next author starts.

```bench-completion-plan
{"version":1,"chunks":[{"id":"SJ-C1","tickets":["1-policy-document.md"],"verification":[{"id":"policy","command":"bench test --package ./internal/releasepreflight"},{"id":"owner","command":"bench test --package ./internal/jsonfile"},{"id":"mutation","command":"bench probe internal/releasepreflight/vulnerability.go --swap 'jsonfile.DecodeDocument(data, &exceptions)' --with 'jsonfile.DecodeDocument([]byte(\"[]\"), &exceptions)' --package ./internal/releasepreflight --run TestVulnerabilityPolicyDocumentContract","probe":"Replace the real policy bytes with an empty array at the shared decode call. The suffix-refusal case must fail."}]},{"id":"SJ-C2","tickets":["2-jsonfile-gate.md"],"verification":[{"id":"owner","command":"bench test --package ./internal/jsonfile"},{"id":"gate","command":"bench test --package ./internal/gate"},{"id":"census","command":"bench test --package ./internal/conformance --run TestBranchNativeArchitectureCensus"},{"id":"mutation","command":"bench probe internal/jsonfile/decode.go --swap 'seen[name] = true' --with 'seen[name] = false' --package ./internal/jsonfile --run TestDocumentStructuralContract","probe":"Disable the shared duplicate-memory update. The nested duplicate assertion must fail."}]},{"id":"SJ-C3","tickets":["3-review-documents.md"],"verification":[{"id":"review","command":"bench test --package ./internal/reviewrecord"},{"id":"fixture","command":"bench test --package ./internal/reviewrecord/recordtest"},{"id":"mutation","command":"bench probe internal/reviewrecord/parse.go --swap 'jsonfile.DecodeDocument(data, value)' --with 'jsonfile.DecodeDocument([]byte(\"null\"), value)' --package ./internal/reviewrecord --run TestReviewJSONCompatibility","probe":"Replace the real review bytes with null at the shared decode call. The valid review-document assertion must fail."}]},{"id":"SJ-C4","tickets":["4-release-documents.md"],"verification":[{"id":"release","command":"bench test --package ./internal/releaseevidence"},{"id":"mutation","command":"bench probe internal/releaseevidence/component_manifest.go --swap 'if len(object) != len(fields) {' --with 'if false && len(object) != len(fields) {' --package ./internal/releaseevidence --run TestComponentJSONCompatibility","probe":"Disable the component extra-key refusal while preserving all required keys. The extra-key assertion must fail."}]},{"id":"SJ-C5","tickets":["5-ownership-closure.md"],"verification":[{"id":"owner-check","command":"bench test --check strict-json-owner"},{"id":"conformance","command":"bench test --package ./internal/conformance"},{"id":"census","command":"bench test --package ./internal/conformance --run TestBranchNativeArchitectureCensus"},{"id":"mutation","command":"bench probe internal/conformance/registry/checks.go --swap 'Name: \"strict-json-owner\"' --with 'Name: \"strict-json-owner-disabled\"' --package ./internal/conformance --run TestStrictJSONOwnerRegisteredPath","probe":"Rename the registry entry while leaving its executable binding unchanged. The registered-path assertion must fail."}]}],"final_verification":[{"id":"integration","command":"bench test --changed --base 37ac80b8adeff4ecbd57f3d813ff8937f2bcbc68"},{"id":"conformance","command":"bench test --package ./internal/conformance"},{"id":"whole-project","command":"bench gate"},{"id":"coverage","command":"bench coverage --check specs/strict-json-ownership/spec.md"}]}
```

### Ticket slice verification

The ticket graph is serial: SJ-C1, SJ-C2, SJ-C3, SJ-C4, then SJ-C5. Each ticket names its immediate predecessor as a blocker.
The spec fence equals the ticket Writes union apart from implicit planning paths and the review pickup.

The first ticket uses the existing owner. The second ticket supplies shared error categories with its first gate consumer.
Later tickets migrate complete review and release caller families. The final ticket grades sole ownership after those migrations.

The source reads for closure are internal/preflight/closure.go, internal/tickets/registry_data.go, and internal/anchors/references.go.
The anchor query reads projects/benchkit.md and CHANGELOG.md. Their existing anchors remain unchanged.
A clean committed checkpoint permits the authoritative per-ticket writes proposals. A dirty planning source cannot supply that proposal evidence.

This phase records document checks only. It does not start implementation or claim a runtime or mutation result.

### Ticket review result

Independent Sol 6.1/xhigh review accepted all five tickets at `e886c8b1c22b2c2878e341f3e287b9704d00bb6c`.
It found no material Standards, Spec, or Coverage findings.
The review verified the complete fifty-row allocation and exact ownership union.
The final closure includes fourteen canary directories, four anchor files, and five bound registries.

Clean-source preflight passed with fourteen green checks and two not-applicable checks.
All five ticket proposals had no missing paths or ordering requirements.
The required prose lane passed at the planning checkpoint.
These are planning checks; runtime preservation, baseline capture, and mutation reds remain implementation obligations.
