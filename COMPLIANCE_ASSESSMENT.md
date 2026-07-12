# Bench repository-controlled compliance assessment — 2026-07-11

Assessment target: the current local Bench repository at `cd43632` and the
clean working tree that existed when review began.

Decision: **NO-GO for bank deployment or unrestricted regulated-project use.**

This is a technology-risk assessment of controls that can be implemented,
tested, documented, and shipped **inside this repository**. It is not a legal
opinion or certification against a named regulation.

## Scope boundary

In scope:

- source code, scripts, skills, hooks, adapters, and repository configuration;
- local process-launch behavior, environment handling, data written by Bench,
  worktree/commit behavior, failure recovery, and exit/status semantics;
- the project gate, tests, canary, CI/release workflow definitions, dependency
  pins, package composition, reproducible-build evidence, and release metadata;
- repository-owned documentation such as security policy, supported versions,
  data handling, threat model, dependency notices, SBOM, checksums, rollback,
  and incident/recovery instructions.

Explicitly out of scope and neither scored nor recommended as remediation:

- workstation/container sandboxing, endpoint security, operating-system
  accounts, filesystem policy outside paths Bench creates, and device rollout;
- enterprise IAM, credential brokers, secret-store administration, and user
  entitlements;
- network firewalls, proxies, DNS, perimeter allowlists, and egress monitoring;
- server-side branch protection, required-review settings, hosted-repository
  administrators, and central CI runner administration;
- central SIEM/log retention, incident platforms, artifact-registry
  organization settings, publisher-account controls, and external signing-key
  custody.

Those external controls matter in a real bank deployment, but including them
here would violate the requested local-repository boundary. This report asks a
narrower question: **does the repository ship defensible secure defaults,
enforcement, evidence, and failure behavior that downstream bank controls can
rely on?** Today it does not.

## Executive control decision

Bench is privileged developer automation. The repository shows unusually good
attention to tests and guardrails, but four local facts prevent approval:

1. a normal SessionStart can remove a user-owned Git worktree and strand its
   unique commit;
2. a cached green verdict can authorize a commit after the configured gate has
   changed to a failing oracle;
3. the release workflow bypasses the local gate, while its own vulnerability
   scanner currently reports 25 reachable standard-library findings;
4. the distributable package cannot run adoption commands and a release-shaped
   wrapper selects the build host's binary before the correct platform package.

These are repository-controlled integrity failures, not deployment-environment
gaps. They must be fixed in code and contract tests before any bank pilot or
production approval is technically supportable.

| Repo-controlled domain | Rating | Basis |
|---|---|---|
| User/repository state integrity | **red** | Foreign worktree deletion; shift failure erases successful changes |
| Oracle/change integrity | **red** | Gate cache omits oracle identity; release does not run the gate |
| Secure build and vulnerability management | **red** | Release scanner fails; scanner absent from local oracle; no PR/push workflow |
| Artifact/package integrity | **red** | Packed adoption broken; wrong binary precedence; no artifact manifest/checksums |
| Data minimization and secret exposure | **red** | Full environment inherited; objectives exposed through argv/logs/commits |
| Code-execution trust | **red** | SessionStart executes every hook script's `--describe` body |
| Auditability and traceability | **red** | No durable structured run/security-event record; stale intents are misleading |
| Licensing and software composition | **amber/red** | Root MIT file exists; platform packages omit it; no notices/SBOM/license report |
| Availability and recovery | **red** | Unbounded fetch/agent/model paths and destructive failure cleanup |
| SDLC verification | **amber** | Strong gate/race/canary, but semantic and release seams are outside it |
| Security/support governance | **amber/red** | SECURITY exists but lacks supported versions, severity/response, data flows, and recovery policy |

## Blocking findings

### C-01 — Automatic cleanup can delete unowned work and strand commits

Control objectives: data integrity, recoverability, change authorization,
least destructive behavior.

Evidence:

- SessionStart silently invokes `resume-clean`
  (`.bench/hooks/session-start.sh:37`).
- `ConservativeCleanup` removes every clean, unlocked out-of-pool worktree
  (`internal/worktree/resume.go:26-66`).
- The classifier defines out-of-pool only by path; it has no Bench ownership
  proof (`internal/worktree/classifier.go:50-73`).
- Explicit `worktree clean` also removes all such worktrees without the
  confirmation advertised by `bin/bench.sh:269`
  (`internal/worktree/clean.go:99-150`).
- Two independent throwaway probes removed a detached worktree with a unique
  commit; no ref contained it afterward and a no-reflog `git fsck` found it
  unreachable.
- Existing tests expect broad cleanup, so the current green gate enforces the
  unsafe contract.

Repository-controlled remediation:

1. Create one ownership registry/marker for every Bench-created worktree and
   require a matching identity before automatic cleanup.
2. Treat every unowned worktree as read-only informational state.
3. Never remove detached work unless its commit is durably referenced by a
   Bench-owned recovery ref.
4. Make explicit cleanup path-scoped, dry-run first, and require an explicit
   acknowledgement for destructive disposition.
5. Add black-box preservation contracts and a canary that fails if the
   ownership check is removed.

Approval evidence: ordinary branch and detached-unique user worktrees survive
SessionStart, `resume-clean`, and default explicit cleanup; only marker-owned
fixtures are removed.

### C-02 — Gate verdict reuse is not bound to the active oracle

Control objectives: change integrity, trustworthy approval evidence,
fail-closed enforcement.

Evidence:

- Commit reuse checks only cached status and working-tree hash
  (`internal/commit/commit.go:105-114`;
  `internal/status/status.go:175-219`).
- It omits resolved gate kind, command, script digest, and relevant
  configuration.
- A probe ran `BENCH_GATE=true bench gate`, then changed the oracle to
  `BENCH_GATE=false`; `bench commit` reused green and committed.
- Verdict writes ignore `os.WriteFile` errors
  (`internal/gate/gate.go:187-208`), so a failed red write can leave an older
  green record.

Repository-controlled remediation:

1. Version the verdict schema and include a deterministic fingerprint of the
   resolved gate and relevant configuration.
2. Reuse only when tree hash, oracle fingerprint, schema, and freshness match.
3. Write through a same-directory temporary file with sync/atomic promotion.
4. If a red verdict cannot replace prior state, invalidate/remove reuse state
   and fail the requested action.
5. Add contracts for gate-command change, gate-script content change,
   auto-detected-gate change, malformed cache, and write failure.

Approval evidence: no green verdict from oracle A is reusable under oracle B,
and an unrecordable red can never leave an actionable green.

### C-03 — Repository release controls do not enforce a releasable build

Control objectives: secure SDLC, vulnerability management, controlled release,
evidence consistency.

Evidence:

- `.github/workflows/release.yml:11-66` publishes on any `v*` tag without
  running `.bench/gate.sh`, the Go test/race/vet suite, ShellCheck, the canary,
  or installed-artifact setup tests.
- It does not enforce exact SemVer, tag/package/changelog agreement, or that
  the selected commit is on the intended release line.
- Release conformance checks only textual tokens
  (`internal/conformance/package_core_checks_test.go:265-290`).
- The local gate has no vulnerability phase.
- The exact pinned release scanner currently reports 25 reachable
  standard-library vulnerability paths for Go 1.25.0 and exits 3. Fixes extend
  through Go 1.25.12.
- There is no pull-request or push workflow in the repository, so the release
  tag is the first automated repository workflow.

Repository-controlled remediation:

1. Make a single repository script the release preflight and call the full
   gate, vulnerability scan, artifact build, package inspection, and installed
   smoke from it.
2. Call that script from both PR/push CI and tag release; add a canary/static
   structural test that proves every required phase remains wired.
3. Update/pin the fixed Go patch toolchain and make scanner red a release
   blocker with documented exception schema (no silent ignore).
4. Enforce version/tag/changelog/binary/manifest equality before publishing.
5. Save machine-readable test summaries and bind their hashes to the release
   manifest produced by the repository.

Approval evidence: one commit produces green gate, race/vet, scanner, package
inspection, and clean-room installed-smoke records through the same preflight
used by release.

### C-04 — The release artifact is unusable and platform identity is unsound

Control objectives: software supply-chain integrity, artifact identity,
reproducible deployment.

Evidence:

- The package omits root `AGENTS.md`, while `adoption_route` requires it;
  packed and git-installed `link`/`init`/`doctor` fail
  (`package.json:13-28`; `bin/bench.sh:201-209`).
- `prepare` builds `dist/bench`, package files include `dist/`, the workflow
  creates all platform packages below that directory, and the launcher prefers
  root `dist/bench` (`package.json:9-16`;
  `.github/workflows/release.yml:42-66`; `bin/bench.sh:119-147`).
- A release-shaped pack contained the Linux x64 host binary plus all four
  platform binaries; it was approximately 33.6 MB compressed/61.8 MB unpacked.
- The four platform packages contain only `package.json` and `bin/bench`.
- No tag, published package, checksum manifest, SBOM, third-party notice, or
  offline release archive exists.

Repository-controlled remediation:

1. Build each package from an empty, explicit staging allowlist.
2. Assert wrapper and platform tarball contents, executable mode, version, and
   binary target before publish.
3. Install and execute the exact tarballs in isolated prefixes; exercise setup,
   doctor, relink, fresh clone, operational command, and unlink.
4. Generate SHA-256 manifests, SBOM, dependency/license notice, build metadata,
   and an index binding artifacts to commit/toolchain/gate evidence.
5. Produce network-independent per-platform archives and document their local
   verification/install path.
6. Make sequential publication resumable by verifying an existing version's
   digest, staging under a non-default tag, publishing wrapper last, and only
   then promoting.

Approval evidence: the exact indexed artifacts install offline, select the
correct binary on each advertised target, contain legal/evidence files, and
produce the recorded digest.

## High findings

### H-01 — Agent and gate subprocesses inherit an overbroad environment

Control objectives: least privilege, secret minimization, confidential-data
handling.

The shift adapter receives `os.Environ()` plus `BENCH_SHIFT=1`
(`internal/shift/loop.go:271-280`). Gate execution removes only two
wrapper-routing variables and otherwise passes the environment
(`internal/gate/gate.go:125-155`). Repo-local adapters/gates therefore receive
every credential or sensitive value present in the parent process.

Objectives are also placed in the process argument list, stdout, a scratch
file, intent records, and commit subjects (`internal/shift/loop.go:98-130` and
`:159-180`). This creates repository-controlled leakage into logs and Git
history even if external infrastructure is perfectly configured.

Repository-controlled remediation:

- launch agents and gates from separate minimal passlists with documented
  opt-in additions;
- pass prompt content via stdin or a mode-0600 file, not argv;
- store an objective identifier in durable records and redact/sanitize commit
  subjects and terminal summaries;
- ship `DATA_HANDLING.md` describing every repository-controlled prompt,
  environment, file, log, network, cache, and retention path;
- add sentinel tests proving denied variables do not reach default subprocesses
  and control characters cannot enter commits/structured output.

### H-02 — Ambient guard reporting executes unreviewed repository code

Control objectives: trusted code execution, safe inspection, bounded startup.

`guards.Rows` enumerates every `.bench/hooks/*.sh` and runs it with Bash and
`--describe` (`internal/guards/guards.go:47-74,112-135`). SessionStart invokes
that reporting path (`.bench/hooks/session-start.sh:39-41`). An unwired
sentinel script executed during a probe. Per-script 5s+3s bounds are sequential
and have no aggregate cap.

Repository-controlled remediation: make descriptions static manifest data;
execute only an exact managed allowlist when execution is genuinely necessary;
report unknown files without running them; use an aggregate startup deadline;
and add a sentinel canary for non-execution.

### H-03 — Durable, accurate audit evidence is absent

Control objectives: traceability, accountability, incident reconstruction.

Git commits and the mode-0600 intent ledger provide partial evidence, but they
do not record a complete shift identity, resolved line, gate fingerprint,
adapter result, file/recovery disposition, or artifact relationship. Gate
cache and scratch state are mutable. Normal subshell completion leaves stale
intents: `bench status --all` showed 12 “resume interrupted work” entries with
no active process or lease (`internal/worktree/subshell.go:30-56`;
`internal/intent/intent.go:298-308`). Misleading evidence is worse than a
clearly absent field.

Repository-controlled remediation:

- define a versioned, redacted event schema for session/shift start and end,
  resolved adapter/line, gate fingerprint/verdict, commit/recovery reference,
  cleanup decision, and release evidence;
- append atomically to a repository-controlled local log with rotation and
  explicit retention configuration;
- distinguish completed, failed, interrupted, recovered, and abandoned states;
- complete intents on normal exit and recover crashes only from matching lease
  identity;
- document that local logs are mutable and are evidence inputs, not a
  tamper-proof central audit system (the latter is out of this report's scope).

### H-04 — Runtime repair and implicit network behavior are not secure defaults

Control objectives: deterministic execution, supply-chain control, egress
transparency, availability.

Positive change since the prior assessment: `BENCH_NO_REPAIR` exists and the
repair path verifies the digest announced by registry metadata. Remaining
repo-controlled gaps:

- repair is enabled by default and artifact plus digest originate from the
  same configurable registry metadata;
- fetch has no response-size, decompression-size, or total deadline bound;
- concurrent failed repair can remove a target another process installed;
- worktree acquisition implicitly executes unbounded, possibly interactive
  `git fetch origin` and ignores the result
  (`internal/worktree/lifecycle.go:162-164`);
- model lookups lack subprocess/response bounds and query sequentially;
- no single no-network mode covers repair, model discovery, and Git fetch.

Repository-controlled remediation: default release artifacts to no repair;
require explicit repair opt-in; pin expected artifacts through a shipped
manifest independent of transport metadata; bound/download to a temporary
file and atomically promote; add `BENCH_OFFLINE=1` with contracts proving zero
repository-initiated network attempts; make Git refresh explicit and
noninteractive; bound every process and response.

### H-05 — Shift failure semantics do not preserve evidence or work

Control objectives: recoverability, availability, accurate processing status.

A green agent iteration is erased when the subsequent commit fails because
deferred release hard-resets and cleans the pool (`internal/shift/loop.go:172-182`;
`internal/worktree/lifecycle.go:225-256`). Adapter failure is ignored and can
produce exit 0 with “objective likely met.” Missing objectives become the
unbounded “improve the codebase”; invalid caps and cap exhaustion do not yield
an honest incomplete state.

Repository-controlled remediation: preserve a locked worktree or recovery ref
after every post-mutation failure; print and event-log the recovery reference;
propagate staging/adapter/commit errors; require an objective; validate caps;
add a wall timeout; and define distinct complete/incomplete/failed/interrupted
exit/schema states.

### H-06 — Licensing, support, and software-composition records are incomplete

Control objectives: legal attribution, ownership, vulnerability response,
supported-software governance.

Positive evidence: the repository has a root MIT `LICENSE`; `SECURITY.md`
documents that hooks are honest-mistake controls and lists network behavior;
the direct Go dependency is pinned and checksummed; package metadata declares
MIT.

Gaps controlled by this repository:

- platform packages omit the license text;
- there is no third-party notice, SBOM, durable license report, dependency
  update policy, or license-change gate;
- `SECURITY.md` lacks supported versions, severity intake, response targets,
  release/rollback guidance, data/credential flow, and recovery expectations;
- it names a personal email as the only contact and is absent from npm package
  contents;
- package metadata lacks repository, bugs, homepage, and support fields;
- there is no threat model or data-handling document.

Repository-controlled remediation: generate and ship SPDX or CycloneDX SBOM,
dependency/license notice and machine-readable license inventory; fail release
on unreviewed license classes; include LICENSE/SECURITY/notices in every
artifact; complete supported-version, reporting, response, rollback, and EOL
policy; add a repo threat model and data-flow inventory.

### H-07 — Release traceability and reproducibility evidence are incomplete

Control objectives: artifact provenance, repeatability, release rollback,
configuration identification.

Two local builds were byte-identical and the binary embedded clean VCS
metadata, which is good evidence. The release workflow pins action SHAs and
requests provenance. It does not emit a repository-owned manifest binding:

- exact source commit and package version;
- Go/Node/npm versions and build flags;
- dependency and platform manifest digests;
- gate/test/race/vet/vulnerability results;
- each tarball/binary checksum and package content inventory;
- rollback/deprecation target.

The current Linux artifact is dynamically linked to glibc and retains debug
data, while docs claim Linux generally.

Repository-controlled remediation: add a deterministic release-manifest
generator, static builds where supported, exact platform contract, stripped
release flags, reproducibility comparison, checksum verification, and
rollback/deprecation metadata. External signature trust is out of scope; the
repo must still produce all bytes and identities a signer/verifier would bind.

### H-08 — Validation/read errors can be reported as compliant empty state

Control objectives: complete processing, error detection, trustworthy status.

Coverage can pass with no acceptance map and can validate nonexistent story
numbers. Learnings/maps/roadmap/status paths convert non-absence read failures
to “empty” or exit 0. These behaviors let corrupted or unreadable control
records disappear from the dashboard/gate.

Repository-controlled remediation: distinguish absent, empty, malformed,
unreadable, and unsupported-schema states; fail closed in the gate; use exact
story sets/ranges; and add hostile filesystem/parse fixtures to contracts and
canary coverage.

## Additional material findings

### M-01 — Permissions hardening is partial

Fresh and existing pool directories are targeted to mode 0700 and lease/intent
files to 0600. However pool chmod failures are ignored and tests permit the
flow to continue; ownership/symlink assumptions are not uniformly rejected.

Repository-controlled remediation: propagate permission-setting failures,
reject non-owned/symlinked pool roots before use, revalidate after creation,
and contract-test permissive pre-existing paths. This applies only to paths
Bench creates or selects, within the report's scope.

### M-02 — Availability bounds are inconsistent

Model API calls have some HTTP timeout handling, but agent execution, gate
execution, implicit Git fetch, repair, total SessionStart guard inspection,
and several reads/outputs lack coherent bounds. `bench outline` emitted 169 KB
by default. Invalid/huge iteration caps are accepted.

Repository-controlled remediation: define one timeout/size/cap policy,
propagate context cancellation through every subprocess, bound response and
output sizes with truncation metadata, and return a distinct timeout/incomplete
status.

### M-03 — Consumer and maintainer capabilities are not separated

Every linked project receives `bench-assess`, `bench-update-kit`, and
`craft-synthesis` because link copies the complete command/skill trees
(`internal/adopt/link.go:104-107`). These kit-maintenance surfaces can interpret
or replace a consumer's root assessment/README/CHANGELOG. This is unnecessary
privilege and confusing product scope.

Repository-controlled remediation: generate a consumer payload from one
allowlist; keep kit-only maintenance commands out; expose a narrow,
version-pinned, manifest-aware consumer upgrade path; test that forbidden
maintainer assets never enter the package/link plan.

### M-04 — Upgrade/uninstall state can become inaccurate

Relink replaces the manifest without reconciling removed rows, leaving stale
skills active but unowned. `unlink` can return 0 and claim success while
modified managed files and its manifest remain. The documented raw `rm` shim
removal can select an unrelated executable.

Repository-controlled remediation: transactional old/new manifest reconcile;
retain ownership for unresolved files; machine-readable partial status and
nonzero exit; marker-verifying shim removal; downgrade/schema contracts.

### M-05 — Decision records and product claims drift from behavior

ADR 0003 says fixtures always run the full inner gate, while current nested
fixtures select an owning phase. ADR 0004 says publish does not carry
`dist/bench`, contradicted by npm `prepare`. README claims paths that artifact
tests disprove. Six closed decision maps remain in the live directory.

Repository-controlled remediation: amend ADRs to current resulting state;
derive package/CLI claims from executable artifact contracts; retire closed
maps with specs; add reference/inventory conformance where a canonical source
exists.

## Positive controls worth retaining

- The full local gate, conformance packages, and deliberately broken canary are
  substantive and passed.
- Full `go test -race -count=1 ./...`, `go vet ./...`, `go mod verify`, and
  `git diff --check` passed.
- Build outputs were byte-identical across two clean local builds.
- Release actions are SHA-pinned; workflow permissions are narrowly declared;
  provenance is requested; the vulnerability scanner version is pinned.
- The dependency footprint is small: one direct external Go dependency,
  pinned and checksummed, with no runtime service dependency.
- Link conflict preflight, path containment, manifest fingerprints, and
  symlink defenses are thoughtful.
- `bench commit` scopes staged paths and refuses unrelated dirty paths; the
  flaw is cache authority, not an absence of commit discipline.
- Pool and intent designs attempt private modes and atomic ledger writes.
- Agent and gate process groups have cancellation machinery; several hooks
  have bounded execution.
- `SECURITY.md` honestly states that interactive hooks are not an
  evasion-resistant security boundary. That prior issue is closed; it should
  not be reclassified as a missing security boundary.

These controls reduce remediation cost. They do not compensate for the
critical paths above, because a reliable test oracle that asserts unsafe
behavior is still unsafe.

## Prior compliance-report reconciliation

The previous file was misspelled `COMPLIACE_ASSESSMENT.md`; this assessment
replaces it with `COMPLIANCE_ASSESSMENT.md`.

Implemented, correctly bounded, or out-of-repo items have been removed from
the active finding set:

- the old H-01 claim that advisory hooks must be an anti-evasion security
  boundary is closed by the explicit `SECURITY.md` trust statement;
- the old M-02 configured-gate finding is not independently actionable: a
  project-selected gate is intentionally trusted code. H-02 above is different
  because merely reporting unwired scripts executes them;
- the old endpoint-install governance item and recommendations for OS sandbox,
  IAM, firewalling, hosted branch settings, SIEM, and central retention are
  outside the requested repository-only scope and are removed;
- prior claims that root LICENSE/SECURITY were absent, the checkout was dirty,
  default-branch pre-push logic was baked incorrectly, and gate decision
  records were stale are implemented or factually obsolete;
- `BENCH_NO_REPAIR` and digest verification partially remediate the old repair
  finding, so only the remaining repository-controlled gaps are retained;
- 0700 pool and 0600 intent/lease changes partially remediate the old
  permission finding; ignored failure/ownership checks remain M-01.

Still current and re-evidenced rather than copied: environment inheritance
(H-01), repair/network defaults (H-04), audit trail (H-03),
license/governance (H-06), vulnerability/build controls (C-03), permissions
(M-01), and availability (M-02).

## Repository-only remediation roadmap

| Order | Required repository change | Proof required for closure |
|---:|---|---|
| 1 | Ownership-safe worktree cleanup and detached recovery refs | Preservation contracts + cleanup canary |
| 2 | Oracle-fingerprinted, atomic, fail-closed gate verdicts | Gate-switch and write-failure contracts |
| 3 | Preserve shift work/evidence on every failure; truthful result states | Commit/stage/adapter/signal/cap recovery matrix |
| 4 | Static/allowlisted guard metadata; no ambient execution of unknown hooks | Unknown-hook non-execution canary |
| 5 | Fixed Go toolchain and one release preflight used by CI and publish | Gate/race/vet/vuln/package records from same commit |
| 6 | Clean package staging, native artifact smoke, resumable publication | Allowlisted tarball inventories and offline install runs |
| 7 | Environment passlists, private prompt transport, redaction, data-handling doc | Sentinel non-inheritance and control-character tests |
| 8 | Explicit/offline network posture and bounded repair/fetch/model operations | Network-disabled integration contracts and timeout/size tests |
| 9 | Versioned local audit-event and accurate intent lifecycle | Completed/failed/interrupted/recovered event fixtures |
| 10 | SBOM, notices, license inventory, complete SECURITY/support/threat/data docs | Every artifact contains and manifest-binds the records |
| 11 | Transactional upgrade/unlink and consumer/maintainer payload split | Upgrade/downgrade/residual and forbidden-payload tests |
| 12 | Fail-closed coverage/read/status semantics and unified resource bounds | Hostile filesystem/schema/range/output fixtures |

No item in this table requires an out-of-repository control to be implemented
or tested. Native platform execution and package publication are evidence
collection activities, but the scripts, manifests, policies, and assertions
that define success remain repository-owned.

## Approval evidence expected from the repository

Before reassessment, the repository should produce one self-contained evidence
directory or release bundle containing:

1. source commit, semantic version, toolchain/build flags, dependency digests,
   and supported target matrix;
2. full gate, race, vet, vulnerability, canary, artifact-content, and installed
   smoke results tied to that commit;
3. binary and package SHA-256 manifest plus reproducibility comparison;
4. SBOM, third-party/license notice, machine-readable license inventory, and
   package file inventories;
5. completed SECURITY, DATA_HANDLING, THREAT_MODEL, support/version/EOL,
   network/offline, rollback, and recovery documents;
6. examples of versioned audit events for success, failure, interruption, and
   recovery with sensitive fields redacted;
7. contract evidence for foreign-worktree preservation, gate identity change,
   post-agent failure recovery, unknown-hook non-execution, environment
   minimization, offline mode, and transactional upgrade/uninstall;
8. a release manifest binding all preceding files and artifact digests.

An external verifier may later sign, retain, or ingest that evidence, but those
external operations are outside this report. The repository's obligation is to
make the evidence complete, deterministic, and verifiable.

## Local verification performed

| Verification | Result | Interpretation |
|---|---:|---|
| `bin/bench.sh gate` including contracts/conformance/canary | pass | Strong local regression signal; not a release approval because critical unsafe behavior is encoded and release does not call it |
| `go test -race -count=1 ./...` | pass | No race detector findings in exercised Go paths |
| `go vet ./...` | pass | No vet diagnostics |
| `go mod verify` | pass | Downloaded module contents match checksums |
| `git diff --check` before report edits | pass | No whitespace diagnostics |
| Two clean local builds | pass, identical | Same SHA-256: `82a60ccd770de9f203b20e8bbc26d4713af75ff77afe37770043a5e475a11bc9` |
| `govulncheck@v1.6.0 ./...` | **fail** | 25 reachable standard-library findings; release blocker |
| High-confidence credential-pattern path/name scan | no finding | Narrow check only; no full secret scanner was available |
| Source/package and release-shaped npm dry runs | **fail readiness** | Adoption marker missing; wrapper/platform contamination and missing legal/evidence files |
| Foreign/detached worktree cleanup probe | **fail** | User worktree removed; unique commit became unreachable |
| Gate A green → gate B red commit probe | **fail** | Cached A verdict authorized B commit |
| Shift commit-failure recovery probe | **fail** | Successful agent change erased |
| Unknown hook SessionStart/guards probe | **fail** | Unwired shell script executed |
| Clean source VCS metadata | pass | HEAD `cd43632`, `vcs.modified=false` at assessment start |

Unavailable local scanners included gitleaks, gosec, Syft, Trivy, OSV-Scanner,
CycloneDX tooling, and Cosign. Their absence is recorded as missing evidence,
not a pass. Nothing was published, no external repository or registry settings
were assessed, and no claims are made about native platforms that were not
executed.

## Final recommendation

Reject bank deployment from the current repository state. The fastest credible
path is to fix the four blocking integrity/release defects first, then make the
release preflight and exact artifact—not the source-tree gate alone—the unit of
approval. In parallel, reduce subprocess data exposure, stop ambient hook
execution, preserve failed work, and make the repository emit a complete local
evidence bundle.

Reassessment should attach to one immutable version and its generated manifest,
after every repository-only closure test above passes. A green current gate or
manual source checkout is insufficient.
