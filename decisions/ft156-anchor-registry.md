# FT156 anchor mechanism and registry

Status: ready

## Destination

The staged program for the conformance-anchor surface: what strengthens, what
gets tooling, what is stated honestly, and what is deferred. This map is the
decision source for the FT156 spec(s); the registry seam itself is spec work.

## #1: The staged program

Blocked by: none
Type: Grill

### Question

Four moves compete: restate coverage rows to the weaker claim substring
anchors actually prove plus a section-scoped `.bench/BENCH.md` red-capability
fixture; fix the `requireCollapsed` comment-strip asymmetry; replace
substring matching with a stronger mechanism; and a `bench anchors <path>`
pre-edit query. Which, and in what order?

### Answer

Resolved 2026-08-02, as a staged program. The stronger-than-substring
mechanism, the fourth competing move, is deferred — see Out of scope.

0. **The `.bench/BENCH.md` section-scoped red-capability fixture is pulled
   forward, before everything else** — it is independent and cheap, it lands
   beside or before the FT107/FT100 prose passes that edit that file first,
   and FT152's stories 1 and 15 rest on the untested combination. It is
   authored against the current mechanism and re-homes onto the registry
   when extraction lands.
1. **Registry extraction first among the mechanism work, carrying
   `bench anchors <path>`** — the anchor set becomes declarative data with
   one shared matcher, and the query command reads it so an editor sees
   which anchors pin a file before editing. Mechanism-neutral: a later
   stronger mechanism reuses the registry. The anchors live today as inline
   closures in one conformance test file (~180 needles: ~73 collapsed, ~92
   plain require, 11 forbid, 4 in-section), and the extraction collapses the
   plain and collapsed require paths into the same shared matcher.
2. **Comment-strip rides the extracted matcher** — with matching centralized,
   the require-direction matchers stripping HTML comments like `forbid`
   already does becomes a small change on the shared matcher plus a red
   demonstration, not a ~165-anchor sweep.
3. **Honest-rows reduced** — after comment-strip lands, only the rows still
   overclaiming against paraphrase evasion are restated to the weaker claim.

## Not yet specified

## Spec-writer discretion

- Registry representation (table, generated file, or parsed declaration),
  provided anchors stay single-sourced and the existing section-scoping
  facility survives. Placement is not discretion: for `bench anchors` to
  build at all, the registry data must live in a non-test package below the
  conformance import edge (the anchors sit in a `_test.go` file `cmd/bench`
  cannot import, and `internal/conformance/registry` documents that edge);
  the exact package within that constraint is spec work.
- Whether the program lands as one spec or several; the ~165-anchor sweep
  premise behind the roadmap's "own spec" sizing for comment-strip no longer
  holds once the matcher is shared.

## Out of scope

- A stronger-than-substring matching mechanism: deferred by the 2026-08-02
  decision session. Re-price after the FT107 and FT100 prose passes settle
  the surface, informed by whether FT158's falsification passes observe real
  paraphrase evasion. Not to be reopened by a spec.

## Sources
