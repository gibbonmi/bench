# Leave the checkout path out of the oracle identity

Blocked by: none
Writes: internal/gate/subject.go, internal/gate/prospective.go, internal/gate/verdict.go, internal/gate/oracle_identity_test.go, internal/gate/baseline_schedule_identity_test.go (new), .bench/commitment.json
Covers: GR1, GR2, GR3, GR4, GR5, GR6, GR7, GR8, GR9, GR45

## What to build

`buildSubjectForTree` keeps the oracle digest exactly as it is today. The digest keeps every frame, the identity-root frame, and the identity-root spelling of each in-root path, in today's order.

The oracle identity changes in three ways:

- It leaves out the identity-root frame.
- It spells an executable path inside the runtime root relative to the runtime root. Otherwise it spells a path inside the identity root relative to the identity root, and every other path stays absolute.
- It frames the baseline phase schedule after the closure frames and before the posture frame.

Today one byte stream feeds both digests. So the collector writes each path frame to each digest in that digest's spelling, and it writes every other frame to both digests.

The schedule is the answer of `readManifest` at `manifestPath(identityRoot)`, the loader's own read. Present bytes enter as bytes, and a live link enters as the bytes of its target. An absent manifest enters as one fixed value.
An error answer, such as a dangling link, a FIFO, or a directory, enters as one fixed defect value with no path and no message. Do not write a second reader of the file.

The `subject` comment in `internal/gate/verdict.go` states the new identity frames. The `subject.Identity` field keeps its name, because the reader table and ticket 02 read it.

`TestOracleIdentityAndOracleShareTheClosure` drops its identity-root case, because that case now moves only the oracle digest. The new `TestOracleIdentityLeavesOutTheCheckoutPath` grades the identity root. The reviewer approved this change on 2026-10-26.

The reader-table test calls the fixture of `internal/gate/reader_plan_test.go`: `newLandingFixture`, `storeTable`, and `selects`. That file holds 400 lines, so it takes no new line. The sibling identity root is a linked worktree of the fixture repository that `runGitIn` creates.
The architecture census refuses a direct `exec.Command`, or a `Run` or `Output` call with a `worktree` literal, in a gate test.

`internal/gate/subject.go` holds 456 of its 513 granted lines. If the schedule frame does not fit, put its helper in `internal/gate/prospective.go` beside `baselineRunnerIdentity`.
Keep `oracle_identity_test.go` under 400 lines. The schedule rows GR3 to GR6 go to the new `baseline_schedule_identity_test.go`.

`.bench/commitment.json` is on this line for the binding commit of the spec-stage close. Do not edit it.

## Acceptance

- [ ] With an in-root launcher and tool, two identity roots with equal content, one of them with a space in its path, give one oracle identity (GR1) and two oracle digests (GR2).
- [ ] An absent baseline manifest and a present empty one give two identities (GR3).
- [ ] Two present baseline manifests with different bytes give two identities (GR4).
- [ ] A live symbolic link at the baseline manifest gives two identities for two targets with different bytes (GR5).
- [ ] A dangling link and a FIFO at the baseline manifest give one identity, and each subject build returns before the test deadline (GR6).
- [ ] With a `PATH` tool inside the checkout, the ordinary subject of the checkout and a prospective subject of the same tree with the checkout as identity root give one identity (GR45).
- [ ] The reader plan selects the reader set for a table whose identity a linked worktree of the fixture computed (GR9).
- [ ] The posture still changes the identity (GR7), and each other input of `TestOracleIdentityAndOracleShareTheClosure` still changes both digests (GR8).
