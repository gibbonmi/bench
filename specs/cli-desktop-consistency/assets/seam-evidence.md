# Specification seam evidence

The existing doctor, adoption transaction, and SessionStart route provide sufficient integration seams.
The new compatibility owner isolates interface facts and their evidence limits.
The approved source remains one compiled map, including its original probe and research assets.

Date: 2026-10-02
Baseline: `2eb7478db1098ef5d7e01168ba61045f371bdd65`
Consumed by: the staged CLI and desktop consistency specification
Drift: Recheck changed source owners and refresh affected runtime evidence before implementation claims compatibility.
Retire when: Implementation verification and the final qualification record supersede this authoring census.

## Source validity

Both structured source assets were reread before the map moved.
The reports still describe the recorded incident and its explicit unknowns.
Their prior source references match the unchanged implementation at the authoring baseline.
No successful current desktop recovery was observed during specification.
No new claim of full installed compatibility was derived from the sources.

The source move preserves the index, all tickets, and both assets.
The repository search found only three tracked absolute map references, all inside the index.
Those references now name the compiled asset location.
Local capture references are updated at the phase handoff.

## Current owners and concrete scenarios

| Scenario | Existing owner read | Selected evidence |
| --- | --- | --- |
| A healthy local command coexists with a broken desktop tool | `internal/adopt/doctor.go:182`; `.bench/hooks/session-start.sh:40` | CD02, CD08, CD15, CD39, CD40 |
| A consumer needs one shared installation | `internal/adopt/setup.go:21`; `internal/adopt/setup.go:304`; `internal/adopt/setup_report.go:15` | CD18 and CD63 |
| A managed file differs from the installed payload | `internal/adopt/link_transaction.go:131`; `internal/adopt/link_stage.go:17` | CD19 through CD22 |
| A successful repair must remain reversible | `internal/adopt/transaction.go:77`; `internal/adopt/transaction.go:114`; `internal/adopt/transaction.go:164` | CD26 through CD36 |
| A startup probe stalls | `internal/sessioninspect/sessioninspect.go:41`; `internal/sessioninspect/sessioninspect.go:67` | CD37, CD56, CD61 |
| A hook starts from a deeper directory | `.codex/hooks.json:8`; `.bench/hooks/session-start.sh:20` | CD12, CD38, CD63 |
| Two interfaces claim the same work | Existing Bench assignment and landing contract in `.bench/BENCH.md` | CD47 and CD48 |

The present transaction deletes successful backups at its terminal step.
Its rollback also ignores individual restoration errors.
CD28, CD30, and CD36 prohibit carrying those behaviors into the stronger compatibility-repair claim.
This change belongs to the adoption owner, not a second repair installer.

## Test and enforcement precedents

| Boundary | Precedent read | Consequence |
| --- | --- | --- |
| Local diagnostic row | `internal/adopt/doctor_gate_inputs_test.go:15` | Drive the real diagnostic producer with supplied and missing facts |
| Real adoption | `internal/adopt/link_transaction_test.go:79`; `internal/adopt/link_transaction_test.go:90` | Use the existing consumer fixture and actual Link path |
| Missing Bench core | `internal/systemtest/session_start_test.go:343` | Drive the actual hook through the system process owner |
| Outside a repository | `internal/systemtest/session_start_test.go:395` | Preserve informational silence outside the hook's scope |
| Aggregate startup timeout | `internal/sessioninspect/sessioninspect_test.go:20` | Preserve bounded startup and its existing exit posture |
| Safe record reads | `internal/bounds/classify.go:75`; `internal/bounds/classify.go:163` | Reuse bounded classification and link refusal |
| Command grammar | `internal/usage/parse.go:129`; `internal/adopt/doctor.go:176` | Validate flags before the selected operation reads paths |
| Command dispatch | `cmd/bench/command_registry.go:32`; `cmd/bench/main.go:142` | Extend doctor through its existing adoption route |
| Gate execution | `.bench/gate.sh:32`; `internal/gate/phases.go:46` | Ordinary tests and the sealed system suite remain the executed roots |
| Coverage citations | `internal/coverage/citations.go:109` | Planned tests are not citations to existing declarations |
| Fence equality | `internal/preflight/fence_writes.go:18`; `internal/spec/fences.go:25` | The spec fence must equal ticket ownership outside implicit paths |
| Completion plan | `internal/preflight/plan.go:17`; `internal/reviewrecord/plan.go:45` | Preflight reads the committed source tip, not an uncommitted plan |

Planned tests acquire exact existing-test citations after implementation creates their declarations.
The coverage check currently reports those rows as planned, uncited evidence.
That diagnostic is not evidence that the feature already passes tests.

## Reader and writer census

The hidden-file sweep excluded only Git metadata.
It covered commands, hooks, adapters, scripts, workflows, conformance tests, public docs, and fixtures.
No GitHub workflow caller of the changed doctor form appeared in that sweep.
The legacy forms remain unchanged.

| Changed fact | Readers or writers | Ticket |
| --- | --- | --- |
| Doctor grammar and help | `internal/adopt/doctor.go`, `cmd/bench/main.go`, `cmd/bench/help_inventory_test.go` | C1 and C2 |
| Command dispatch contract | `cmd/bench/command_registry.go`, its tests, subcommand routing, and AXI conformance | Derived ticket closure |
| Shared repository assets | `buildLinkPlan`, `transactionalLink`, `convergeSetup`, and `finishSetup` in `internal/adopt` | C2 |
| Managed destination changes | Setup, Link, Upgrade, Unlink, doctorFix, and the new undo path | C2 |
| Startup instruction | `.bench/hooks/session-start.sh`, `internal/sessioninspect`, and hook entry-point parity tests | C3 |
| Shared workflow obligation | `.bench/BENCH.md`, `.bench/BENCH-reference.md`, README, and `internal/anchors` | C3 |
| New durable repair data | Adoption repair record and `DATA_HANDLING.md` | C2 |
| Decision-map location | The moved map index and the local capture handoff | Specification phase |

The command-help needle is `bench doctor [--fix]`.
Its public rendering comes from `cmd/bench/main.go` and its exact inventory assertion lives in `cmd/bench/help_inventory_test.go`.
The existing by-path hint remains valid and need not be rewritten.
The registry closure proposal supplies every bound file before the ticket graph reaches sign-off.

The existing startup hint, `bench CLI:`, keeps its meaning and bytes.
The new obligation adds a separate message rather than renaming that hint.
Entry-point parity tests must cover the added output through both the hook and its existing command route.

The payload already ships the agreement, reference guide, README, and hook tree.
No new payload root or file-count grant is required for documentation.
The implementation must split responsibilities within existing budgets instead of editing reviewer-owned grants.

## Source-clause coverage

| Decision source | Approved clauses | Coverage |
| --- | --- | --- |
| Ticket #1 | Independent chats and repository continuity | CD46, CD57, CD58 |
| Ticket #2 | WSL CLI, desktop WSL agent, preserved Linux and macOS, excluded native Windows | CD01, CD45, CD49, CD50, CD52, CD53 |
| Ticket #3 | Commands, files, skills, hooks, permissions, worktrees, reviews, recovery | CD07, CD39, CD45, CD47 through CD51, CD57 through CD59 |
| Ticket #4 | One setup, automatic start and resume checks, reversible repairs, interruption and policy boundaries | CD18 through CD24, CD26 through CD38, CD61, CD63 |
| Ticket #5 | Concurrent repository use and isolated writers | CD29, CD47, CD48 |
| Ticket #6 | Shared Bench contract, separate personal state, conflict reporting | CD03 through CD08, CD14, CD21 through CD23, CD33, CD46 |
| Ticket #7 | Affected-work refusal, recovery action, unaffected work | CD06, CD07, CD41 through CD44, CD60, CD64 |
| Ticket #8 | Supported repair authority and explicit private-runtime decision | CD23 through CD25, CD60 |
| Ticket #9 | One complete delivered outcome | CD18, CD45, CD54 |
| Ticket #10 | Negative actual-interface evidence and remaining qualification | CD02, CD08, CD15, CD39, CD40, CD45, CD54 |
| Ticket #11 | Supported routes, configuration boundaries, unknown upstream repair | CD03, CD04, CD08, CD23 through CD25, CD43, CD60 |
| Ticket #12 | Capability evidence, invalidation, live retest, required equivalents | CD13, CD15, CD38 through CD45, CD55, CD56, CD62, CD64 |

## Hostile-input dispositions

| Profile edge class | Disposition |
| --- | --- |
| Spaces, glob characters, unquoted arguments, and deep cwd | CD09 and CD12 use real command arguments and decoy paths |
| Control bytes and sink-permitted line separators | CD16 drives TOON with ESC, BEL, tab, newline, and return |
| Numeric-looking cells | CD16 uses the existing TOON encoder with leading-zero identifiers |
| Git C-quoted patch headers | No new patch parser; existing review and landing routes survive under CD48 |
| Unicode whitespace and hand-edited fields | CD09 rejects invalid command operands; CD34 rejects malformed record identifiers |
| A write changes the facts reported | CD20 verifies repeated tracked repair; CD44 requires post-repair retest |
| Relative file references and trailing spaces | CD12 resolves repository identity through the existing Git owner |
| Missing trailing newline | CD10 and CD11 include complete records without a terminal newline |
| Absent versus empty files and directories | CD11, CD26, CD31, and CD34 retain distinct dispositions |
| FIFOs, devices, sockets, and dangling links | CD10 and CD34 refuse unsafe record reads before opening |
| Live input links and symlink parents | CD34 refuses repair-record redirection; CD35 rechecks destination identity |
| Symbolic refs in Bench namespaces | No new ref writer; existing worktree lifecycle survives under CD47 and CD48 |
| Raw Git substituted into lifecycle guidance | CD48 requires the actual Bench lifecycle in the live transcript |
| JSON-escaped shell operators and operator runs | Existing command guard remains authoritative; CD49 and CD50 use its harmless forbidden fixture |
| A flag value mistaken for a positional | CD09 supplies interface-looking and flag-looking values in invalid positions |
| Prose tokens mistaken for live grammar | No new Markdown parser; coverage and completion-plan readers remain authoritative |
| Non-ASCII whitespace in authored Markdown | Prose and existing document-validity checks grade these artifacts |
| Missing PATH tools and symlink invocation | CD07 and CD12 drive by-path and symlinked wrapper forms |
| Every shipped command and hook surface | CD07, CD18, CD38, and CD63 cover kit, consumer, hook, and adapter integration |
| Foreign registrations, reused paths, dirty state, and plan drift | CD32, CD35, CD47, and CD48 preserve the existing ownership boundary |
| Add/add merge conflicts on generated paths | No new Git composer; CD48 uses the unchanged landing owner |
| Fresh-process reload after serialization | CD27 interrupts and resumes through a second process |
| Non-TTY input | CD09 asserts that the new modes never prompt; setup keeps its existing confirmation contract |
| Host filesystem stalls | CD61 preserves the existing aggregate bound; CD26 through CD28 retain failure evidence |
| Unterminated Markdown delimiters | Existing prose and source-validity checks remain authoritative |
| A symlinked temporary root | CD12 uses the canonical repository identity from the existing Git owner |
| Fast-lane private composed checkout | No new branch-range checker; the existing commit lane grades the composed documents |

The implementation checks both present and missing required capabilities.
It checks optional absence separately from required absence.
It drives runtime, configuration, and workspace fingerprint changes independently under CD13.
It uses canned facts for report comparisons instead of comparing two elapsed-time-dependent live reports.

## Pre-review proof checklist

Cited symbols: Doctor, Setup, Link, Upgrade, Unlink, transactionalLink, promoteAll, rollback, Inspect, and Command resolve in the named source files.
Import edges: Existing adoption, session inspection, bounds, TOON, Git, usage, and command packages passed `go list`.
The new compatibility package is planned and does not yet compile as a package.
Its consumers must pass `go list` after implementation creates it.
It must not import either adoption or session inspection, which would create a cycle through its consumers.

Source-row clauses and occurrences: The table above covers all twelve decision tickets in the sole compiled source.
Promised field labels: `context{field,value,source}`, `checks{check,state,action}`, and `live{capability,action}`.
Changed-function callers: The reader and writer census names each affected production caller and its direct dispatch helper.
Copy survival: No production owner is replaced by a copied implementation.
CD30 removes the final retained backup entry to prove that every preimage is necessary.
Rendered-shape readers: The doctor help inventory and startup parity surfaces are named above and closed by preflight.

The checks do not authenticate the active harness or a same-user adversary.
The spec states the existing trusted entry points instead of adding a self-attestation claim.
Live probes qualify behavior only after the actual interface produces evidence.

## Verification and review boundary

The independent review round is the reviewer's spec-and-ticket sign-off, as the project profile requires.
The author retains writes under the user's explicit authorship instruction.
The first author check found planned tests represented as existing declarations.
The corrected map keeps future evidence explicitly planned.
No implementation tests were created to make specification checks pass.

Prose, coverage, the complete ticket preflight, and each derived ownership-closure proposal passed before sign-off.
The reviewer approved the implementation line, seams, coverage, fences, exclusions, and complete ticket graph on 2026-10-02.
