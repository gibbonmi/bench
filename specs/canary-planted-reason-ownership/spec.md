# Canary planted-reason ownership

Status: implemented

Decision source: specs/canary-planted-reason-ownership/decisions/canary-planted-reason-ownership.md

## Problem

`bench canary` validates a fixture inventory but reports that it dispatched owners,
while ordinary tests directly prove only part of the retained kit fixture set. A linked
repo created by `bench init` receives a seed fixture and prose claiming the inventory
command proves the example check bites, even though Bench has no language-agnostic way
to execute that project check. Current-state ADRs still describe the retired nested
canary sweep. Contributors therefore cannot tell which surface proves inventory shape,
which surface proves a planted reason, or which guarantees a linked repo actually has.

## Solution

Make the contract match the branch-native architecture. `bench canary` and the ship
check validate a non-empty inventory with one accepted check binding per fixture and
report only that fact. Ordinary kit tests derive every retained fixture from that
inventory, materialize its mutation, call its exact registered in-process owner, require
the fixture's non-empty `EXPECT` text in the produced diagnostics, restore the subject,
and require that text to disappear. `bench init` creates no example fixture and makes no
planted-reason claim; linked projects own such proof in their native gate tests.
Current-state documentation records this split and deprecates the two nested-sweep
performance decisions.

## User stories

1. As an operator running `bench canary` or the ship inventory check, I want a truthful,
   deterministic inventory verdict, so success cannot be mistaken for owner execution
   or a planted-reason proof. Line: gpt-5.6-terra / medium. The public and ship surfaces
   have known seams, but gate-facing vocabulary and fail-closed validation merit the
   profile's mid-tier conformance line.
2. As a Bench maintainer weakening any retained kit check, I want the ordinary gate to
   turn red for that fixture's own planted reason and return green after restoration, so
   no retained `EXPECT` is ceremonial. Line: gpt-5.6-terra / high. This is a quantified
   oracle repair whose cheapest failure is a silently omitted fixture, so the mid tier
   needs a high-effort falsification pass.
3. As a maintainer initializing a linked repo, I want Bench to scaffold only guarantees
   it owns and leave the repo visibly unconfigured, so project-native tests—not a fake
   generic canary journey—establish whether project checks bite. Line: gpt-5.6-terra /
   medium. The behavior is exact and the existing init seam is known, but it changes the
   fail-closed scaffold contract.
4. As a contributor arriving cold, I want the profile, platform reference, ADRs, public
   help, shipped gate-authoring guidance, compatibility prose, and changelog to describe
   the same current canary contract as the code, so I do not design against a retired
   nested sweep. Line: gpt-5.6-sol / high. Project documentation steers future kit work
   and the profile binds every cold session, so the documented leverage override applies.

## Implementation decisions

- `canary inventory` and `planted-reason proof` retain their meanings from `CONTEXT.md`.
  Inventory is the non-empty set of accepted fixture-to-check bindings. A
  planted-reason proof is a direct mutation test through one exact registered owner,
  including restoration. Inventory is never described as dispatch, invocation, or
  sweep; a planted-reason proof explicitly invokes its registered owner.
- The existing in-process inventory decision remains the production seam. Fixture
  discovery feeds one selection/validation result consumed by the public command and
  ship check. The production callback-based `Dispatch`/`DispatchResult` layer and its
  dispatched-result field leave because no production caller may execute an owner.
  The complete current production package has no other function-typed parameter, so no
  production function or method in `internal/canary` retains or introduces one. Internal
  names below the decision may change, but there is no vestigial invocation callback or
  second inventory derivation.
- A successful public command prints exactly
  `canary inventory ok (<n> fixture bindings)`. `<n>` is the accepted binding count
  derived from the current fixture producer, never a stored constant. The existing root
  positional grammar, help behavior, and exit classes stay intact. Invalid, duplicate,
  empty, unsafe, unreadable, or unbound inventory remains fail-closed.
- Binding has two producer-derived partitions. A kit conformance family resolves to its
  registered check; a linked project's non-empty two-level family resolves to a
  `project:<family>` inventory identity and carries no execution claim. A flat fixture
  with neither an explicit `CHECK` nor a family has no binding and is red. At kit scope,
  the ordinary conformance family-registry check and the direct-proof owner resolution
  additionally reject a family that is not registered; the public inventory command
  does not misrepresent a linked-project family identity as a kit check owner.
- The ship-tier canary step calls the same inventory validator as the public command and
  starts no fixture owner, gate, wrapper, `go test`, or `go run`. The selected-executable
  system journey proves routing to that production inventory path and its complete
  aggregate only.
- The ordinary planted-reason proof derives its universe from `canary.Fixtures`; test
  code carries no independently maintained fixture census. For each produced fixture it
  resolves the exact non-meta check and its declared tier from the registry, reads a non-empty
  normalized `EXPECT`, materializes the mutation in a temporary subject, calls the
  registered in-process owner, checks that the expectation occurs in its diagnostics,
  restores the subject, calls the same owner again, and checks that the expectation is
  absent. A missing or ineligible owner is a red, never a skip.
- The current producer yields 184 retained fixtures, but the producer-derived proof
  runner completes only 181. `release-digest-binding-omitted` and
  `release-package-evidence-omitted` reach the ship-tier `release-evidence-probe`, whose
  authenticated-clone precondition fails on the materialized fixture root before either
  fixture's `EXPECT` diagnostic. `unrouted-subcommand` cannot materialize because its
  `MUTATE.json` names a removed command-map line instead of the current
  `commandRegistry` entry. This observed 181/184 partition, not the source-only owner
  census, is the implementation baseline.
- Apply decision #3's retirement alternative to exactly the two release fixtures.
  `release-digest-binding-omitted` becomes an owning-package mutation test over the
  canonical release-index encoder: the encoded artifact row must carry the non-empty
  component-manifest digest under `component_manifest_sha256`. The existing JSON-tag
  mutation must make that test red. `release-package-evidence-omitted` becomes an
  owning-package test that invokes the existing `build-release-evidence.mjs` seam against
  temporary wrapper/platform staging and requires the registry-derived package evidence
  in the wrapper output. Removing the wrapper's `copyPackageEvidence` call must make that
  test red before any full artifact, authenticated clone, or release preflight journey.
  The two fixture directories and their then-unreferenced shared base asset leave the
  producer in the same atomic landing as their replacement tests; no semantic tripwire
  is dropped or reduced to source-text matching.
- `unrouted-subcommand` remains a retained fixture with the `subcommand-routing` owner and
  reanchors its mutation to the current `commandRegistry` declaration. Its control repair
  lands atomically with the universal owner proof.
- The accepted retained producer therefore has 182 fixtures. The direct test calls each
  retained fixture's exact registered owner in-process at its declared registry tier.
  The independently authored current count in the public-command expectation is 182 and
  remains the omission grader. The initial 31-member gap partitions into the two
  owning-package migrations and 29 retained missing direct comparisons. The 26 retained
  `package-core-guard` members are
  `bounds-duplicate-owner`, `native-reproducibility-handoff-dropped`,
  `native-smoke-workflow-dropped`, `native-trigger-comment-spoof`,
  `offline-network-repair-allowed`, `offline-registry-fallback-allowed`,
  `offline-slice1-operation-omitted`, `offline-stage-interruption-ignored`,
  `preflight-native-call-bypassed`, `preflight-native-upload-bypassed`,
  `preflight-publish-ancestry-omitted`, `preflight-publish-changelog-omitted`,
  `preflight-publish-identity-omitted`, `preflight-publish-needs-bypassed`,
  `preflight-publish-order-bypassed`, `preflight-release-call-bypassed`,
  `preflight-verify-analysis-omitted`, `preflight-verify-artifact-omitted`,
  `preflight-verify-gate-omitted`, `preflight-verify-smoke-omitted`,
  `preflight-verify-vulnerability-omitted`, `reintroduced-bare-skip`,
  `release-future-owner-omitted`, `release-public-profile-omitted`,
  `reproducibility-byte-compare-bypassed`, and `unrouted-subcommand`. The two
  `compliance-hardening` members are `missing-license` and
  `mutable-workflow-action`; the `injected-ports` member is `unregistered-port`.
  This baseline enumeration is review evidence, not a new executable registry.
- The accepted candidate proves each retained fixture through its registered
  owner and keeps the complete 182-member retained producer. Inventory-only retention, a
  hard-coded sibling diagnostic, an empty-tree collision screen, fixture retirement, and
  reliance on the removed sweep are not accepted proof in this candidate; only the two
  exact owning-package migrations above are authorized retirements.
- `bench init` stops writing the example fixture, `DO-NOT-SHIP` marker, example check,
  and planted-reason claims. Its generated gate retains the configuration sentinel and
  its existing repo-local-before-PATH resolver, then invokes `"$bench" canary "$root"`
  only for inventory validation. A new linked repo is red until its project adds real
  checks and, where desired, real fixture bindings plus project-native planted-reason
  tests. Re-running init never recreates the retired seed and preserves project-owned
  checks and fixtures.
- The shipped `.bench/lib/canary-run.sh` filename remains as a compatibility surface for
  already-linked gates, but its comments and failure diagnostic describe inventory
  validation. It calls only the public inventory command.
- ADR 0001 and `craft-gate` state the two-layer current contract: direct ordinary owner
  tests defend kit fixtures, while linked repos receive inventory validation and own
  their native planted-reason proof. ADRs 0003 and 0009 are rewritten as current-state
  records with `Status: deprecated`; they retain no operative nested-sweep latency or
  concurrency policy. The platform reference removes the authenticated inner-canary and
  canary-phase-selector claims from its phase-manifest description. Public wrapper help,
  the profile, and changelog use the same terms without storing a fixture count.
- Ticket slicing may sequence inventory vocabulary, direct fixture proof, linked-repo
  scaffolding, and current-state documentation, but every ticket lands independently
  green and states no guarantee stronger than its landed code proves. The fixture-proof
  ticket lands the two exact migrations first; the truthful-inventory ticket is blocked
  by it because only that tree has the accepted 182-binding aggregate. The final accepted
  candidate closes the complete contract; splitting it into separately shippable specs
  is not authorized.
- Bootstrap authority is unchanged. The tagged system journey enters through the
  already-selected Bench executable, and the inventory command launches no successor.
  No path, fixture record, callback, or candidate-controlled executable is treated as an
  authority for execution.

## Testing decisions

- TDD attaches at three existing seams: the in-process inventory decision plus public
  command, the ordinary `runFixtureBite` owner seam, and the real `Init` scaffold. The
  two already-working release tripwires attach regression tests at their existing
  owning-package encoder and focused builder seams and use mutation probes rather than a
  manufactured initial red. No new production test interface is introduced.
- Inventory decision tests exercise single and complete producer sets, invalid/duplicate
  bindings, unsafe names, absent/empty trees, malformed binding markers, dangling
  symlinks, and special files. An empty-tree diagnostic states only that the fixture
  inventory has no accepted bindings; it makes no claim that a gate did or did not prove
  checks bite. A FIFO or other special `EXPECT` must be rejected by metadata before any
  content read can block. Public command tests assert exact output and exit, while the one tagged
  system journey guards selected-executable routing and the composition of count plus
  wording.
- One producer-derived ordinary test owns the universal planted-reason assertion. A
  test-internal universe runner accepts a canary directory—not fixture names—and derives
  its cases by calling `canary.Fixtures` before invoking a supplied proof callback. The
  live universal test is a thin wrapper that passes the kit canary directory and a proof
  callback to that runner. The callback calls `runFixtureBite` and records an identity only
  after that proof returns; the test then compares the callback-observed identity set with
  an independently fetched `canary.Fixtures` key set. A derivation bite test passes a
  synthetic tree containing both a family-bound dev fixture and an explicit `CHECK`-bound
  `release-evidence-probe` fixture absent from the live repository, and requires the same
  runner callback to observe both. A separate resolution-only assertion passes the
  synthetic ship fixture to `resolveFixtureBite` and requires the registered check plus
  `registry.Ship`; it never invokes the supplied callback or the owner. A test-source
  shape grader rejects fixture-name
  literals, a second proof loop, recording before proof, or any live case source other
  than the one runner call and producer-key comparison. Replacing the runner's producer
  call with a stored list must make the synthetic test red; filtering CHECK-bound or
  ship-tier fixtures must make the synthetic and live set-equality assertions red; and
  replacing the live call with a hand list must make the shape grader red. Its live subtests
  use fixture identity for attribution and drive the exact registered check. The existing
  helper remains the behavioral shape: expectation present after mutation, expectation
  absent after restoration. Adding a retained fixture automatically adds a required
  subtest; no manual fixture-name list can make the gate green.
- The release-index owning-package test encodes production `artifactEvidence`, decodes the
  resulting JSON independently, and asserts its non-empty component-manifest digest under
  the contract key; it does not inspect Go source. Its exact omission-graded symbol is
  `TestReleaseIndexBindsComponentManifestDigest`. The package-evidence owning-package
  test follows the existing Node-process prior art in `internal/releaseevidence`, prepares
  only the temporary staging the production script consumes, runs that script directly,
  and asserts registry-derived evidence in the wrapper package. It does not call
  `build-artifacts.sh`, authenticate a clone, or build/rehearse release artifacts. The
  exact symbol is `TestBuildReleaseEvidenceIncludesRegisteredPackageEvidence`. The
  historical fixture mutations are the required self-probes for these replacements.
- One independently authored retired-fixture replacement census in
  `internal/conformance` is the omission grader for the two migrations. It parses the
  owning-package test declarations and requires each exact replacement test symbol in its
  named file, while also requiring both fixture directories and their now-unreferenced
  shared base asset to be absent. Deleting or renaming either replacement test must make
  the census red. This inventory stores only the two migration obligations needed to
  grade omission; release behavior remains derived and asserted at the owning-package
  seams above.
- The `unrouted-subcommand` fixture materializes its current `commandRegistry` mutation
  and produces its existing routing diagnostic. A source-only owner census or a callback
  record made before successful materialization is not accepted evidence for it.
- Fixture-resolution tests cover zero-byte, ASCII-whitespace-only, and
  non-ASCII-whitespace-only `EXPECT` values as distinct refusal cases. The central
  mutation that replaces normalization with a raw non-empty byte check must make both
  whitespace cases red.
- During implementation, the central-property mutation replaces one previously uncovered
  fixture's `EXPECT` with a unique sentinel and requires the owning package test to fail
  naming that fixture and sentinel. The independent coordinator probe uses a different
  fixture and mutation kind. Both restore source bytes before any landing.
- Init tests drive the real scaffold in a repository path containing spaces and glob
  characters, assert the seed and example check are absent, assert the configuration
  sentinel plus repo-local-before-PATH resolution and quoted `"$bench" canary "$root"`
  invocation remain, and re-run after project-owned checks and fixtures exist to prove
  they are preserved.
- The ordinary `go test -count=1 ./...` gate phase observes the direct proofs. The
  branch-native architecture census adds `internal/conformance/fixture_bite_test.go` to
  its exact `directArchitectureTests` set and rejects a generic process or repository
  constructor there, including nested gates, wrappers, per-fixture `go test`, and
  `go run`. The universal runner and its live wrapper stay in that censused file. The
  ship owner's existing implementation files remain outside that direct-harness set.
  Declared-tier resolution is covered by the synthetic resolution-only assertion; no
  retained live fixture or synthetic proof callback enters the authenticated ship
  journey. No focused fixture/family execution surface is added.
- Skill/ADR/profile semantics are not self-grading. Exact stale-vocabulary searches
  catch retained retired terms in production and in conformance helper names, test names,
  comments, and diagnostics, and a fresh Fable cross-harness review re-derives the decision
  source and checks the final prose. The public wrapper help is exercised as real output,
  and the changelog records the user-visible command and scaffold change under the
  existing Unreleased sections.

### Seam diagrams

    trigger: bench canary or ship inventory validation
        │
        ▼
    fixture tree ──▶ [ discover → select accepted bindings ] ──▶ exact count or refusal
                                  ◀ tests attach here: decision tests and real CLI/system output

    trigger: ordinary Go mutation test
        │
        ▼
    produced fixture ──▶ [ materialize → exact registered owner → restore ] ──▶ planted red, restored green
                                      ◀ tests attach here: runFixtureBite over the producer-derived inventory

    trigger: bench init in a linked repo
        │
        ▼
    empty/project tree ──▶ [ Init scaffold ] ──▶ red configuration sentinel + inventory-only gate call
                                  ◀ tests attach here: real scaffold output, files, and idempotent re-run

### Acceptance coverage map

| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| CI1 | 1 | A successful public run prints exactly `canary inventory ok (<n> fixture bindings)`, with `<n>` derived from the accepted inventory; after the two exact owning-package migrations the retained kit reports 182. | real `bench canary` command | observed red: `test "$(go run -buildvcs=false ./cmd/bench canary .)" = "canary inventory ok (182 fixture bindings)"` exited 1 | Exact output and candidate producer count reject the current dispatch claim, a stored count, partial aggregation, and retention of the two false-proof fixtures. |
| CI2 | 1 | No production function or method in `internal/canary` declares any function-typed parameter, and inventory validation exposes no dispatch result or dispatched-result field. | in-process canary API-shape test plus decision behavior | observed red: `! rg -n '^func [^(]+\([^)]*func\(' internal/canary --glob '*.go' -g '!*_test.go'` exited 1 on the only current callback; the implementation assertion parses every production `FuncDecl`, including methods, and rejects every function-typed parameter | The complete-package shape assertion rejects a callback consuming a binding, fixture, name, or renamed wrapper, plus method/unexported executors and result-field-only wording repairs; the live census shows no legitimate callback seam is banned. |
| CI3 | 1 | Empty selection, duplicate fixture identity, a fixture value with no binding, unsafe binding, and control-bearing identity all refuse deterministically; the empty-tree refusal says only that no accepted inventory binding exists and claims no gate proof. | inventory selection decision | observed red: the empty-tree path currently emits `the gate cannot prove its own checks bite`; existing `TestSelectCanaryOwnersFromImmutableFixtures` covers the empty, duplicate, missing-owner, and hostile refusal classes | The exact selection partitions prevent an unresolved or ambiguous fixture from satisfying the count-only happy path, while the diagnostic assertion catches retained execution-proof vocabulary on the same repaired path. |
| CI4 | 1 | A missing, dangling-symlink, or special-file fixture marker fails closed by metadata before a content read; a FIFO `EXPECT` cannot block or count as a fixture. | fixture discovery and control-record classification | observed red: a synthetic one-fixture tree with FIFO `EXPECT` made `go run -buildvcs=false ./cmd/bench canary <root>` print success, so the refusal probe exited 1 | Metadata-first refusal catches the hostile file types that the current `os.Stat` presence check admits without making inventory execute or interpret a planted reason. |
| CI5 | 1 | The ship step, compatibility shim in both resolved- and missing-Bench branches, unit tests, and tagged selected-executable journey all consume the same inventory decision and use inventory-only vocabulary and diagnostics. | ship library call plus compatibility and real selected-executable system seams | observed red: the current source sweep finds `SweepShip`, two `canary sweep failed` branches, `TestCanaryDispatchAndAggregation`, and `fixture owners dispatched` | Cross-consumer assertions catch a locally correct command whose ship or shipped compatibility path still claims execution, including when the required Bench command is missing. |
| CI6 | 1 | The existing optional root positional, `help`/`--help`/`-h`, too-many-argument usage, and exit 0/1/2 classes remain unchanged, while top-level `bench` help describes fixture inventory validation rather than running a gate. | command grammar, command registry route, and real wrapper help | observed red: top-level `bench` help currently says `bench canary [root]` will `run the gate against known-broken fixtures`; `TestRunCanaryDispatchesToCommand` and shared usage grammar tests already cover the unchanged subcommand grammar | Real wrapper output catches a repaired subcommand whose primary command list still advertises the retired sweep. |
| PB1 | 2 | Every fixture returned by the accepted 182-member retained producer, including the 29 retained members of the initial gap, receives one completed subtest through its exact registered non-meta owner at that owner's declared tier. | producer-derived ordinary fixture-bite test plus exact public aggregate | observed red: the producer-derived runner completed 181/184; after removing the two release fixtures, `unrouted-subcommand` still fails materialization on its removed anchor, so only 181/182 complete; deleting any other fixture violates CI1's independently expected 182 aggregate | The universal proof makes omission, record-before-proof, stale-control, and wrong-owner execution red, while the independent candidate aggregate makes a shrink-and-ignore repair red without creating a production registry. |
| PB2 | 2 | Adding either a family-bound dev fixture or an explicit `CHECK`-bound ship fixture automatically creates a required planted-reason subtest; no manual fixture list, tier filter, binding-form filter, or stored fixture count controls completeness. | test-internal universe runner from canary directory through `canary.Fixtures` to proof callback | identified red: a synthetic tree adds one registered-family fixture and one `CHECK: release-evidence-probe` fixture, both absent from the repository, and requires the shared runner callback to observe both; replacing the producer call with today's literal names or filtering ship/CHECK bindings must fail | The two producer extensions grade derivation and both binding partitions, so freezing all 182 retained names or quietly excluding a future ship owner cannot satisfy the row. |
| PB6 | 2 | The live universal proof consumes the same runner once, records an identity only after `runFixtureBite` returns for it, and asserts that observed set equals a fresh `canary.Fixtures` key set, with no fixture-name literals, manual case source, or second proof loop. | live callback-observed set equality plus AST shape of the universal test and shared runner | identified red: filter CHECK-bound or ship-tier entries inside the runner, or replace the live runner call with the current hand-listed passing subset while leaving PB2's synthetic runner test intact; set equality must fail on the filter and the source-shape grader must fail on the hand list or record-before-proof shape | This couples the producer-derived cases to completed live proofs and makes any observed subset—including 181—distinguishable from the accepted 182 even when inventory and the synthetic helper remain correct. |
| PB3 | 2 | For each retained fixture, its normalized non-empty `EXPECT` occurs after materialization and is absent after restoration through the same owner. | `runFixtureBite` behavioral seam | observed red: the `bounds-duplicate-owner` sentinel remained unobserved because no current subtest reaches that fixture; existing direct subtests already demonstrate the paired predicate for their subset | The paired assertion rejects collision-only checks, hard-coded sibling diagnostics, and a test that never proves restored green. |
| PB4 | 2 | An absent, zero-byte, ASCII-whitespace-only, or non-ASCII-whitespace-only `EXPECT`, missing fixture, unbound or unregistered check, or meta check fails with fixture attribution rather than skipping; resolution returns a registered owner's declared dev or ship tier before invocation. | fixture-bite resolution | identified red: extend `TestFixtureBiteResolutionRefusesInvalidInputs` with ASCII and U+2003-only files, then replace normalization with a raw non-empty byte predicate; both new cases must fail. Replace the current dev-only refusal with a resolution-only synthetic `CHECK: release-evidence-probe` assertion that expects `registry.Ship`; retaining the dev hard-code must fail without invoking that owner | These failure partitions prevent an ineligible fixture from disappearing from the universal claim, independently grade Unicode-aware normalization, and make declared-tier resolution red-capable without an authenticated ship journey. |
| PB5 | 2 | Direct proof invokes only the registered owner as an in-process function and introduces no generic repository/process constructor, gate, wrapper, `go test`, or `go run` in `internal/conformance/fixture_bite_test.go`; an owner may perform the work already assigned to its declared tier. | exact `directArchitectureTests` membership plus ordinary conformance owner | identified red: `TestBranchNativeArchitectureCensus` currently omits the fixture-proof file, so adding an `exec.Command` there stays green; add that exact file to the census set, then the same constructor must be reported while existing release-owner files remain accepted | Per-file census ownership catches restoration of the retired nested sweep in the harness without prohibiting an exact registered owner's native implementation. |
| PB7 | 2 | The two false-proof release fixtures and their shared base asset leave the canary producer only with omission-graded owning-package mutation tests preserving their semantic tripwires: canonical index encoding binds the component-manifest digest, and focused production builder staging includes registry-derived package evidence in the wrapper. `unrouted-subcommand` mutates the current `commandRegistry` entry and reaches its existing diagnostic. | `internal/releaseevidence` encoder and focused builder seams, retired-fixture replacement census, plus real fixture materialization through the shared bite helper | not TDD-able: both release behaviors already exist but lack owning-package regression tests; applying their historical JSON-tag and omitted-copy mutations must make the replacements red. Identified omission red: deleting or renaming either replacement test must make the independently authored census fail. Observed baseline: the universal runner completes 181/184 because both release fixtures fail before `EXPECT` and `unrouted-subcommand` has no live anchor | The exact migration/control partition preserves the release semantics and grades replacement omission while preventing an authenticated ship journey, source-text behavior check, or stale mutation control from being counted as direct proof. |
| LR1 | 3 | `bench init` writes no example fixture, `EXPECT`, `DO-NOT-SHIP` marker, example check, or output claiming that Bench proved a project check bites. | real Init scaffold | observed red: a real init probe exited 1 because it created `tests/canary/example/example/EXPECT` and printed the seed-canary claim | File and output assertions catch both the false evidence and prose-only removal of the claim. |
| LR2 | 3 | The generated gate keeps the fail-closed configuration sentinel and repo-local-before-PATH resolution, then invokes `"$bench" canary "$root"` only for inventory validation and is red on absent/empty inventory until the project configures real checks and bindings; a project-owned two-level family remains an accepted `project:<family>` inventory binding without becoming an executable owner. | generated gate text plus real inventory command | observed red: `! rg -n -e 'const seedCanaryPath' -e DO-NOT-SHIP -e 'seed canary' -e 'canary sweep failed' internal/adopt/init.go` exited 1; current `TestScaffoldGateUsesCanarySubcommand` covers the resolver and current `TestInitScaffoldsTwoLevelSeedCanary` covers the project-family binding that their replacements preserve | The combined assertion preserves selected executable routing, the only linked-project binding form, and a surviving inventory caller while refusing a disguised fixture journey or silently green empty scaffold. |
| LR3 | 3 | Re-running init in a path containing spaces and glob characters never recreates the retired seed and preserves project-owned checks and fixture bindings. | Init re-run against serialized scaffold state | observed red: after moving the current seed aside, a second init in `/tmp/bench init [x].*` recreated it and the absence probe exited 1 | A second real invocation catches a first-run-only deletion and unquoted-root regressions without inventing a generic owner protocol. |
| DOC1 | 4 | ADR 0001, `craft-gate`, and the project profile state that kit ordinary tests own direct planted-reason proof, while linked repos receive Bench inventory validation and own native proof. | current-state documentation and shipped guidance reviewed against code and decision map | not TDD-able: semantic truth spans prose and code; the exact-candidate Opus review re-derives all three, while stale-term searches provide only a deletion tripwire | Re-derived semantic review catches polished wording—or future gate-authoring guidance—that still assigns execution or a seed proof to the inventory command. |
| DOC2 | 4 | ADRs 0003 and 0009 carry `Status: deprecated` and describe no live nested-sweep latency, worker, inner-gate, or `GOMAXPROCS` policy. | ADR current-state review | observed red: the stale-ADR negation probe exited 1 on `canary sweep`, `complete inner gate`, and the old gate-runs-itself claim | Exact retired-mechanism terms make the stale operational decisions visible; review checks the replacement is a current state rather than history. |
| DOC3 | 4 | Production comments, result fields, conformance helper/test rationale, system expectations, public wrapper help, compatibility prose, shipped gate-authoring guidance and platform reference, profile, ADRs, and the Unreleased changelog use inventory/proof terms consistently and store no independent fixture count beyond the omission-grading test expectation. | repository-wide current-claim census plus semantic review | observed red: the current census finds dispatch/sweep claims in canary, conformance helper names/comments/diagnostics, system, adopt, pre-release, wrapper, compatibility, `craft-gate`, `.bench/BENCH-reference.md`, profile, and ADR surfaces | The cross-surface census catches a composition in which each code path passes locally while one current-state consumer or direct-proof test still reasons from the old sweep. |
| CC1 | 1, 2, 3, 4 | The final candidate simultaneously reports the 182-binding inventory only, directly proves that complete retained producer, preserves the two migrated release tripwires in owning-package tests, scaffolds no false linked-repo proof, and documents exactly those guarantees. | real CLI and init seams plus owning-package release tests, complete ordinary gate, and exact-candidate review | observed red: the public vocabulary probe, linked scaffold probe, 181/184 runner, and 31-member source census fail on the current tree | This is the composition-degenerate row: per-fence tests cannot justify a final candidate whose code, scaffold, tests, and docs describe different contracts. |

### Edge inventory

- Error path — CI3, CI4, PB4, and LR2 cover invalid inventory, unsafe control
  records, ineligible owners, and the intentionally red unconfigured scaffold.
- Empty or absent input — CI3 refuses empty inventory; PB4 refuses missing/empty
  expectations; LR2 keeps an initialized but unconfigured repo red.
- Boundary values — CI1 covers one derived count contract and the complete accepted
  182-member aggregate without storing either count in production; PB1 uses that
  independently authored expectation to reject retirement of a current member.
- Malformed input — CI3, CI4, and PB4 cover malformed markers, control-bearing
  identities, special files, dangling links, empty `EXPECT`, and invalid owner classes.
- Interrupted or partial state — **Won't handle:** init transactionality is unchanged and
  this repair removes writes; it adds no durable execution or fixture mutation outside a
  test-owned temporary subject.
- Re-run idempotency — CI1 requires deterministic inventory output, PB3 repeats the owner
  after restoration, and LR3 drives a fresh process over serialized init state.
- Process-boundary lifecycle — CI5 uses the selected-executable system journey and LR3
  uses a second init process; PB5 keeps each planted-reason proof in-process.
- Hostile environment — CI4 covers special files and dangling links; LR2 preserves the
  generated linked-repo by-path CLI before global PATH; LR3 covers spaces and glob
  characters; CI3 covers control bytes; CI5 covers the compatibility shim's resolved and
  missing-Bench branches plus selected-executable system routing. TTY input, network,
  signals, worktree destruction, and host `fsync` pressure are **Won't handle** because
  none is consumed or introduced by these non-interactive, network-free, bounded
  inventory and in-process test surfaces.
- Hand-edited line endings and whitespace — PB3 accepts a non-empty `EXPECT` with or
  without a trailing newline after normalization; PB4 adds explicit ASCII and U+2003-only
  refusal cases plus a raw-byte-predicate mutation, rather than treating the inherited
  zero-byte case as evidence for whitespace handling.
- Command self-write — LR1 and LR3 assert post-write truth and repeated application in a
  tracked-style repository; the read-only inventory command changes no reported fact.
- Argument grammar and quoting — CI6 preserves `--`/help/arity behavior and LR2/LR3 keep
  the multi-word root as one quoted argument.

## Ownership fences

- Inventory command and decision: `internal/canary`
- Ship inventory consumer: `internal/preprelease`
- Selected-executable system expectation: `internal/systemtest`
- Top-level command route test: `cmd/bench/main_test.go`
- Public wrapper help: `bin/bench.sh`
- Shipped compatibility shim: `.bench/lib/canary-run.sh`
- Shipped gate-authoring guidance: `.agents/skills/bench-craft-gate/SKILL.md`
- Shared platform reference: `.bench/BENCH-reference.md`
- Direct planted-reason proof and fixture classification: `internal/conformance`
- Release-tripwire owning-package replacements:
  `internal/releaseevidence/release_index_test.go` and
  `internal/releaseevidence/package_artifact_test.go`
- Exact false-proof fixture changes:
  `tests/canary/package-core-guard/release-digest-binding-omitted/`,
  `tests/canary/package-core-guard/release-package-evidence-omitted/`,
  `tests/canary/package-core-guard/release-evidence-probe-base.txt`, and
  `tests/canary/package-core-guard/unrouted-subcommand/MUTATE.json`
- Linked-repo scaffold and its tests: `internal/adopt`
- Current profile: `projects/benchkit.md`
- Current canary ADR: `docs/adr/0001-working-tree-gate-tripwire.md`
- Deprecated latency ADR: `docs/adr/0003-canary-latency-budget.md`
- Deprecated concurrency ADR: `docs/adr/0009-canary-concurrency-budget.md`
- User-visible release note: `CHANGELOG.md`

These are the complete implementation write envelopes. They authorize no edit to
`.bench/gate.sh`, `internal/gate`, `ROADMAP.md`, `CONTEXT.md`, another decision map or
spec, assessment/capture artifacts, or assignment worktrees. Reviewer disposition:
pre-approved for this staged spec; implementation still derives independently-green
tickets inside these fences.

## Out of scope

- A production fixture-journey owner, nested gate, wrapper, per-fixture `go test`, or
  per-fixture `go run`: rejected architecture, not deferred remainder. Reintroducing it
  would be a separate 7-edit, 3-gate-run capability and requires a new reviewer decision.
- FT168 focused fixture/family execution: separate 5-edit, 2-gate-run capability after
  FT153 lands and its premise is revalidated. It must select direct ordinary proofs
  without becoming a second oracle or receiving gate credit.
- Gate performance pricing: separate 3-edit, 2-gate-run assessment and repair after the
  correctness owner is complete; this spec adds no concurrency or latency policy.
- FT153/FT168 roadmap reconciliation: workflow-owned maintenance, 1 edit and 1 gate run,
  not part of the canary behavior candidate.
- Retiring or replacing any of the remaining 182 fixtures: separate 4-edit, 2-gate-run
  migration covering the fixture, an equivalent owning-package mutation test, its
  omission grader, and the public aggregate/current-state claim. This candidate already
  applies that complete migration to the two named release fixtures and CI1 deliberately
  holds the resulting aggregate; another failed fixture must reopen decision #3's
  alternative before its control enters an ownership fence.
- A generic language-agnostic linked-project check protocol: separate 6-edit, 3-gate-run
  product capability with no current trust or execution contract. The surviving linked
  caller is inventory validation.
