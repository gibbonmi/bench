# Managed integration repair verification

## Scope

C2 implements CD18 through CD36 in the original approved author session.
The native review and C2 checkpoint remain pending.
No post-review repair cycle has been consumed for C2.
One preservation hardening pass extended the approved cases.

The tests use private repositories and configuration homes.
No test changes the active CLI or desktop configuration.
The desktop normal shell still fails before process creation.
Elevated diagnostic execution does not qualify that interface.

## Verification

| Route | Result | Time in milliseconds |
| --- | --- | --- |
| Complete adoption package | Pass; no skips | 115642 |
| Complete public repair package | Pass; no skips | 90236 |
| Complete transaction package | Pass; no skips | 476 |
| Compatibility package and fixture support | Pass; no skips | 4 |
| Complete command package | Pass; no skips | 26255 |
| Sealed system suite | Pass; no skips | 98007 |
| Conformance package | Pass; three capability skips | 57271 |
| Structure growth from C1 checkpoint | Pass | not recorded |

The conformance package tests do not replace a live-tree gate.
Two capability skips report unavailable Unix socket bindings.
The third reports unavailable permission to create a character device.
No environment skip occurred.

The sealed system suite tests every competing adoption writer with a held destination lock.
It covers separate repositories, aliased hook paths, different Bench homes, and an unrelated writer.
Its interrupted-repair case starts a fresh process to undo a partially published transaction.

## Undo omission

The probe changed the undo loop start from the final entry to the preceding entry.
This omits one retained destination while keeping the mutation compilable.
TestCompatibilityUndo failed its exact bytes and mode comparison in 1759 milliseconds.
The probe reported bit and restored the source.
The restored test passed in 1667 milliseconds with no skips.

## Preservation defects and repairs

| Test | Observed red | Verified repair |
| --- | --- | --- |
| TestCompatibilityKitUpdatesCanonicalBroker | An older canonical manifest was refused | The existing manifest producer validates the old binding before replacement |
| TestRepairPreservesRepositoryRefs | Repair changed a remote reference outside its record | Compatibility repair leaves remote references unchanged |
| TestUndoRefusesInvalidRecord | Unknown state and public recovery files admitted undo | The record reader validates state and private modes |
| TestUndoRestoresCompleteMode | Undo lost a special permission bit | Publication restores the complete mode after the data write |
| TestRepairRefusesPublicNamespace | Repair used a public recovery namespace | An existing unsafe namespace is refused without permission changes |
| TestCompatibilityKitModifiedConflict | Changed shim and hook modes lacked conflict reports | Shared eligibility checks preserve and report those assets |

The public repair suite verifies these repairs together.
The broker update also verifies exact restoration of its previous manifest.
The retained transaction suite verifies incomplete recovery reporting and preservation of available backups.

## Fixture corrections

The repeated-repair fixture commits before both repair calls.
Thus, both calls observe the same branch-resolution prerequisites.
The writer census preserves an earlier test's modified gate before preparing its second repository.
Its assertions read the complete bounded command report, including the spill.
These corrections keep the original pass criteria and supply no behavioral-red claim.

## Live-tree gate

Run 20261002T232002.729534954Z-2816038 completed with a red test phase.
Its only failed assertions required canonical citations for the realized C2 coverage rows.
The formatting, vet, race, and system phases passed.
Every behavioral package also passed.

The map now uses the required path and test-name citations.
The coverage owner accepts all 64 rows.
C3 still has 13 planned hermetic seams without realized citations.

A later alias regression failed in 36 milliseconds.
The selected shim path differed from its transaction's canonical destination.
Repair now uses the destination identity already captured by the shared transaction owner.
The regression also requires undo to restore the previous shim bytes through that alias.
All kit cases passed in 936 milliseconds after the correction.
The live-root documentation check passed in 1614 milliseconds with no skips.

## Final focused source

The final adoption-family run passed with no failures or skips.
Adoption took 57894 milliseconds, public repair took 47680, and transaction tests took 664.
This run includes the canonical shim-identity correction.
The earlier completed checks retain their own source and scope.

## Remaining work

Freeze the C2 source with the current verification record.
Then obtain the three independent native reviews and the C2 checkpoint.
C3 still owns session integration and actual CLI and desktop qualification.
