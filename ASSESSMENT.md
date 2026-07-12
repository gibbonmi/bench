# Bench release-readiness assessment — 2026-07-11

Assessment target: the current local repository at `cd43632` plus the clean
working tree that existed when the review began. This replaces the 2026-07-08
assessment. It covers adoption, packaging, release automation, first-hour user
experience, CLI/core safety, workflow commands, all shipped skills, gate
authority, tests, records, and live operational state.

Severity in this report:

- **critical** — normal use can lose or strand user state, or the release
  cannot be trusted to preserve it;
- **high** — an advertised invariant, install route, or release guarantee is
  false on a supported path;
- **medium** — a reachable defect, misleading success, or material adoption
  obstacle;
- **low** — bounded friction, maintainability drift, or defense in depth.

## Decision

**NO-GO for deployment into other projects and NO-GO for a public npm
release.** Bench has a substantial, fast local verification system, but its
normal session-start path can delete Git worktrees it does not own; its gate
cache can authorize a commit after the gate has changed to a failing oracle;
its shift loop can erase successful agent work after a commit failure; and the
only advertised distributable cannot execute `link`. The release workflow
would also put the build host's Linux binary ahead of the selected platform
package.

The source checkout is useful for informed maintainer experimentation, but it
is not safe as a default tool in an arbitrary repository. Until R-01 is fixed,
do not let Bench automatically clean worktrees, and do not use the worktree
cleanup command in a repository with user-managed worktrees.

| Surface | Readiness | Assessment |
|---|---:|---|
| Source-tree development | conditional | Strong test feedback, but unsafe worktree ownership and shift recovery semantics |
| Install into a new project | 0/5 | Registry package absent; advertised git package cannot run adoption commands |
| First-hour usability | 1/5 | The first documented command is unavailable until after linking; setup assets and restart guidance are incomplete |
| Existing-project safety | 1/5 | Conflicts are mostly preserved, but relink strands assets and SessionStart can remove foreign worktrees |
| CLI truthfulness/recovery | 2/5 | Good structured surfaces, undermined by false success, false-empty reads, and destructive teardown |
| Skills/workflow coherence | 2/5 | Strong conceptual decomposition, but several phase contracts contradict one another and maintainer-only skills ship to consumers |
| Local quality engineering | 4/5 | Full gate, canary, race tests, vet, and reproducible builds pass; important semantic paths are outside the oracle |
| Release engineering | 1/5 | Pinned actions and provenance exist, but release bypasses the gate, fails vulnerability scan, and produces the wrong wrapper shape |

Current inventory: **3 critical, 12 high, 30 medium, and 8 low** findings.
Counts include the detailed adoption, CLI, and skill tables below; overlapping
compliance consequences are not counted again here.

There is no honest “one or two command” journey today. The right target is
clear and has been added to `ROADMAP.md` as FT76: one mechanical command plus,
only where project facts require judgment, one harness-native setup
conversation.

## Release blockers

### R-01 — Session start deletes worktrees Bench does not own

Severity: **critical**

`.bench/hooks/session-start.sh:37` silently runs `resume-clean` whenever a
session starts. `ConservativeCleanup` removes every clean, unlocked
`ClassOutOfPool` worktree (`internal/worktree/resume.go:26-66`), but
`ClassOutOfPool` means only “not the root and not inside Bench's pool”
(`internal/worktree/classifier.go:50-73`). It carries no ownership proof.

Two independent throwaway-repository probes reproduced removal of an ordinary
user-managed detached worktree containing a unique commit. After cleanup, no
ref contained that commit and `git fsck --unreachable --no-reflogs` reported
it unreachable. Git may retain the object until garbage collection, but Bench
has destroyed the checkout and its durable reference. Existing tests require
the broad removal behavior, so the green gate currently protects the defect.

The explicit `bench worktree clean` path is broader still: it removes every
out-of-pool worktree without asking (`internal/worktree/clean.go:99-150`),
although CLI help promises removal “after confirmation” (`bin/bench.sh:269`).

Required release condition: automatic cleanup may touch only a worktree with a
verifiable Bench ownership record. Foreign worktrees must be informational.
Detached work must never be removed unless its commit is durably referenced.
Add black-box contracts for an ordinary branch worktree and a detached unique
commit, plus a canary that proves ownership checks bite.

### R-02 — A green verdict from one gate authorizes a commit under another

Severity: **critical**

Commit reuse binds a cached green verdict only to the working-tree hash
(`internal/commit/commit.go:105-114`; `internal/status/status.go:175-219`). It
does not bind the verdict to the resolved gate kind, command, script digest,
or relevant configuration. Reproduced in a throwaway repository:

```text
BENCH_GATE=true  bench gate
BENCH_GATE=false bench commit -m should-have-failed change.txt
gate: green (fresh verdict reused for this tree)
committed 1 path(s)
```

The verdict writer also ignores write failures
(`internal/gate/gate.go:187-208`), so a failed attempt to record red can leave
an older green record available. This breaks the central claim that the active
gate is the oracle.

Required release condition: version the cache record and include a fingerprint
of the fully resolved oracle. Reuse only an exact tree-and-oracle match. Write
atomically, and fail closed if a red verdict cannot replace prior state.

### R-03 — Shift teardown can erase a successful agent iteration

Severity: **critical**

When the gate is green but `git commit` fails, the shift loop returns 1
(`internal/shift/loop.go:172-182`). Deferred teardown releases the pool entry,
whose `Release` hard-resets and cleans before dropping the lease
(`internal/worktree/lifecycle.go:225-256`). A probe with a successful adapter,
green gate, and missing Git identity ended with no commit containing the new
file and a cleaned pool worktree. There was no recovery path.

`stageTouched` also ignores staging failures (`internal/shift/shift.go:61-71`).
Any failure after agent mutation must preserve a locked recovery worktree or a
durable recovery ref and print its exact path. Cover missing identity, failing
commit hook, staging failure, signal, and teardown error.

### R-04 — The advertised package installs but cannot adopt a repository

Severity: **high**

The package intentionally omits root `AGENTS.md`, while the launcher's
`adoption_route` refuses `link`, `init`, and `doctor` unless both
`.agents/commands` and root `AGENTS.md` exist (`package.json:13-28`;
`bin/bench.sh:201-209`). Both a packed install and the documented git-`npx`
route reproduced:

```text
benchkit 0.2.0 (linux/amd64)
bench: link/init/doctor must run from the real Bench kit
```

The postinstall doctor fails for the same reason, so it does not establish a
working durable shim. `redbench` and the checked `@redbench/*` packages are not
published, and the repository has no release tags. A stranger therefore has
no working supported entry point.

Required release condition: test the artifact, not only the source tree. The
gate must pack, install in an isolated prefix, and exercise `version`, setup or
`link --init`, `doctor`, a local operational command, relink, and unlink.

### R-05 — Release packing selects the build host's binary on other platforms

Severity: **high**

`npm publish` runs `prepare`, which creates `dist/bench`; `files` ships all of
`dist/`; the release workflow also generates four packages under
`dist/packages`; and the launcher checks root `dist/bench` before the selected
optional platform dependency (`package.json:9-16`; `.github/workflows/release.yml:42-66`;
`bin/bench.sh:119-147`).

A release-shaped dry run produced a 33.6 MB compressed wrapper containing the
Linux x64 host binary and all four nested platform binaries. On macOS or Linux
arm64, the launcher would attempt the Linux x64 host executable first. The
presence of that executable also suppresses repair. ADR 0004's statement that
the publish checkout has no root binary is false because npm runs `prepare`.

Required release condition: construct every package in a clean staging tree;
make wrapper contents an explicit allowlist; assert no `dist/bench` or nested
packages in the wrapper; and execute native artifact smokes for every claimed
OS/architecture.

### R-06 — Release bypasses the project oracle and currently fails its scanner

Severity: **high**

The only workflow triggers on any `v*` tag and publishes without running
`.bench/gate.sh`, `go test`, `go vet`, ShellCheck, the canary, or installed
artifact adoption. It does not enforce exact SemVer, tag/package/changelog
agreement, or default-branch ancestry (`.github/workflows/release.yml:11-66`).
The conformance check merely looks for release tokens
(`internal/conformance/package_core_checks_test.go:265-290`).

The local gate is green, but the scanner used by release is red. This exact
command reports 25 reachable standard-library vulnerability paths from the
Go 1.25.0 toolchain, with fixes through Go 1.25.12:

```text
go run golang.org/x/vuln/cmd/govulncheck@v1.6.0 ./...
Your code is affected by 25 vulnerabilities from the Go standard library.
exit status 3
```

Examples include `crypto/tls` GO-2026-5856 and `html/template`
GO-2026-4980. This is scanner evidence, not a claim that every reachable
symbol is exploitable in Bench, but it is an objective release failure.

Required release condition: update to a fixed toolchain, run the authoritative
gate plus release-only security/package checks before any publish, and add a
canary proving those steps cannot be deleted or weakened unnoticed.

### R-07 — Opening a session executes every shell file in the hooks directory

Severity: **high**

`bench guards` discovers every `.bench/hooks/*.sh` and executes it with Bash
and `--describe` (`internal/guards/guards.go:47-74,112-135`). SessionStart runs
`guards --brief`, so an unwired or user-added script executes simply because a
session opened (`.bench/hooks/session-start.sh:39-41`). A sentinel hook probe
confirmed execution. Only a foreign pre-push hook receives a non-execution
safeguard. Timeouts are per file and sequential, with no aggregate deadline.

Required release condition: descriptions come from static metadata or a
manifest-owned allowlist. Unknown scripts are reported, never executed. Add an
unknown-hook sentinel contract and an aggregate time budget.

### R-08 — Shift can report likely success after the adapter failed

Severity: **high**

Adapter start and exit failures are deliberately ignored
(`internal/shift/loop.go:271-290`). With `/bin/false` as the adapter and a green
gate, a probe exited 0 with “objective likely met,” zero changes, and zero
commits. The default for a missing objective is the unbounded “improve the
codebase” (`internal/shift/shift.go:133-140`), invalid or negative iteration
caps are accepted, and cap exhaustion still exits successfully.

Required release condition: distinguish adapter failure, no-op completion,
objective predicate success, and cap exhaustion with honest exit codes and
structured results. Require an objective, validate positive bounded caps, and
set a wall deadline.

### R-09 — The installed shim routes maintenance commands into a known refusal

Severity: **high**

The generated stable shim chooses the repo-local wrapper first
(`internal/adopt/doctor.go:107-142`). That wrapper intentionally refuses
`link`, `init`, and `doctor`, so the advertised installed `bench link` and
`bench doctor --fix` paths self-loop. This was reproduced after generating the
shim. Route adoption/upgrade/doctor commands to the installed kit and
operational commands to the project-local runtime, and pin the interaction in
one end-to-end contract.

## The one- or two-command target

The desired experience should be:

```sh
# no global install: one shell command
npx --yes redbench@<immutable-version> setup

# or a durable install: two shell commands
npm install --global redbench@<immutable-version>
bench setup
```

If repository-specific decisions cannot be inferred safely, the mechanical
command ends by asking for exactly one harness-native action after a reload:

```text
/bench-setup-repo       # Claude command surface
$bench-setup-repo       # Codex/AGENTS skill surface
```

FT76 records the fuller, Matt-Pocock-style interaction: explore the repository,
present detected facts and missing decisions, ask one unresolved question at a
time, preview the write plan, then write. `setup` must compose the current
`link`, `init`, and setup skill; it must not create a second installer.

Definition of done for that journey:

- works from the exact packed artifact and an approved offline artifact;
- selects and proves the correct native binary;
- performs a transactional preflight before changing the repository;
- preserves surrounding user content in `AGENTS.md` or `CLAUDE.md` and updates
  one marker-owned block without duplicates;
- composes existing JSON settings and Git hooks or stops with a precise,
  non-partial conflict report;
- seeds the gate, canary, project profile, context/ADR skeletons, journal, and
  repo-local executable path from shipped templates;
- prints an honest intentionally-red state when project checks are not yet
  configured, never a false success;
- is idempotent, downgrade-aware, and reversible by `unlink`;
- leaves fresh clones runnable without a global `bench`;
- succeeds in empty repos, established repos, monorepos, paths with spaces,
  existing agent files, existing hooks/settings, and no-network mode;
- requires at most one shell command plus one prompt-driven setup conversation.

## Adoption, packaging, and first-hour findings

| ID | Sev | Finding and evidence | Improvement |
|---|---|---|---|
| A-01 | high | No registry package or immutable tag exists; the README's fallback uses mutable `#main` while all builds report `0.2.0`. `CHANGELOG.md` still has a large Unreleased section and an obsolete 0.2.0 command vocabulary. | Define dev SHA/dirty and immutable release versions; publish only after artifact gates pass. |
| A-02 | medium | README starts with `/bench-setup-repo`/`$bench-setup-repo` before link has installed those surfaces, and neither link output nor quickstart clearly requires a harness reload. | Lead with install; print the exact reload and next invocation. |
| A-03 | medium | `init` and setup point at `projects/<name>.md` examples that consumers cannot access; linked `BENCH-reference` names a missing `.bench/skills-index.sh`; shipped examples do not satisfy the current setup schema. | Ship one canonical generic profile/template set and validate every linked pointer. |
| A-04 | medium | A pre-existing project `CLAUDE.md` is preserved byte-for-byte, but link still declares all-harness success even if Claude never imports `AGENTS.md` or `.bench/BENCH.md`. | Add a marker-owned import block policy or an explicit red doctor row. |
| A-05 | medium | Relink writes the new plan and replaces the manifest without pruning old rows (`internal/adopt/link.go:308-330`). Removed kit skills remain active, become unowned, and cannot be unlinked. | Reconcile old/new manifests transactionally; preserve modified files and retain ownership until disposition. |
| A-06 | medium | Link copies an ignored `.bench/dist/bench`. A fresh clone loses it, and the linked wrapper suggests `scripts/go-build.sh`, which does not exist in a consumer repo. | Make the tracked local launcher self-bootstrap from a pinned artifact/version or ship a self-contained repo-local runtime. |
| A-07 | medium | Link preflights conflicts, but later sequential writes have no rollback on I/O failure. | Stage writes, fsync where relevant, then atomically promote or roll back. |
| A-08 | medium | Platform packages contain only `package.json` and `bin/bench`; they omit LICENSE, SECURITY, dependency notices, SBOM, and checksums. | Put required legal/security/evidence files in every independently published artifact. |
| A-09 | medium | Four platform packages publish sequentially before the wrapper. A partial immutable publish is not normally rerunnable. | Preflight all names, publish to a staging tag, verify already-present digests, poll dependencies, publish wrapper last, then promote. |
| A-10 | medium | The documented uninstall may run `rm -f "$(command -v bench)"` after npm removal and delete a foreign executable; doctor can recommend the same for a path it classified foreign. | Add marker-verifying `bench shim remove`; never print raw deletion for unowned paths. |
| A-11 | medium | Linux x64 is dynamically linked to glibc although docs promise Linux generally. A static `CGO_ENABLED=0` build succeeds. Release binaries retain debug data. | Publish/document an exact matrix, build static where supported, strip release symbols, and run native/musl smokes. |
| A-12 | medium | There is no cold/offline archive, checksum set, or internal-registry example. Runtime repair is the only fallback. | Produce self-contained per-platform archives and npm tarballs that install with network disabled. |
| A-13 | low | `unlink` returns 0 and says Bench was removed even when modified managed files and the manifest remain. | Use partial/nonzero status and a machine-readable residual list. |
| A-14 | low | Repair has no fetch timeout, response-size/decompression limit, or concurrency-safe cleanup; package metadata lacks repository/homepage/bugs/support fields; naming is Bench/redbench/benchkit. | Bound repair, make cache promotion atomic, complete metadata, and settle one user-facing identity. |

## CLI, safety, and operational findings

| ID | Sev | Finding | Improvement |
|---|---|---|---|
| C-01 | medium | Coverage validation exits green when a spec has no coverage map. It stores only the maximum story number, so stories 1 and 3 make nonexistent story 2 valid; story 0 and reversed ranges also pass. | Require a map or explicit historical marker; validate exact story membership, positive IDs, and range direction. |
| C-02 | medium | Learnings, maps, roadmap, outline, and status helpers turn non-`ENOENT` read errors into authoritative empty states. Probes with directories/dangling entries exited 0. | Only absence means empty; permission, type, parse, and traversal errors must be structured unknown/failure states. |
| C-03 | medium | Default-branch resolution has two answers. `ResolvedDefault` handles a sole `master`; `DefaultBranch` fabricates `main`, and diff/roadmap use the latter. A master-only repo made `roadmap --context` fail. | Single-source resolved-default facts and make callers handle “unknown.” |
| C-04 | medium | Intent lock stale takeover is age-based even for a live PID; competing reclaimers and unconditional release can remove a successor's lock. | Use identity-safe locking (`flock` or nonce/rename verification); never age out a proven-live owner. |
| C-05 | medium | Normal subshell completion releases its warm pool worktree but does not complete its intent. Existing paths are treated as live, so `status --all` showed 12 stale “resume interrupted work” rows. | Complete on normal exit and correlate crash recovery with the active lease identity. |
| C-06 | medium | Default `bench outline` emitted 169,050 bytes/2,193 lines with no cap or truncation metadata. | Return a bounded summary with counts and `truncated`; reserve `--full` for deliberate use. |
| C-07 | medium | Worktree acquisition silently runs unbounded `git fetch origin`, may prompt, and ignores failure. | Make fetch explicit or bounded, noninteractive, observable, and disableable. |
| C-08 | medium | Shift objectives appear in process arguments, stdout, a scratch file, intent records, and durable commit subjects. Control characters/classified text have no policy or redaction. | Transport prompts through stdin/private files, separate objective ID from text, sanitize commit subjects, and document data handling. |
| C-09 | medium | `bench gate nonsense` and several query commands accept trailing garbage; coverage slug lookup changes with current directory. | Centralize exact arity and anchor repo-relative resolution at the root. |
| C-10 | medium | Model discovery has no subprocess deadline, runs providers sequentially, and uses unbounded `io.ReadAll`. | Add bounded reads/deadlines and query independent providers concurrently. |
| C-11 | low | Leading-dash commit paths have no clean `--` interface; naming a directory does not authorize its changed children; help exits as an error. | Define a conventional argument grammar and contract-test hostile filenames. |
| C-12 | low | Capability-dependent symlink/hook tests skip on unsupported filesystems; the one-minute concurrency deadline is coupled to a similarly sized shell timeout. | Report skipped security classes as explicit evidence and separate the deadlines. |

## Skills and workflow consistency

The skill set has good conceptual seams: gate authoring, specs, TDD, review,
delegation, line selection, design systems, and ADRs largely point to one
discipline owner. The problems are semantic integration and consumer surface
area, not a need for more skills.

| ID | Sev | Inconsistency or redundancy | Improvement |
|---|---|---|---|
| S-01 | high | `link` copies every command and skill (`internal/adopt/link.go:104-107`), including kit-only `bench-assess`, `bench-update-kit`, and `craft-synthesis`. In a consumer repo these can replace root assessments, pull upstream HEAD, and interpret the consumer's README/CHANGELOG as Bench provenance. | Define one canonical consumer payload manifest. Keep maintainer-only surfaces in the kit; give consumers a versioned `upgrade/relink` path. |
| S-02 | high | README says a clear idea can skip shaping, while `bench-write-spec` refuses every feature without a closed map (`bench-write-spec.md:26-35`) and shape says even a clear idea must create one (`bench-shape-idea.md:22-26`). | Either make shape mandatory everywhere or let write-spec create a zero-fog map inline. |
| S-03 | high | `craft-delegate` permits inline code only for one source-line (`SKILL.md:15-23`), but implement says any no-spec lighter change may remain inline and calls it the sole exception (`bench-implement-spec.md:35-49`). Harnesses without write subagents must stop for ordinary work. | Choose one capability-aware delegation policy and state it only in craft-delegate; adapters reference it. |
| S-04 | high | Implement runs a finishing `bench commit --spec`, then final-check says it owns the landing gate/commit/status flip. A clean reviewed branch can arrive at final-check with nothing left to commit. | Give the state transition and landing commit to one phase; make the other a green no-op/reporting handoff. |
| S-05 | high | BENCH invariant 4 says commit only on green, but debug instructs committing a failing repro before a shift (`bench-debug.md:98-104`). | Keep the red observation uncommitted and preserve it safely, or explicitly reopen/change the invariant with reviewer approval. |
| S-06 | medium | What-next requires one `roadmap --context` snapshot but does not request `--full`, then rereads the journal directly. The default context truncates. | Use one `--context --full` snapshot and consume its schema as the one input. |
| S-07 | medium | Copyable final-check/implement examples show `bench commit -m "<msg>"` without the required paths. The CLI exits 2. | Make examples executable and add argv conformance tests for documented commands. |
| S-08 | medium | `craft-line`'s unquoted YAML description contains `#2`; YAML treats it as a comment. The live catalog description was truncated, losing important triggers. | Quote frontmatter scalars and parse real YAML in conformance/canary tests. |
| S-09 | medium | Setup requires concrete profiles, hostile inputs, and model bindings, but the referenced templates are not linked and current shipped examples lag the schema. | FT76 should own a validated, shipped seed set rather than prose pointing outside the consumer tree. |
| S-10 | medium | Skills-index derivation exists in shell and Go, while linked prose points at a shell helper that is not linked. This violates the repo's one-source-per-fact standard. | Pick one generator/query implementation and derive every index and check from it. |
| S-11 | medium | Design-it-twice says to paste briefs verbatim, but those briefs omit model, effort, iteration cap, paths, and return shape required by craft-delegate. | Make briefs task bodies embedded inside a complete delegation charge, with an inline fallback when delegation is unavailable. |
| S-12 | medium | Craft-synthesis recognizes only upstream- and learnings-sourced candidates; reviewer-directed assessment findings have no legal origin even though they can change the kit. | Add an assessment/local-reviewed origin or route every such item through one explicit synthesis entry. |
| S-13 | medium | `CONTEXT.md` defines the line with a rough token cap, while BENCH/craft-line define an iteration cap and say token budgets are not a stopping rule. | Use the canonical line definition everywhere. |
| S-14 | medium | Regroup's README example suggests a broad shift even though write-spec permits shift only when every story is cheap and fully gate-observable. | Use an example that satisfies routing or direct it to interactive implementation. |
| S-15 | low | Shape's bootstrap wording can be read as stopping after one frontier question, while craft-grill says continue until fog is gone. Design-system's trigger fires for all UI work, but its body assumes a design system already exists. | Clarify phase-vs-skill termination and add the no-design-source branch. |
| S-16 | low | Stale references remain: skills index location, README's `.claude/skills` tree, a nonexistent `lines-env.sh` port comment, and the AXI query inventory omits `outline`/`roadmap --context`. | Run a semantic reference/inventory sweep from canonical registries. |
| S-17 | low | Six closed decision maps remain in `decisions/` while `bench maps` reports zero; one carries superseded delegation wording. `skills-assessment.md` is an unreferenced 228-line historical report with already-implemented recommendations stated as current. | Retire/delete closed maps with their spec lifecycle and remove historical assessment artifacts from the live tree; Git is the archive. |
| S-18 | low | Craft-skills' warning against contrastive examples is not consistently dogfooded in grill/line/spec prose. The mutable `https://axi.md` reference also weakens reproducibility. | Apply the authoring lint to first-party skills and pin external normative references. |

Efficiency gains should therefore come from **removing surfaces and duplicate
derivations**, not from adding another umbrella skill: split consumer and kit
maintenance payloads, single-source the skills/CLI inventories, bound large
query output, consume one full roadmap snapshot, delete retired maps/history,
and make setup reuse the existing link/init/setup seams.

## Gate, tests, records, and release evidence

### What is strong

- `bin/bench.sh gate` passed every phase, contract package, conformance suite,
  and canary.
- `go test -race -count=1 ./...`, `go vet ./...`, `go mod verify`, and
  `git diff --check` passed.
- Two clean local builds were byte-identical with SHA-256
  `82a60ccd770de9f203b20e8bbc26d4713af75ff77afe37770043a5e475a11bc9`.
- The current binary records commit `cd43632` and `vcs.modified=false`.
- The release workflow pins action SHAs, uses minimal declared permissions,
  includes OIDC provenance, pins the vulnerability scanner, and publishes the
  wrapper after platform packages.
- Link/unlink use manifest fingerprints and path containment; source-tree
  conflict preflight is conservative; the dependency set is small and pinned.

### Why green is not release approval

- The gate has no vulnerability phase and is not called by release.
- Artifact tests do not install and operate the exact tarball produced by the
  workflow.
- Existing contracts intentionally encode foreign-worktree removal and
  arbitrary hook description execution.
- No contract changes the gate identity between gate and commit.
- No contract forces a post-agent commit failure and proves recovery.
- Release conformance is textual rather than behavioral.
- ADR 0003 still says every fixture runs the full inner gate although fixtures
  now select an owning phase. ADR 0004's clean-publish-tree claim is false.
- The release emits no checksum manifest, SBOM, third-party notice, durable test
  report, or record binding tag, commit, toolchain, gate result, and package
  digests.

## Prior-assessment reconciliation

The previous backlog was checked against current behavior before this report
was written. Implemented items were removed rather than carried as findings:
the package identity moved to `redbench`; clone build instructions, link-owned
`CLAUDE.md` reversal, coverage missing-argument text, live default-branch
pre-push resolution, package exclusion of `projects/benchkit.md`, symlink-hop
limits, Codex hook timeout, salvage status guidance, cache/marker constants,
worktree classification errors, CLI unknown-command exit, canary success
output, ADR 0002's cache wording, July 7 CHANGELOG entries, and the earlier
lease-reclaimer identity defect are implemented.

The old claims that tarball adoption, no-global relink, gate-cache reuse, and
worktree cleanup were sound do **not** survive executable artifact/adversarial
testing; they are replaced by R-01, R-02, R-04, R-05, and R-09. The earlier
profile-template concern remains as A-03. Capability-skipped security tests
remain as C-12.

## Ranked improvement plan

Order is based on user-state safety first, then oracle integrity, then a
releasable artifact, then onboarding polish.

| Priority | Work | Evidence required | Rough agent time |
|---:|---|---|---:|
| 1 | Ownership-safe worktree cleanup; disable broad SessionStart cleanup; detached-commit protection | Two black-box preservation contracts and a biting canary | 1–2 days |
| 2 | Gate-cache oracle fingerprint and atomic/fail-closed verdict writes | Gate A green → gate B red must refuse commit; write-failure test | 1 day |
| 3 | Shift failure recovery and truthful completion; staging/adapter/cap handling | Missing identity/hook/stage/adapter/cap tests retain state and return honest codes | 1–2 days |
| 4 | Stop executing unknown hooks; static managed guard metadata and total deadline | Unknown sentinel never runs under guards or SessionStart | 0.5–1 day |
| 5 | Clean package staging and command-aware runtime selection | Four-host resolution tests; wrapper allowlist; pack/install/setup/relink/unlink flow | 2–3 days |
| 6 | Make release consume the gate; update Go; enforce version/changelog/tag and artifact smoke | Current scanner green; release-step canary; dry-run manifest | 1–2 days |
| 7 | Fix installed shim/adoption routing and fresh-clone local runtime | Global and npx install interaction tests with network disabled after install | 1 day |
| 8 | Shape/build FT76 one-command bootstrap | Empty/established/monorepo fixtures; idempotency, rollback, preservation, reload handoff | 3–5 days |
| 9 | Transactional relink/upgrade/unlink and managed instruction composition | Upgrade/downgrade plus modified/stale asset matrix | 2–3 days |
| 10 | Consumer-vs-maintainer payload split and skill phase reconciliation | Payload manifest contract; no kit-only commands in consumer; workflow scenario tests | 2–3 days |
| 11 | Coverage/read-error/default-branch/intent correctness batch | Focused red-to-green AXI and runtime contracts | 2–3 days |
| 12 | Offline/reproducible release evidence and package governance | Per-platform archives, checksums, SBOM/notices, support/security docs | 2–3 days |
| 13 | Output/resource/CLI consistency and stale-doc cleanup | Bounded output schemas, arity corpus, docs command probes | 1–2 days |

Minimum credible release-blocking work is approximately **9–14 agent-days**
plus native-platform execution and publisher access. A polished one-command,
offline-capable, established-repository experience is more realistically
**18–28 agent-days**. Those ranges assume the architecture is retained; a
reviewer decision that changes instruction-file ownership or distribution may
move them.

## Exit criteria for reassessment

Do not publish merely because the current gate is green. Reassess when all of
the following are true:

1. R-01 through R-09 have executable regression tests and are closed.
2. `govulncheck` and the full repository gate pass on the same tagged checkout.
3. The exact release-shaped wrapper contains only allowed files and installs
   the correct binary on every supported platform.
4. Pack/install/setup/doctor/relink/unlink and fresh-clone flows pass from an
   isolated prefix without a source checkout.
5. The npm version, tag, changelog, binary version, commit, and artifact
   manifest agree exactly.
6. Publishing is staged, resumable, and wrapper-last with digest verification.
7. FT76's one-command flow preserves existing instruction/config/hook content
   and is idempotent and reversible.
8. Consumer installations exclude kit-maintainer-only skills and commands.
9. Security, data-handling, support, license/notice, SBOM, checksum, rollback,
   and supported-platform artifacts ship with the release.
10. A clean room user can complete setup from README alone in one shell command
    plus at most one harness-native conversation.

## Verification and limits

Executed during this assessment:

- full Bench gate and canary;
- full Go race suite, vet, module verification, and diff check;
- the exact release `govulncheck` version;
- packed and release-shaped npm dry runs;
- packed/git install, global lifecycle, stable-shim, relink-with-removed-asset,
  and fresh-clone probes;
- all four cross-builds, binary resolution/format inspection, static and
  stripped comparison, and byte-for-byte reproducibility comparison;
- throwaway reproductions for foreign/detached worktree deletion, gate-cache
  authority change, shift commit failure, failed adapter false success,
  unwired hook execution, coverage gaps, false-empty reads, and default-branch
  divergence;
- static review of every craft skill, every workflow command and Codex adapter,
  the launcher/core, tests, package generator, workflow, README, project
  profile, ADRs, roadmap, and prior assessments.

Limits: nothing was published; publisher credentials and registry organization
configuration were not assessed; GitHub Actions itself was not executed;
native macOS/arm64, Linux arm64, musl/NixOS, and Windows execution were not
available. Registry/native-runner claims therefore remain release evidence to
collect, not assumptions. The bank-style repo-controlled compliance decision
and remediation evidence are in `COMPLIANCE_ASSESSMENT.md`.
