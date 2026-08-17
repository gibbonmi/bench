# Repair: attribute a degraded roadmap/ directory, stop the wrapped-heading double diagnostic, refuse a live symlink row file

Blocked by: repair-fixture-harness-dedup-and-comments.md, repair-typed-diagnostic.md
Writes: internal/roadmap/tree.go, internal/roadmap/tree_validation.go, internal/roadmap/occurrences.go, internal/roadmap/tree_test.go, internal/roadmap/context_test.go, tests/canary/roadmap-detail-integrity/, specs/progressive-roadmap/spec.md

## What to build

Review findings Coverage C1, C2, C3, C5, C6, C8 (`reviews/progressive-roadmap.md`).

- **C1 (worst issue).** `LoadTree` already captures `tree.DirState`/`DirReason`
  (`internal/roadmap/tree.go:44`) but `ParseDocument`/`ValidateRoadmapTree`
  never read it — a `roadmap/` that is a regular file, FIFO, or unreadable
  directory currently either produces 67 misleading "missing detail owner"
  diagnostics (index present) or is silently clean (index absent). Add a
  diagnostic naming the actual cause when `DirState` is degraded (anything but
  parsed or, when the index has rows, empty), sourced at `roadmap/` itself.
  Add a canary fixture for this class in
  `tests/canary/roadmap-detail-integrity/` (mirror an existing fixture's
  shape) and a unit test in `tree_test.go` that plants a degraded `roadmap/`
  and asserts the new diagnostic — this also proves the existing
  `RoadmapDir + "/"` entry in `occurrenceSequenceTrusted`
  (`internal/roadmap/occurrences.go:94`) is actually reachable; today deleting
  that entry doesn't red anything.
- **C2.** Add a test planting an *empty* `roadmap/FT7.md` (not absent) and
  assert it reports the heading-mismatch diagnostic per the existing design
  comment at `tree.go:155-156`, not silence.
- **C3.** `tree.go:95-98`'s wrapped-heading branch `continue`s before
  `indexed[id] = true` (line 103), so a wrapped heading whose row file *is*
  present also fires a spurious `orphan detail file` diagnostic from
  `listingDiagnostics`. Fix so a wrapped heading emits exactly its one
  diagnostic — mark the ID as claimed by the index (without creating a row)
  before continuing, so the directory pass doesn't also call it orphaned.
  Update the `roadmap-wrapped-heading` canary fixture's `EXPECT` if it was
  silently tolerating two diagnostics (check whether its `EXPECT` is an exact
  match or substring match first).
- **C5.** `tree.go:46`'s `bounds.Classify` follows symlinks, so a row file
  that is a live symlink to an in-tree-but-non-row-file target is read as
  authoritative — contradicting `spec.md:237`'s "Won't handle … the
  classifier's wrong-type state reports it" claim, which is false for a live
  symlink to a regular file. Check how `bounds.Classify` is used elsewhere in
  this codebase for a similar authoritative-content read (e.g. `ROADMAP.md`
  itself, decision maps) to see whether symlink refusal already has a
  convention to match; if none exists, add symlink detection scoped to the
  roadmap row-file read path (not a global `bounds` package change — that
  package is shared broadly and a behavior change there needs its own
  review). A live symlink row file should report the same class of
  diagnostic C1 or the existing "unrecognized/wrong-type" class reports —
  pick whichever is the closer semantic fit and say why. Then correct
  `spec.md:237`'s Won't-handle claim to match what the code actually does
  (either "now refused" or a narrowed, accurate claim — do not leave prose
  contradicting code).
- **C6.** Add a test asserting the `--context` `sources` block's `roadmap/`
  row reports `absent` when the directory doesn't exist and `empty` when it
  exists with no files (today only `parsed` is pinned).
- **C8.** `TestOccurrenceLedgerMalformedAndLineEndings`
  (`internal/roadmap/context_test.go`) discards the diagnostics return
  (`valid, failures, _`) for its mixed-LF/CRLF case, silently accepting a
  `heading does not match` fault. Capture and assert on that return instead
  of discarding it — either the test's fixture needs to avoid producing that
  fault (if unintended) or the assertion needs to expect it (if it's the
  correct, if surprising, outcome — say which).

## Acceptance

- [ ] A degraded `roadmap/` directory produces one diagnostic naming the real cause, canaried, and the sequence-trust entry it feeds is reachable by a red-then-green test.
- [ ] An empty row file reports heading-mismatch, tested.
- [ ] A wrapped heading with its row file present emits exactly one diagnostic.
- [ ] A live symlink row file is refused (or its handling is explicitly and correctly documented if the reviewer's read of `bounds.Classify` shows refusing it is out of scope) — `spec.md:237` matches the code either way.
- [ ] `--context`'s `roadmap/` sources row is tested at `absent` and `empty`, not only `parsed`.
- [ ] The line-ending test's diagnostic return is asserted, not discarded.
- [ ] `go test ./...` and `bench gate` stay green.
