# Learnings — usage journal

- 2026-08-21 — `bench commit --spec <slug>` exited
  `error: landed-but-checkout-incomplete: exit status 1` on a light-path
  tickets-only spec. The commit itself was correct: gate green, the three real
  files landed, and the retired ticket folder was excluded from history. Only
  the working-tree deletion did not complete, so `specs/<slug>/` stayed on disk
  as untracked residue and the session had to `rm -rf` it by hand and re-gate.
  What happened: a green landing reported a red-looking error, which reads as a
  failed commit until you inspect `git log`. Right behavior: the retirement's
  checkout step either completes or the command names which half succeeded and
  what the caller must remove. Proposed rule change: none for the agent — this
  is a `bench commit` defect, and the fix needs a reviewer call on whether the
  checkout is retried, or the residue is reported as a named next action.

- 2026-08-21 — The first home for the wrapper's pool-home test was
  `internal/conformance`, which red the gate's own meta-check
  (`conformance meta unregistered live-tree assertion`). That package admits a
  live-tree test only through its executable registry or a
  fixture-construction classification, and a policy assertion on the live
  wrapper is neither. Moving it to `internal/systemtest`, beside the adoption
  journey that already execs the real wrapper, was the seam. Why it was missed:
  the author picked the package by "which one already execs bash" instead of
  reading the package's own admission rule first. Proposed rule change:
  `craft-seams` gains one line — before placing a test in a package that owns a
  registry or inventory, read that registry's admission rule; a package that
  enumerates its own tests has a seam contract, not just a location.
