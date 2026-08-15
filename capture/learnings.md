# Learnings — usage journal

## 2026-08-15 — a decision map compiled under a spec is deleted with that spec [open]

**What happened.** The deepening map (`/bench-shape-idea` over the 2026-08-15 survey)
was compiled into `specs/gate-run-transaction/decisions/`, and its ticket #5 told
later specs and light-path tickets to cite that path. `bench spec retire
gate-run-transaction` deleted the spec directory whole, so the only in-tree copy of a
map still owning six unbuilt lanes was gone; the next session recovered it from git
history and re-homed it at `decisions/deepening-2026-08.md`.

**Right behavior.** A map that outlives one spec lives in `decisions/`, never under a
spec's own directory; only a map consumed entirely by that one spec may travel with it.

**Proposed rule change.** `/bench-shape-idea` exit: a map whose tickets route to more
than one spec or light-path lane compiles to `decisions/<map>.md`; `/bench-write-spec`
cites, and does not copy, that path.
