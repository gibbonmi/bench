# Consistent Bench behavior in CLI and desktop

Status: staged

Decision source: `specs/cli-desktop-consistency/decisions/cli-desktop-consistency.md` (ready compiled map).

Verification log: 0 iteration(s) to accept — independent reviewer sign-off is pending.

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

Line: gpt-5.6-sol / high.
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
| CD01 | 2 | The compatibility report names the explicitly selected CLI or desktop interface | planned TestCompatibilityInterface in internal/adopt/compatibility_test.go | Swapping the selected interface changes the expected context |
| CD02 | 2 | An unknown active runtime remains unknown when a launcher version is available | planned TestCompatibilityRuntimeProvenance in internal/compatibility/inspect_test.go | A PATH launcher cannot supply the active server identity |
| CD03 | 3 | The report names the active configuration home without assuming the other interface shares it | planned TestCompatibilityHomes in internal/compatibility/inspect_test.go | Different Linux and mounted Windows homes must remain distinct |
| CD04 | 3 | An unreadable effective setting reports unknown with its source | planned TestCompatibilityUnknownConfig in internal/compatibility/inspect_test.go | Dropping a failed source would fabricate agreement |
| CD05 | 4 | Diagnostics leave both user configuration homes byte-identical | planned TestCompatibilityReadOnly in internal/adopt/compatibility_test.go | Independent before-and-after snapshots catch silent synchronization |
| CD06 | 5 | An absent required repository asset reports its exact restoration action | planned TestCompatibilityMissingAsset in internal/adopt/compatibility_test.go | Omitting one required asset must produce a failed check |
| CD07 | 6 | An absent global Bench command reports the existing by-path route | planned TestCompatibilityMissingPath in internal/systemtest/compatibility_test.go | A fixture without global Bench must still reach the local wrapper |
| CD08 | 7 | A declared hook remains unverified until its live behavior is observed | planned TestCompatibilityDeclaredHook in internal/compatibility/inspect_test.go | A configuration file alone cannot produce a live pass |
| CD09 | 8 | An invalid interface operand exits 2 before configuration reads | planned TestCompatibilityGrammar in internal/adopt/compatibility_test.go | A read sentinel catches validation after inspection |
| CD10 | 8 | Special files and dangling configuration links return a bounded diagnostic | planned TestCompatibilityFileKinds in internal/compatibility/inspect_test.go | FIFO and dangling-link fixtures cannot become empty successful inputs |
| CD11 | 8 | An empty required configuration file receives a distinct result from an absent file | planned TestCompatibilityAbsentEmpty in internal/compatibility/inspect_test.go | Both fixture shapes must retain their different evidence states |
| CD12 | 8 | A path containing spaces and glob characters reaches the intended repository | planned TestCompatibilityPaths in internal/systemtest/compatibility_test.go | Decoy paths expose splitting or glob expansion |
| CD13 | 29 | A changed relevant context produces a different compatibility fingerprint | planned TestCompatibilityFingerprint in internal/compatibility/inspect_test.go | One-field mutations invalidate the previous observation context |
| CD14 | 29 | A personal model preference alone does not create a Bench configuration conflict | planned TestCompatibilityPersonalSettings in internal/compatibility/inspect_test.go | Two permitted preferences must not require synchronization |
| CD15 | 2 | The inspection result never labels the whole active chat qualified | planned TestCompatibilityInspectionBoundary in internal/adopt/compatibility_test.go | Successful local checks still require live tool evidence |
| CD16 | 8 | Control-bearing report fields cannot inject a terminal control sequence | planned TestCompatibilityOutput in internal/compatibility/inspect_test.go | The output seam must refuse unsafe cells through the existing renderer |
| CD17 | 26 | Legacy doctor invocations retain their prior exit and output behavior | planned TestCompatibilityLegacyDoctor in internal/adopt/compatibility_test.go | A differential fixture compares legacy modes before and after integration |
| CD18 | 1 | One setup installs the shared Bench integration without copying either user configuration home | planned TestCompatibilitySetup in internal/systemtest/compatibility_test.go | Fresh linked-repository fixtures keep unrelated home sentinels intact |
| CD19 | 9 | The compatibility repair restores an unmodified managed asset from the canonical payload | planned TestCompatibilityManagedRepair in internal/adopt/compatibility_repair_test.go | Removing one managed hook must become a working installed hook |
| CD20 | 9 | A second identical repair leaves all managed destination bytes and modes unchanged | planned TestCompatibilityRepairIdempotence in internal/adopt/compatibility_repair_test.go | Tracked repeated application exposes self-induced drift |
| CD21 | 10 | A modified managed asset remains unchanged and reports a conflict | planned TestCompatibilityModifiedConflict in internal/adopt/compatibility_repair_test.go | A user edit must survive the automatic repair path |
| CD22 | 10 | A foreign file at a repair destination remains unchanged | planned TestCompatibilityForeignConflict in internal/adopt/compatibility_repair_test.go | A marker-free sentinel cannot become Bench-owned |
| CD23 | 11 | A repair requiring security-policy change returns a decision action without applying it | planned TestCompatibilityPolicyBoundary in internal/adopt/compatibility_repair_test.go | Permission and trust-record sentinels must remain untouched |
| CD24 | 11 | A repair requiring process interruption returns a decision action without sending a signal | planned TestCompatibilityInterruptionBoundary in internal/adopt/compatibility_repair_test.go | A supervised live process must survive the repair attempt |
| CD25 | 11 | A private runtime path receives no automatic mutation | planned TestCompatibilityRuntimeBoundary in internal/adopt/compatibility_repair_test.go | An absent launcher path must stay outside the repair plan |
| CD26 | 12 | A backup-publication failure leaves all repair destinations unchanged | planned TestCompatibilityBackupFailure in internal/adopt/compatibility_fault_test.go | A denied record directory prevents the first destination write |
| CD27 | 12 | An interrupted repair retains enough state for a fresh process to recover each touched target | planned TestCompatibilityInterruptedRepair in internal/systemtest/compatibility_test.go | A second process must recover after an interrupted multi-target promotion |
| CD28 | 12 | A failed terminal record publication reports incomplete repair and preserves recovery data | planned TestCompatibilityTerminalFailure in internal/adopt/compatibility_fault_test.go | Successful destination writes cannot conceal a missing terminal record |
| CD29 | 13 | Concurrent Bench writers cannot interleave changes to the same repair destination | planned TestCompatibilityConcurrentWriters in internal/systemtest/compatibility_test.go | A barrier holds one writer while the competing writer attempts publication |
| CD30 | 14 | Undo restores every recorded preimage byte and mode | planned TestCompatibilityUndo in internal/adopt/compatibility_repair_test.go | Deleting the final backup entry must fail the complete restoration comparison |
| CD31 | 14 | Undo removes a target that the repair created from absence | planned TestCompatibilityUndoCreated in internal/adopt/compatibility_repair_test.go | An absent preimage must not become an empty surviving file |
| CD32 | 14 | Undo refuses a destination changed after repair | planned TestCompatibilityUndoConflict in internal/adopt/compatibility_repair_test.go | A later edit must survive the stale undo request |
| CD33 | 15 | Repair records contain only managed preimages and their necessary restoration metadata | planned TestCompatibilityRepairPrivacy in internal/adopt/compatibility_repair_test.go | Credential and raw-environment sentinels must not occur in the record |
| CD34 | 8 | A hostile repair identifier cannot select a path outside its private record directory | planned TestCompatibilityRepairPaths in internal/adopt/compatibility_fault_test.go | Traversal and symlink fixtures cannot redirect reads or writes |
| CD35 | 9 | A repair rechecks destination identity immediately before its publication | planned TestCompatibilityPlanDrift in internal/adopt/compatibility_fault_test.go | A target replaced after inspection must prevent its planned write |
| CD36 | 12 | A failed restore reports each unresolved target and preserves its available backup | planned TestCompatibilityRestoreFailure in internal/adopt/compatibility_fault_test.go | A restore error cannot become a successful rollback |
| CD37 | 16 | SessionStart emits the compatibility check obligation before dependent work | planned TestCompatibilityStartup in internal/sessioninspect/compatibility_test.go | Removing the lifecycle call loses the required next action |
| CD38 | 17 | A resumed session repeats the active-interface checks | planned TestCompatibilityResume in internal/systemtest/compatibility_test.go | A prior success followed by a changed context cannot suppress the resume check |
| CD39 | 18 | A hook-process success cannot replace the actual chat-tool probe | review-owned: `assets/qualification.md`, actual-tool transcript | A hook-only success beside a failed tool must remain unqualified |
| CD40 | 18 | An escalated diagnostic success cannot replace normal-permission evidence | review-owned: `assets/qualification.md`, permission comparison | The observed incident supplies the cheap wrong-success control |
| CD41 | 19 | A missing optional desktop tool leaves an unrelated Bench operation available | planned TestCompatibilityOptionalTool in internal/compatibility/inspect_test.go | A disabled preview must not block a shell-only diagnostic |
| CD42 | 19 | A missing required capability stops its dependent operation before a mutation | review-owned: `assets/qualification.md`, missing-required-capability transcript | A mutation sentinel must remain absent when its required check fails |
| CD43 | 20 | The shipped instructions expose a recovery action when no Bench command can start | review-owned: shipped README and active-session failure transcript | A command-only error route cannot satisfy the failed-shell scenario |
| CD44 | 21 | A repair result requires a retest through the previously failed interface | planned TestCompatibilityRetestRequired in internal/sessioninspect/compatibility_test.go | Printing repair success must not print recovered compatibility |
| CD45 | 21 | The qualification record includes successful recovery through normal desktop tools | review-owned: `assets/qualification.md`, recovery transcript | A still-failing desktop leaves the complete outcome unaccepted |
| CD46 | 22 | Independent CLI and desktop chats obtain the same repository Bench rules | review-owned: `assets/qualification.md`, two fresh chats | Shared conversation context cannot supply the missing integration |
| CD47 | 23 | Concurrent CLI and desktop writers use distinct Bench assignments | review-owned: `assets/qualification.md`, concurrent worktree exercise | Both chats must produce separate owned paths before their writes |
| CD48 | 23 | Both interfaces complete the existing Bench review and landing route | review-owned: `assets/qualification.md`, disposable repository lifecycle | Raw Git publication cannot satisfy the recorded lifecycle |
| CD49 | 24 | The actual CLI interface rejects the harmless forbidden-command fixture | review-owned: `assets/qualification.md`, CLI guard transcript | The rejected command sentinel must never execute |
| CD50 | 24 | The actual desktop interface rejects the harmless forbidden-command fixture | review-owned: `assets/qualification.md`, desktop guard transcript | The rejected command sentinel must never execute |
| CD51 | 25 | A verified equivalent review route produces the same review outcome | review-owned: `assets/qualification.md`, review comparison | A declared equivalent without an observed outcome is insufficient |
| CD52 | 26 | Existing macOS and Linux contract fixtures retain their prior outcomes | planned TestCompatibilityOtherHosts in internal/systemtest/compatibility_test.go | Platform-specific branches must not change established non-WSL behavior |
| CD53 | 27 | Native Windows Bench execution remains outside the new qualification claim | planned TestCompatibilitySupportBoundary in internal/compatibility/inspect_test.go | A Windows-native fixture must not report WSL qualification |
| CD54 | 28 | The complete outcome remains unaccepted while any required live qualification row is missing | review-owned: full acceptance reconciliation before landing | A green package suite cannot erase missing actual-interface evidence |
| CD55 | 29 | A runtime version difference alone does not fail a capability check | planned TestCompatibilityVersionDifference in internal/compatibility/inspect_test.go | Different version labels with identical observed capabilities must not conflict |
| CD56 | 29 | Unknown context prevents reuse of an earlier live check | planned TestCompatibilityUnknownContext in internal/sessioninspect/compatibility_test.go | A missing runtime or policy observation cannot inherit a prior pass |
| CD57 | 30 | Each actual interface writes and reads the expected bytes in its own scratch file | review-owned: `assets/qualification.md`, CLI and desktop file probes | A subprocess fixture alone cannot qualify the actual file tools |
| CD58 | 30 | Each actual interface invokes the repository skill from its installed path | review-owned: `assets/qualification.md`, skill invocation transcripts | An available file without a usable harness skill does not qualify |
| CD59 | 24 | A required permission conflict remains visible without broadening the configured policy | planned TestCompatibilityPermissionConflict in internal/systemtest/compatibility_test.go | A policy sentinel exposes an automatic bypass |
| CD60 | 20 | An unknown upstream repair remains explicitly unresolved in the recovery guidance | planned TestCompatibilityUnknownRepair in internal/sessioninspect/compatibility_test.go | A fabricated repair action cannot turn unknown into recovered |
| CD61 | 16 | A compatibility inspection timeout preserves the existing informational startup exit | planned TestCompatibilityStartupTimeout in internal/sessioninspect/compatibility_test.go | Slow discovery cannot prevent the chat from opening |
| CD62 | 17 | A relevant configuration change invalidates the current session observation | review-owned: `assets/qualification.md`, configuration drift transcript | A changed hook or permission context must trigger another affected check |
| CD63 | 1 | A linked repository receives the compatibility instructions through the real payload installer | planned TestCompatibilityPayload in internal/systemtest/compatibility_test.go | Kit-only documents cannot satisfy a linked consumer |
| CD64 | 19 | An equivalent route counts only after that operation succeeds through the route | planned TestCompatibilityEquivalentEvidence in internal/compatibility/inspect_test.go | Naming an alternate tool cannot satisfy the required capability |

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

- `CONTEXT.md`
- `.bench/BENCH-reference.md`
- `.bench/BENCH.md`
- `.bench/hooks/session-start.sh`
- `CHANGELOG.md`
- `DATA_HANDLING.md`
- `README.md`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/main.go`
- `internal/adopt`
- `internal/adopt/compatibility.go`
- `internal/adopt/compatibility_test.go`
- `internal/adopt/doctor.go`
- `internal/anchors`
- `internal/compatibility`
- `internal/conformance/entry_point_parity_bite_test.go`
- `internal/conformance/entry_point_parity_test.go`
- `internal/sessioninspect`
- `internal/systemtest/compatibility_test.go`

Reviewer disposition: proposed for approval with the ticket graph.
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
The reviewer approves these engineering choices with the seams and fences.

Review round: one independent reviewer sign-off round, with at most two author correction passes.
The current artifact has not received that sign-off.
No implementation is authorized by specification authorship alone.

## Completion plan

```bench-completion-plan
{"version":1,"chunks":[{"id":"C1","tickets":["1-diagnose-interface.md"],"verification":[{"id":"compatibility","command":"bench test --package ./internal/compatibility"},{"id":"adopt","command":"bench test --package ./internal/adopt"},{"id":"system","command":"bench test --check system"}]},{"id":"C2","tickets":["2-repair-managed-integration.md"],"verification":[{"id":"adopt","command":"bench test --package ./internal/adopt"},{"id":"system","command":"bench test --check system"}]},{"id":"C3","tickets":["3-check-and-qualify-sessions.md"],"verification":[{"id":"compatibility","command":"bench test --package ./internal/compatibility"},{"id":"adopt","command":"bench test --package ./internal/adopt"},{"id":"system","command":"bench test --check system"},{"id":"live-qualification","command":"review assets/qualification.md against the source-bound live acceptance rows"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/cli-desktop-consistency/spec.md"},{"id":"system","command":"bench test --check system"},{"id":"live-qualification","command":"review actual CLI and desktop transcripts against every live acceptance row"}]}
```
