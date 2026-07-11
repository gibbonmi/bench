# Compliance hardening

Status: staged

## Problem

A bank-posture assessment (`COMPLIACE_ASSESSMENT.md`, 2026-07-10) surfaced defects
and undocumented trust posture in the kit. A future assessor — or any security-minded
user — has to *infer* the kit's trust model from source, several fixable defects ship
(loose pool permissions, unbounded model HTTP queries, an un-opt-out-able self-repair
download, mutable workflow action tags, a missing license), and there is no written
security contact or disclosure path.

## Solution

A hardening pass that extracts the genuine engineering value from the assessment
without chasing certification. Fix the findings that are defects for *any* user, and
write the trust model down so posture is *read*, not inferred. Deployment-environment
controls (OS sandboxing, audit infrastructure, signed provenance, endpoint management)
stay out of scope — they belong to a deploying organization's platform, not this kit.

Concretely: tighten worktree pool/lease permissions; bound model-query HTTP with a
deadline; add a `BENCH_NO_REPAIR` kill-switch and a version+digest announcement to the
self-repair path; pin every workflow action to a commit digest and add release-time
`govulncheck`; add root `LICENSE` and `SECURITY.md`.

## User stories

1. As a user whose pool holds checkouts of possibly-private code, I want the worktree
   pool created `0700` and lease files `0600` so another local user cannot read my
   working trees.
   Line: gpt-5.6-luna / medium. A known seam (`internal/worktree`) with the correctness
   fully observable by stat, but the tighten-existing-pool and unowned-directory edges
   need care beyond a trivial constant swap.

2. As a user with a pool created under the old `0777` mode, I want it tightened to
   `0700` on the next acquisition so I get the fix without a manual migration step.
   Line: gpt-5.6-luna / medium. Same seam and effort as story 1; the idempotent
   re-tighten and not-owned-directory cases are the reason this is medium not low.

3. As a user running `bench models` against a hung or slow provider, I want the query
   to fail within a bounded time and render the existing structured-unavailable row,
   not hang the command.
   Line: gpt-5.6-luna / medium. The `doHTTP` seam already exists, but choosing the
   deadline as a documented constant and driving a real hanging server through the
   dedicated client is more than a one-line change.

4. As an operator in an environment that forbids the self-repair download, I want to
   set `BENCH_NO_REPAIR` so the repair path refuses with a structured error and makes
   zero network calls, letting me distribute the binary internally instead.
   Line: gpt-5.6-luna / medium. A shell short-circuit at a known seam, black-box
   observable through the existing repair contract harness.

5. As a user watching a repair run, I want it to announce the package version and
   digest it is about to install so an unexpected version or artifact is visible before
   it lands.
   Line: gpt-5.6-luna / low. A single stderr line added to the `.mjs` at the point the
   integrity digest is already in hand.

6. As a maintainer, I want every workflow `uses:` pinned to a full commit digest (with
   the version in a comment) so a compromised or force-moved action tag cannot silently
   change what CI runs.
   Line: gpt-5.6-luna / low. Mechanical substitution; the acceptance signal is a static
   grep the local gate enforces.

7. As a maintainer, I want a gate check that fails if any workflow `uses:` is not
   digest-shaped so pinning cannot silently regress on a later edit.
   Line: gpt-5.6-terra / medium. Oracle-authoring under `craft-gate`: a wrong gate is
   the worst bug class here, and the check plus its canary must be authored carefully.

8. As a maintainer, I want `govulncheck` run as a required release-gating step so a
   known vulnerability in the dependency surface blocks a release.
   Line: gpt-5.6-terra / medium. A `craft-gate` fail-posture decision (where it attaches
   and how it fails), not a mechanical edit.

9. As a downstream consumer, I want a root `LICENSE` with MIT text matching the
   `package.json` declaration so the license I am granted is unambiguous.
   Line: gpt-5.6-luna / low. Verbatim canonical MIT text with a verified copyright line.

10. As a maintainer, I want a gate check that fails if root `LICENSE` is absent so the
    governance file cannot silently disappear.
    Line: gpt-5.6-terra / medium. Oracle-authoring under `craft-gate`, kit-self scoped,
    with a canary proving it bites.

11. As a security-minded user or assessor, I want a root `SECURITY.md` stating the trust
    model, the egress inventory, a security contact, and a disclosure expectation so I
    read the posture instead of inferring it from source.
    Line: gpt-5.6-terra / high. Decided-state governance prose that must be precise and
    keep one-source-per-fact with the `gitguard` source comment; the profile's
    doc-authoring bucket would permit top, but this is transcription of a closed map, so
    mid at high effort is the honest call.

## Implementation decisions

**`internal/worktree` (stories 1–2).** Pool directories are created `0700` and lease
files `0600`; all other worktree behavior is unchanged. Because `MkdirAll` is subject
to umask and a pool may already exist under the old mode, acquisition both creates with
`0700` *and* explicitly re-`Chmod`s the pool root to `0700` on every acquisition — an
idempotent tighten that repairs a pre-existing loose pool without a migration step. The
re-tighten is best-effort against a directory the process does not own: a `Chmod`
failure must not corrupt pool state or abort acquisition (log-and-continue posture,
matching the existing best-effort `git fetch`). The lease create at `lifecycle.go`
switches its `O_EXCL` mode argument from `0o644` to `0o600`.

**`internal/models` (story 3).** A dedicated `*http.Client` with a bounded total
deadline replaces `http.DefaultClient`; `doHTTP` is rebound to that client's `Do`. The
deadline is a named package constant with a comment stating the tradeoff — too short a
value misreports a live-but-slow provider as unavailable — set to a value generous for a
models-list endpoint (10s). On timeout the request error flows through the *existing*
`apiInventory` path, yielding the current `query failed` unavailable row and unchanged
exit code (0). No new row type, no new exit code.

**Repair path — `bin/bench.sh` + `bin/bench-repair-binary.mjs` (stories 4–5).** The
kill-switch is checked in `repair_binary()` in `bench.sh`, *before* node is invoked:
when `BENCH_NO_REPAIR` is non-empty, print a structured refusal to stderr
(`bench: repair disabled by BENCH_NO_REPAIR`) and return non-zero without spawning node,
guaranteeing zero network I/O. The outer command still falls through to the existing 127
build-remedy path. The announcement is added in the `.mjs` once `dist.integrity` is in
hand and before the write: one stderr line naming the package, version, and a digest
fragment (`bench: installing <pkg>@<version> sha512:<prefix>`). Mint `BENCH_NO_REPAIR`
as a new variable — no existing offline convention exists to reuse.

**`.github/workflows/release.yml` (stories 6, 8).** Every `uses:` is rewritten to a
40-char commit SHA with the human version in a trailing comment (`# vX.Y.Z`). A new
release step installs a pinned `govulncheck` (`go install
golang.org/x/vuln/cmd/govulncheck@<pinned-version>`) and runs `govulncheck ./...` before
publish; a finding fails the release (fail-closed). The step is CI-only: it needs
network (vuln DB) and a tool install, which the local gate must not depend on.

**Gate checks — `internal/conformance` (stories 7, 10).** Two new checks join the
kit-self conformance suite (`TestRootConformance`, run as the `ln` gate phase). Both are
keyed on `kitRoot`, not the graded `root` — following the `checkConformanceCanaryFamilies(kitRoot)`
precedent — so linked repos are **not** forced to carry a `LICENSE` or workflow
directory. Check A: root `LICENSE` exists and its first line matches the canonical MIT
opening. Check B: every `uses:` under `.github/workflows/*.yml` is digest-shaped
(`@[0-9a-f]{40}`, optionally followed by a comment). Both are pure static reads — no
network — so they belong on the local gate. Each ships a `tests/canary/` case proving it
goes red when the invariant is violated.

**Root `LICENSE` (story 9).** Canonical MIT text, `Copyright (c) 2026 gibbonmi`,
matching the `package.json` `"license": "MIT"` declaration.

**Root `SECURITY.md` (story 11).** One document consolidating: the trust model (process
boundary is not the security boundary; hooks and guards are honest-mistake layers; gates
are trusted code; the real boundaries are the harness/OS sandbox and server-enforced
branch protection); the egress inventory (`bench models` → OpenAI/Anthropic APIs; repair
→ npm registry; worktree acquisition → best-effort `git fetch origin`); a security
contact (gibs.mikej@gmail.com); and a one-paragraph disclosure expectation. To keep
one-source-per-fact, `SECURITY.md` states the *posture* of the advisory-not-boundary
framing but does not restate the `internal/gitguard` mechanism — that source comment
stays canonical.

No new Go package. All code changes are thin edits inside existing deep modules.

## Testing decisions

- **A good test here** exercises external behavior at the module's existing seam: stat
  the resulting file modes, observe a bounded failure against a hanging stub, observe a
  refusal with zero registry hits, grep the rendered artifact. Not internals.
- **Prior art:** `internal/worktree/lifecycle_test.go` (mode and umask assertions),
  `internal/models/models_test.go` (the `doHTTP` swap and unavailable-row assertions),
  `internal/contract/surface/binary_repair_test.go` (real `bench.sh`→`.mjs` driven
  against a stub registry with hit-counting), and any existing `internal/conformance`
  root-file check plus its `tests/canary/` case.
- **Gate command:** `bench gate` (the project gate; the `ln` phase runs
  `TestRootConformance`, `contract` runs the surface contracts).

### Seam diagram

Worktree permissions (stories 1–2):

    trigger: bench worktree / shift acquires a pool worktree
        │
        ▼
    root, resetRef  ──▶  [ worktree.Acquire → MkdirAll+Chmod 0700 ]  ──▶  leased worktree
    pre-existing 0777 pool ──▶ [ tryCreate lease O_EXCL 0600 ]     ──▶  lease file 0600
                                   ◀ tests attach here: stat pool dir + lease file modes
                                     after Acquire (set umask explicitly; seed a loose pool)

Model-query deadline (story 3):

    trigger: bench models (provider key set)
        │
        ▼
    url, key  ──▶  [ fetchDataIDs → dedicated client.Do (bounded) ]  ──▶  ids | error
                       ◀ tests attach here: point url at a hanging httptest server;
                         assert an unavailable row returns within the deadline, no hang

Repair kill-switch + announce (stories 4–5):

    trigger: bench <porcelain> with no platform binary (--repair path)
        │
        ▼
    BENCH_NO_REPAIR set ──▶ [ repair_binary(): refuse before node ]  ──▶ stderr refusal, exit 127, 0 hits
    unset               ──▶ [ .mjs: fetch meta → announce v+digest → install ] ──▶ stderr announce, binary
                              ◀ tests attach here (binary_repair_test.go): stub registry,
                                assert Hits()==0 on refusal; assert announce line on install

License + workflow-pin gate checks (stories 7, 10):

    trigger: bench gate (ln phase → TestRootConformance, kitRoot)
        │
        ▼
    kitRoot ──▶ [ check A: LICENSE present + MIT first line ]  ──▶ diag | ok
    kitRoot ──▶ [ check B: every workflow uses: is @<40-hex> ] ──▶ diag | ok
                   ◀ tests attach here: tests/canary case removes LICENSE / unpins a
                     uses:, asserts the gate goes red (EXPECT names the dropped invariant)

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | fresh pool dir is `0700`, lease is `0600` | `worktree` pkg test | stat after `Acquire` asserts `0700`/`0600`; fails on `0777`/`0644` | a mode regression changes the observed permission bits directly |
| 2 | pre-existing `0777` pool tightened on acquire | `worktree` pkg test | seed a `0777` pool dir, `Acquire`, stat → red if still `0777` | only the re-`Chmod` produces `0700`; a create-only fix leaves the seeded dir loose |
| 2 | tighten tolerates an unowned/failing dir | `worktree` pkg test | `Acquire` still returns a lease when `Chmod` cannot apply | a fix that aborts on `Chmod` error fails to return a worktree |
| 3 | hung provider fails within the deadline | `models` pkg test | hanging httptest server → unavailable row returns bounded; a hang trips the test timeout | without a client deadline the request never returns |
| 3 | timeout renders the existing unavailable row, exit 0 | `models` pkg test | assert `query failed` row and exit 0 on timeout | a new error path or nonzero exit would diverge from the asserted output |
| 4 | `BENCH_NO_REPAIR` refuses with zero network | `binary_repair_test.go` | set var, run porcelain → `registry.Hits()==0`, refusal on stderr, exit 127 | any network call increments hits; a missing short-circuit spawns node and fetches |
| 5 | repair announces version + digest before install | `binary_repair_test.go` | assert stderr contains the version and a `sha512:` fragment before the wrote lines | a silent install omits the announce line |
| 6,7 | every workflow `uses:` is digest-pinned | `internal/conformance` + canary | canary unpins a `uses:` to `@v4` → gate red; a live grep over `release.yml` | a mutable tag fails the `@[0-9a-f]{40}` shape check |
| 8 | govulncheck gates a release | release workflow | not TDD-able locally — CI-only step, no hermetic red signal (see edge inventory) | verified by workflow presence + manual pre-tag smoke, per the release comment |
| 9,10 | root `LICENSE` present with MIT first line | `internal/conformance` + canary | canary removes `LICENSE` → gate red; live check asserts presence + first line | a missing or wrong-license file fails the presence/first-line assertion |
| 11 | `SECURITY.md` states trust model + egress + contact | — | not TDD-able — prose content; presence is assertable, correctness is review (see edge inventory) | reviewed against the map's #3/#5/#7 decisions; presence can be gate-checked if desired |

### Edge inventory

- **Empty / malformed input** — model timeout: covered (row 3). Malformed registry
  metadata and truncated tar already covered by existing repair contracts; unchanged.
- **Permission / ownership** — pool tighten against a not-owned or Chmod-failing
  directory: covered (row "tighten tolerates an unowned/failing dir").
- **Concurrency / idempotency** — re-acquiring an already-`0700` pool re-applies `0700`
  harmlessly: covered by the idempotent tighten (row 2); existing lease-race contracts
  unchanged.
- **umask variance** — pool/lease mode tests set umask explicitly so a permissive umask
  cannot mask a mode regression: covered (rows 1–2, per Handoff item 6).
- **Absent tool / dependency** — `govulncheck` absent locally: **Won't handle** — the
  check is CI-only where the tool is installed deterministically; it is never part of the
  local gate, so there is no local-absent case. Node absent for repair: already covered by
  the existing no-node contract.
- **Kill-switch set-but-empty vs unset** — `BENCH_NO_REPAIR=` (empty) must behave as
  unset (repair proceeds); non-empty refuses: the refusal test sets a non-empty value;
  empty-equals-unset is the shell `-n` semantics, **Won't handle** as a separate row —
  standard `[[ -n ]]` idiom, asserting it would test bash, not the kit.
- **Digest bump procedure** — pinned actions silently freeze without a stated bump path:
  handled in prose — `SECURITY.md`/the workflow comment names how to refresh a pin; not a
  test row (Handoff item 9).
- **Linked-repo false positive** — the LICENSE and workflow-pin checks keyed on `root`
  would fire in every linked repo: handled by keying both on `kitRoot` (implementation
  decision); the canary runs kit-self so it exercises the intended scope.
- **Spaces / globs / control bytes** — n/a: no new user-supplied string surface (Handoff
  item 6). Existing repair space-path contract (`WithSpacePath`) stays green.

## Out of scope

- **A dedicated PR/push CI workflow** — the repo deliberately enforces via the local gate
  plus the `pre-push` hook, with `release.yml` as the only workflow; adding a PR-CI
  surface is a separate capability and a larger architectural decision than this pass. It
  means `govulncheck` and the gate run only at release time, not per-PR. ~1 workflow file,
  3 gate runs to build later. *(Flagged for veto: if you want per-PR scanning, this is the
  cut to reopen.)*
- **Signature verification of the repair download against an org trust root** — the
  kill-switch already lets an org disable repair and distribute internally; a trust-root
  verifier is separate machinery. ~1 module, several gate runs.
- **Commit/release signing (H-02), tamper-evident audit logging (H-03), endpoint
  management (M-04), SBOM / signed provenance, OS sandboxing, environment allowlists** —
  reviewer/process or deployment-platform concerns the kit does not own (map #1, #3, Out
  of scope). A kit-sized shift-session log is parked in `IDEAS.md`.
- **Fuzzing of manifest/tar/envelope parsing** — parked pending evidence of real parser
  rot; the existing hostile-input suites are the current posture (map #8).
- **Optional `SECURITY.md` presence gate check** — presence is assertable but the map did
  not require enforcing it; left as a reviewer call. ~1 check + canary, 2 gate runs.
