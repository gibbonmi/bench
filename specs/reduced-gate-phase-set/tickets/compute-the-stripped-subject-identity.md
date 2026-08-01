# Compute the stripped subject identity

Blocked by: Declare the reduced gate scope

Ownership fence: `internal/gate/subject.go`, `internal/gate/stripped_subject_test.go`
Assumptions: `benchgit.TreeHash` builds a throwaway index, so the stripped variant is the same construction with the allowlisted paths dropped from that index

## What to build

A second subject identity computed over the tree with the allowlisted paths
excluded, alongside the existing whole-tree identity. The whole-tree identity keeps
governing included phases and the overall subject record; the stripped identity is
what decides whether an excludable phase's evidence still answers for the tree.

Everything the reuse decision rests on is this arithmetic, and it is dangerous in
one direction more than the other. An identity that moves too readily makes the
feature inert — no reduction ever happens. An identity that moves too rarely reuses
evidence for a tree that really changed, which is the failure that puts ungraded
code behind a green verdict. Build it so an over-broad strip is visible.

## Acceptance

- [ ] [R06] An edit confined to allowlisted paths leaves the stripped identity unchanged while the whole-tree identity moves.
- [ ] [R07] Any edit outside the allowlist moves the stripped identity, including an edit under a parent directory of a declared path.
