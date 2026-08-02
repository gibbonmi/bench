# FT183 per-component-scoping review residuals

Status: shaping

## Destination

Dispositions for the two faces left by the shipped per-component-gate-scoping
review: the unreachable whole-changeset reduced fallback in `internal/gate`'s
component scoping, and the unbound derivation-source check in
`internal/conformance` that grades `internal/gate`'s registry.

## #1: Retire the whole-changeset reduced fallback

Blocked by: none
Type: Grill

### Question

After per-component scoping, `reducedInheritance` answers only for a kit root
with no `go.mod` — a shape production never has (review S1; the comment-truth
half already landed at `7936b70`). Remove the path, or keep it for linked-repo
futures?

### Answer

Resolved 2026-08-02: remove it. The removal simplifies `internal/gate` and
retires `reduced_run_test`'s synthetic fixture; if a linked-repo shape ever
genuinely needs a whole-changeset reduced path, it is re-added against that
real shape rather than carried as an unreachable branch.

## #2: Observe which derivation a registry row resolves through

Blocked by: none
Type: Research

### Question

The derivation-source check proves an entry is derivation-sourced but not
that it binds to its NAMED derivation, so a wrong-but-derived resolver swap
passes (review Sp1, disclosed in `derivation_source_test.go`'s header). What
observation design would expose which function a registry row actually
resolves through, and at what cost?

### Answer

— (open: produce a short cited design summary of candidate observation
mechanisms over `internal/gate`'s registry, with the cost and blast radius of
each; parked as research by the 2026-08-02 decision session — no ruling until
the summary exists)

## Not yet specified

- The reduced-verdict record class's disposition under #1: `reducedInheritance`
  is the sole writer of `Reduced` verdicts, so removal orphans the persisted
  record fields and their readers (`inherits`, `validateInheritance`, the
  status and prep-release consumers). Whether the record class retires with
  the path or keeps a reader for legacy on-disk records is open for the spec
  to surface.

## Spec-writer discretion

## Out of scope

- Reintroducing any whole-changeset reduced path for the kit root; #1 closed
  that direction.

## Sources
