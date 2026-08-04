# Validate the persisted checkpoint ref before reclamation trusts it

Blocked by: none
Ownership fence: `internal/specbuild`
Contracts: the stored `CheckpointRef` crosses `internal/specbuild/checkpoint.go`→`internal/specbuild/state.go`'s residue enumeration as authorization to delete the ref it names, and is asserted by CR1 and CR2 through the service's own save and load rather than a hand-mangled state file the loader would refuse anyway; the checkpoint ref's canonical form crosses `checkpoint.go`→`validCore` and is asserted by CR5 as one construction rather than two
Assumptions: `refs/bench/specbuild/checkpoint/` plus `digest(Run + "\x00" + assignment ID)` in `checkpoint.go:195` is the one construction of a checkpoint ref name and stays that way; terminality still comes from the owning record's own flag and this ticket does not touch that judgment; claims re-derived from the tree at pickup

## What to build

A critical authority defect found by the authoritative review of candidate `64fda745`,
reproduced against a checkout of that commit before this ticket was written.

`validCore` in `internal/specbuild/state.go` validates a record's identity, spec, branch,
base, candidate tip, and operation table — and accepts an assignment's `CheckpointRef` as
an arbitrary string. `provisionalResidue` then resolves that stored string with
`s.refObject(assigned.CheckpointRef)`, checking neither the namespace it sits in nor
whether it matches the ref the run would have written, and claims it with the owning
record's disposition. For a terminal record that disposition is `reclaimable`, and
`ApplyReclaim` deletes every reclaimable ref.

The observed reproduction: a terminal run whose assignment `CheckpointRef` is set to
`refs/heads/main`, written through the service's own `save` and read back through its own
`load` — so `record.valid` is the trust boundary that accepted it, not a mangled file —
plans `refs/heads/main` as `reclaimable` and `ApplyReclaim` deletes it, returning a nil
error and a receipt naming the default branch among the refs it spent.

Validate the persisted ref where the record's other identity facts are validated: a stored
`CheckpointRef` is either empty or exactly the ref that run and assignment would have
written, derived from the one construction in `checkpoint.go` rather than a second copy of
it. A record failing that check is incomplete prior state, which is what `load` already
says about every other identity violation.

The enumeration is then entitled to the premise its own comment already states. The
comment on `residueRecords` argues that record order decides nothing "because no two
records of one slug name the same ref", and reasons about the checkpoint ref being written
"under that same run identity digested with an assignment ID" — which is true of refs the
checkpoint path writes and was not true of the string the enumeration actually read. That
gap is the mistaken premise this defect rests on. Once the value is validated the comment
describes an enforced invariant rather than a hoped-for one; say so in the comment, and do
not leave it asserting a guarantee no code makes.

## Acceptance

- [ ] [CR1] a record whose assignment `CheckpointRef` names a ref outside the checkpoint namespace is refused by the loader as incomplete prior state, asserted by writing it through the service's own save.
- [ ] [CR2] a record whose assignment `CheckpointRef` sits in the checkpoint namespace but does not match this run and assignment's own digest is refused the same way.
- [ ] [CR3] reclamation of a terminal run never plans or deletes a ref outside the namespaces it owns, asserted by driving plan and apply and observing the default branch survive.
- [ ] [CR4] an assignment that has not checkpointed yet — an empty `CheckpointRef` — still loads, assigns, and integrates, so the validation does not break the normal path.
- [ ] [CR5] the validation derives the expected ref from the same construction `checkpoint.go` writes, rather than restating the prefix or the digest recipe.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CR1 | drop the namespace half of the checkpoint-ref validation | the foreign-checkpoint-ref loader test | remove that clause from `validCore`, run `go test ./internal/specbuild -run Reclaim -count=1 -timeout 300s`, expect the loader-refusal assertion to fail |
| CR2 | validate only the namespace prefix and not the digest | the wrong-assignment-checkpoint-ref test | weaken the check to a prefix test, run `go test ./internal/specbuild -run Reclaim -count=1 -timeout 300s`, expect the digest-refusal assertion to fail |
| CR3 | restore the unvalidated `refObject` claim in `provisionalResidue` | the default-branch-survival test | revert the enumeration guard, run `go test ./internal/specbuild -run Reclaim -count=1 -timeout 300s`, expect the surviving-branch assertion to fail |
| CR4 | require a non-empty `CheckpointRef` on every assignment | the pre-checkpoint assignment test | make the empty case invalid, run `go test ./internal/specbuild -count=1 -timeout 600s`, expect the assign-then-load assertion to fail |
| CR5 | restate the checkpoint ref prefix and digest recipe inside `validCore` | the single-source construction test | inline the literal and the digest, run `go test ./internal/specbuild -count=1 -timeout 600s`, expect the single-source assertion to fail |
