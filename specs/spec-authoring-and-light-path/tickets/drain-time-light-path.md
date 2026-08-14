# Let a what-next drain implement light-path items on the spot

Blocked by: none
Writes: .agents/commands/bench-what-next.md, .bench/BENCH.md, internal/anchors/registry_data.go, internal/canary/inventory_test.go, tests/canary/workflow-guidance-anchors/

## What to build

The `/bench-what-next` drain verdict set gains "implement now" for an item
meeting the light-path observables (the observables already exist in
BENCH.md, so no sibling gates this — BENCH.md, registry, and fixture paths
are shared with siblings and land serially). When the reviewer chooses it,
the session writes the one ticket file, spawns a write-delegate charged with
it (`craft-delegate` isolation, `craft-line` routing), and verifies the
returned diff against the ticket's acceptance rows and the gate. The
implemented item lands as its own commit on green, written into
bench-what-next.md as the second named exception to the drain's
one-batch-commit rule beside the existing two-spec-retire exception — a second
Require row pins the commit exception and a **new** fixture backs it (the
existing `one uncommitted batch diff` needle stays as-is; no fixture backs
that rule today), and the file's "Two constraints shape that single commit"
count word updates to match the added third. Added BENCH.md prose must avoid
the file's whole-file Forbid on the bare substring `thorough`.
Items needing a reviewer decision, a new seam, or spec-level design still
graduate to `ROADMAP.md`. The `.bench/BENCH.md` edit is one clause: the
Capture paragraph's graduate-only sentence gains "or close by implementation
during that same reviewed drain" — graduation remains drain-only, an
implemented item closes instead of graduating, and BENCH.md stays within its
180-line budget. One long-needle Require anchor pins the
route sentence (ticket file, write-delegate, main-session verification) with
a new fixture; a second Require pins the commit exception.
## Acceptance

- [ ] bench-what-next.md offers implement-now for light-path-eligible drained
      items and its route sentence names ticket file, write-delegate, and
      main-session verification, pinned by one long needle whose fixture bites
      both halves (covers WF14)
- [ ] the implemented item's own-commit-on-green lands as the second named
      exception to the one-batch-commit rule, pinned by its own second
      Require row with a new fixture biting both halves (covers WF17)
- [ ] `.bench/BENCH.md` carries the close-by-implementation clause and remains
      within its 180-line budget (covers WF13)
- [ ] the ticket's fixture additions update the canary binding count in this
      ticket's green commit (covers WF22)
