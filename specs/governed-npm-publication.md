# FT83 slice 3 — governed npm publication

Status: implemented

Decision map: `decisions/governed-offline-release-bundle.md`; publication research
`decisions/governed-offline-release-publication-research.md`.

## Problem

Bench can now build, reproducibly compare, and offline-verify the complete immutable
release set — five npm tarballs and four per-target archives — and the authoritative
preflight indexes them under one deterministic release index. What it cannot do is
publish that set to npm under governance. There is no state machine that stages every
name, waits for platform dependencies, publishes the wrapper last, verifies live
integrity before promoting, and rolls back by deprecation rather than unpublish. There
is no durable publication record bound to the approved release, no reviewer-authorized
handoff for interactive 2FA approval, and no fail-closed proof that a selected release
profile's external-owner records exist before real publication is allowed. Until this
slice lands, "governed release" stops at the artifact boundary and the FT83 row cannot
retire.

## Solution

Add one deep publication module owning a resumable npm state machine and an npm-registry
port. The state machine builds and locally verifies the complete immutable set, submits
every package under a version-specific candidate tag, verifies or approves platform
packages first and the wrapper last, reverifies the complete live set, then promotes
platform `latest` tags first and the wrapper last. First publication uses
reviewer-authorized direct publish because the five identities do not yet exist; later
releases use stage-only OIDC submission with interactive 2FA approval. Retries accept an
already-live package only when registry integrity exactly matches the approved local
tarball; after any package is live, failure preserves the old `latest`, removes candidate
tags, deprecates the bad version, and requires a new version — automation never
unpublishes.

The port has two adapters: the existing hermetic registry fixture is the gate adapter and
public npm is the runbook adapter. Publication produces one durable, non-deterministic
`publication-record.json` that references the immutable release-index digest and carries
ordered registry transitions, stage IDs, authentication mode, tag and integrity
observations, timestamps, result, and provenance digests. It never rewrites the
reproducible index or the package bytes and it is not a second deterministic manifest.
The authoritative preflight remains the only publication authority: publish mode requires
an explicit profile and stays red for any selected profile whose required FT87, FT88, or
(bank) FT71 producer record is absent. New `bench release
prepare|submit|status|promote|rollback` operations are idempotent, non-interactive, emit
compact TOON with an exact `next_action` for external approval, and never ingest
credentials into evidence.

## User stories

1. As a release maintainer of still-unpublished identities, I want first publication to
   direct-publish the four platform packages under a version-specific non-default tag and
   the wrapper last, verifying each live integrity before advancing, so that a partial or
   out-of-order publication cannot leave a default-installable release. Line:
   `gpt-5.6-luna` / low. Ordering, tag, and per-package integrity checks are fully
   observable at the registry port through the hermetic fixture.

2. As a maintainer of an existing identity, I want subsequent releases to stage every
   package under a version-specific non-default tag with stage-only OIDC and to hand off
   an exact `next_action` for interactive 2FA approval of platform packages first and the
   wrapper last, so that automation submits but only a present maintainer approves. Line:
   `gpt-5.6-luna` / medium. The staged path cannot run against real npm in the gate, so
   its submit/wait/approve logic is exercised through the fixture and the emitted handoff.

3. As a release operator resuming after a partial success, I want an already-live package
   accepted only when its registry integrity byte-matches the approved local tarball and
   any mismatch to stop the release, so that a retry never advances past substituted or
   drifted bytes. Line: `gpt-5.6-luna` / medium. Resume, duplicate-transition, and
   mismatch classification are subtle but fully fixture-observable.

4. As a release maintainer, I want promotion to move platform `latest` tags first and the
   wrapper `latest` last, and only after the complete live set reverifies, so that no
   consumer can resolve a wrapper whose platform dependency is not yet default. Line:
   `gpt-5.6-luna` / low. Promotion order and gating on complete reverification are
   observable as ordered dist-tag calls.

5. As a release maintainer of a failed post-live release, I want rollback to remove
   candidate tags, preserve the prior `latest`, and deprecate the bad version with a
   recovery message while never unpublishing, so that a broken release is retired without
   burning recovery state or pretending an immutable version can be removed. Line:
   `gpt-5.6-luna` / low. Deprecate-and-tag-restore calls and the absence of any unpublish
   call are directly assertable.

6. As a release reviewer, I want one durable `publication-record.json` bound to the
   immutable release-index digest that carries every ordered transition — package,
   version, local and registry integrity, stage ID, authentication mode, tag state,
   timestamp, result, and provenance digests — and contains no credential material, so
   that the publication history is auditable from one ledger without a second manifest.
   Line: `gpt-5.6-luna` / low. Exact-field and no-secret assertions run against the
   promoted record.

7. As a bank-track reviewer, I want publish mode to require an explicit profile and stay
   red for any selected profile whose required FT87, FT88, or (bank) FT71 producer record
   is absent, so that the machinery ships now without ever declaring the repository
   release- or bank-ready before its producers exist. Line: `gpt-5.6-luna` / low. The
   requirement registry already encodes profile requiredness; missing-record and
   conforming fixtures make the fail-closed verdict red or green.

8. As the reviewer, I want behavior contracts driving the built `bench release` CLI
   against the hermetic fixture plus biting canaries for publication-order bypass,
   premature wrapper promotion, integrity-mismatch acceptance, and an unpublish attempt,
   each with a distinct message, so that governed publication is an enforced claim rather
   than documentation, and every enforcement dependency fails closed. Line:
   `gpt-5.6-luna` / medium. The mutations are explicit and gate-observable, and the
   profile's cached routing spends medium effort on changes to the oracle.

9. As a release reviewer, I want the runbook to require live-registry integrity
   re-verification, first-publication name-ownership as an explicit manual precondition,
   and interactive 2FA approval that the hermetic gate cannot perform, so that unproved
   publication authority never authorizes a release. Line: `gpt-5.6-terra` / medium.
   These semantics are the ones the gate structurally cannot grade, so they are asserted
   in the runbook and preflight preconditions rather than faked green.

Stories 1 and 3–8 are separable gate-observable slices against the fixture. Story 2's
staged path and story 9's live-registry semantics supply the runbook proof the hermetic
gate cannot; the completed integration uses the highest required line, `gpt-5.6-terra` /
medium, for the runbook and live-registry boundary while the remaining implementation
stays on its per-story line.

## Implementation decisions

- **One new deep module; the rest composes existing seams.** The publication module owns
  the resumable state machine, transition classification, and an npm-registry port.
  Callers do not know registry transition rules, retry classification, or tag ordering.
  The `bench release` shell subcommands, npm invocations, and TOON rendering are thin and
  carry no duplicate policy. The authoritative preflight remains the only publication
  authority and composes the module's verdict; it is not reimplemented.
- **Two port adapters, two evidence classes.** The port has exactly two adapters: the
  existing hermetic registry fixture is the gate adapter and public npm is the runbook
  adapter. The fixture gains staged-submit/approve simulation, read-only integrity query,
  dist-tag transition, and deprecation endpoints so it can exercise the full state
  machine; it remains a test adapter, not a second package builder and not a claim about
  public npm behavior.
- **Two publication paths.** First publication direct-publishes the four platform
  packages under one version-specific non-default tag, verifies each live integrity,
  direct-publishes and verifies the wrapper last, then promotes. Subsequent releases pin
  npm 11.15+/Node 22.14+, stage every package under a version-specific non-default tag
  with stage-only OIDC, and wait for maintainer approval. Reviewer presence is required
  for first publication, staged approval or rejection, and final promotion; verification
  and bounded registry polling are automatic.
- **Integrity-aware, resumable transitions.** A retry treats an already-live package as
  complete only when its registry SHA-512 SRI equals the approved local tarball's; any
  mismatch is terminal. Before approval, failed stages are safe cleanup. After any package
  is live, failure preserves the old `latest`, removes candidate tags, deprecates the bad
  version with a recovery message, and requires a new version. Automation never
  unpublishes. A duplicate or resumed transition is idempotent, not a second live change.
- **Promotion order.** Platform `latest` tags move first and the wrapper `latest` last,
  and only after the complete live set reverifies. A non-default candidate tag is public
  but not default; an exact version or that tag remains installable, so a candidate tag is
  never treated as private.
- **One publication ledger, bound by digest, never a second manifest.** Publication
  produces one durable `publication-record.json` that references the immutable
  release-index digest and carries ordered transitions with package, version, local and
  registry integrity, stage ID when applicable, authentication mode, tag state, timestamp,
  result, and provenance digests. It never rewrites the release index, `SHA256SUMS`, or
  package bytes. npm provenance and receipts are non-reproducible, created after the
  reproducible build, excluded from the reproducibility comparison, and bound only here.
  No credential ever enters the record.
- **Profile-gated fail-closed authority.** Publish mode requires `--profile public|bank`
  and its full green preflight evidence is the only thing that authorizes release. The
  existing requirement registry is the single source of profile requiredness; publish
  stays red for any selected profile whose required producer record is absent, and a
  focused run never authorizes publication. No second requiredness policy is introduced.
- **CLI contract.** `bench release prepare|submit|status|promote|rollback` consume the
  release directory and durable record, are idempotent and non-interactive, emit compact
  TOON on stdout including structured errors and an exact `next_action` for external 2FA
  approval, keep progress on stderr, and exit 0 on success or no-op, 1 on unsatisfied
  release intent, and 2 on usage.

Publication-record envelope (decision-bearing shape, not the full schema):

    { "schema_version": 1, "release_index_sha256": "<immutable digest>",
      "path": "public|first", "profile": "public|bank",
      "transitions": [ { "package": "...", "version": "...", "action":
        "stage|publish|verify|promote|deprecate|tag-remove", "auth_mode":
        "oidc-stage|direct|approval", "stage_id": "...", "local_integrity": "sha512-...",
        "registry_integrity": "sha512-...", "tag_state": "...", "result": "...",
        "timestamp": "..." } ],
      "provenance": [ { "package": "...", "sha256": "..." } ], "result": "..." }

## Testing decisions

- Tests attach to the built `bench release` CLI, the built authoritative preflight and its
  promoted evidence, and the promoted `publication-record.json`. They drive real bytes and
  real registry calls through the hermetic fixture; they do not mock private Go
  collaborators, duplicate transition policy, or infer registry behavior from source.
- Publication behavior extends the existing artifact-surface, preflight-evidence, and
  hermetic-fixture contracts. The fixture's request log is the observation surface for
  call order, tag state, integrity, and the absence of forbidden calls.
- The feature must pass the project gate: `.bench/gate.sh`. Real public-registry
  transactions, first-publication name ownership, and interactive 2FA are runbook and
  manual-precondition evidence; the gate proves only the fixture-exercisable state machine
  and fails closed on any missing tool, runner, or record rather than treating it as a
  skipped step.
- Honest classification: rows whose behavior the gate cannot exercise (real publish,
  live integrity, interactive approval, first-name claim) are marked runbook-only with a
  reason and do not fake a red signal.
- Pre-implementation probes run on 2026-07-19 are named `P1` through `P4`:
  - `P1` searched `internal/` and `bin/bench.sh`: only `release-preflight` routes to the
    binary; no `bench release` publication subcommand or state machine exists.
  - `P2` read `scripts/offline-registry.mjs`: it serves tarball PUT upload, GET metadata
    with `dist-tags`/`integrity`, and tarball download, but has no dist-tag transition,
    deprecation, staged-submit/approve, or integrity-mismatch rejection endpoint.
  - `P3` read the release index schema (`Index`) and evidence finalizer: neither carries a
    publication section, and no `publication-record.json` is produced today.
  - `P4` per the research asset, local npm is 11.11.0 with no `stage` command, confirming
    the staged path is exercised only through the fixture simulation and the runbook, and
    the first-publication direct path is the only gate-plausible live shape.

### Seam diagrams

Seam 1 — publication state machine and registry port:

    trigger: bench release submit|promote|rollback with approved release directory
        │
        ▼
    approved set + release index ──▶ [ resumable publication state machine ] ──▶ registry
    reviewer authorization        ──▶ [ integrity/order/tag classification  ] ──▶ calls
                                              │
                                              ▼
                                  hermetic fixture request log + durable record
                    ◀ tests attach here: drive real calls against the fixture,
                      observe PUT/verify/dist-tag order, integrity retry, no unpublish

Seam 2 — publish-mode authority and record binding:

    trigger: bench release-preflight --mode publish --profile plus a publication run
        │
        ▼
    requirement registry + records ──▶ [ authoritative preflight authority ] ──▶ verdict
    approved release-index digest   ──▶ [ publication record assembly       ] ──▶ record
                                              │
                                              ▼
                                  release-index-bound publication-record.json
                    ◀ tests attach here: withhold a required producer record (red),
                      verify record references the immutable digest and hides secrets

Seam 3 — CLI hybrid contract and resume:

    trigger: bench release <op> invoked twice against the same state
        │
        ▼
    args + durable record ──▶ [ idempotent non-interactive CLI ] ──▶ TOON + exit code
                          ──▶ [ next_action / status renderer   ] ──▶ stderr progress
                                          │
                                          ▼
                              second run is a no-op with identical state
                    ◀ tests attach here: assert 0/1/2 exit codes, next_action for 2FA,
                      structured stdout errors, and idempotent rerun

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | First publication direct-publishes the four platform packages under a version-specific non-default tag and the wrapper last, verifying each live integrity before advancing | publication state machine | M1: reorder to publish the wrapper before a platform package; fixture request-log assertion is red | Ordered request observation rejects wrapper-first or default-tag publication a happy-path publisher would allow |
| 2 | The staged path submits every package under a non-default tag with stage-only OIDC and emits an exact `next_action` for platform-first then wrapper-last 2FA approval | publication state machine plus CLI contract | M2: drop the approval-wait or approve the wrapper before platforms; handoff/order assertion is red | Fixture-simulated staging plus the emitted handoff catch an implementation that publishes where it must only stage |
| 3 | A retry accepts an already-live package only on exact registry-vs-local integrity match and stops terminally on mismatch | publication state machine | M3: accept a live package whose fixture integrity differs; resume assertion is red | Integrity-gated resume defeats a retry that advances past substituted or drifted registry bytes |
| 4 | Promotion moves platform `latest` first and the wrapper `latest` last, only after the complete live set reverifies | publication state machine | M4: promote the wrapper `latest` before a platform or before reverification; dist-tag order assertion is red | Observed tag-transition order rejects a wrapper resolvable before its default platform dependency |
| 5 | Rollback removes candidate tags, preserves the prior `latest`, deprecates the bad version, and issues no unpublish call | publication state machine | M5: emit an unpublish call or drop the deprecate; absence/deprecate assertion is red | Asserting a forbidden call is absent and a deprecate is present catches rollback-by-unpublish and lost prior `latest` |
| 6 | `publication-record.json` references the immutable release-index digest and carries every ordered transition with integrity, stage ID, auth mode, tag state, timestamp, result, and provenance, with no credentials | publish-mode authority and record binding | M6: alter the referenced digest, drop a transition field, or leak a token; record assertion is red | Exact-field and no-secret assertions catch an unbound, incomplete, or credential-bearing ledger and a second manifest |
| 7 | Publish mode requires an explicit profile and stays red for a selected profile whose required FT87, FT88, or (bank) FT71 record is absent, green when all are present | publish-mode authority | M7: remove one required producer fixture record for the selected profile; preflight publish is red with attributed message | Missing-record fixtures prove the fail-closed slot bites without duplicating producer content; conforming fixtures prove composition |
| 8 | The built CLI drives real fixture journeys and canaries reject order bypass, premature wrapper promotion, integrity-mismatch acceptance, and an unpublish attempt, each with a distinct message | CLI contract plus project gate canary | M1/M3/M4/M5: mutate each surface; gate reports its distinct red | Independent mutations prove every new authorization dependency bites rather than trusting a green fixture |
| 8 | `bench release` operations return exit 0 success/no-op, 1 unsatisfied intent, 2 usage; emit TOON stdout with `next_action`; a second identical run is an idempotent no-op | CLI contract | M2: make a resumed run repeat a live transition or return a wrong exit code; idempotency/exit assertion is red | The hybrid-contract and rerun assertions catch a non-idempotent or exit-code-lying CLI a documentation-only surface would ship |
| edge of 1–8 | Control bytes in registry-sourced names/tags/integrity, duplicate or resumed transitions, a partial live set, and plan-vs-apply drift fail closed and preserve recovery state | publication state machine | M3/M6: inject a hostile registry response or a drifted plan; classification assertion is red | Hostile-registry fixtures exercise publication's owned edge classes at the port rather than a private parser |
| edge of 1–8 | SIGINT during submit, wait, promote, deprecate, or record write leaves no promoted partial ledger, no orphan candidate tag beyond recovery, and a clean idempotent resume | publication state machine plus CLI contract | M5: interrupt at a transition; prior-record and tag-state assertion is red on non-atomic recovery | Interrupt-at-stage plus a second run catch leaked state and non-atomic durable recovery jointly owned with preflight |
| 9 | Live public first publication, real-registry integrity re-verification, and interactive 2FA approval are runbook and manual-precondition steps | runbook | Cannot go red in gate — runbook-only: no network, no credentials, no existing package, and no interactive approver in CI | Honest classification keeps the gate from faking a signal it cannot produce; these are asserted in the runbook and preflight preconditions |
| 9 | Name ownership and publish authority are an explicit manual precondition, since an E404 cannot distinguish an available public name from an inaccessible identity | runbook | Cannot go red in gate — runbook-only: E404 is ambiguous per the research probe and proves neither availability nor authority | Recording the precondition prevents the machinery from inferring authority the registry cannot confirm |

Cheapest wrong implementations checked against the map: a publisher that ignores order
fails the request-log order rows; a retry that trusts liveness over integrity fails M3; a
promoter that moves the wrapper first fails M4; a rollback that unpublishes fails M5; a
record that omits the index digest or a transition field fails M6; a machine that skips
the profile gate fails M7's missing-record fixture; and a CLI that is not idempotent fails
the rerun row.

### Edge inventory

- **Error path:** name/version already live with mismatched integrity, stage submission
  rejected, approval timeout, promotion of an unverified set, deprecation failure, and a
  registry 4xx/5xx are coverage rows.
- **Empty/absent input:** absent approved release directory, absent required producer
  record, absent prior publication record, and an empty registry response distinguish
  absent from present-empty in coverage rows.
- **Boundary values:** exactly four platform packages plus one wrapper; one version-specific
  candidate tag; platform-first/wrapper-last on both publish and promote; and zero unpublish
  calls are coverage rows.
- **Malformed input:** malformed registry metadata, unknown transition/status, control bytes
  in registry-sourced names/tags/integrity, and an inconsistent release-index reference are
  coverage rows.
- **Interrupted/partial state:** SIGINT at submit, wait, promote, deprecate, and record
  write, plus prior-record preservation, are coverage rows jointly owned with preflight.
- **Re-run idempotency:** resumed submit, resumed promote, duplicate transition, and a
  second no-op run are coverage rows.
- **Hostile environment:** hostile registry responses, denied non-fixture egress, missing
  npm/node tools, and control bytes are coverage rows; multi-word argument quoting is owned
  by the thin shell caller per the profile checklist.
- **Public npm registry live behavior:** the real publish, live integrity, and interactive
  2FA approval are runbook rows, not gate rows, because the hermetic fixture cannot claim a
  name, hold credentials, or perform proof-of-presence.
- **External signing and provenance trust:** **Won't handle** — signing-key custody and
  transparency-log trust are a separate supply-chain capability excluded by the decision
  map; this slice only records provenance digests it is handed.
- **Credential custody and OIDC issuance:** **Won't handle** — trusted-publisher setup, 2FA
  enrolment, and secret storage are host/registry administration excluded by the map; the
  module consumes an already-authorized session and stores no secret.

## Out of scope

- **FT87 general offline/network controls and FT88 data-handling and FT71 local-event
  record production:** these producer capabilities are already represented by the
  fail-closed registry slots this slice gates on; producing their content is a separate
  feature. Estimated at 10–18 edits and 3 gate runs per roadmap feature.
- **Real public-registry publication and name claiming:** performing a live publish,
  claiming the five identities, or configuring trusted publishers is a credentialed
  operational act the decision map excludes from this work. Estimated at 6–10 edits and 3
  gate runs of runbook and workflow wiring.
- **External signature and provenance trust roots:** signing-key custody, transparency
  logs, and receipt-signature verification are a separate supply-chain capability outside
  the decision map. Estimated at 18–28 edits and 5 gate runs.
- **Native Windows publication:** Windows package identities, selection, and publication
  semantics are a separate platform capability; WSL2 continues to consume Linux artifacts.
  Estimated at 20–30 edits and 5 gate runs.
