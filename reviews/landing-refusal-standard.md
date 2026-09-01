# Review: landing-refusal-standard

Range: `f470a3c53e9084393065c07c349fc4bdaf536c16..f7456ba81e5d6a42a90f16f72385b9125edcb122`
Raw findings: 15 (Standards 5, Spec 5, Coverage 5). Repair targets after
collapse: 5. Reviewer-decision parks: 3.

## Standards

Count: 5. Worst: the reference doc's paragraph makes two universal claims the
code does not keep.

- **S1 — no-op.** "Required argument" versus empty-string call sites. The spec
  refutes the concern: "An empty route string still compiles, so the registry
  walk grades the value the face prints" (spec, Implementation decisions). The
  registry walk is the enforcement, and it runs.
- **S2 — auto-fix (repair R4).** `.bench/BENCH-reference.md:187` says each
  refusal "prints its failed paths"; the dirty-source face prints no paths by
  design (LRS13). Reword the second sentence; keep the anchored needle line
  verbatim; sync the canary fixture copy.
- **S3 — ask-user (parked).** `faceResumeDestinationResidue` routes six
  distinct destination states through one repair sentence
  (`land_refusal.go:930`, `land_resume.go:77`). Route granularity is a design
  call.
- **S4 — no-op.** The unregistered-face fallback in
  `landingRefusalFaceByName` is a defensive default; kept.
- **S5 — no-op.** LRS tags in test comments match the tree's idiom.

## Spec

Count: 5. Worst: the conflict face lives outside the registry the spec names
as the authoritative inventory.

- **P1 — ask-user (parked).** Rows LRS1/LRS2 name
  `TestIdentityComponentRegistryHasAProducingFixture`; the delivered walk is
  `TestLandingRefusalRegistryHasAProducingFixture`. The behavior exists; the
  row's test name is wrong, and the build may not edit acceptance rows.
- **P2 — auto-fix (repair R4).** `landingConflictRefusal` composes the
  conflict route outside `landingRefusalFaces` (spec: "One constructor turns
  an entry into a refusalError... The registry is the authoritative
  inventory"). Register the conflict face; the registry now supports a face
  with a dynamic detail, so the original blocker is gone.
- **P3 — no-op with note.** `land_specless_test.go` is absent from the fence
  list, but `bench preflight review` reports `paths-authorized` green over
  this exact range, so the oracle accepts the path. The fence row stays a
  spec-hygiene note for the reviewer.
- **P4 — no-op (flagged).** LRS17's "observed tip" is the tip the tree holds
  (`wanted=`), the only value that makes the re-run work. Non-behavioral
  wording contradiction, flagged for reviewer veto.
- **P5 — auto-fix (repair R5).** LRS3's re-run-tail assertion samples two
  faces. Widen the loop to every first-run preflight face in
  `landingRefusalFaces`.

## Coverage

Count: 5. Worst: a reachable fence-face branch prints a refusal with no
`next=` field.

- **C1 — auto-fix (repair R1).** The non-typed branch in
  `landingSourceRange` (`land_identity.go:198-207`) wraps the sentence, so
  `landingFaceByDetail` misses it and no route attaches.
  `TestLandCommandAbsentSpecFolderKeepsTheUnreadableRefusal` pins that record
  and asserts nothing about `next=`. Route the branch; add the assertion.
- **C2 — auto-fix (repair R2).** `resumeRerun`'s output is never observed; a
  mutation that returns `""` stays green. Add a resume fixture that drives a
  source-side face through `resumeAssignment`.
- **C3 — auto-fix (repair R3).** `residueRemovalCommand`'s two placeholder
  branches (`!lineSafe`, empty list) execute unasserted. Add the assertions.
- **C4 — no-op.** `sourceMergePending`'s `rev-parse` error branch is
  uncovered; both error branches return the same fallback.
- **C5 — no-op.** Defensive branches unreachable by construction; the
  registry/fixture bijection is pinned both ways.

## Repair tickets

- R1 (C1): route the fence face's wrapped-error branch; assert `next=` in
  `land_tickets_only_test.go`. Writes: `internal/worktree/land_identity.go`,
  `internal/worktree/land_refusal.go`,
  `internal/worktree/land_tickets_only_test.go`.
- R2 (C2): resume fixture that observes `resumeRerun`. Writes:
  `internal/worktree/identity_component_test.go`,
  `internal/worktree/land_resume_refusal_test.go`.
- R3 (C3): assert both placeholder branches of `residueRemovalCommand`.
  Writes: `internal/worktree/land_release_refusal_test.go`.
- R4 (S2+P2): register the conflict face in `landingRefusalFaces` with a
  producing fixture; reword the doc's second sentence; sync the canary
  fixture. Writes: `internal/worktree/land_refusal.go`,
  `internal/worktree/land.go`,
  `internal/worktree/identity_component_test.go`,
  `.bench/BENCH-reference.md`,
  `tests/canary/workflow-guidance-anchors/reference-refusal-route-shape/`.
- R5 (P5): widen LRS3's re-run-tail loop to every first-run preflight face.
  Writes: `internal/worktree/land_surface_test.go`.
