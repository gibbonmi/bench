# Prospective artifact recovery

Status: implemented

Roadmap: FT251

Decision source: roadmap/FT251.md, the reviewed artifact from the 2026-08-24 stable-owner-landing review

Verification log: 2 Terra xhigh iterations; all cited blockers repaired, and the user stopped further review before a final acceptance pass.

## Problem

A prospective gate creates a private Git worktree and an owner-authored run binary.
Normal completion removes both resources through deferred cleanup.

SIGKILL prevents the deferred cleanup from running.
The dead landing then leaves the worktree, its Git registration, and its run binary in temporary storage.
No later prospective operation can identify and remove that residue.

## Solution

One prospective artifact owner manages the private worktree and every owner-authored run binary as one bundle.
The bundle records the owner PID and repository identity before it creates the private worktree.

Each prospective operation first sweeps recognized bundles for its repository in the current operating-system temporary directory.
The sweep removes a bundle only when the owner PID is definitely absent.
It retains live, invalid, and foreign bundles.

Normal completion closes the bundle in process.
A later prospective operation closes a bundle whose owner died before cleanup.

## User stories

### Recover dead-owner resources

Line: gpt-5.6-terra / high.

This group introduces the new lifecycle seam and its process recovery proof.

1. As an operator, I want the next prospective operation to remove a dead owner's bundle, so that SIGKILL does not leave durable residue.
2. As a reviewer, I want each bundle bound to one repository, so that a sweep cannot remove another repository's resources.
3. As an operator, I want bundle cleanup to remove the Git registration, so that Git does not retain a stale worktree entry.
4. As an operator, I want repeated sweeps to be idempotent, so that a retry cannot fail on an already-removed bundle.
5. As an operator, I want recovery before new artifact creation, so that every new run starts after known dead residue is removed.

### Preserve live and unproved resources

Line: gpt-5.6-terra / high.

This group makes automatic deletion conservative under concurrency and hostile state.

6. As an operator, I want a live owner's bundle retained, so that a concurrent prospective operation cannot destroy active work.
7. As a reviewer, I want an invalid owner record to deny deletion, so that malformed state cannot grant destructive authority.
8. As an operator, I want foreign temporary entries retained, so that a Bench prefix cannot grant ownership.
9. As an operator, I want an answering PID retained regardless of age, so that PID reuse causes a safe leak instead of deletion.
10. As an operator, I want each concurrent bundle classified independently, so that one dead owner cannot expose one live owner's resources.
11. As an operator, I want hostile bundle paths confined exactly, so that shell text cannot widen deletion scope.

### Preserve prospective execution

Line: gpt-5.6-terra / medium.

This group keeps the accepted gate and landing contracts while it changes resource ownership.

12. As an operator, I want normal terminal outcomes to remove the current bundle, so that recovery does not weaken ordinary cleanup.
13. As a reviewer, I want inherited run binaries excluded from bundle ownership, so that cleanup cannot remove a trusted baseline executable.
14. As a reviewer, I want every prospective checkout producer to use the owner convention, so that no shared caller keeps the original leak.
15. As an operator, I want exact green evidence reuse unchanged, so that artifact recovery does not add another full gate run.
16. As a reviewer, I want recovery to execute no old candidate bytes, so that residue cannot become cleanup authority.

## Implementation decisions

A narrow prospective artifact module owns the bundle lifecycle.
The gate package consumes this module instead of coordinating temporary paths itself.
The module is not a generic cleanup framework and exposes no public command.

The bundle root uses one dedicated prefix in the current operating-system temporary directory.
It contains a fixed checkout child and every owner-authored run-binary directory.
The layout makes the bundle root the one source for teardown scope.

The owner writes one strict record before it creates the checkout child.
The record contains a schema, a positive owner PID, and the canonical Git common-directory identity.
The owner publishes the record atomically as a regular 0600 file.

The sweep scans only the dedicated bundle prefix.
It accepts only a canonical bundle root with a valid record for the current repository.
It does not follow a symbolic link or accept a special file.

The process probe treats success and permission refusal as live.
It treats only the operating system's absent-process result as dead.
All other probe results retain the bundle.

The sweep removes a recognized dead bundle in a fixed order.
It first removes a present Git worktree registration.
It then removes the one bundle root and its owner-authored binary children.
A cleanup error refuses the new prospective operation and retains the unresolved bundle.

The prospective gate supplies the bundle root to the existing run-binary factory.
The existing source-selection and executable verification rules do not change.
An inherited selection remains outside the bundle and keeps its existing lifetime.

The full gate, evidence inspection, and fast lane use the same bundle owner.
The sweep runs before each of these producers creates a prospective checkout.
Post-publication landing resume does not create a prospective checkout and remains unchanged.

The installed wrapper rejects routing overrides and validates the installation manifest's broker path, version, platform, and digest.
It then executes the installed broker under the existing operating-system Git and toolchain trust assumption.
Fresh candidate control starts only when that broker intentionally builds and runs the new prospective checkout.

The broker executes no recovered artifact while it classifies or removes a dead bundle.
Candidate bytes do not authenticate the broker or the cleanup decision.

## Testing decisions

The highest seam is a real Git child-process journey through prospective authorization.
The journey blocks after the checkout and run binary exist.
The test kills the owner process group and starts a fresh authorization process.

The fresh process proves removal of the old bundle, Git registration, and run binary.
It also proves that the new prospective subject reaches its ordinary result.
A second fresh process proves idempotency across the process boundary.

The module seam receives table tests for record shape, repository identity, PID liveness, and partial state.
The tests use real temporary directories and inject only the operating-system process probe.
Existing stable-owner cleanup tests remain the prior art for normal terminal outcomes.

An ordering seam observes the published record before checkout creation starts.
Partial-state tests cover a record-only bundle, a checkout without a binary, and a stale Git registration.
A recovered candidate executable carries a planted marker that must never appear.

### Seam diagram

    prospective gate, inspection, or fast lane
                      │
                      ▼
    [ prospective artifact owner ] ──▶ dead-owner sweep
                      │                      │
                      │                      └──▶ remove old Git registration and bundle
                      ▼
       owner record + private checkout + owner-authored run binary
                      │
                      ▼
             prospective execution

    process tests attach at prospective authorization
    lifecycle tests attach at the prospective artifact owner

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| PAR01 | 1 | A fresh authorization removes a dead owner's checkout and owner-authored run binary before it runs. | child-process prospective authorization | A defer-only implementation leaves both resources after SIGKILL. |
| PAR02 | 2 | A bundle record for another Git common directory cannot authorize deletion. | prospective artifact owner | A prefix-only sweep crosses repository ownership. |
| PAR03 | 3 | A dead-bundle sweep removes the exact Git worktree registration before it removes the bundle root. | real Git lifecycle seam | A filesystem-only cleanup leaves stale Git administration state. |
| PAR04 | 4 | A second fresh-process sweep of the removed bundle changes no path. | child-process prospective authorization | A non-idempotent cleanup fails after the first removal. |
| PAR05 | 5 | The old dead bundle is absent before the new prospective checkout is materialized. | prospective artifact owner | A post-run sweep leaves residue throughout the next execution. |
| PAR06 | 1 | The owner record is visible before the checkout-creation seam runs. | prospective artifact ordering seam | A checkout-first implementation exposes an unowned private tree. |
| PAR07 | 1 | The published owner record has the exact schema, PID, common-directory identity, regular-file type, and 0600 mode. | prospective artifact ordering seam | An incomplete record cannot prove the deletion scope. |
| PAR08 | 1 | A dead record-only bundle is removed by the next sweep. | prospective artifact owner | Cleanup that assumes a checkout exists leaves interrupted creation residue. |
| PAR09 | 1 | A dead bundle with a checkout and no run binary is removed by the next sweep. | real Git lifecycle seam | Cleanup that assumes a completed binary build leaves the checkout. |
| PAR10 | 3 | A dead bundle with only a stale Git registration is removed by the next sweep. | real Git lifecycle seam | Cleanup that requires a present checkout leaves Git administration state. |
| PAR11 | 6 | A blocked live owner keeps its checkout and owner-authored run binary during a concurrent authorization. | concurrent child-process authorization | An unconditional prefix sweep deletes active resources. |
| PAR12 | 7 | An unsupported owner-record schema retains the bundle root unchanged. | prospective artifact owner | A permissive decoder grants deletion authority to unknown data. |
| PAR13 | 7 | A missing, empty, or malformed owner record retains the bundle root unchanged. | prospective artifact owner | A parse failure becomes false deletion authority. |
| PAR14 | 7 | A symbolic-link or special-file owner record retains the bundle root unchanged. | prospective artifact owner | A file-shape failure becomes false deletion authority. |
| PAR15 | 7 | A symbolic-link or special-file bundle-root candidate remains unchanged. | prospective artifact owner | A scanner that follows or opens the candidate escapes its namespace. |
| PAR16 | 8 | A foreign same-prefix directory remains byte-for-byte unchanged after a sweep. | prospective artifact owner | Directory naming alone grants destructive authority. |
| PAR17 | 6 | A permission-refused process probe retains the bundle unchanged. | injected process-probe seam | Treating EPERM as dead deletes a possibly live owner's resources. |
| PAR18 | 7 | A non-ESRCH process-probe error retains the bundle unchanged. | injected process-probe seam | An unknown probe result becomes false death proof. |
| PAR19 | 9 | A PID that answers live retains its bundle even when its record is old. | injected process-probe seam | An age fallback can delete a live or reused PID's bundle. |
| PAR20 | 10 | One sweep removes a dead-owner bundle and retains an independently live-owner bundle. | concurrent child-process authorization | A global cleanup conflates two artifact owners. |
| PAR21 | 11 | A bundle path with spaces and glob characters stays confined to its exact root. | real Git lifecycle seam | Shell expansion widens cleanup beyond the recorded bundle. |
| PAR22 | 12 | A green prospective execution leaves no current bundle. | prospective execution lifecycle | The new owner can omit cleanup on success. |
| PAR23 | 12 | A red prospective execution leaves no current bundle. | prospective execution lifecycle | The new owner can omit cleanup on a gate refusal. |
| PAR24 | 12 | A timed-out prospective execution leaves no current bundle. | prospective execution lifecycle | The new owner can omit timeout cleanup. |
| PAR25 | 12 | A cancelled prospective execution leaves no current bundle. | prospective execution lifecycle | The new owner can omit interrupt cleanup. |
| PAR26 | 12 | A run-binary build refusal leaves no current bundle. | prospective run-binary selection | The new owner can omit pre-selection cleanup. |
| PAR27 | 13 | A sweep never removes the inherited run binary selected from the baseline. | prospective run-binary selection | Broad binary cleanup deletes executable state it does not own. |
| PAR28 | 14 | A full-gate checkout publishes a valid owner record. | prospective checkout producer | The full gate can remain on the defer-only helper. |
| PAR29 | 14 | An evidence-inspection checkout publishes a valid owner record. | prospective checkout producer | Inspection can remain on the defer-only helper. |
| PAR30 | 14 | A fast-lane checkout publishes a valid owner record. | prospective checkout producer | The fast lane can remain on the defer-only helper. |
| PAR31 | 15 | An identical retry after dead-bundle recovery reuses the same exact green evidence. | gate evidence seam | Cleanup coupled to evidence forces another full gate run. |
| PAR32 | 16 | Recovery removes a planted old candidate executable without running its marker. | child-process prospective authorization | Recovered candidate bytes execute before the fresh subject starts. |
| PAR33 | 3 | A Git registration-removal failure refuses before new checkout creation and retains the recognized dead bundle root. | prospective artifact owner | An ignored Git cleanup error leaves residue and starts another producer. |

### Edge inventory

- The owner is live, absent, permission-protected, or represented by a reused PID.
- The record is missing, empty, malformed, unsupported, nonregular, or a symbolic link.
- The bundle root is regular, a symbolic link, or a special file.
- The record names the current repository, another repository, or a missing repository.
- The bundle has no checkout, no run binary, a stale registration, or a complete artifact set.
- One dead bundle exists beside one live bundle.
- The bundle path contains spaces or glob characters.
- The gate ends green, red, timed out, cancelled, or refused during the run-binary build.
- The prospective producer is the full gate, evidence inspection, or fast lane.
- The selected run binary is owner-authored or inherited.
- A sweep runs once or runs again after complete removal.
- Git registration removal succeeds or returns a controlled failure.

**Won't handle** a legacy unrecorded prospective tree — a new prospective operation is the surviving in-scope owner-record producer.

**Won't handle** a changed temporary-directory root between failure and recovery — the unchanged root is the surviving in-scope discovery surface.

**Won't handle** immediate reclamation while an unrelated process reuses the owner PID — a later sweep after that PID exits is the surviving in-scope recovery path.

**Won't handle** descendants that survive the killed owner — FT201 owns process-group parent-death behavior for the surviving in-scope command path.

## Ownership fences

- `internal/gate/prospectiveartifact/`
- `internal/gate/engine.go`
- `internal/gate/lane.go`
- `internal/gate/prospective.go`
- `internal/gate/prospective_owner_test.go`
- `internal/gate/lane_test.go`
- `internal/worktree/land_freshness_test.go`
- `internal/systemtest/owner_artifact_recovery_test.go`
- `projects/benchkit.md`
- `CHANGELOG.md`

The two tickets are serial because both edit the prospective artifact owner.

## Out of scope

- A public artifact-cleanup command is a separate operator surface: 6 edits, 2 gate runs.
- A repository registry that spans temporary-directory changes is separate durable state: 8 edits, 3 gate runs.
- A process-start identity probe for immediate PID-reuse cleanup is a separate platform capability: 10 edits, 3 gate runs.
- Parent-death enforcement for descendant process groups is FT201: 12 edits, 3 gate runs.

## Further notes

FT251 is a learnings-origin Map candidate under the kit synthesis discipline.
No existing owner joins the prospective checkout and its owner-authored binary lifecycle.
The narrow module fills that gap without changing public command grammar.
