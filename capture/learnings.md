# Learnings — usage journal

## 2026-08-15 — FT208 authoring scoped producer hardening to the consolidated reader [open]

**What happened.** FT208's authoring pass locked its first coverage map around
`internal/skillsindex`. Advisory reviews found that the producer-reader inventory and
composition fences omitted live conformance, package-surface, AXI-registry, and anchor
openers; they also replaced a nondeterministic temp-file signal sentinel and a story
that bundled SIGINT lifecycle with missing-Git diagnosis. The miss happened because I
scoped the work to the consolidated semantic reader and inferred consumers from
filenames instead of enumerating every live opener and registered execution path.

**Right behavior.** Producer hardening includes every executable consumer that can
reopen the producer, while process-interruption coverage controls the exact vulnerable
interval through a production-reached handshake and keeps unrelated command recovery
in its own story.

**Proposed rule change.** Before locking hostile-input rows, derive the complete
producer-reader set from executable registries and the call graph, then place one real
composition row through each producer class. Every interrupt row must use a
deterministic production-reached handshake rather than filesystem polling, a sleep, or
an appearance sentinel.
