# Compute the tree-free oracle identity beside the oracle digest

Blocked by: none
Writes: internal/gate/subject.go, internal/gate/verdict.go, internal/gate/oracle_identity_test.go (new)
Covers: NG1, NG2, NG3, NG4

## What to build

`buildSubjectForTree` computes one more digest in the same pass. The identity digest receives every frame that the oracle digest receives, except the tree. It then receives the `BENCH_REQUIRE_CAPABILITIES` posture of the process. The `subject` type carries the result as one more field.

The oracle digest keeps every frame and its order, so its value does not change. The identity digest is the value that later tickets compare between a whole green and a landing. It is the oracle identity of map decision #8.

`internal/gate/subject.go` holds 446 of its 513 granted lines. Keep the growth small, or move the second writer into a new file.

## Acceptance

- [ ] Two subjects that differ only in their tree have one identity digest and two oracle digests.
- [ ] Setting `BENCH_REQUIRE_CAPABILITIES=1` changes the identity digest of one subject.
- [ ] Changing the bytes of a declared tool changes both the identity digest and the oracle digest.
- [ ] The diff of `buildSubjectForTree` adds frames only to the identity digest, and the oracle frames keep their order.
