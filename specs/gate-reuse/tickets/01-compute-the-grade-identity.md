# Compute the grade identity from the verdict inputs only

Blocked by: none
Writes: internal/gate/subject.go, internal/gate/prospective.go, internal/gate/gate.go, internal/gate/run_transaction.go, internal/gate/oracle_identity_test.go, internal/gate/grade_identity_test.go (new), internal/gate/baseline_schedule_identity_test.go (new), .bench/commitment.json
Covers: GR1, GR2, GR3, GR4, GR5, GR6, GR7, GR8, GR9, GR45, GR54, GR55

## What to build

`buildSubjectForTree` computes the tree-free identity digest from the verdict inputs only. That digest is the grade identity.
It frames the gate resolution and its launcher closure, `PATH`, the gate-input manifest, the declared environment values, the declared path bytes, and the declared tools.
It also frames the bytes of the baseline phase manifest, the run binary member, the running executable digest, and the capability posture last.
It frames no identity root path, no policy version, and no baseline runner identity.

An executable inside the runtime root enters by its path relative to the runtime root. Otherwise, an executable inside the identity root enters by its path relative to the identity root.
Every other executable keeps its absolute path. So an ordinary subject and a prospective subject of one tree from one checkout keep one identity.

The schedule member is the answer of `readManifest` at the identity root's `.bench/phases.json`, which is the baseline of a prospective run. Present bytes enter as bytes, and a live link enters as its target's bytes.
An absent manifest enters as one value. A defect answer, such as a dangling link or a FIFO, enters as one stable value with no path. `readManifest` keeps its signature.

The run binary member is `freshness.Digest` over the run binary's source root, computed before any build. A source root that declares no build inputs frames one stable value.
Move the source-root rule of `runBinarySource` in `gate.go` into one helper that takes the resolution kind, the runtime root, and the storage root.
The selection at `run_transaction.go:159` calls that helper. The identity calls it where the prospective script's presence is known, in `hashProspectivePreparation` or after it, with the identity root as the storage root.

The running executable digest comes from `brokermanifest.Digest` over the running executable. A package-level seam supplies that executable and computes its digest once for each process, so a test can vary it.

The oracle digest keeps each current frame in its current order, the baseline runner identity included, and it gains no frame. So `baselineRunnerIdentity` stays as an oracle frame only, and the evidence store keeps its current name until ticket 02.
`gate.go` holds 373 lines and `run_transaction.go` holds 370 lines, and each stays at or under its 400-line limit.

## Acceptance

- [ ] Two identity roots with equal content and a space in one path give one grade identity and two oracle digests.
- [ ] An absent manifest, an empty manifest, and two manifests with different bytes give distinct identities, and a link gives its target's identity.
- [ ] A dangling link and a FIFO at the manifest give one identity, and the subject build returns.
- [ ] Each of the posture, the tool bytes, `PATH`, an environment value, the declared path bytes, the run binary member, and the running executable digest changes the identity.
- [ ] A prospective subject of one tree gives one identity from a worktree whose build inputs equal the tree and from a primary checkout whose build inputs differ.
- [ ] The reader plan admits a table that a sibling identity root stored over equal content.
