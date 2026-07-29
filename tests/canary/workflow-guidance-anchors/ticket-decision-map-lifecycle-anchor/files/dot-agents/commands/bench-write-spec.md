Top-level `decisions/` holds pre-spec working maps. A new spec is compiled from a
named top-level `decisions/<topic>.md`; run `bench maps` and confirm it shows no row
for the map. Settled provenance there is deliberately outside the top-level `bench maps`
query.

The source map and any map-owned assets from top-level `decisions/` move
into `specs/<slug>/decisions/`, preserving their useful relative layout, and update every
reference to the moved paths. A re-run reads the already-compiled spec-local map; it never
recreates a top-level copy. Whole-folder retirement removes the compiled maps and map-owned
assets under `specs/<slug>/decisions/` with the spec and its tickets; there is no separate
decision-map cleanup step.
