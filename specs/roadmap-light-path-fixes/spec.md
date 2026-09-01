# Roadmap light-path fixes

Status: staged

Roadmap: FT92, FT94, FT104, FT117, FT164, FT166, FT178, FT179, FT268, FT270

Decision source: the ten named roadmap rows, checked against the live tree

Verification log: 3 iteration(s) to accept — the first review rejected the
multi-ticket light-path folder; the second split aggregate slices and corrected
scope; the third serialized the final shared registry owner. The reviewer
approved the breakdown and bench worktree shell token.

## Problem

Ten roadmap rows describe bounded fixes that are ready for implementation.
Several rows contain more than one independently green change. Treating each
row as one light-path ticket creates broad fences, hidden merge conflicts, and
acceptance that cannot identify one owner.

Some roadmap behavior has also landed in part. The learning writer, system
check, and full finding loop already exist. Rebuilding those facts would create
duplicate ownership.

## Solution

Keep one reviewed umbrella for scheduling and retirement. Split every roadmap
row at its executable seam. Each ticket owns one small behavior, focused
evidence, and an exact write fence.

Preserve behavior already present in the live tree. Tickets implement only the
remaining delta. Blockers order behavioral dependencies. The coordinator
explicitly serializes each otherwise-frontier pair whose `Writes:` entries
overlap, including required registry and fixture closure paths, even when no
blocker edge exists. Nonconflicting frontier tickets remain eligible for
parallel authoring.

## User stories

1. As a gate operator, I can identify drift and receive declared-input protection in linked repositories.
2. As a maintainer, I can change one resume test golden without duplicating its expected format.
3. As a coordinator, I stop bounded flaky retries and avoid aggregate grading while delegated tests remain live.
4. As a CLI caller, parser-owned verbs provide one consistent usage contract.
5. As a reviewer, repair charges and done claims resolve to concrete evidence.
6. As a contributor, sanctioned writers preserve primary-local capture grammar.
7. As an automation caller, worktree dispatch cannot fall into an interactive shell accidentally.
8. As a maintainer, high-stakes APIs and comment rules state durable contracts.
9. As a gate maintainer, prospective artifact failure edges have direct tests and one owner per fact.
10. As a test caller, named checks and prose findings are discoverable and complete.

## Implementation decisions

- FT92 splits transaction attribution from the consumer gate scaffold.
- FT94 shares only the test expectation. Production rendering remains an independent subject.
- FT104 separates coordinator guidance from the runtime fold retry.
- FT117 gives spec parsing, doctor routing, and commit usage separate owners.
- FT164 separates the repair template, done-claim proof, and installed-lane fallback.
- FT166 adds only the missing retrospective writer. The learning writer remains unchanged.
- FT178 preserves the human shell as bench worktree shell. Bare dispatch becomes usage-only.
- FT179 names each documented API group. No repository-wide comment sweep is allowed.
- FT268 gives each residual one owner. Shared test files create explicit blocker edges.
- FT270 preserves the existing system check and full finding loop.
- FT270 adds only discovery, prose routing, skip disclosure, and gate-prose diagnostics.
- The FT281 files remain outside every fence. In particular, no ticket owns tag_census files.
- Each ticket runs focused checks. Integration folds and whole-tree gates remain serial.

## Testing decisions

- CLI changes use parser-level tests and built-binary journeys where dispatch is the seam.
- Guidance changes update their canonical anchor tests without copying the whole rule.
- Process lifecycle changes cover normal exit, interrupt, and termination.
- Prospective artifact changes use the existing owner and system-test families.
- Prose behavior uses the canonical prose grader and exact stdout or stderr contracts.
- Existing behavior receives preservation checks only where its caller could regress.

### Seam diagram

    roadmap row
        |
        v
    focused ticket -> package or CLI seam -> focused evidence
        |
        v
    retained integration source -> serial gate -> reviewed landing

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| LF1 | 1 | drift identifies tree movement or the exact manifest path | gate transaction outcome tests | a generic refusal cannot satisfy the expected subject |
| LF2 | 1 | linked consumers reject ignored declared inputs through an explicit scaffolded consumer-only scope | setup consumer-gate journey and named-check environment ownership | path-derived scope or benchkit-only conformance cannot make the consumer red |
| LF3 | 2 | two resume tests share one expected-format helper | worktree unit and runtime tests | a second literal violates the one-source rule |
| LF4 | 3 | coordination stops after two proven flaky refusals and waits for quiet delegates | delegation anchors and behavior tests | prose without the stop and quiet conditions remains incomplete |
| LF5 | 3 | an empty-reason infrastructure fold gets one verified retry | worktree merge tests | a second retry or an unverified retry fails the counter |
| LF6 | 4 | spec arguments use parser-standard help and refusals | internal/spec tests | hand parsing diverges on help or unknown flags |
| LF7 | 4 | doctor uses the parser and nested routing is graded accurately | adopt and routing tests | the old switch or stale exemption remains observable |
| LF8 | 4 | commit usage errors print one flat contract | internal/commit tests | nested usage text remains byte-visible |
| LF9 | 5 | every repair charge uses the named five-field template | delegation anchors | an improvised charge omits an enforced field |
| LF10 | 5 | done claims resolve owners and preserve ticket attribution | delegation anchors | a missing artifact or umbrella ledger remains acceptable |
| LF11 | 5 | installed-lane repair uses the candidate commit core and rebuilds after landing | delegation anchors | the fallback can bypass snapshot grading |
| LF12 | 6 | a retrospective writer validates and round-trips one primary-local artifact | roadmap and retros tests | malformed input or overwrite behavior fails directly |
| LF13 | 7 | bare or unknown worktree dispatch acquires nothing | command registry journey | accidental fallback creates observable state |
| LF14 | 7 | the explicit worktree shell releases or leaves a reclaimable lease on signals | worktree lifecycle tests | a killed shell leaves an unreclaimable owner |
| LF15 | 8 | comment and review skills enforce the remaining FT179 rules | skill anchor tests | new provenance or reviewer register escapes review |
| LF16 | 8 | release-evidence boundaries state their contracts | review-owned: enumerate the functions named by the LF16 ticket | an omitted named function fails reviewer enumeration |
| LF17 | 8 | preflight and marker-wait boundaries state grammar and failure contracts | review-owned: enumerate the functions named by the LF17 ticket | caller history cannot substitute for each named contract |
| LF18 | 8 | gate resolution and execution boundaries state contracts | review-owned: enumerate the functions named by the LF18 ticket | an omitted named oracle function fails reviewer enumeration |
| LF19 | 8 | worktree lifecycle and wrapper dispatch state contracts | review-owned: enumerate the functions named by the LF19 ticket | an omitted named automation boundary fails reviewer enumeration |
| LF20 | 9 | prospective execution has one root path and proves bundle confinement | prospective owner test | the test-only branch or empty-root setup escapes |
| LF21 | 9 | one checkout-name owner serves production and system tests | prospective artifact system journey | a retyped child name can drift |
| LF22 | 9 | Open refuses absent and dangling temporary roots | prospectiveartifact tests | both filesystem states reach the same direct seam |
| LF23 | 9 | ReadPublished has a negative caller | prospectiveartifact tests | a failure path without a caller remains unproved |
| LF24 | 9 | permission tests state their live edge without retired citations | prospectiveartifact tests | stale row provenance remains in edited comments |
| LF25 | 10 | test help and unknown checks show one named-check inventory | testreport command tests | trial-only discovery or generic usage remains |
| LF26 | 10 | the prose named check reaches the canonical grader and prints every finding | testreport tests | truncation or a second grader loses findings |
| LF27 | 10 | a skipped root-conformance test cannot read as root green | testreport package result tests | a package pass masks the skip |
| LF28 | 10 | gate-prose help succeeds and findings include each sentence | gate-prose command tests | line-only output still requires another file read |

### Edge inventory

- Error paths: unknown flags, malformed capture, invalid roots, and failed folds remain fail-closed.
- Empty input: bare worktree, missing check names, and empty retrospective fields receive explicit behavior.
- Boundary values: the second flaky refusal stops; the first verified infrastructure refusal retries once.
- Interrupted state: an explicit shell signal leaves no unreclaimable worktree.
- Re-run idempotency: one retry stays one retry, and capture writers never overwrite prior entries.
- Hostile paths: manifest paths and prospective roots with spaces or glob characters remain literal.
- Partial implementation: existing learning, system-check, and full-finding behavior remains owned by current code.

**Won't handle** — FT249 durable capture storage — that roadmap row owns the storage decision.

**Won't handle** — untouched legacy provenance comments — FT179 limits removal to edited lines.

**Won't handle** — test filters inside the system suite — FT270 concerns named check routing.

## Ownership fences

- `specs/roadmap-light-path-fixes/`
- `reviews/roadmap-light-path-fixes.md`
- `internal/gate/run_transaction.go`
- `internal/gate/run_outcomes_test.go`
- `internal/adopt/setup.go`
- `internal/adopt/setup_test.go`
- `internal/adopt/adopt_test.go`
- `internal/conformance/validity_checks_test.go`
- `internal/conformance/registry/scope.go`
- `internal/conformance/subcommand_routing_test.go`
- `internal/testreport/command.go`
- `internal/testreport/check_test.go`
- `internal/worktree/`
- `internal/spec/`
- `internal/adopt/doctor.go`
- `internal/commit/`
- `internal/roadmap/`
- `internal/retros/`
- `internal/learnings/`
- `internal/testreport/`
- `internal/gate/gate_prose.go`
- `internal/gate/gate_prose_test.go`
- `internal/prose/`
- `internal/gate/gate.go`
- `internal/gate/decision.go`
- `internal/gate/engine.go`
- `internal/gate/phases.go`
- `internal/gate/prospective.go`
- `internal/gate/prospective_owner_test.go`
- `internal/gate/prospectiveartifact/`
- `internal/systemtest/owner_artifact_recovery_test.go`
- `internal/systemtest/adoption_test.go`
- `internal/releaseevidence/`
- `internal/preflight/`
- `internal/contract/`
- `cmd/bench/main.go`
- `cmd/bench/main_test.go`
- `cmd/bench/command_registry_test.go`
- `.agents/skills/bench-craft-delegate/`
- `.agents/skills/bench-craft-line/SKILL.md`
- `.agents/skills/bench-craft-comments/SKILL.md`
- `.agents/skills/bench-craft-review/SKILL.md`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`

## Out of scope

- Shared capture storage remains FT249.
- A general command-registry redesign is not required.
- Existing comment debt outside edited lines remains untouched.
- FT281 citation package scope remains independently owned.
