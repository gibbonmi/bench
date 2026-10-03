# Consistent Bench behavior in CLI and desktop

Status: staged

Decision source: `specs/cli-desktop-consistency/decisions/cli-desktop-consistency.md` (ready compiled map).

Verification log: 1 iteration(s) to accept — the reviewer approved the specification and ticket graph on 2026-10-02.

## Problem

The CLI can run Bench while the desktop fails before its shell starts.
The approved probe records that difference under normal permissions.
Bench's current startup hook can print useful context without proving that the active chat tool works.

The probe and research remain in the compiled decision source.
Their negative result does not identify the deletion actor or prove the cause of the desktop failure.
The complete workflow and successful desktop recovery remain unqualified.

## Solution

One repository setup supplies the shared Bench contract to independent CLI and desktop chats.
Each chat checks the capabilities needed for its current operation.
Checks distinguish inspected configuration from actual tool behavior.
Supported reversible repairs preserve settings, and unresolved failures identify the affected work and the next supported action.

Qualify Codex CLI in WSL2 and the Windows desktop app with its agent in WSL2.
Preserve existing macOS and Linux behavior.
Setup, automatic checks, and recovery form one complete delivery.
A green unit suite without the required live evidence cannot complete this specification.

## User stories

Line: gpt-6-astra / high.
Implementation-line reason: Recovery and active-interface proof are the hardest boundaries. The approved behavior is precise, but live evidence remains weak and fixture coverage needs new cases.
Harder chunks: C2, C3.

### Setup and identity

1. As a Bench user, I want one repository setup for both interfaces, so that I do not maintain separate Bench installations.
2. As a Bench user, I want the active interface identified, so that a terminal result cannot stand in for a chat result.
3. As a Bench user, I want configuration sources identified, so that conflicting settings have a clear owner.
4. As a Bench user, I want personal settings preserved, so that Bench consistency does not synchronize my accounts or preferences.
5. As a Bench user, I want missing Bench assets reported, so that I can restore the required integration.
6. As a Bench user, I want absent tools diagnosed, so that a missing PATH entry has a concrete recovery route.
7. As a Bench user, I want unknown effective settings reported honestly, so that a file declaration cannot imply active enforcement.
8. As a Bench user, I want malformed inputs refused, so that diagnostics cannot hang or inspect an unintended path.

### Repair and preservation

9. As a Bench user, I want reversible automatic repairs, so that routine damage does not require repeated setup.
10. As a Bench user, I want changed user files preserved, so that a repair cannot overwrite my work.
11. As a Bench user, I want repairs bounded by authority, so that Bench cannot restart work or weaken permissions without my decision.
12. As a Bench user, I want failed repairs recoverable, so that interruption cannot destroy the previous state.
13. As a Bench user, I want concurrent repairs serialized, so that two chats cannot overwrite the same shared asset.
14. As a Bench user, I want an undo action, so that I can restore a completed automatic repair.
15. As a Bench user, I want private repair records, so that reversible repairs do not publish local data.

### Session behavior and qualification

16. As a Bench user, I want automatic startup checks, so that a new chat discovers missing capabilities before dependent work.
17. As a Bench user, I want automatic resume checks, so that old success cannot authorize a changed environment.
18. As a Bench user, I want actual chat-tool evidence, so that a hook or escalated shell cannot falsely qualify normal execution.
19. As a Bench user, I want only affected work blocked, so that an optional desktop failure does not stop unrelated work.
20. As a Bench user, I want recovery guidance without a working shell, so that process-start failure does not create a dead end.
21. As a Bench user, I want recovery verified in the failed interface, so that a successful repair command cannot conceal continuing failure.
22. As a Bench user, I want independent chats, so that Bench behavior does not depend on synchronized conversations.
23. As a Bench user, I want isolated concurrent writers, so that CLI and desktop work share the existing assignment rules.
24. As a Bench user, I want consistent hooks and permissions, so that both interfaces enforce the same Bench contract.
25. As a Bench user, I want equivalent review outcomes, so that optional desktop presentation tools do not control completion.
26. As a Bench user, I want preserved macOS and Linux behavior, so that WSL qualification does not regress existing installations.
27. As a Bench user, I want explicit support boundaries, so that native Windows execution is not implied.
28. As a Bench user, I want one complete delivered outcome, so that setup alone cannot count as compatibility.
29. As a Bench user, I want current capability evidence, so that runtime version numbers do not substitute for working operations.
30. As a Bench user, I want consistent file access and skills, so that both chats can complete the full workflow.

## Implementation decisions

### Diagnostics and evidence

Extend the existing doctor command with an opt-in compatibility mode.
Keep legacy doctor behavior unchanged, including its existing exit contract.
The new mode has these planned forms:

```text
bench doctor --compat <codex-cli|codex-desktop>
bench doctor --compat <codex-cli|codex-desktop> --fix
bench doctor --compat <codex-cli|codex-desktop> --undo <repair-id>
```

The first form inspects and writes no configuration or health cache.
The mutation flags are mutually exclusive.
An unknown interface, missing value, or invalid flag combination exits 2 before inspection.
A completed local inspection exits 0 only when its required local checks pass.
A failed or unknown required local check exits 1 with its exact next action.
A zero exit never certifies the entire active chat.

The report separates context, local checks, and live obligations.
Its TOON tables are `context{field,value,source}`, `checks{check,state,action}`, and `live{capability,action}`.
The check states are `ok`, `failed`, `unknown`, and `not-required`.
An explicit interface argument is a caller declaration, not an authenticated runtime identity.

Context names the repository, execution environment, configuration home, active runtime, and relevant policy provenance.
It distinguishes an observed active runtime from an installed launcher.
Unknown active metadata remains unknown.
No hard-coded Windows username or shared Codex home enters the contract.
Configuration inspection must not infer effective trust or permissions from one configuration file.
Unsupported configuration formats or unavailable effective sources remain unknown.

A single compatibility owner defines capability requirements, evidence provenance, and recovery classifications.
Doctor and the startup integration consume that owner.
The owner reuses the existing bounded file readers, process bounds, Git root resolution, and TOON renderer.
The adoption owner retains configuration writes and its canonical payload plan.
No second installer, shell parser, or configuration-precedence implementation is introduced.

The capability inventory comes from the approved workflow classes and the selected operation.
It includes commands, files, rules, skills, hooks, permissions, worktrees, reviews, and recovery.
Desktop-specific tools enter only when the operation has no verified equivalent route.
The producer inventory also supplies report rows, so counts and advertisements cannot drift.

### Session checks

The existing SessionStart route requests an active-interface check on startup and resume.
Shared instructions also state the obligation when the hook itself is unavailable.
The startup hook stays informational and preserves its existing aggregate deadline.
A timeout prints a diagnostic and leaves the required live evidence unknown.

Before dependent work, the agent runs a normal-permission shell probe through the actual chat tool.
It then invokes the repository's Bench wrapper through that same tool path.
The active interface also supplies file, skill, hook, and permission evidence needed by the operation.
A fresh subprocess, integrated terminal, or escalated shell is diagnostic evidence only.

The agent retains successful live evidence only for the current session and its observed context.
No reusable green certificate is written to disk.
Start, resume, runtime replacement, workspace change, or relevant policy change invalidates the affected observations.
When context cannot be compared, the agent repeats the required live checks before dependent work.
Different version numbers alone are not a failure.

A failed required check stops only dependent work.
Unrelated read-only discussion and operations with their own verified capabilities remain available.
The instructions name the blocked operation, failed check, and supported next action.
These are workflow obligations, not a new security sandbox against a malicious agent.

### Setup and reversible repair

One setup installs the shared repository assets through the existing adoption transaction.
The setup result distinguishes completed local installation from pending interface qualification.
Each interface may retain a separate user configuration home.
Bench reports conflicts and never synchronizes credentials, personal preferences, models, or chat history.

Compatibility repair uses canonical Bench payload provenance to select eligible destinations.
Only missing or unchanged Bench-managed assets qualify for automatic replacement.
When a check finds an eligible repair, the agent invokes the compatibility fix before repeating the affected checks.
Modified managed files and foreign files remain unchanged with a conflict action.
The kit source checkout uses its existing doctor repair route and never links consumer copies over source files.

Every repair preserves its preimages before its first destination mutation.
The adoption transaction remains the publication owner.
Its compatibility path retains restoration metadata after success, instead of discarding all preimages.
Undo restores only destinations whose current identity still matches that repair's postimage.
An originally absent destination returns to absence.

Repairs serialize on the canonical destination identities they share.
The writer census includes setup, link, upgrade, unlink, doctor repair, and compatibility undo.
Those paths must use the same exclusion protocol when they touch a protected destination.
Unrelated destinations can proceed concurrently.
A destination change after inspection refuses the stale plan before that destination's mutation.

Private repair records live below the resolved Bench home, in a compatibility-repair namespace keyed by repository identity.
Directories use mode 0700, and files use mode 0600.
Records contain managed preimages, modes, destination identities, and the metadata needed for undo.
They exclude credentials, raw environment values, chat content, and copied user configuration homes.
Records remain local until explicit removal after their recovery purpose ends.

A write, close, persistence, or restore failure reports incomplete repair and retains the available recovery state.
An interrupted repair can be inspected and recovered by a fresh process.
The implementation must not declare rollback success after ignoring a restore error.
A successful repair still requires a live retest through the failed interface.

Supported harness repair routes are eligible only when their authority, reversibility, and installed behavior are established.
No such automatic upstream repair is established by the current incident.
Process interruption, security-policy changes, trust approval, and private-runtime edits require the existing explicit reviewer decision.
An unavailable supported repair stays unresolved.

### Recovery without command execution

The shipped agreement includes a short failure instruction that does not require a successful Bench process.
The reference guide owns the full recovery procedure, and README links to it for the user.
The procedure identifies the failing interface and preserves the original startup error.
It permits supported external diagnostics without treating their success as recovered compatibility.
It directs the user to preserve active work before an authorized app restart or upstream escalation.
No automatic restart or private launcher repair enters this build.

### Executable authority

The current harness, operating-system permissions, and installed Bench launcher remain the trusted entry points.
Compatibility checks do not authenticate their own caller or replace the installed broker trust chain.
The wrapper launches the selected Bench binary through its current resolution contract.
Doctor delegates writes to the existing adoption owner, and lifecycle guidance delegates worktree publication to the existing broker.
This specification adds no claim of refusal before the harness can execute its first command.

## Implementation chunks

Each chunk contains one serial ticket and a review checkpoint.
C1 introduces the shared diagnostic owner and receives review before its consumers start.
C2 supplies safe repair to C3.
These checkpoints do not authorize separate release of incomplete compatibility.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| C1 / `1-diagnose-interface.md` | Diagnose the selected interface | CD01, CD02, CD03, CD04, CD05, CD06, CD07, CD08, CD09, CD10, CD11, CD12, CD13, CD14, CD15, CD16, CD17 | Named tests and live evidence in the owned coverage rows | no |
| C2 / `2-repair-managed-integration.md` | Repair managed integration reversibly | CD18, CD19, CD20, CD21, CD22, CD23, CD24, CD25, CD26, CD27, CD28, CD29, CD30, CD31, CD32, CD33, CD34, CD35, CD36 | Named tests and live evidence in the owned coverage rows | yes |
| C3 / `3-check-and-qualify-sessions.md` | Check and qualify independent sessions | CD37, CD38, CD39, CD40, CD41, CD42, CD43, CD44, CD45, CD46, CD47, CD48, CD49, CD50, CD51, CD52, CD53, CD54, CD55, CD56, CD57, CD58, CD59, CD60, CD61, CD62, CD63, CD64 | Named tests and live evidence in the owned coverage rows | yes |

## Testing decisions

Tests drive the production doctor entrypoint, the adoption transaction, and the real SessionStart script.
The compatibility owner receives typed facts through the same boundary used by the real collectors.
Tests do not substitute a package variable in one process and claim that a separate process observed it.
TDD applies to repair preservation and context invalidation.
The implementation records a targeted omission or swap red for every independent policy expectation.

Ordinary Go tests execute in the existing gate test phase.
System-tagged tests run through `bench test --check system`, which supplies `BENCH_KIT` and the sealed subject binary.
No external account, network request, or active desktop runtime becomes a hermetic gate dependency.
Live qualification remains review-owned and mandatory before the implementation can land.

Existing precedents and the enforcement-reader census are in [seam evidence](assets/seam-evidence.md).
The same-owner transaction tests expose persistence failures.
The real system-test process owner exposes wrapper, hook, and linked-repository integration.

### Doctor-route mutation

The probe owner excludes system checks from its automatic targets.
Its documented manual route preserves the same mutation and system assertion.
This procedure implements the `doctor-route-probe` requirement.
The author performs each step separately and retains its result.

1. Copy `cmd/bench/main.go` to a unique private backup path outside the repository.
2. Verify that the backup exists and matches the source bytes and mode.
3. Replace the unique `Run: adoptCommand("doctor")` with `Run: adoptCommand("setup")` in the source file.
4. Run `bench test --check system` and retain the terminal result.
5. Restore the preserved bytes and mode, even if the test fails unexpectedly.
6. Verify that the restored source matches the backup.
7. Run `bench test --check system` and require a green result.

The mutated run must compile and fail TestCompatibilityMissingPath for the doctor route.
Another failure does not prove this requirement.
The author records the mutation outcome, failing test, exit code, and verified restoration.
The preserved copy remains available until restoration succeeds.
This procedure changes no production requirement or system-test obligation.

### Seam diagram

```text
actual CLI or desktop chat
  | normal tool call                         failed tool startup
  v                                          |
repository wrapper                           v
  |                                  loaded instruction + user guide
  v                                          |
doctor --compat -> compatibility owner        supported external diagnosis
  |                  ^                       |
  |              typed fact tests            live retest in failed chat
  v
adoption transaction -> managed destination -> retained undo preimage
  ^                         ^
repair fault tests      concurrent writer tests

SessionStart -> session-inspect -> live-check obligation -> actual chat tools
      ^                ^                                   ^
real hook tests   deadline tests                    review-owned evidence
```

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| CD01 | 2 | The compatibility report names the explicitly selected CLI or desktop interface | `internal/systemtest/compatibility_test.go` (`TestCompatibilityCollectorInterfaces`) | Swapping the selected interface changes the expected context |
| CD02 | 2 | An unknown active runtime remains unknown when a launcher version is available | `internal/compatibility/inspect_test.go` (`TestCompatibilityRuntimeProvenance`) | A PATH launcher cannot supply the active server identity |
| CD03 | 3 | The report names the active configuration home without assuming the other interface shares it | `internal/compatibility/inspect_test.go` (`TestCompatibilityHomes`); `internal/adopt/compatibility_test.go` (`TestCompatibilityMissingConfigurationHome`) | Different Linux and mounted Windows homes must remain distinct |
| CD04 | 3 | An unreadable effective setting reports unknown with its source | `internal/compatibility/inspect_test.go` (`TestCompatibilityUnknownConfig`) | Dropping a failed source would fabricate agreement |
| CD05 | 4 | Diagnostics leave both user configuration homes byte-identical | `internal/systemtest/compatibility_test.go` (`TestCompatibilityCollectorReadOnly`) | Independent before-and-after snapshots catch silent synchronization |
| CD06 | 5 | An absent required repository asset reports its exact restoration action | `internal/adopt/compatibility_test.go` (`TestCompatibilityMissingAsset`) | Omitting one required asset must produce a failed check |
| CD07 | 6 | An absent global Bench command reports the existing by-path route | `internal/systemtest/compatibility_test.go` (`TestCompatibilityMissingPath`) | A fixture without global Bench must still reach the local wrapper |
| CD08 | 7 | A declared hook remains unverified until its live behavior is observed | `internal/compatibility/inspect_test.go` (`TestCompatibilityDeclaredHook`) | A configuration file alone cannot produce a live pass |
| CD09 | 8 | An invalid interface operand exits 2 before configuration reads | `internal/adopt/compatibility_test.go` (`TestCompatibilityGrammar`) | A read sentinel catches validation after inspection |
| CD10 | 8 | Special files and dangling configuration links return a bounded diagnostic | `internal/compatibility/inspect_test.go` (`TestCompatibilityFileKinds`) | FIFO and dangling-link fixtures cannot become empty successful inputs |
| CD11 | 8 | An empty required configuration file receives a distinct result from an absent file | `internal/compatibility/inspect_test.go` (`TestCompatibilityAbsentEmpty`) | Both fixture shapes must retain their different evidence states |
| CD12 | 8 | A path containing spaces and glob characters reaches the intended repository | `internal/systemtest/compatibility_test.go` (`TestCompatibilityPaths`) | Decoy paths expose splitting or glob expansion |
| CD13 | 29 | A changed relevant context produces a different compatibility fingerprint | `internal/compatibility/inspect_test.go` (`TestCompatibilityFingerprint`); `internal/systemtest/compatibility_test.go` (`TestCompatibilityCollectorPolicyFingerprint`) | One-field mutations invalidate the previous observation context |
| CD14 | 29 | A personal model preference alone does not create a Bench configuration conflict | `internal/compatibility/inspect_test.go` (`TestCompatibilityPersonalSettings`) | Two permitted preferences must not require synchronization |
| CD15 | 2 | The inspection result never labels the whole active chat qualified | `internal/adopt/compatibility_test.go` (`TestCompatibilityInspectionBoundary`) | Successful local checks still require live tool evidence |
| CD16 | 8 | Control-bearing report fields cannot inject a terminal control sequence | `internal/compatibility/inspect_test.go` (`TestCompatibilityOutput`) | The output seam must refuse unsafe cells through the existing renderer |
| CD17 | 26 | Legacy doctor invocations retain their prior exit and output behavior | `internal/adopt/compatibility_test.go` (`TestCompatibilityLegacyDoctor`) | A differential fixture compares legacy modes before and after integration |
| CD18 | 1 | One setup installs the shared Bench integration without copying either user configuration home | `internal/systemtest/compatibility_test.go` (`TestCompatibilitySetup`) | Fresh linked-repository fixtures keep unrelated home sentinels intact |
| CD19 | 9 | The compatibility repair restores an unmodified managed asset from the canonical payload | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityManagedRepair`) | Removing one managed hook must become a working installed hook |
| CD20 | 9 | A second identical repair leaves all managed destination bytes and modes unchanged | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityRepairIdempotence`) | Tracked repeated application exposes self-induced drift |
| CD21 | 10 | A modified managed asset remains unchanged and reports a conflict | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityModifiedConflict`) | A user edit must survive the automatic repair path |
| CD22 | 10 | A foreign file at a repair destination remains unchanged | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityForeignConflict`) | A marker-free sentinel cannot become Bench-owned |
| CD23 | 11 | A repair requiring security-policy change returns a decision action without applying it | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityPolicyBoundary`) | Permission and trust-record sentinels must remain untouched |
| CD24 | 11 | A repair requiring process interruption returns a decision action without sending a signal | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityInterruptionBoundary`) | A supervised live process must survive the repair attempt |
| CD25 | 11 | A private runtime path receives no automatic mutation | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityRuntimeBoundary`) | An absent launcher path must stay outside the repair plan |
| CD26 | 12 | A backup-publication failure leaves all repair destinations unchanged | `internal/adopt/transaction/transaction_test.go` (`TestCompatibilityBackupFailure`) | A denied record directory prevents the first destination write |
| CD27 | 12 | An interrupted repair retains enough state for a fresh process to recover each touched target | `internal/systemtest/compatibility_test.go` (`TestCompatibilityInterruptedRepair`) | A second process must recover after an interrupted multi-target promotion |
| CD28 | 12 | A failed terminal record publication reports incomplete repair and preserves recovery data | `internal/adopt/transaction/transaction_test.go` (`TestCompatibilityTerminalFailure`) | Successful destination writes cannot conceal a missing terminal record |
| CD29 | 13 | Concurrent Bench writers cannot interleave changes to the same repair destination | `internal/systemtest/compatibility_test.go` (`TestCompatibilityConcurrentWriters`) | A barrier holds one writer while the competing writer attempts publication |
| CD30 | 14 | Undo restores every recorded preimage byte and mode | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityUndo`) | Deleting the final backup entry must fail the complete restoration comparison |
| CD31 | 14 | Undo removes a target that the repair created from absence | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityUndoCreated`) | An absent preimage must not become an empty surviving file |
| CD32 | 14 | Undo refuses a destination changed after repair | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityUndoConflict`) | A later edit must survive the stale undo request |
| CD33 | 15 | Repair records contain only managed preimages and their necessary restoration metadata | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityRepairPrivacy`) | Credential and raw-environment sentinels must not occur in the record |
| CD34 | 8 | A hostile repair identifier cannot select a path outside its private record directory | `internal/adopt/transaction/transaction_test.go` (`TestCompatibilityRepairPaths`) | Traversal and symlink fixtures cannot redirect reads or writes |
| CD35 | 9 | A repair rechecks destination identity immediately before its publication | `internal/adopt/transaction/transaction_test.go` (`TestCompatibilityPlanDrift`); `internal/adopt/transaction/identity_test.go` (`TestIdentityAfterReplacementPreparation`) | A target replaced after inspection must prevent its planned write |
| CD36 | 12 | A failed restore reports each unresolved target and preserves its available backup | `internal/adopt/transaction/transaction_test.go` (`TestCompatibilityRestoreFailure`); `internal/adopt/transaction/identity_test.go` (`TestUndoReportsEveryPreflightFailure`) | A restore error cannot become a successful rollback |
| CD37 | 16 | SessionStart emits the compatibility check obligation before dependent work | `internal/sessioninspect/compatibility_test.go` (`TestCompatibilityStartup`) | Removing the lifecycle call loses the required next action |
| CD38 | 17 | A resumed session repeats the active-interface checks | `internal/systemtest/compatibility_session_test.go` (`TestCompatibilityResume`) | A prior success followed by a changed context cannot suppress the resume check |
| CD39 | 18 | A hook-process success cannot replace the actual chat-tool probe | review-owned: `assets/qualification.md`, actual-tool transcript | A hook-only success beside a failed tool must remain unqualified |
| CD40 | 18 | An escalated diagnostic success cannot replace normal-permission evidence | review-owned: `assets/qualification.md`, permission comparison | The observed incident supplies the cheap wrong-success control |
| CD41 | 19 | A missing optional desktop tool leaves an unrelated Bench operation available | `internal/compatibility/inspect_test.go` (`TestCompatibilityOptionalTool`) | A disabled preview must not block a shell-only diagnostic |
| CD42 | 19 | A missing required capability stops its dependent operation before a mutation | review-owned: `assets/qualification.md`, missing-required-capability transcript | A mutation sentinel must remain absent when its required check fails |
| CD43 | 20 | The shipped instructions expose a recovery action when no Bench command can start | review-owned: shipped README and active-session failure transcript | A command-only error route cannot satisfy the failed-shell scenario |
| CD44 | 21 | A repair result requires a retest through the previously failed interface | `internal/sessioninspect/compatibility_test.go` (`TestCompatibilityRetestRequired`) | Printing repair success must not print recovered compatibility |
| CD45 | 21 | The qualification record includes successful recovery through normal desktop tools | review-owned: `assets/qualification.md`, recovery transcript | A still-failing desktop leaves the complete outcome unaccepted |
| CD46 | 22 | Independent CLI and desktop chats obtain the same repository Bench rules | review-owned: `assets/qualification.md`, two fresh chats | Shared conversation context cannot supply the missing integration |
| CD47 | 23 | Concurrent CLI and desktop writers use distinct Bench assignments | review-owned: `assets/qualification.md`, concurrent worktree exercise | Both chats must produce separate owned paths before their writes |
| CD48 | 23 | Both interfaces complete the existing Bench review and landing route | review-owned: `assets/qualification.md`, disposable repository lifecycle | Raw Git publication cannot satisfy the recorded lifecycle |
| CD49 | 24 | The actual CLI interface rejects the harmless forbidden-command fixture | review-owned: `assets/qualification.md`, CLI guard transcript | The rejected command sentinel must never execute |
| CD50 | 24 | The actual desktop interface rejects the harmless forbidden-command fixture | review-owned: `assets/qualification.md`, desktop guard transcript | The rejected command sentinel must never execute |
| CD51 | 25 | A verified equivalent review route produces the same review outcome | review-owned: `assets/qualification.md`, review comparison | A declared equivalent without an observed outcome is insufficient |
| CD52 | 26 | Existing macOS and Linux contract fixtures retain their prior outcomes | `internal/systemtest/compatibility_session_test.go` (`TestCompatibilityOtherHosts`) | Platform-specific branches must not change established non-WSL behavior |
| CD53 | 27 | Native Windows Bench execution remains outside the new qualification claim | `internal/compatibility/inspect_test.go` (`TestCompatibilitySupportBoundary`) | A Windows-native fixture must not report WSL qualification |
| CD54 | 28 | The complete outcome remains unaccepted while any required live qualification row is missing | review-owned: full acceptance reconciliation before landing | A green package suite cannot erase missing actual-interface evidence |
| CD55 | 29 | A runtime version difference alone does not fail a capability check | `internal/compatibility/inspect_test.go` (`TestCompatibilityVersionDifference`) | Different version labels with identical observed capabilities must not conflict |
| CD56 | 29 | Unknown context prevents reuse of an earlier live check | `internal/sessioninspect/compatibility_test.go` (`TestCompatibilityUnknownContext`) | A missing runtime or policy observation cannot inherit a prior pass |
| CD57 | 30 | Each actual interface writes and reads the expected bytes in its own scratch file | review-owned: `assets/qualification.md`, CLI and desktop file probes | A subprocess fixture alone cannot qualify the actual file tools |
| CD58 | 30 | Each actual interface invokes the repository skill from its installed path | review-owned: `assets/qualification.md`, skill invocation transcripts | An available file without a usable harness skill does not qualify |
| CD59 | 24 | A required permission conflict remains visible without broadening the configured policy | `internal/systemtest/compatibility_session_test.go` (`TestCompatibilityPermissionConflict`) | A policy sentinel exposes an automatic bypass |
| CD60 | 20 | An unknown upstream repair remains explicitly unresolved in the recovery guidance | `internal/sessioninspect/compatibility_test.go` (`TestCompatibilityUnknownRepair`) | A fabricated repair action cannot turn unknown into recovered |
| CD61 | 16 | A compatibility inspection timeout preserves the existing informational startup exit | `internal/sessioninspect/compatibility_test.go` (`TestCompatibilityStartupTimeout`) | Slow discovery cannot prevent the chat from opening |
| CD62 | 17 | A relevant configuration change invalidates the current session observation | review-owned: `assets/qualification.md`, configuration drift transcript | A changed hook or permission context must trigger another affected check |
| CD63 | 1 | A linked repository receives the compatibility instructions through the real payload installer | `internal/systemtest/compatibility_session_test.go` (`TestCompatibilityPayload`) | Kit-only documents cannot satisfy a linked consumer |
| CD64 | 19 | An equivalent route counts only after that operation succeeds through the route | `internal/compatibility/inspect_test.go` (`TestCompatibilityEquivalentEvidence`, `TestCompatibilityFailureAfterSuccess`) | Naming an alternate tool cannot satisfy the required capability |

### Edge inventory

The [edge dispositions](assets/seam-evidence.md#hostile-input-dispositions) walk the complete project checklist.
Every behavior serves both the kit checkout and linked repositories unless its row states a live-host boundary.
The repair owner distinguishes absent, empty, malformed, unreadable, and wrong-type records.
A managed symlink uses the existing adoption provenance rules; an unknown repair-record symlink is refused.
A valid symlinked repository root resolves through the existing Git identity owner.

Won't handle: conversation migration or history synchronization — independent chats continue through repository state.
Won't handle: native Windows Bench execution — the qualified Windows app runs its agent in WSL2.
Won't handle: all desktop operating systems — existing macOS and Linux behavior remains preserved.
Won't handle: every optional desktop tool — operations use verified required capabilities or verified equivalent routes.
Won't handle: automatic private-runtime edits — affected operations remain blocked with the reviewer-owned repair route.
Won't handle: authentication against a malicious same-user process — existing trust and permission boundaries remain authoritative.

## Ownership fences

- `.bench/BENCH-reference.md`
- `.bench/BENCH.md`
- `.bench/hooks/session-start.sh`
- `CHANGELOG.md`
- `CONTEXT.md`
- `DATA_HANDLING.md`
- `README.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/main.go`
- `internal/adopt`
- `internal/adopt/compatibility.go`
- `internal/adopt/compatibility_test.go`
- `internal/adopt/doctor.go`
- `internal/anchors`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_decision_maps.go`
- `internal/anchors/registry_decision_maps_test.go`
- `internal/compatibility`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/entry_point_parity_bite_test.go`
- `internal/conformance/entry_point_parity_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/preflight/evidencecmd/evidence_file_reconstruction_test.go`
- `internal/runbinary/runbinary_test.go`
- `internal/sessioninspect`
- `internal/systemtest/compatibility_test.go`
- `internal/systemtest/compatibility_session_test.go`
- `tests/canary/data-handling-derivation/undocumented-passlist-var`
- `tests/canary/docs-currency-token-diet/benchref-imported`
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped`
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated`
- `tests/canary/docs-currency-token-diet/dogfood-referent-shipped`
- `tests/canary/docs-currency-token-diet/missing-cli-inventory`
- `tests/canary/docs-currency-token-diet/readme-command-first`
- `tests/canary/docs-currency-token-diet/signal-vocabulary-drift`
- `tests/canary/docs-currency-token-diet/stale-cli-doc-reference`
- `tests/canary/load-validity-metadata/readme-shared-rule-drift`
- `tests/canary/load-validity-metadata/shared-rule-drift`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key`
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy`
- `tests/canary/skills-index-command-adapters/dangling-index`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/skills-index-command-adapters/missing-index-field`
- `tests/canary/skills-index-command-adapters/stale-index-wording`
- `tests/canary/skills-index-command-adapters/unindexed-skill`
- `tests/canary/workflow-guidance-anchors/agents-handoff-section-rule`
- `tests/canary/workflow-guidance-anchors/capture-sink-anchor`
- `tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns`
- `tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-coverage-map-term`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-parts`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-decision-map-term`
- `tests/canary/workflow-guidance-anchors/context-reader-sweep-term`
- `tests/canary/workflow-guidance-anchors/context-ticket-vocabulary`
- `tests/canary/workflow-guidance-anchors/core-absent-bench-operational-layer`
- `tests/canary/workflow-guidance-anchors/core-absent-category-context`
- `tests/canary/workflow-guidance-anchors/core-absent-category-oracle`
- `tests/canary/workflow-guidance-anchors/core-absent-category-setup`
- `tests/canary/workflow-guidance-anchors/core-absent-category-work`
- `tests/canary/workflow-guidance-anchors/core-absent-gate-authority`
- `tests/canary/workflow-guidance-anchors/core-absent-kit-only-ship`
- `tests/canary/workflow-guidance-anchors/core-absent-no-path-fallback`
- `tests/canary/workflow-guidance-anchors/core-absent-retained-integration-source`
- `tests/canary/workflow-guidance-anchors/core-absent-retro-capture-owner`
- `tests/canary/workflow-guidance-anchors/core-absent-retro-drain-owner`
- `tests/canary/workflow-guidance-anchors/core-absent-skills-guidance`
- `tests/canary/workflow-guidance-anchors/core-absent-upgrade-route`
- `tests/canary/workflow-guidance-anchors/core-cli-bounded-complete-response`
- `tests/canary/workflow-guidance-anchors/core-cli-call-pair`
- `tests/canary/workflow-guidance-anchors/core-progressive-loading-term`
- `tests/canary/workflow-guidance-anchors/delegated-chunk-tip-review`
- `tests/canary/workflow-guidance-anchors/delegated-green-predecessor`
- `tests/canary/workflow-guidance-anchors/delegated-opt-in-entry`
- `tests/canary/workflow-guidance-anchors/delegated-pending-predecessor-stop`
- `tests/canary/workflow-guidance-anchors/delegated-prerequisite-checkpoint`
- `tests/canary/workflow-guidance-anchors/delegated-repair-ownership`
- `tests/canary/workflow-guidance-anchors/fix-dont-park-section-relocated`
- `tests/canary/workflow-guidance-anchors/light-path-inline`
- `tests/canary/workflow-guidance-anchors/readme-shaping-skip`
- `tests/canary/workflow-guidance-anchors/reference-agent-push-rule`
- `tests/canary/workflow-guidance-anchors/reference-bench-operational-layer`
- `tests/canary/workflow-guidance-anchors/reference-category-context`
- `tests/canary/workflow-guidance-anchors/reference-category-oracle`
- `tests/canary/workflow-guidance-anchors/reference-category-setup`
- `tests/canary/workflow-guidance-anchors/reference-category-work`
- `tests/canary/workflow-guidance-anchors/reference-gate-authority`
- `tests/canary/workflow-guidance-anchors/reference-kit-only-ship`
- `tests/canary/workflow-guidance-anchors/reference-no-path-fallback`
- `tests/canary/workflow-guidance-anchors/reference-progressive-loading-term`
- `tests/canary/workflow-guidance-anchors/reference-refusal-route-shape`
- `tests/canary/workflow-guidance-anchors/reference-retro-capture-owner`
- `tests/canary/workflow-guidance-anchors/reference-retro-drain-owner`
- `tests/canary/workflow-guidance-anchors/reference-skills-guidance`
- `tests/canary/workflow-guidance-anchors/reference-upgrade-route`
- `tests/canary/workflow-guidance-anchors/structured-phase-progress-anchor`
- `tests/canary/workflow-guidance-anchors/ticket-light-path-anchor`


Destination composition scope:

The user directed completion and landing on 2026-10-03.
C3 includes the changes inherited from main at b4a2fe693be05b997c3f30ae3ec5c86a4f0e01cf.
The following exact paths extend the fence and ticket writes for this composition.
Their committed contents and modes equal that destination commit.
The continuation preserves these paths and reviews their integration with C3.
Any new edit to these paths needs its own behavioral justification.

- `.agents/commands/bench-debug.md`
- `.agents/commands/bench-final-check.md`
- `.agents/commands/bench-implement-spec.md`
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md`
- `.agents/skills/bench-craft-line/SKILL.md`
- `.agents/skills/bench-craft-line/references/bounded-repair-policy.md`
- `.agents/skills/bench-implement-spec/references/staleness-pass.md`
- `AGENTS.md`
- `ASSESSMENT.md`
- `ROADMAP.md`
- `decisions/cost-follows-project-size.md`
- `decisions/cost-follows-project-size/tickets/1.md`
- `decisions/cost-follows-project-size/tickets/10.md`
- `decisions/cost-follows-project-size/tickets/2.md`
- `decisions/cost-follows-project-size/tickets/3.md`
- `decisions/cost-follows-project-size/tickets/4.md`
- `decisions/cost-follows-project-size/tickets/5.md`
- `decisions/cost-follows-project-size/tickets/6.md`
- `decisions/cost-follows-project-size/tickets/7.md`
- `decisions/cost-follows-project-size/tickets/8.md`
- `decisions/cost-follows-project-size/tickets/9.md`
- `internal/bounds/bounds.go`
- `internal/canary/inventory.go`
- `internal/census/census_test.go`
- `internal/conformance/implementation_continuation_test.go`
- `internal/gate/gate_go.go`
- `internal/gate/kit_source.go`
- `internal/gate/kit_value_test.go`
- `internal/gate/lane_select.go`
- `internal/gate/lane_test.go`
- `internal/gate/phases.go`
- `internal/gate/record_time_skew_test.go`
- `internal/gate/tag_census_test.go`
- `internal/gate/verdict_registry.go`
- `internal/preflight/decision.go`
- `internal/preflight/evidencecmd/operations.go`
- `internal/preflight/proposal.go`
- `internal/prose/walk.go`
- `internal/prose/walk_test.go`
- `internal/racetests/racetests.go`
- `internal/releaseevidence/release_evidence.go`
- `internal/releaseevidence/release_requirements.go`
- `internal/releaseevidence/requirement_inspection.go`
- `internal/releasepreflight/command.go`
- `internal/releasepreflight/decision.go`
- `internal/releasepreflight/types.go`
- `internal/reviewrecord/check.go`
- `internal/reviewrecord/completion.go`
- `internal/reviewrecord/coverage.go`
- `internal/reviewrecord/delegated.go`
- `internal/reviewrecord/parse.go`
- `internal/reviewrecord/plan.go`
- `internal/reviewrecord/record.go`
- `internal/reviewrecord/recordcmd/amendment.go`
- `internal/reviewrecord/recordcmd/command.go`
- `internal/reviewrecord/recordcmd/command_test.go`
- `internal/reviewrecord/recordcmd/completion.go`
- `internal/reviewrecord/recordcmd/completion_test.go`
- `internal/reviewrecord/recordtest/fixture.go`
- `internal/reviewrecord/source.go`
- `internal/reviewrecord/write.go`
- `internal/tickets/registry_data.go`
- `internal/tickets/registry_data_test.go`
- `internal/worktree/build.go`
- `internal/worktree/build_test.go`
- `internal/worktree/classifier.go`
- `internal/worktree/clean.go`
- `internal/worktree/clean_discard.go`
- `internal/worktree/clean_discard_test.go`
- `internal/worktree/clean_discard_transaction_test.go`
- `internal/worktree/clean_landed.go`
- `internal/worktree/clean_landed_apply_test.go`
- `internal/worktree/clean_landed_hostile_test.go`
- `internal/worktree/clean_landed_test.go`
- `internal/worktree/clean_set.go`
- `internal/worktree/clean_set_apply.go`
- `internal/worktree/clean_set_apply_test.go`
- `internal/worktree/clean_set_outcomes_test.go`
- `internal/worktree/clean_set_refusal_test.go`
- `internal/worktree/clean_set_test.go`
- `internal/worktree/clean_unclaimed_test.go`
- `internal/worktree/delegated_integration_test.go`
- `internal/worktree/effect_census_test.go`
- `internal/worktree/effects.go`
- `internal/worktree/identity_component_test.go`
- `internal/worktree/joins.go`
- `internal/worktree/journey_race_test.go`
- `internal/worktree/journey_test.go`
- `internal/worktree/land.go`
- `internal/worktree/land_broker_notice_test.go`
- `internal/worktree/land_effects.go`
- `internal/worktree/land_effects_cleanup_test.go`
- `internal/worktree/land_effects_test.go`
- `internal/worktree/land_empty_sibling_test.go`
- `internal/worktree/land_facts_test.go`
- `internal/worktree/land_fixtures_test.go`
- `internal/worktree/land_flags_test.go`
- `internal/worktree/land_freshness_test.go`
- `internal/worktree/land_identity.go`
- `internal/worktree/land_identity_test.go`
- `internal/worktree/land_local_capture_test.go`
- `internal/worktree/land_marker_fixture_test.go`
- `internal/worktree/land_prunes_landed_siblings_test.go`
- `internal/worktree/land_reauthorization_test.go`
- `internal/worktree/land_release_refusal_test.go`
- `internal/worktree/land_resume.go`
- `internal/worktree/land_resume_refusal_test.go`
- `internal/worktree/land_resume_test.go`
- `internal/worktree/land_specless_test.go`
- `internal/worktree/land_surface_test.go`
- `internal/worktree/land_tickets_only_test.go`
- `internal/worktree/lifecycle.go`
- `internal/worktree/lifecycle_acquire_test.go`
- `internal/worktree/list.go`
- `internal/worktree/list_actions_test.go`
- `internal/worktree/list_selected.go`
- `internal/worktree/list_selected_test.go`
- `internal/worktree/live_binary.go`
- `internal/worktree/live_binary_test.go`
- `internal/worktree/merge.go`
- `internal/worktree/merge_caller_root_test.go`
- `internal/worktree/merge_test.go`
- `internal/worktree/orphan_render_test.go`
- `internal/worktree/orphan_sweep_test.go`
- `internal/worktree/ownership.go`
- `internal/worktree/ownership_test.go`
- `internal/worktree/parallel_census_test.go`
- `internal/worktree/path_identifier_test.go`
- `internal/worktree/pool_reclaim.go`
- `internal/worktree/pool_reclaim_test.go`
- `internal/worktree/reauthorize.go`
- `internal/worktree/reauthorize_test.go`
- `internal/worktree/recovery_retry_test.go`
- `internal/worktree/reset.go`
- `internal/worktree/reset_apply.go`
- `internal/worktree/reset_apply_test.go`
- `internal/worktree/reset_plan_test.go`
- `internal/worktree/reset_refusal_test.go`
- `internal/worktree/reset_repair_test.go`
- `internal/worktree/reset_restore_refusal_test.go`
- `internal/worktree/resume.go`
- `internal/worktree/resume_clean_ambient_test.go`
- `internal/worktree/resume_test.go`
- `internal/worktree/serial_ceiling_test.go`
- `internal/worktree/single_read_census_cases_test.go`
- `internal/worktree/single_read_census_test.go`
- `internal/worktree/snapshot.go`
- `internal/worktree/snapshot_test.go`
- `internal/worktree/subshell.go`
- `internal/worktree/subshell_test.go`
- `internal/worktree/table_name_census_test.go`
- `internal/worktree/verb_fixture_test.go`
- `internal/worktree/verb_result_route_test.go`
- `internal/worktree/verb_runner_check_test.go`
- `internal/worktree/verb_runner_test.go`
- `internal/worktree/worktree.go`
- `internal/worktree/worktree_test.go`
- `projects/benchkit.md`
- `roadmap/FT101.md`
- `roadmap/FT284.md`
- `roadmap/FT293.md`
- `roadmap/FT318.md`
- `roadmap/FT327.md`
- `roadmap/FT352.md`
- `roadmap/FT356.md`
- `roadmap/FT361.md`
- `roadmap/FT367.md`
- `roadmap/FT368.md`
- `roadmap/FT369.md`
- `roadmap/FT373.md`
- `roadmap/FT374.md`
- `roadmap/FT375.md`
- `roadmap/FT376.md`
- `roadmap/FT377.md`
- `specs/ft370-comment-only-evidence/spec.md`
- `specs/ft370-comment-only-evidence/tickets/1-move-tree-change-reader.md`
- `specs/ft370-comment-only-evidence/tickets/2-prove-comment-only-gaps.md`
- `specs/ft370-comment-only-evidence/tickets/3-accept-proven-gaps-at-checkpoint.md`
- `specs/ft370-comment-only-evidence/tickets/4-state-comment-only-correction-rule.md`
- `specs/markdown-block-reader/spec.md`
- `specs/markdown-block-reader/tickets/1-add-block-reader.md`
- `specs/markdown-block-reader/tickets/10-read-agents-marker-blocks.md`
- `specs/markdown-block-reader/tickets/11-forbid-block-rule-copies.md`
- `specs/markdown-block-reader/tickets/2-read-spec-and-coverage-blocks.md`
- `specs/markdown-block-reader/tickets/3-read-roadmap-blocks.md`
- `specs/markdown-block-reader/tickets/4-read-field-scan-blocks.md`
- `specs/markdown-block-reader/tickets/5-read-handoff-blocks.md`
- `specs/markdown-block-reader/tickets/6-read-journal-blocks.md`
- `specs/markdown-block-reader/tickets/7-read-review-record-blocks.md`
- `specs/markdown-block-reader/tickets/8-read-anchor-blocks.md`
- `specs/markdown-block-reader/tickets/9-read-skills-frontmatter.md`
- `specs/worktree-seam-reduction/decisions/worktree-seams.md`
- `specs/worktree-seam-reduction/decisions/worktree-seams/assets/seam-classification.md`
- `specs/worktree-seam-reduction/decisions/worktree-seams/tickets/1.md`
- `specs/worktree-seam-reduction/decisions/worktree-seams/tickets/10.md`
- `specs/worktree-seam-reduction/decisions/worktree-seams/tickets/11.md`
- `specs/worktree-seam-reduction/decisions/worktree-seams/tickets/2.md`
- `specs/worktree-seam-reduction/decisions/worktree-seams/tickets/3.md`
- `specs/worktree-seam-reduction/decisions/worktree-seams/tickets/4.md`
- `specs/worktree-seam-reduction/decisions/worktree-seams/tickets/5.md`
- `specs/worktree-seam-reduction/decisions/worktree-seams/tickets/6.md`
- `specs/worktree-seam-reduction/decisions/worktree-seams/tickets/7.md`
- `specs/worktree-seam-reduction/decisions/worktree-seams/tickets/8.md`
- `specs/worktree-seam-reduction/decisions/worktree-seams/tickets/9.md`
- `specs/worktree-seam-reduction/spec.md`
- `specs/worktree-seam-reduction/tickets/1-read-the-kit-value-once-in-gate.md`
- `specs/worktree-seam-reduction/tickets/10-fault-the-merge-reconcile-with-a-stale-index-lock.md`
- `specs/worktree-seam-reduction/tickets/11-lift-each-read-below-a-census-entry.md`
- `specs/worktree-seam-reduction/tickets/12-refuse-a-read-below-a-census-entry.md`
- `specs/worktree-seam-reduction/tickets/13-make-the-serial-ceiling-exact.md`
- `specs/worktree-seam-reduction/tickets/14-name-each-table-in-production.md`
- `specs/worktree-seam-reduction/tickets/15-require-the-joins-route-of-a-not-called-stub.md`
- `specs/worktree-seam-reduction/tickets/2-carry-an-ambient-value-below-each-verb-entry.md`
- `specs/worktree-seam-reduction/tickets/3-pass-the-kit-value-to-merge-and-land.md`
- `specs/worktree-seam-reduction/tickets/4-interrupt-the-landing-marker-with-a-gate-script.md`
- `specs/worktree-seam-reduction/tickets/5-fault-the-landing-follow-on-steps-with-real-fixtures.md`
- `specs/worktree-seam-reduction/tickets/6-land-the-stubbed-landing-tests-for-real.md`
- `specs/worktree-seam-reduction/tickets/7-fault-the-cleanup-reads-with-real-fixtures.md`
- `specs/worktree-seam-reduction/tickets/8-fault-the-reauthorize-unlock-with-a-denied-admin-directory.md`
- `specs/worktree-seam-reduction/tickets/9-fault-the-reset-move-with-real-fixtures.md`

The inherited guidance paths require the following exact fixture dependencies.
Their source bytes remain unchanged.

- `tests/canary/docs-currency-token-diet/introduces-undeclared-command`
- `tests/canary/docs-currency-token-diet/stale-command-reference`
- `tests/canary/guidance-prose-budgets/over-budget-skill`
- `tests/canary/line-routing/line-binding-prose-drift`
- `tests/canary/package-core-guard/bounds-classify-limit-restated`
- `tests/canary/package-core-guard/bounds-discovery-window-unwrapped`
- `tests/canary/package-core-guard/bounds-dot-import-package-alias`
- `tests/canary/package-core-guard/bounds-dot-import-wait`
- `tests/canary/package-core-guard/bounds-duplicate-owner`
- `tests/canary/package-core-guard/bounds-intent-window-fixed`
- `tests/canary/package-core-guard/bounds-multiple-dot-import-wait`
- `tests/canary/package-core-guard/bounds-parenthesized-wait`
- `tests/canary/package-core-guard/bounds-raw-elapsed-wait`
- `tests/canary/package-core-guard/bounds-raw-injected-wait`
- `tests/canary/package-core-guard/bounds-raw-wait-deadline`
- `tests/canary/package-core-guard/bounds-raw-wait-duration`
- `tests/canary/package-core-guard/bounds-read-limit-restated`
- `tests/canary/package-core-guard/bounds-reassigned-wait-duration`
- `tests/canary/package-core-guard/bounds-redeclared-wait-duration`
- `tests/canary/package-core-guard/bounds-worktree-window-unwrapped`
- `tests/canary/roadmap-detail-integrity/roadmap-degraded-row-directory`
- `tests/canary/roadmap-detail-integrity/roadmap-duplicate-row`
- `tests/canary/roadmap-detail-integrity/roadmap-heading-mismatch`
- `tests/canary/roadmap-detail-integrity/roadmap-inline-body`
- `tests/canary/roadmap-detail-integrity/roadmap-missing-detail-owner`
- `tests/canary/roadmap-detail-integrity/roadmap-next-duplicate-line`
- `tests/canary/roadmap-detail-integrity/roadmap-next-missing-line`
- `tests/canary/roadmap-detail-integrity/roadmap-next-unanchored-line`
- `tests/canary/roadmap-detail-integrity/roadmap-next-unknown-token`
- `tests/canary/roadmap-detail-integrity/roadmap-orphan-detail`
- `tests/canary/roadmap-detail-integrity/roadmap-unreadable-detail`
- `tests/canary/roadmap-detail-integrity/roadmap-unrecognized-file`
- `tests/canary/roadmap-detail-integrity/roadmap-wrapped-heading`
- `tests/canary/skill-description-budgets/budget-table-missing`
- `tests/canary/skill-description-budgets/description-folded`
- `tests/canary/skill-description-budgets/description-missing`
- `tests/canary/skill-description-budgets/over-budget-command`
- `tests/canary/skill-description-budgets/over-budget-description`
- `tests/canary/workflow-guidance-anchors/agents-system-suite-route`
- `tests/canary/workflow-guidance-anchors/benchkit-hostile-input-heading`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-owner`
- `tests/canary/workflow-guidance-anchors/benchkit-review-round-routing`
- `tests/canary/workflow-guidance-anchors/benchkit-spec-ownership`
- `tests/canary/workflow-guidance-anchors/benchkit-system-suite-route`
- `tests/canary/workflow-guidance-anchors/calibration-abstained-row`
- `tests/canary/workflow-guidance-anchors/calibration-abstention-final`
- `tests/canary/workflow-guidance-anchors/calibration-claim-schema`
- `tests/canary/workflow-guidance-anchors/calibration-coordinator-label`
- `tests/canary/workflow-guidance-anchors/calibration-declaration-line`
- `tests/canary/workflow-guidance-anchors/calibration-no-free-text`
- `tests/canary/workflow-guidance-anchors/calibration-retro-aggregate-duty`
- `tests/canary/workflow-guidance-anchors/calibration-retro-table-duty`
- `tests/canary/workflow-guidance-anchors/calibration-status-definitions`
- `tests/canary/workflow-guidance-anchors/debug-archaeology-anchor`
- `tests/canary/workflow-guidance-anchors/debug-phase1-stop-gate-softened`
- `tests/canary/workflow-guidance-anchors/debug-red-commit`
- `tests/canary/workflow-guidance-anchors/debug-reproduction-economics-deleted`
- `tests/canary/workflow-guidance-anchors/delegate-anchor-probe-owning-check`
- `tests/canary/workflow-guidance-anchors/delegate-grammar-fence-inventory`
- `tests/canary/workflow-guidance-anchors/delegate-live-tree-inventory-fence`
- `tests/canary/workflow-guidance-anchors/delegate-out-of-fence-write`
- `tests/canary/workflow-guidance-anchors/delegate-probe-mutated-bytes`
- `tests/canary/workflow-guidance-anchors/delegate-root-conformance-pass`
- `tests/canary/workflow-guidance-anchors/delegate-serial-ceiling-fence`
- `tests/canary/workflow-guidance-anchors/delegate-skip-ownership-check`
- `tests/canary/workflow-guidance-anchors/delegated-account-inventory`
- `tests/canary/workflow-guidance-anchors/delegated-account-reconciliation`
- `tests/canary/workflow-guidance-anchors/delegated-author-limit`
- `tests/canary/workflow-guidance-anchors/delegated-cap-exhaustion-trigger`
- `tests/canary/workflow-guidance-anchors/delegated-dispatch-declaration`
- `tests/canary/workflow-guidance-anchors/delegated-entry-refusals`
- `tests/canary/workflow-guidance-anchors/delegated-final-verification`
- `tests/canary/workflow-guidance-anchors/delegated-lost-session-trigger`
- `tests/canary/workflow-guidance-anchors/delegated-mid-review-route`
- `tests/canary/workflow-guidance-anchors/delegated-no-progress-trigger`
- `tests/canary/workflow-guidance-anchors/delegated-paid-comparison`
- `tests/canary/workflow-guidance-anchors/delegated-resumption-contents`
- `tests/canary/workflow-guidance-anchors/delegated-terminal-failure-trigger`
- `tests/canary/workflow-guidance-anchors/delegated-tier-authorization`
- `tests/canary/workflow-guidance-anchors/delegated-unbound-model-stop`
- `tests/canary/workflow-guidance-anchors/delegated-writer-termination`
- `tests/canary/workflow-guidance-anchors/dg-1`
- `tests/canary/workflow-guidance-anchors/dg-1-contradiction`
- `tests/canary/workflow-guidance-anchors/dg-2`
- `tests/canary/workflow-guidance-anchors/dg-2-contradiction`
- `tests/canary/workflow-guidance-anchors/dg-25`
- `tests/canary/workflow-guidance-anchors/dg-26`
- `tests/canary/workflow-guidance-anchors/dg-29`
- `tests/canary/workflow-guidance-anchors/dg-29-verification-target`
- `tests/canary/workflow-guidance-anchors/dg-3`
- `tests/canary/workflow-guidance-anchors/dg-3-blanket-ban`
- `tests/canary/workflow-guidance-anchors/dg-3-contradiction`
- `tests/canary/workflow-guidance-anchors/dg-30`
- `tests/canary/workflow-guidance-anchors/dg-31`
- `tests/canary/workflow-guidance-anchors/dg-31-contradiction-trigger`
- `tests/canary/workflow-guidance-anchors/dg-4`
- `tests/canary/workflow-guidance-anchors/dg-4-contradiction`
- `tests/canary/workflow-guidance-anchors/dg-5`
- `tests/canary/workflow-guidance-anchors/dg-5-contradiction`
- `tests/canary/workflow-guidance-anchors/dg-5-dirty`
- `tests/canary/workflow-guidance-anchors/dg-5-dirty-contradiction`
- `tests/canary/workflow-guidance-anchors/dg-6`
- `tests/canary/workflow-guidance-anchors/dg-6-command`
- `tests/canary/workflow-guidance-anchors/dg-6-digest`
- `tests/canary/workflow-guidance-anchors/dg-6-dirty`
- `tests/canary/workflow-guidance-anchors/dg-6-surface`
- `tests/canary/workflow-guidance-anchors/final-check-bare-leftover-clean-retired`
- `tests/canary/workflow-guidance-anchors/final-check-census-read-before-land`
- `tests/canary/workflow-guidance-anchors/final-check-landed-worktree-sweep`
- `tests/canary/workflow-guidance-anchors/final-check-light-path-changelog-heading`
- `tests/canary/workflow-guidance-anchors/final-check-scratch-branch-clean`
- `tests/canary/workflow-guidance-anchors/implement-spec-adoption-freshness`
- `tests/canary/workflow-guidance-anchors/implement-spec-coverage-task-seeding`
- `tests/canary/workflow-guidance-anchors/implement-spec-cross-harness-pointer`
- `tests/canary/workflow-guidance-anchors/implement-spec-entry-validation`
- `tests/canary/workflow-guidance-anchors/implement-spec-incapable-harness`
- `tests/canary/workflow-guidance-anchors/implement-spec-inline-exception`
- `tests/canary/workflow-guidance-anchors/implement-spec-mandatory-delegation-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-prose-owner-transfer`
- `tests/canary/workflow-guidance-anchors/implement-spec-read-only-helper`
- `tests/canary/workflow-guidance-anchors/implement-spec-red-preflight-route`
- `tests/canary/workflow-guidance-anchors/implement-spec-review-charge-omitted`
- `tests/canary/workflow-guidance-anchors/implement-spec-review-charge-order`
- `tests/canary/workflow-guidance-anchors/implement-spec-review-charge-reversed`
- `tests/canary/workflow-guidance-anchors/implement-spec-status-flip-anchor`
- `tests/canary/workflow-guidance-anchors/implement-spec-worktree-before-preflight`
- `tests/canary/workflow-guidance-anchors/implement-spec-write-delegation`
- `tests/canary/workflow-guidance-anchors/line-anchor-missing`
- `tests/canary/workflow-guidance-anchors/prepared-build-approval`
- `tests/canary/workflow-guidance-anchors/prepared-build-freshness`
- `tests/canary/workflow-guidance-anchors/prepared-triage-authority`
- `tests/canary/workflow-guidance-anchors/prepared-triage-cli-boundary`
- `tests/canary/workflow-guidance-anchors/prepared-triage-input-contract`
- `tests/canary/workflow-guidance-anchors/ticket-stage-routing-anchor`

Reviewer disposition: approved with the ticket graph on 2026-10-02.
The fence is the union of ticket writes outside the spec and capture folders.
The review pickup is `reviews/cli-desktop-consistency.md`.
The implementation must request an in-scope plan expansion before editing any missing closure path.
No grant to the reviewer-owned structure budget is included.

## Out of scope

The edge inventory records the approved exclusions.
No new capability is deferred from the approved outcome.
These exclusions require 0 edits and 0 gate runs in this build.
A future expansion requires its own reviewed scope and estimate.

## Further notes

### C3 test placement

C3 uses a separate system test file for session behavior.
The existing fixture remains the shared owner.
The acceptance rows, chunk IDs, dependencies, and system check remain unchanged.

### Qualification and completion

C3 writes `assets/qualification.md` with a transcript reference for each live row.
Each entry names its interface, active runtime evidence, configuration sources, repository identity, permission mode, invocation, and observed result.
An unknown value stays unknown and cannot authorize evidence reuse.
The report distinguishes direct observation from reviewer-supplied transcripts.
It includes a capability inventory with each desktop tool's required, equivalent, or optional disposition.

The qualification runs independent fresh chats and resumed chats in the actual installed interfaces.
It includes harmless hook rejection, ordinary command success, file access, skill invocation, concurrent writers, review, landing, and failed-capability recovery.
The reviewer decides any disruptive recovery action before that action runs.
A fixture or a fresh sandbox process cannot substitute for these live rows.
The present desktop failure must clear through a supported, authorized route before complete qualification succeeds.

### Source clauses and additions

[Seam evidence](assets/seam-evidence.md) maps every approved source clause to its acceptance rows.
The same artifact records source validity, reader sweeps, execution roots, and the fixed pre-review proof checklist.

Flagged additions: the diagnostic mode, context fingerprint, repair journal, and guarded undo implement the approved checks and reversibility.
They add no new supported environment or repair authority.
The shared writer exclusion is necessary to preserve the approved concurrent-use contract.
The reviewer approved these engineering choices with the seams and fences.

Review round: one independent reviewer sign-off round, with at most two author correction passes.
The reviewer approved this specification, its ticket graph, the implementation line, and the ownership fences on 2026-10-02.
Ticket 1 has completed author verification, independent review, and its checkpoint.
Its behavioral evidence is in [first-ticket verification](assets/first-ticket-verification.md).
The additional artifact repair is recorded in [cycle 3 evidence](assets/repair-cycle-3.md).

### Author transfer

The prior ticket author stopped before the original session took ownership.
The prior coordinator confirmed no worktree processes and verified restoration of the doctor-route mutation.
All eight uncommitted implementation files remain preserved.
The prior delegated assignment and its source remain recorded at commit 0adf2a61b588dd6c2ba4f8f465fac0eb8f7bdc1f.
Earlier verification retains its original author and source; the current session performs fresh verification.
The author session is codex:01a0fc58-0f59-7f52-9c4c-a25247d5fa21.

The scoped comment repair is recorded in [fourth-cycle verification](assets/repair-cycle-4.md).

## Completion plan

The pre-build audit found no implementation drift from the approved source.
The independent amendment review found no blocking issue.

The user directed implementation in the original authoring session on 2026-10-02.
The version 1 plan records this supported single-session execution form.
The user-directed session retains its declared Astra / high line.
Independent review still covers all three axes after each chunk.
The initial implementation has no numeric iteration cap.
The bounded repair policy still limits each chunk to two post-review repair cycles.

Ticket 1 verifies the doctor route with the Doctor-route mutation procedure.
The system check must reject that swap through TestCompatibilityMissingPath.
After the probe restores the source, the system check must pass.

The reviewer approved the inherited predecessor-base assertion repair on 2026-10-02.
C1 uses the shared TOON encoder for the expected revision row.
The assertion still rejects the wrong base and an earlier chunk in the diff.
The evidence command package and a wrong-base mutation verify this repair.

The reviewer extended C1 by one repair cycle for unexpected live artifact writes.
The package check uses the installed Node runtime that satisfies the declared floor.
The manifest-preservation test uses the existing private kit-copy fixture.
Its byte and absence assertions remain unchanged.
The manifest-directory swap must make that same test fail.


### Final integration continuation

The reviewer directed completion and landing on 2026-10-03.
The coordinator composed current main through Bench after resolving one changelog insertion conflict without changing its text.
The composed source passed the whole-project gate.

The current version 2 plan records the fresh C3 successor and the actual integration coordinator.
C1 and C2 keep their accepted version 1 source records and original author identities.
Their empty current histories indicate no new author dispatch for those completed tickets.

The C3 successor reconciles retained live evidence and repeats verification on the composed source.
Fresh independent axes review that source before its checkpoint.
The three chunk IDs, acceptance rows, requirements, and behavior remain unchanged.
C3 retains one consumed repair cycle; this integration and evidence update adds no product repair.
The original author history above describes the earlier version 1 execution.

```bench-completion-plan
{"version":2,"chunks":[{"id":"C1","tickets":["1-diagnose-interface.md"],"verification":[{"id":"compatibility","command":"bench test --package ./internal/compatibility","ticket":"1-diagnose-interface.md"},{"id":"adopt","command":"bench test --package ./internal/adopt","ticket":"1-diagnose-interface.md"},{"id":"system","command":"bench test --check system","ticket":"1-diagnose-interface.md"},{"id":"doctor-route-probe","command":"doctor-route-system-swap: follow the Doctor-route mutation procedure","probe":"swap","ticket":"1-diagnose-interface.md"},{"id":"evidence-command","command":"bench test --package ./internal/preflight/evidencecmd","ticket":"1-diagnose-interface.md"},{"id":"run-binary","command":"bench test --package ./internal/runbinary","ticket":"1-diagnose-interface.md"},{"id":"manifest-directory-probe","command":"bench probe internal/runbinary/runbinary.go --swap 'return runBuildScript(ctx, sourceRoot, output, filepath.Dir(output))' --with 'return runBuildScript(ctx, sourceRoot, output, \"\")' --package ./internal/runbinary --run '^TestBuildLeavesTheWrapperManifestUntouched$'","probe":"swap","ticket":"1-diagnose-interface.md"}]},{"id":"C2","tickets":["2-repair-managed-integration.md"],"verification":[{"id":"compatibility","command":"bench test --package ./internal/compatibility","ticket":"2-repair-managed-integration.md"},{"id":"adopt","command":"bench test --package ./internal/adopt","ticket":"2-repair-managed-integration.md"},{"id":"repair","command":"bench test --package ./internal/adopt/repairtest","ticket":"2-repair-managed-integration.md"},{"id":"transaction","command":"bench test --package ./internal/adopt/transaction","ticket":"2-repair-managed-integration.md"},{"id":"system","command":"bench test --check system","ticket":"2-repair-managed-integration.md"}]},{"id":"C3","tickets":["3-check-and-qualify-sessions.md"],"verification":[{"id":"compatibility","command":"bench test --package ./internal/compatibility","ticket":"3-check-and-qualify-sessions.md"},{"id":"adopt","command":"bench test --package ./internal/adopt","ticket":"3-check-and-qualify-sessions.md"},{"id":"system","command":"bench test --check system","ticket":"3-check-and-qualify-sessions.md"},{"id":"live-qualification","command":"review assets/qualification.md against the source-bound live acceptance rows","ticket":"3-check-and-qualify-sessions.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/cli-desktop-consistency/spec.md"},{"id":"system","command":"bench test --check system"},{"id":"live-qualification","command":"review actual CLI and desktop transcripts against every live acceptance row"}],"execution":{"mode":"delegate","run_id":"cli-desktop-consistency-20261002","orchestrator_session":"codex:01a101a2-3138-7cc3-a799-cb270ed0b460","author_limit":1,"assignments":{"1-diagnose-interface.md":[],"2-repair-managed-integration.md":[],"3-check-and-qualify-sessions.md":[{"session":"codex:01a10131-751b-71a2-b998-8f7cadba1da0","assignment":"f7123d5b2ad592acc6e5a239c1f1f389","model":"gpt-6-astra","effort":"high","source":"fa1dc27d7730ca3889dad592dd51dff6d7b8a899","native_ref":"reviews/cli-desktop-consistency.md: C3 confirming review and author verification"},{"session":"codex:01a102ad-b534-7561-a10f-aa45d97e554b","assignment":"f7123d5b2ad592acc6e5a239c1f1f389","model":"gpt-6-astra","effort":"high","source":"5ca38033cb2dc72e9d109dc440bdf89636942ee9","native_ref":"native:/root/c3_closeout_author; ready in Desktop session 01a101a2-3138-7cc3-a799-cb270ed0b460","predecessor":"codex:01a10131-751b-71a2-b998-8f7cadba1da0","trigger":"user-directed","stopped":"The prior author completed its fa1dc27 handoff; no prior author transaction remains. The successor waited until the coordinator merge terminated green at native 69995d.","preserved":"fa1dc27d7730ca3889dad592dd51dff6d7b8a899 and all original native review and verification occurrences"}]}}}
```
