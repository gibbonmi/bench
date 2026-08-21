# Land executable freshness

Status: staged

Decision source: reviewer-confirmed current conversation, 2026-08-20 (FT242, re-scoped). The reviewer confirmed three closures in this session: FT242's original ask — a spec amendment reaches the destination through one sanctioned step — already shipped as FT225 (`cb1462a6`, `33aa5258`, retired `bef52480`), so the amendment path is closed and stays closed; the residual defect is `bench worktree land` executing from a stale dev binary, which is what enforced the retired identical-bytes refusal during the ft230 landing and turned a two-line spec amendment into roughly twenty minutes of detour; the fix scopes to the land command, not to a wrapper-wide or porcelain-wide sweep.

Verification log: 2 iteration(s) to accept — one opus/medium round returned one blocking finding and five folds. Blocking: story 3's ordering promise had no red-capable row, since every proof between the parse and `landReviewed` is read-only, so a check inserted just before publication would have passed the whole map while emitting the wrong remedy; folded as row LF9 and the exact placement rule. Folds: stories 7 and 8 (owner-inherited promises beyond the decision source) removed with rows LF6/LF7 moved to the edge inventory as cited dispositions; row LF10 added so the fenced registry test is exercised (closure forwards `Command.Executable`); the `canonicalPath` hop named in the placement decision; the Bootstrap authority section extended with the unverified wrapper hop and the hand-`freshness-publish` seal-forgery residual, the latter as its own Won't handle. The second iteration verified every fold against the tree and accepted.

## Problem

`bench worktree land` executes whatever executable the wrapper resolves. In the
kit source repository that is the dev build `dist/bench`, which only a hand run
of `scripts/go-build.sh` refreshes. The gate refuses a stale executable
(`.bench/gate.sh` runs `freshness-check` before trusting `BENCH_RUN_BINARY`),
but the land command performs its own refusals — the landing contract — before
any gate runs. A stale binary therefore enforces whatever contract was current
when it was built. On 2026-08-20 a `dist/bench` built one day earlier enforced
the identical-spec-bytes refusal that FT225 had retired that afternoon, and the
operator paid a hand mirror commit plus two full destination gate runs for a
two-line spec amendment.

## Solution

The first-run landing proves its own executable is fresh before it enforces any
landing contract. Where the repository declares Go build inputs
(`scripts/go-build.inputs` at the root), `bench worktree land` runs the existing
freshness owner against its own executable and refuses on any error, reusing
the owner's message — which already names the exact rebuild command. Where the
repository declares no Go build inputs (every linked repository), the check
does not apply. The resume path is exempt: it only completes an
already-published landing, and that landing may itself have changed the Go
sources.

## User stories

Line: opus / medium. One small, well-scoped Go change that composes an existing
owner at an existing command seam, with dense prior art in both packages.

1. As a reviewer, I want the landing to refuse when its executable was not
   built from the current sources, so that a retired contract is never enforced
   on my landing.
2. As a reviewer, I want the refusal to name the exact rebuild command, so that
   recovery is one paste rather than an investigation.
3. As an agent, I want the refusal to fire before any landing proof or
   mutation, so that a stale binary never half-executes a landing.
4. As a user of a linked repository, I want the landing to skip the freshness
   proof when the repository declares no Go build inputs, so that installed
   release binaries keep landing without a dev-only proof.
5. As a reviewer, I want the resume of a published landing to complete without
   the freshness proof, so that a landing whose own diff changed the Go sources
   can still finish its marker, reconcile, and release steps.
6. As a reviewer, I want a dev executable without an intact adjacent seal
   refused, so that an unsanctioned plain `go build` never lands.

## Implementation decisions

- The check consumes `freshness.Verify` — the same owner `freshness-check` and
  the gate consume — never a second staleness derivation. Land reaches it
  through a package-var seam (prior art: `landReviewed`,
  `authorizeLandingSource` in the same file) so command tests can substitute a
  deterministic fault.
- Applicability is decided by `Lstat` of `scripts/go-build.inputs` at the
  landing root. Only not-exist skips the check. A present manifest of any form
  — empty, symlink, special — routes to the owner, whose own reading discipline
  refuses what it cannot trust. Presence, not content, is the predicate, so a
  broken link is never classified as an authoritative absence. That predicate
  lives beside the manifest path it reads, in the freshness owner, because the
  path already has exactly one source there; land asks the owner rather than
  repeating the literal.
- Placement: in `LandCommand`, after the resume dispatch, the grammar parse,
  and the `canonicalPath` argument proof, and before `landingDestination` —
  the first repository proof. Input-shape refusals (grammar, canonical path)
  stay first; every repository proof comes after the freshness proof, so a
  stale binary in a repo that would also refuse for a dirty destination or a
  moved tip still emits the rebuild remedy, not the earlier proof's message.
  The resume path never runs the check.
- The refusal reuses the existing `refused{detail=...}` line and exit 1, with
  the owner's message (which carries the copy-paste rebuild remedy, proven
  hostile-path-safe by the owner's own tests) as the detail.
- The command learns its own path the same way `freshness-check` does: the
  registry passes `Command.Executable` into `LandCommand`, whose signature
  grows that one parameter. The registry closure in `cmd/bench/main.go` is the
  only production caller.

### Bootstrap authority

The claim is refusal-before-execution: the landing refuses to act from an
untrusted executable. The hops: the wrapper resolves an executable path and
executes it; that executable then authenticates itself against its adjacent
seal (executable digest proves the running bytes are the sealed bytes) and
against the tree (source digest proves the sealed bytes were built from the
current build inputs). The seal is written only by the sanctioned build —
`scripts/go-build.sh` promotes through `freshness-publish` from the
just-built staged binary — so the trust root is the sanctioned build script.
The first hop is itself unverified: `bin/bench.sh` is plain tracked shell, and
a wrapper-side check is deliberately deferred (see Out of scope). Two more
residuals are named, not hidden. The check runs inside the subject, so an
executable that predates this feature, or one patched to skip its own check,
cannot be caught by itself. And the seal cannot authenticate its own writer:
`freshness-publish` is a public plumbing verb, so a hand `go build` plus a
hand publish produces a seal-matched executable the check accepts. The
independent roots that survive these residuals are the gate — which builds and
grades a private exact-source binary and never reuses `dist/bench` — and the
operator's sanctioned rebuild. This mechanism defends against staleness
accidents, not adversaries; see the Won't handle lines.

## Testing decisions

- A good test drives `bench worktree land` at the command surface and observes
  the refusal line, the exit code, and the untouched git state — never the
  internal call.
- Wiring and ordering tests live in `internal/worktree/land_test.go` (prior
  art: the fixture journeys already there), substituting the check seam where
  a real stale digest would need a Go toolchain in the fixture.
- Binding tests (the real owner reached through the land surface) use a
  sealless fixture executable, which the owner refuses before it ever needs a
  Go module — no toolchain in the fixture.
- Hostile artifact classes (symlinks, special files, mtime ties, torn bytes)
  are the owner's, already proven in `internal/freshness/freshness_test.go`;
  the edge inventory cites the exact functions, read this session.
- The gate seam that observes the feature: the Go test phase of `bench gate`.

### Seam diagram

    reviewer: bench worktree land <path> --request … --base … --source-tip … --spec … -m …
        │
        ▼
    invoked executable + landing root
        │
        ▼
    [ LandCommand first-run entry ]
        │  manifest at root?  ──no──▶  landing proofs proceed (linked repo)
        │ yes
        ▼
    [ verifyLandingExecutable = freshness.Verify(root, executable) ]  ──error──▶  refused{detail=<owner message with rebuild remedy>}, exit 1, no state touched
        │ nil
        ▼
    landing proofs → LandReviewed → publication
        ◀ tests attach at the command surface: fixture repo, fixture executable, substituted or real owner

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| LF1 | 1, 2 | with the manifest present and the freshness owner refusing, `bench worktree land` exits 1 with a `refused{detail=...}` line carrying the owner's message | land command test in `internal/worktree/land_test.go`, fixture repo plus erroring check seam | an implementation with no check lands from the stale executable anyway |
| LF2 | 3 | the LF1 refusal leaves the destination ref, the green marker, and the assignment state untouched | the LF1 test asserts git refs and intent state after the refusal | a check placed after publication mutates before it refuses |
| LF3 | 4 | with no `scripts/go-build.inputs` at the root, the landing never consults the check and reaches its next landing proof | land command test with a panicking check seam and no manifest in the fixture | an unconditional check refuses every linked repository |
| LF4 | 5 | `--resume` completes its marker, reconcile, and release steps while the check seam errors | the existing resume journey `TestLandCommandPublicResumeCompletesPublishedReleaseWithoutRepublishing` extended with an erroring check seam | a check placed before the resume dispatch strands post-publication recovery |
| LF5 | 6 | with the manifest present and no seal beside the fixture executable, the real owner refuses through the land surface | land command test binding the production `freshness.Verify` with a sealless fixture executable | wiring that substitutes a stub but never binds the real owner |
| LF8 | 4 | a present-but-empty manifest still routes the landing to the check | land command test, empty-manifest fixture plus erroring check seam | a content-sniffing applicability predicate silently skips the dev context |
| LF9 | 3 | with the manifest present, an erroring check seam, and a destination state that would independently refuse, the emitted `refused{detail=...}` is the owner's message | land command test with a dirty-destination fixture, prior art `TestLandCommandRefusesDestinationAndSourceStateBeforeGate` and its call counter over the seam | a check inserted just before `landReviewed` passes every other row and emits the wrong remedy |
| LF10 | 1 | the registry closure forwards the invoked `Command.Executable` into the land check | command registry test asserting the forwarded value, precedent the `freshness-check` closure | a closure passing a literal or an empty path verifies the wrong executable |

### Edge inventory

- Owner-inherited coverage, cited rather than rowed — these tests predate this
  spec and stay green across the diff, so they are dispositions, not
  acceptance rows: tied-mtime staleness and torn or patched bytes
  (`TestVerifyUsesContentRatherThanMtime`,
  `TestVerifyRefusesUntrustedArtifactStates`), live and dangling symlinks at
  the seal or the executable
  (`TestVerifyRefusesLiveAndDanglingSymlinkArtifacts`) — all read this
  session in `internal/freshness/freshness_test.go`.
- Manifest absent versus present-but-empty: distinct behaviors, both rowed
  (LF3, LF8).
- Manifest as a dangling or live symlink, a directory, or an empty file:
  `Lstat` reports presence, so the landing routes to the owner rather than
  skipping; the owner's reading discipline decides from there. Never classified
  as absent. Graded directly by
  `TestDeclaresBuildInputsReadsPresenceRatherThanContent`, since the owner's
  own symlink suite covers the seal and the executable, never the manifest.
- Seal or executable as FIFO or other special file: the owner refuses before
  reading (`TestVerifyRefusesSpecialArtifactsBeforeReading`, named from the
  owner's suite this session).
- Hostile repository paths inside the refusal detail: the owner's remedy line
  is copy-paste safe (`TestRefusalRebuildActionIsCopyPasteSafeForHostilePaths`,
  read this session as a name in the owner's suite); the land refusal passes
  the detail through the existing `sanitize.Controls` line.
- **Won't handle:** a stale executable at `--resume` — resume only completes a
  landing already published by a fresh-checked first run, and that landing's
  own diff may change the Go sources; the surviving in-scope caller is the
  operator who rebuilds and reruns anything pre-publication.
- **Won't handle:** auto-rebuild as the remedy — the landing never executes
  what it just built; the remedy stays the named hand command the owner's
  message already carries.
- **Won't handle:** an executable predating this feature or patched to skip
  its own check — self-attestation cannot catch it; the gate's private
  exact-source build and the operator's sanctioned rebuild are the surviving
  independent checks (see Bootstrap authority).
- **Won't handle:** a hand `go build` followed by a hand `freshness-publish`,
  which produces a seal-matched executable the check accepts — the seal
  cannot authenticate its own writer; the same two independent checks
  survive it.
- **Won't handle:** freshness proofs for other porcelain verbs (`shift`,
  `commit`, `status`, …) — a separate capability, priced in Out of scope; the
  land command is where a stale contract is irreversible.

## Ownership fences

- `internal/freshness/freshness.go`
- `internal/worktree/land.go`
- `internal/worktree/land_test.go`
- `cmd/bench/main.go`
- `cmd/bench/command_registry_test.go`
- `ROADMAP.md`
- `roadmap/FT242.md`

## Out of scope

- Porcelain-wide executable freshness (every verb checks before acting): 3
  edits, 2 gate runs.
- A `bench status` row surfacing a stale dev binary ambiently: 2 edits, 2 gate
  runs.
- A wrapper-side check at dev-build resolution in `bin/bench.sh`: 2 edits, 2
  gate runs.

## Further notes

The staging commit also rewrites `roadmap/FT242.md` (and the `ROADMAP.md` index
line) to the re-scoped decision, so the board records that the original ask
closed as FT225 and this spec carries the residual. The first landing of this
spec necessarily runs the pre-feature executable; the protection begins with
the first rebuild after it lands.
