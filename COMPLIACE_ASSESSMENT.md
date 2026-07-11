# Bench Local Software Compliance Assessment

Assessment date: 2026-07-10  
Assessment target: the current local repository and working tree  
Assessment posture: bank software compliance and technology risk review  
Decision: **do not approve for unrestricted bank use in the current state**

## Executive decision

**Approval likelihood: 20% — Grade D / remediation required.**

Bench has a stronger internal correctness discipline than most developer utilities:
its Go tests, vet checks, ShellCheck checks, targeted race tests, conformance suite,
behavior contracts, and deliberately broken canaries are all substantive. The code
also contains good path-containment, symlink, lifecycle, and commit-scope defenses.

Those strengths do not resolve the principal bank risks. Bench is not merely a
passive wrapper. It launches autonomous coding agents, inherits the caller's full
environment, permits workspace writes and arbitrary configured gate execution,
changes Git state, installs hooks and PATH shims, performs implicit network calls,
and contains a repair path that downloads and executes a binary. Its own destructive
Git guard explicitly states that it is an honest-mistake control rather than an
evasion-resistant security boundary.

The current checkout is also not an immutable approval candidate: it contains 30
staged, modified, or untracked entries; HEAD is unsigned; and the locally built
binary identifies itself as built from a modified tree. A bank approval must attach
to a specific reviewed commit and artifact digest, not a mutable developer checkout.

I would consider a **restricted, non-production pilot approximately 65% likely to be
approved** if it ran in an isolated disposable environment with synthetic or
non-confidential code, no bank credentials, deny-by-default egress, an internally
built pinned binary, server-enforced change controls, and retained audit logs. That
is not approval of the current unrestricted operating model.

## Scope and limitations

This assessment covers only files and behavior evidenced by the local repository,
including current uncommitted changes. It does not assess any externally hosted or
previously released artifact, external service account, remote repository setting,
or distribution channel.

The review evaluates common bank control objectives:

- least privilege and separation of duties;
- secrets and confidential-data protection;
- network egress governance;
- software supply-chain integrity;
- traceability and audit retention;
- controlled change and reproducible artifacts;
- secure development and vulnerability management;
- licensing and ownership governance;
- resilience, safe failure, and recovery.

This is a technical compliance recommendation, not a legal opinion or formal
certification against a named regulatory framework.

## System and inherent-risk classification

Bench should be classified as **privileged developer automation / code-execution
orchestration**, not as a low-risk command-line wrapper.

Its effective capabilities include:

- launching a caller-selected executable agent with a supplied prompt;
- giving that agent write access to a Git worktree;
- inheriting process environment variables, including credentials present there;
- running repository-provided and environment-provided shell gates;
- staging and committing changes after a local gate verdict;
- creating, resetting, and cleaning pooled Git worktrees;
- installing repository hooks and a stable-PATH executable shim;
- querying external model services when credentials are present;
- fetching from a Git remote during worktree acquisition;
- downloading, caching, marking executable, and launching a replacement binary.

The inherent risk is therefore **high** even if the implementation is bug-free. A
bank control decision must address what the utility is authorized to do, not only
whether it does those things correctly.

## Blocking findings

### C-01 — Full environment inheritance exposes credentials to autonomous subprocesses

Severity: **Critical**  
Control domains: least privilege, secrets management, data loss prevention

The shift runner executes the configured agent and passes `os.Environ()` with only
one additional variable. There is no allowlist or credential scrubbing
(`internal/shift/loop.go:258-268`). The standard adapter then launches an autonomous
coding process with workspace-write access (`.bench/adapters/codex:14-18`); the
other adapters likewise launch their corresponding agent processes
(`.bench/adapters/claude:14-18`, `.bench/adapters/opencode:14-18`).

The gate path has the same trust problem: it removes only two wrapper-routing
variables and passes every other environment variable to a repository executable or
configured shell command (`internal/gate/gate.go:75-86`,
`internal/gate/gate.go:103-138`, `internal/gate/gate.go:141-155`).

On a bank workstation, environment variables commonly include source-control,
artifact, cloud, database, proxy, or model-provider credentials. A prompt-injected
agent, compromised adapter, or modified gate can read and transmit them with the
same OS authority as the user. Worktree isolation does not isolate environment
secrets or network access.

Required remediation:

1. Launch agents and gates with a minimal explicit environment allowlist.
2. Remove long-lived credentials from process environments; use short-lived,
   audience-bound credentials from an approved broker where access is necessary.
3. Run the agent in an OS-enforced sandbox or container with an explicit filesystem
   allowlist, no access to the developer home directory, and controlled egress.
4. Add tests proving sensitive sentinel variables and files are unavailable to the
   agent and gate processes.

### C-02 — Runtime repair trusts mutable remote metadata to install executable code

Severity: **Critical**  
Control domains: software supply chain, code integrity, execution control

The repair path selects a user-configurable remote service, fetches metadata, follows
the metadata's artifact URL, downloads bytes, extracts a binary, writes it with
executable permissions, and atomically installs it in the Bench cache
(`bin/bench-repair-binary.mjs:31-59`,
`bin/bench-repair-binary.mjs:82-93`,
`bin/bench-repair-binary.mjs:109-142`).

The SHA-512 integrity value is obtained from the same metadata source as the binary.
The code accurately documents that this protects against corruption and transport
tampering, not compromise of that source (`bin/bench-repair-binary.mjs:95-107`).
There is no bank-controlled digest in the repository, no signature verification
against a bank trust root, and no approval manifest binding version, platform, and
hash before execution.

Required remediation:

1. Disable runtime repair in the bank build.
2. Distribute an internally built binary through an approved internal channel.
3. Bind each platform artifact to a bank-approved digest and signed provenance.
4. Refuse execution when the installed binary does not match the approved manifest.
5. Make any update an explicit, logged administrative action rather than an implicit
   command fallback.

### H-01 — Safety hooks are advisory and bypassable, not a security boundary

Severity: **High**  
Control domains: separation of duties, change control, defense in depth

The destructive-Git analyzer explicitly defines its threat model as an
honest-mistake layer and not an evasion-resistant boundary; wrapper inspection is
only one level deep (`internal/gitguard/gitguard.go:1-10`). A malformed or incomplete
hook envelope becomes an allow decision (`internal/gitguard/gitguard.go:55-68`). The
Stop hook intentionally fails open when its resolver or core is unavailable
(`.bench/hooks/stop.sh:33-60`). The pre-push protection is a local client hook
(`internal/adopt/prepush.sh:1-57`), so it cannot substitute for server-enforced
branch protection or required review.

The worktree and gate controls are valuable reliability measures, but they do not
constrain a malicious or prompt-injected process running with the user's account.

Required remediation:

- Treat hooks as advisory controls in all documentation and risk models.
- Enforce protected branches, required reviews, required independent checks, and
  signed changes on the authoritative Git service.
- Deny direct modification of protected repositories at the OS/token level.
- Enforce agent restrictions outside the agent process using sandbox, identity, and
  network policy controls.

### H-02 — The assessed source is not an immutable, attributable release candidate

Severity: **High**  
Control domains: change management, configuration management, non-repudiation

The local checkout contains 31 staged, modified, or untracked entries. HEAD is commit
`f4a96c1105867fb4c73ab3247d63a46f36293a39`, and Git reports no commit signature.
The built binary reports that revision with `vcs.modified=true` and a `+dirty`
version. There is no tag on HEAD.

The current tests therefore validate a mutable combination of HEAD, index, working
tree, generated binary, and local tool cache. Another reviewer cannot recreate the
exact approval object from a commit identifier alone.

Required remediation:

- Produce a clean reviewed commit containing the intended fixes.
- Require an approved author/reviewer trail and signed commit or signed release
  attestation.
- Build from that clean commit in a controlled environment.
- record the source commit, toolchain, platform, dependency lock/checksum data,
  artifact digest, and test evidence in the approval record.

### H-03 — No durable security audit trail exists for agent activity

Severity: **High**  
Control domains: auditability, incident response, accountability

Bench prints operational output and creates Git commits, but the repository contains
no tamper-evident event log for prompts, tool invocations, subprocess identities,
network destinations, policy decisions, approvals, or denied actions. Shift scratch
notes and objectives are deliberately excluded from changes and removed on cleanup
(`internal/shift/shift.go:24-29`, `internal/shift/shift.go:115-119`,
`.gitignore:16-19`). The gate cache is a mutable local file, not an audit record
(`internal/gate/gate.go:187-208`).

Git history shows resulting code changes, which is useful, but it cannot answer which
commands ran, which data left the host, which credentials were accessible, why a
policy allowed an action, or whether a user bypassed a local hook.

Required remediation:

- Emit structured security events to an approved append-only audit destination.
- Record session identity, source commit, agent and model identity, policy version,
  prompts or approved redacted hashes, command/tool decisions, egress destinations,
  gate results, commits, and reviewer approvals.
- Define retention, access, redaction, incident response, and correlation with
  endpoint and Git-service logs.

### H-04 — Network egress is implicit and insufficiently governed

Severity: **High**  
Control domains: data protection, third-party risk, resilience

Worktree acquisition performs a best-effort `git fetch origin` without a command
option that makes network access explicit (`internal/worktree/lifecycle.go:148-162`).
The model inventory reads credentials from the environment and calls two external
model endpoints (`internal/models/models.go:56-70`,
`internal/models/models.go:92-100`). It uses `http.DefaultClient` and constructs
requests without a deadline or context (`internal/models/models.go:22-27`,
`internal/models/models.go:142-160`). The binary-repair path is an additional network
and executable-ingress surface.

There is no repository policy declaring allowed destinations, proxy enforcement,
certificate requirements, retry and timeout limits, offline mode, data
classification restrictions, or approval workflow for transmitting source-derived
content through agent adapters.

Required remediation:

- Default to offline operation in the bank build.
- Make each network operation explicit and policy-controlled.
- Enforce destination allowlists through the platform, not application convention.
- Add bounded connect, response, and total timeouts.
- Document exactly what data each adapter may transmit and obtain information-security
  and third-party-risk approval before confidential source is in scope.

### H-05 — Licensing and security governance artifacts are incomplete

Severity: **High**  
Control domains: legal compliance, ownership, vulnerability disclosure

The project metadata declares an MIT license (`package.json:53`), but the repository
contains no root `LICENSE` file. The source has one direct third-party Go dependency
pinned to a pseudo-version (`go.mod:7`) and checksummed in `go.sum`, but the repository
contains no SBOM, third-party notice, dependency-license report, or documented legal
approval. The locally cached dependency carries an MIT license, but cache inspection
is not durable repository evidence.

The repository also has no `SECURITY.md`, `CODEOWNERS`, documented security contact,
supported-version policy, vulnerability disclosure procedure, or security-specific
threat model.

Required remediation:

- Add the project's complete license text and verified copyright ownership.
- Generate and retain a machine-readable SBOM for every approved build.
- Add third-party license attribution and automated license-policy enforcement.
- Add security ownership, vulnerability intake, supported-version, severity,
  remediation-SLA, and coordinated-disclosure documentation.

### H-06 — Vulnerability-management and build-pipeline controls are not demonstrated

Severity: **High**  
Control domains: secure SDLC, supply-chain assurance

The repository has excellent behavioral testing, but no checked-in SAST, dependency
vulnerability, secret-scanning, SBOM, container/artifact, or malware-scanning policy.
The available local environment did not contain `govulncheck`, `gosec`, `syft`,
`trivy`, or `osv-scanner`, so none of those checks could be credited. Workflow actions
are referenced by mutable major-version tags rather than immutable commit digests
(`.github/workflows/release.yml:20-33`).

Required remediation:

- Pin workflow actions and build images by immutable digest.
- Add required SAST, dependency vulnerability, secret, license, and SBOM checks.
- Define severity thresholds and time-bound remediation exceptions.
- Retain scanner versions, databases, inputs, and results with the approved artifact.
- Add fuzzing and adversarial tests for manifest parsing, tar parsing, hook envelopes,
  path handling, and shell-command classification.

## Additional material findings

### M-01 — Shared worktree storage permissions are too broad for confidential code

Severity: **Medium**

The worktree pool is created with mode `0777` (`internal/worktree/lifecycle.go:154-157`),
subject only to the process umask. Lease files are created with mode `0644`
(`internal/worktree/lifecycle.go:70-80`). On a multi-user host or permissive umask,
another local account may be able to list or read confidential worktree content, or
interfere with pool state.

Use a bank-controlled private directory, require ownership by the invoking identity,
create directories as `0700`, files as `0600`, reject symlinked or foreign-owned pool
components, and test hostile multi-user filesystem states.

### M-02 — Arbitrary configured gate execution expands the trusted computing base

Severity: **Medium**

An executable repository gate takes precedence, and an environment-provided gate is
executed through `bash -c` (`internal/gate/gate.go:75-86`,
`internal/gate/gate.go:103-119`). This is appropriate for a general development tool,
but in a bank deployment it means repository content and caller environment are
trusted code. The canary protects the Bench repository's own gate; it does not make
every consumer repository gate trustworthy.

Approve gates through policy, execute them in the same restricted environment as the
agent, bind their digest to the repository approval record, and prevent unreviewed
gate changes from accessing credentials or network resources.

### M-03 — Availability controls are incomplete

Severity: **Medium**

The shift loop has iteration caps and process-group cancellation, which is positive,
but agent execution has no wall-clock deadline (`internal/shift/loop.go:258-277`),
standalone gates have no deadline (`internal/gate/gate.go:141-163`), and model HTTP
queries have no timeout. A hung child, gate, endpoint, decompression, or oversized
response can consume workstation or runner resources indefinitely.

Add configurable hard deadlines, response-size and artifact-size limits, disk quotas,
CPU/memory/process limits, bounded decompression, and deterministic cleanup after
forced termination.

### M-04 — Endpoint installation changes require managed deployment controls

Severity: **Medium**

`bench link` installs hooks and managed repository files, while `bench doctor --fix`
selects a writable stable PATH directory and atomically creates an executable shim
(`internal/adopt/doctor.go:49-65`, `internal/adopt/doctor.go:285-340`). Conflict and
ownership checks are good, but there is no bank software inventory registration,
administrator approval, endpoint-management package, or central uninstall evidence.

Deploy through managed endpoint tooling, restrict installation paths, record software
inventory and version, and provide centrally tested rollback.

## Positive controls and engineering strengths

The following materially reduce risk and should be preserved:

1. **Substantive gate and canary design.** The project runs root conformance,
   behavior contracts, ShellCheck, Go build/vet/test/cross-compile checks, and a canary
   suite that proves targeted checks fail when sabotaged (`projects/benchkit.md:119-153`).

2. **Path-scoped gated commits.** `bench commit` refuses unexplained dirty files,
   requires a green gate, stages literal named paths, and commits only afterward
   (`internal/commit/commit.go:86-139`). This is a strong accidental-change control.

3. **Safe adoption behavior.** Link preflight detects symlink parents, foreign files,
   and modified managed files before writing (`internal/adopt/link.go:178-225`).

4. **Safe unlink containment.** Unlink refuses absolute paths, traversal, symlinked
   parents, non-regular targets, and modified managed content; it removes the manifest
   last (`internal/adopt/unlink.go:73-151`).

5. **Manifest integrity.** Managed files and symlinks are SHA-256 fingerprinted, and
   manifest replacement is atomic (`internal/adopt/manifest.go:27-56`,
   `internal/adopt/manifest.go:63-123`).

6. **Worktree race awareness.** Lease creation uses `O_EXCL`, verifies stale-takeover
   identity, protects live owners, and cleans before releasing the lease
   (`internal/worktree/lifecycle.go:70-120`,
   `internal/worktree/lifecycle.go:222-253`).

7. **Process cleanup.** Agent and cancellable gate children run in separate process
   groups so interrupts can terminate descendants before worktree release
   (`internal/shift/loop.go:258-301`).

8. **Small dependency surface.** The compiled Go core has one direct third-party
   dependency, pinned by version and checksum (`go.mod:7`, `go.sum`).

9. **Reproducibility intent.** The Go toolchain is pinned, build flags are centralized,
   and builds use `-trimpath` (`go.mod:3-5`, `scripts/go-build.sh:1-30`). This is a good
   base for reproducible builds once the source and build environment are immutable.

10. **Explicit hostile-input testing.** The project profile covers spaces, glob
    characters, control bytes, newline-less files, absent/empty distinctions, missing
    tools, symlinks, interruption, idempotency, and deep working directories
    (`projects/benchkit.md:85-108`).

## Local verification performed

| Verification | Result | Compliance interpretation |
|---|---:|---|
| `go test -count=1 ./...` | Pass | All local Go packages and contract suites passed. |
| `go vet ./...` | Pass | No vet diagnostics. |
| ShellCheck at warning severity over shipped shell, hooks, adapters, and pre-push asset | Pass | No ShellCheck warning-or-higher diagnostics. |
| Race detection over adopt, commit, gate, gitguard, shift, and worktree packages | Pass | No races detected in the security-sensitive package subset. |
| `git diff --check` | Pass | No whitespace-error diagnostics in the current diff. |
| Go module inventory | Pass with caveat | One direct external Go dependency, pinned and checksummed. No repository SBOM or durable license report. |
| Vulnerability/SAST/SBOM scanners | Not performed | Required tools were unavailable locally; absence is not a pass. |
| Commit signature | Fail | HEAD is unsigned. |
| Clean immutable source state | Fail | 31 staged, modified, or untracked entries; local binary reports a modified VCS tree. |
| Root license/security/ownership artifacts | Fail | No root license text, security policy, ownership file, SBOM, or third-party notice found. |

Passing tests demonstrate implementation quality; they do not by themselves satisfy
bank approval requirements for privilege, data movement, audit, legal, and supply-chain
controls.

## Required approval package

The following evidence should be mandatory before reassessment:

1. A clean, reviewed, signed source commit and immutable artifact digest.
2. An internally controlled build record with pinned toolchain and dependencies,
   reproducibility comparison, SBOM, license report, and signed provenance.
3. A bank build with runtime repair disabled and implicit network activity disabled.
4. An agent execution profile with a minimal environment, credential isolation,
   filesystem sandbox, resource limits, and deny-by-default egress.
5. A data-flow diagram identifying every file, environment variable, prompt, response,
   log, and external destination, with data classification and retention decisions.
6. A threat model covering prompt injection, malicious repositories, compromised
   adapters, local multi-user attacks, supply-chain compromise, hook bypass, credential
   theft, data exfiltration, and denial of service.
7. Server-enforced branch protections, independent required checks, code-owner review,
   and a prohibition on direct agent credentials for protected branches.
8. Tamper-evident audit events integrated with bank monitoring and incident response.
9. Required SAST, dependency, secret, license, SBOM, malware, and artifact-integrity
   checks with documented thresholds and exception governance.
10. Security ownership, vulnerability disclosure, supported-version, patch-SLA,
    rollback, and end-of-life procedures.
11. Penetration testing or adversarial assurance focused on the executable repair path,
    environment leakage, prompt injection, hook bypass, path/symlink attacks, archive
    parsing, worktree races, and configured shell gates.

## Restricted pilot conditions

If the business wants evaluation before the blocking findings are remediated, approval
should be limited to all of the following conditions:

- disposable single-user VM or container;
- synthetic, public, or explicitly non-confidential source only;
- no production, customer, payment, authentication, security, or regulated data;
- no bank credentials or general developer home directory mounted;
- dedicated unprivileged identity with no protected-branch write permission;
- internally built and hash-pinned binary; runtime repair disabled;
- deny-by-default outbound network with individually approved destinations;
- manual review of every resulting diff and no autonomous merge or deployment;
- endpoint, process, Git, and network logs retained under the pilot record;
- time-limited exception with named owner, user population, repositories, and exit
  criteria.

## Final recommendation

**Reject unrestricted deployment now.** The application's correctness controls are
credible, but the current design trusts the user environment, repository, agent,
network, and remote executable source more broadly than a bank can accept. The fastest
credible route to approval is not another round of unit tests; it is a bank-specific
execution profile that removes runtime executable repair, isolates credentials and
data, enforces offline or allowlisted networking, produces durable audit evidence,
and binds approval to a clean signed artifact with complete governance records.
