# Compliance hardening

*Grill answers were reviewer-pre-authorized recommendations ("go with all your
recommendations"), recorded as decided state; the map is the veto surface.*

## Destination

Turn `COMPLIACE_ASSESSMENT.md` (bank-posture compliance review, 2026-07-10) into a
decided remediation scope for the kit: which findings become code changes, which
become a written security posture, and which are deployment-environment concerns the
kit explicitly does not own.

## #1: What posture does the kit take toward a bank-grade assessment?

Blocked by: none
Type: Grill

### Question

Is the goal certification-readiness (build the bank controls), or extraction of the
genuine engineering value from the findings?

### Answer

Extraction, not certification. Bench is a general developer kit; bank deployment
controls (sandboxing, audit infrastructure, signed provenance, endpoint management)
belong to a deploying organization's platform, not the kit. The kit's response is a
hardening pass: fix the findings that are defects for *any* user, and write the trust
model down so a future assessor reads posture instead of inferring it.

## #2: Which findings become code, which become documentation, which are out of scope?

Blocked by: #1
Type: Grill

### Question

The assessment lists 2 critical, 6 high, and 4 medium findings. What is the triage?

### Answer

- **Code:** M-01 pool/lease permissions; M-03 + H-04 model-query timeouts; C-02
  repair kill-switch; H-06 workflow digest pinning and vulnerability scanning; H-05
  root `LICENSE`.
- **Documentation:** C-01 environment inheritance, H-01 advisory hooks, M-02 gate
  trust, and the H-04 egress inventory — one written trust model in a root
  `SECURITY.md` (see #3, #7).
- **Out of scope:** H-02 (commit/release signing — reviewer process decision, not kit
  code), H-03 (tamper-evident audit infrastructure), M-04 (endpoint management), OS
  sandboxing and environment allowlists, SBOM/signed provenance, the restricted-pilot
  machinery.

## #3: Trust boundary (C-01, H-01, M-02) — build isolation features or document the model?

Blocked by: #2
Type: Grill

### Question

Should the kit add an environment allowlist for agent/gate launches, or state that
the process boundary is not the kit's security boundary?

### Answer

Document, don't build. Agents and gates running with the caller's full authority is
the tool's design point: adapters pass provider credentials through the environment
by design, and an allowlist the user must populate with every credential is the same
trust surface with extra failure modes. The real boundaries are the harness/OS
sandbox and server-enforced branch protection, both outside the kit. `SECURITY.md`
states this explicitly — hooks and guards are honest-mistake layers (as
`internal/gitguard` already declares), and repository gates are trusted code.

## #4: Repair path (C-02) — disable, gate, or keep the self-repair download?

Blocked by: #2
Type: Grill

### Question

The repair path downloads and installs an executable trusted to remote metadata.
Remove it, require confirmation, or bound it?

### Answer

Keep it — it is the distribution mechanism — but bound it: add an explicit opt-out
kill-switch (environment variable) under which the repair path refuses with a
structured error and performs no network activity, and make a repair announce the
version and digest it is about to install. Signature verification against an
organizational trust root stays out of scope; an organization that needs it disables
repair and distributes internally, which the kill-switch enables.

## #5: Network and resource bounds (H-04, M-03) — which bounds get built?

Blocked by: #2
Type: Grill

### Question

Missing timeouts exist on model HTTP queries, agent execution, and standalone gates.
Which are defects and which are accepted posture?

### Answer

Model queries get hard deadlines (bounded connect and total time via request
context; a dedicated client replaces `http.DefaultClient`), with timeout rendered
through the same structured-unavailable posture as other provider errors. Agent and
gate execution get **no** wall-clock deadline: long runs are legitimate, and the
iteration cap plus process-group interrupt is the accepted availability posture. The
best-effort `git fetch` in worktree acquisition stays, documented in the
`SECURITY.md` egress inventory rather than flagged per call.

## #6: Filesystem permissions (M-01) — how tight do pool and lease modes go?

Blocked by: #2
Type: Grill

### Question

The worktree pool is created `0777` and lease files `0644`. What modes are correct?

### Answer

Pool directories `0700`, lease files `0600`. No legitimate cross-user consumer of
the pool exists, so there is no sharing use case to preserve; the pool holds full
checkouts of possibly-private code. Existing pools created with the old mode are
tightened on next acquisition rather than requiring manual migration.

## #7: Governance artifacts (H-05) — which files does the repo add?

Blocked by: #2
Type: Grill

### Question

`package.json` declares MIT but no root license text exists; there is no security
policy, ownership, or disclosure documentation. What lands?

### Answer

Two root files: `LICENSE` (MIT text matching the `package.json` declaration, with
verified copyright line) and `SECURITY.md` (the trust model from #3, the egress
inventory from #5, a security contact, and a one-paragraph disclosure expectation
sized for a small open project). SBOM generation, third-party notice automation, and
license-policy tooling are out of scope until the dependency surface grows past its
current single pinned module.

## #8: Pipeline hardening (H-06) — pinning and scanning?

Blocked by: #2
Type: Grill

### Question

Workflow actions use mutable major-version tags, and no vulnerability scanning is
checked in. What is adopted?

### Answer

Pin every workflow `uses:` reference to a full commit digest (version noted in a
comment), and add `govulncheck` as a required check. Whether it attaches to CI only
or also to the local gate — and its fail posture when the tool is absent — is a
gate-authoring decision the spec-writer settles under `craft-gate` (flagged in the
Handoff). Fuzzing of manifest/tar/envelope parsing is parked pending evidence of
real parser rot; the existing hostile-input suites are the current posture.

## Not yet specified

- None.

## Out of scope

- Tamper-evident audit logging and monitoring integration (H-03) — deployment
  infrastructure; a kit-sized shift-session log is parked in `IDEAS.md`.
- Commit and release signing (H-02) — reviewer git/process configuration, not kit code.
- Endpoint-management packaging, inventory registration, central rollback (M-04).
- OS sandboxing, credential brokering, environment allowlists for agent launches.
- SBOM, signed provenance, artifact attestation, malware scanning.
- Restricted-pilot conditions and certification against any named framework.

## Handoff

1. **Module boundaries.** `internal/worktree` owns pool/lease creation modes.
   `internal/models` owns the provider HTTP client and its deadlines. The repair
   path (`bin/bench-repair-binary.mjs` and whatever invokes it) owns the
   kill-switch refusal. `.github/workflows/` owns digest pins and the
   vulnerability-scan job. Repo root owns `LICENSE` and `SECURITY.md`. No new Go
   package.
2. **Contracts.** Pool directories are created (and existing ones tightened to)
   `0700`, lease files `0600`; all other worktree behavior unchanged. Every model
   provider request carries a bounded deadline; timeout yields the existing
   structured-unavailable row and unchanged exit codes. With the kill-switch set,
   repair performs no network I/O, writes a structured error to stderr, and exits
   nonzero; unset, repair announces version + digest before install. `LICENSE`
   exists at root with MIT text. Every workflow `uses:` is a 40-char commit SHA.
3. **Deep vs thin.** All changes are thin edits inside existing deep modules
   (worktree lifecycle, models client, repair script); none introduces a new seam.
   `SECURITY.md` and `LICENSE` are documents with no code seam.
4. **Black-box assertables.** Stat modes on a freshly acquired and a pre-existing
   pool; a hanging stub server observing bounded model-query failure; a
   kill-switch run observing refusal with zero network side effects; `LICENSE`
   presence and first-line match; a grep over workflows proving every `uses:` is
   digest-shaped.
5. **Gate attachment.** Worktree modes and model timeouts attach as package tests
   in their owners. `LICENSE` presence attaches as a root conformance check
   (authored under `craft-gate`, with a canary proving it bites). Repair
   kill-switch and workflow-pinning attachment are flagged under item 7.
6. **Hostile-input owners.** `internal/worktree` owns umask variance (tests set
   umask explicitly), pre-existing loose-mode pools, idempotent re-acquisition,
   and symlinked pool components (existing checks). `internal/models` owns
   hanging, slow, and oversized responses and absent credentials. The repair path
   owns kill-switch set-but-empty vs unset, missing runtime, and mid-download
   interruption (existing atomic install). Pipeline checks own the
   missing-scanner-tool case. Spaces/globs/control bytes: n/a — no new
   user-supplied string surfaces.
7. **Uncertainty flags.** (a) Where `govulncheck` attaches — CI-only vs local gate
   — and its fail posture when absent locally; settle under `craft-gate`. (b) The
   test seam for the `.mjs` repair path — the grill did not verify what harness
   covers that script. (c) The kill-switch variable name and whether an existing
   offline convention should be reused rather than a new variable minted.
8. **Rejected alternatives.** Environment allowlisting for agent/gate launches;
   removing the repair path entirely; agent or gate wall-clock deadlines; building
   audit logging into the kit; SBOM/provenance tooling at the current dependency
   count; per-call confirmation prompts on network operations.
9. **Domain watch-outs.** Tightening an existing pool must handle directories the
   process does not own without corrupting pool state. A too-short model deadline
   misreports a live-but-slow provider as unavailable — the threshold is a
   documented constant, not folklore. Digest-pinned actions need a stated bump
   procedure or they silently freeze. The gitguard advisory-not-boundary framing
   in `SECURITY.md` must match the source comment, one fact one source.

Dependency order: `LICENSE` + `SECURITY.md` first (cheapest, and #3's documented
posture is referenced by the rest); then worktree permissions and model timeouts
(independent package-local fixes); then the repair kill-switch; then workflow
pinning + scanning last (needs the item-7(a) gate decision). Recommended slices,
not a slicing decision — that stays with the reviewer.
