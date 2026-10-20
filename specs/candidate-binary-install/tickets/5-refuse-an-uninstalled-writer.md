# Refuse a declared build's guarded entries on an uninstalled executable

Blocked by: 1-declare-persisted-fields.md, 2-install-the-candidate.md, 3-hold-the-refresh.md
Writes: internal/commitment/repository/window.go, internal/commitment/repository/install_check_test.go (new), internal/brokermanifest/brokermanifest.go, internal/refusalroute/faces_gate.go, internal/refusalroute/registry.go, internal/refusalroute/route_test.go, internal/gate/gate.go, internal/gate/checkpoint.go, internal/gate/complete_checkpoint.go, internal/gate/refusal_route_test.go, internal/gate/checkpoint_persisted_test.go (new), internal/reviewrecord/recordcmd/assignment.go, internal/reviewrecord/recordcmd/command.go, internal/reviewrecord/recordcmd/refusal.go, internal/reviewrecord/recordcmd/assignment_persisted_test.go (new), internal/conformance/refusal_route_bypass_test.go, internal/conformance/refusal_route_guard_test.go
Covers: CB26, CB27, CB28, CB29, CB30, CB31, CB32, CB33, CB34, CB35, CB36, CB57, CB60

## What to build

Deliver CB-C3: one digest-only candidate check guards the chunk checkpoint, the complete checkpoint, and the gap close of `bench record <spec> assignment` of a declared build.

Put the check owner in `internal/commitment/repository/window.go`, beside the open-build answers from ticket 3. Both writer packages already import this package. The check takes the declaration, the running executable's path, and the primary broker manifest path. It reads the manifest with `brokermanifest.Read` and digests the executable with `brokermanifest.Digest`. The caller derives the manifest path from the primary checkout's root and `brokermanifest.Name`, never from `BENCH_WRAPPER`, because a worktree's own `bin/` can hold a manifest of the worktree executable. The check applies only where the repository declares Bench build inputs. It refuses when the digests differ. It does not test field knowledge.

Register one face, `candidate-not-installed`, in the gate face list. The registry records the owner of the rule, so `recordcmd` may print this gate face. The route's first step is an instruction, not a command step: run the install form from ticket 2 from the primary checkout, with the target fact. The merge face list holds the prior art for an instruction step. The second step is the registry's generic composed rerun, `rerun` (`internal/refusalroute/registry.go:168`), so each entry supplies its own rerun fact, as `internal/commit/commit.go:248-249` does with `CommitRerun`. The checkpoint entries render that fact from the existing gate rerun over their label and arguments, through one exported renderer in `faces_gate.go`. The record entry supplies its own `bench record <spec> assignment ...` command. Thus a refused record command never prints `bench gate`. An agent command step for the install would make the route walk run the install inside the gate test, which cannot work there. When the ledger reads, the target is the integration target from ticket 3's declared-build answer, whatever worktree the caller runs in. When the ledger does not read, the route prints the literal `<target>` slot. Add a fact or a text part to the registry for that slot. Add one producing fixture row to the gate face fixtures. That fixture carries the instruction by rebinding the primary manifest to the running test binary's digest before the rerun. Keep every route literal inside the registry so the route-bypass test stays green.

Each entry prints ticket 1's `persisted-fields-invalid` error when the declaration scan refuses.

Call the check at three entries, before any ledger read or write:

- `RunCommand` for `--chunk`, after the verb context opens and before the owner lookup, the review-record read, and the gap write.
- The `--complete` path, at the same point.
- `recordcmd/assignment.go`, before its record read and its gap-close transaction. The record verb has no route printer today, so add the printing path that prints the face and its `next=` line.

An undeclared spec makes the check a no-op, so every existing gate and record fixture keeps its result.

## Acceptance

- [ ] `TestCandidateCheck`: equal digests pass even when the declaration names a path that the current ledger type does not hold, different digests refuse, an undeclared spec passes with any manifest, and a worktree manifest beside `BENCH_WRAPPER` does not replace the primary manifest (CB26, CB27, CB28, CB60).
- [ ] A declared chunk checkpoint and a declared `--complete` checkpoint on a mismatched digest refuse with `candidate-not-installed` and leave the ledger bytes unchanged (CB29, CB30).
- [ ] A declared `bench record <spec> assignment` on a mismatched digest refuses and leaves the ledger bytes unchanged (CB31).
- [ ] With an unreadable ledger, the chunk checkpoint prints `candidate-not-installed` and not the decode error (CB32).
- [ ] The route prints the install command with the integration target label, prints the literal `<target>` slot when the ledger does not read, prints the rerun after the install, and names the integration target from an author worktree (CB33, CB34, CB35, CB57).
- [ ] One route test row shows that a refused `bench record <spec> assignment` prints a rerun that names `bench record`, and a refused checkpoint prints a rerun that names `bench gate` (CB35).
- [ ] A checkpoint that runs a version-cache executable after an install refuses (CB36).

## Checkpoint verification

Run CB-C3's verification commands from the spec's completion plan. Probe the digest comparison, the manifest locator, and each entry's call order with one mutant each, and record the reds. The orchestrator then freezes the CB-C3 delta for independent review.
