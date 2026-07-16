# Release preflight (FT82)

## Destination

One repository-owned release preflight, mirroring the gate-phases shape, that PR/push
verification and tag publication both call — fail-closed identity, ancestry, and scan
checks, with machine-readable phase records. Publication stays blocked until it is green.

## #1: What is the preflight, architecturally?

Type: Grill

### Question
Where does the one authoritative implementation live — script, Go core, or workflow jobs?

### Answer
A Go-core-owned plumbing subcommand (`release-preflight`), entered via a thin repo
script the way `.bench/gate.sh` execs `gate-phases`. It composes the existing gate,
race, vet, govulncheck, `build-artifacts.sh`, and `smoke-artifacts.sh` as named phases
and emits the machine-readable records natively. Both workflows reduce to one call each.

## #2: FT82/FT83 manifest boundary

Type: Grill

### Question
FT82's equality set names a "release manifest," but FT83 owns the full deterministic
manifest. Minimal manifest now, or defer the leg?

### Answer
Minimal manifest now: the preflight's machine-readable output is the v1 manifest — an
identity block (tag, package version, source commit, binary-stamped version, changelog
entry) plus per-phase verdict records. FT83 extends it; it does not replace it.

## #3: One phase set or modes?

Blocked by: #1
Type: Grill

### Question
A PR has no tag, so identity phases can't run there; how do both paths call one preflight?

### Answer
One phase registry, two explicit modes. `verify` (PR/push): full gate, race, vet,
govulncheck, artifact build + inspection, clean-room installed smoke. `publish` (tag):
everything in verify plus the identity phases (SemVer/equality, ancestry, changelog).
The structural contract asserts publish is a **strict superset** of verify. This puts
the full gate in CI on every PR for the first time — accepted at ~50s.

## #4: Changelog leg of the equality set

Blocked by: #2
Type: Grill

### Question
`CHANGELOG.md` has only a rolling `## Unreleased` section — nothing to match a tag against.

### Answer
Version the existing file. Cutting a release renames `## Unreleased` to
`## vX.Y.Z (date)`; the publish phase fails unless a heading exactly matching the tag
exists and no `## Unreleased` content is stranded above it. One file stays the single
source for synthesis history and release notes.

## #5: Scanner-exception schema

Type: Grill

### Question
What does a documented, reviewable exception to a govulncheck finding look like?

### Answer
Tracked `scripts/vuln-exceptions.json` (the `platforms.json` precedent for release
inputs): entries of `id`, `reason`, `expires`. The scan phase parses `govulncheck
-json` and fails on any finding without a matching unexpired entry — and also on an
expired or no-longer-matching entry. Exceptions are loans, not waivers.

## #6: Cross-platform smoke composition

Blocked by: #1, #3
Type: Grill

### Question
The per-platform smoke matrix runs on separate runners and can't live inside one
preflight process; publish today smokes nothing.

### Answer
The preflight owns a host clean-room smoke phase (install the exact built tarballs
into a fresh npm prefix, run `smoke-artifacts.sh`); the CI matrix is a fan-out that
invokes that same phase per runner. `release.yml` gains the matrix verify job as a
`needs:` gate before publish. Workflow-structure conformance asserts both workflows
wire through the preflight.

## #7: Release-line ancestry

Type: Grill

### Question
What ancestry must hold for a tag to be publishable?

### Answer
The tagged commit must be an ancestor of (or equal to) `origin/main` at publish time;
anything else fails closed. No release-branch grammar yet — when a `release/*` branch
first exists, that is a new decision.

## #8: Phase-record retention

Blocked by: #2
Type: Grill

### Question
Where do the machine-readable records live, and for how long?

### Answer
Records plus the v1 manifest under `dist/preflight/` (gitignored); both workflows
upload the directory as a CI artifact — default retention on verify, platform maximum
on publish. Durable in-bundle retention is FT83's manifest work.

## #9: Pre-release tags

Blocked by: #7
Type: Grill

### Question
Does a `v0.3.0-rc.1` tag have a legal path through the publish preflight?

### Answer
No. Only exact `vMAJOR.MINOR.PATCH` publishes; pre-release and build-metadata tags
fail closed until FT83's staged-tag publication legalizes staging.

## #10: Toolchain pin enforcement

Type: Task

### Question
FT82 says the Go patch toolchain is fixed and pinned; `go.mod` already carries
`toolchain go1.25.0`.

### Answer
Resolved inline: no new policy. The preflight enforces that the `toolchain` directive
is an exact patch version, and the manifest records the toolchain actually used, so
drift or removal goes red.

## Not yet specified

- Exact JSON field schema of the per-phase record and v1 manifest (spec-level).
- Whether the preflight's gate phase reuses the durable gate cache or always runs cold
  in CI (spec-level; correctness either way — cache is keyed on tree hash).
- npm-prefix isolation details of the clean-room smoke (temp prefix vs container).

## Out of scope

- Full deterministic manifest, SBOM, reproducibility comparison, staged/dist-tag
  publication, rollback records — FT83.
- Live-registry post-publish verification (launcher self-repair smoke) — cannot run
  pre-publish; FT83's staged flow owns it.
- Release-branch (`release/*`) ancestry grammar — deferred until one exists.
- Subprocess env passlists and data-handling inventory — FT88.

## Handoff

1. **Module boundaries.** (a) `release-preflight` plumbing subcommand in the Go core
   (new internal package mirroring gate-phases): phase registry, verify/publish modes,
   record + manifest emission. (b) Thin shell entry script that execs it. (c)
   `scripts/vuln-exceptions.json` schema + its parser inside the scan phase. (d) The
   two workflows, reduced to preflight calls plus the matrix fan-out and publish steps.
   (e) Conformance/canary extension covering phase-registry integrity and workflow wiring.
2. **Contracts.** Preflight: exit 0 green / non-zero red; `--mode verify|publish`
   (publish requires a tag ref); writes `dist/preflight/<phase>.json` per phase plus
   `dist/preflight/manifest.json` (identity block: tag, package version, source commit,
   binary-stamped version, changelog heading, toolchain; per-phase verdicts). Errors
   structured, fail closed on any unreadable input. Exceptions file: absent = no
   exceptions; present-but-malformed = red, never skip.
3. **Deep vs thin.** The preflight core is the deep module (phase orchestration,
   equality, exceptions, records). The shell entry and workflow YAML are thin
   pass-throughs — no seam of their own beyond structure checks. Existing scripts
   (`build-artifacts.sh`, `smoke-artifacts.sh`) stay owned seams the phases call.
4. **Black-box assertables.** Exit code per fixture; record/manifest JSON contents;
   red on: tag≠package version, missing/mismatched changelog heading, stranded
   `## Unreleased`, non-ancestor tag, pre-release tag, uncovered finding, expired or
   unused exception, non-exact toolchain directive, deleted/bypassed phase (canary);
   green baseline on a conformant fixture.
5. **Gate attachment.** Conformance phase: workflow-structure + phase-registry checks.
   Runtime contracts: exercise the built binary's preflight modes in throwaway fixture
   repos. Canary: broken fixtures per red class above. **Gate-blind:** actual GitHub
   Actions execution (runner matrix, artifact upload, npm publish wiring) — manual
   verify on the first real tag.
6. **Hostile-input owners.** Paths with spaces/globs → fixture repos exercise the
   entry script and record paths. Control bytes in git-sourced text (tag names, commit
   subjects entering the manifest) → manifest writer rejects/redacts. Absent vs empty
   vs malformed exceptions file → scan phase (absent=none, empty/malformed=red).
   Missing tool (govulncheck, npm) → the owning phase fails closed, never skips
   silently. Symlinked invocation and cwd-deeper-than-root → entry script. SIGINT
   mid-preflight → no partial `dist/preflight/` mistaken for complete (write-then-rename).
7. **Uncertainty flags.** Publish-artifact retention: "platform maximum" — the exact
   day cap depends on repo visibility (public repos cap lower than private); spec
   should encode "maximum allowed", not a literal number.
8. **Rejected alternatives.** Separate release-notes file (second source); exceptions
   without expiry (rot into suppression); no-exception hard block; preflight
   orchestrating runners itself; ancestry-free or pre-release publishing; committing
   publish records back to the repo.
9. **Domain watch-outs.** npm publication is immutable — a defect found post-publish
   cannot be reliably unpublished, which is why every check runs pre-publish.
   govulncheck verdicts drift with database updates: a commit green today can scan red
   tomorrow with no code change — the expiry-bearing exception schema exists for that.
   GitHub artifact retention caps differ by repository visibility.

Dependency order: preflight core (phases, records, manifest) → workflow wiring
(both callers, matrix `needs:` gate, superset assertion) → conformance/canary
extension. Sliceable at either boundary; slicing is the reviewer's call.
