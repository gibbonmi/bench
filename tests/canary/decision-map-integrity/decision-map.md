# Multi-harness line binding

Status: ready

## Destination

Make the tier-to-model binding symmetric across harnesses, so each harness
reads, reports, and recommends its own model family. Today `BENCH_TIER_*` holds
Codex ids as the privileged canonical, and `BENCH_ALIAS_*` is a Claude-only
translation. So every advisory surface leads with a Codex id, even inside a
Claude Code session. After this change no family is canonical: the tier
(top/mid/cheap) is the only abstract identity, and each reader resolves it to
the harness that asks.

Scope confirmed with reviewer: the bite is presentation only. No run is
mis-launching today.

## Notes

The migration is a hard cut with no dual-read shim. `bench doctor` finds the
retired keys and reports the exact rewrite.

The recommendation rule lives in the advisory layer, because enforcement stays
permissive. The abstract rule names no family; it resolves per harness.

## Decisions so far

- [Data model](multi-harness-line-binding/tickets/1.md): The tier is the only canonical identity, and each family is a peer.
- [lines.env schema](multi-harness-line-binding/tickets/2.md): One `BENCH_<HARNESS>_<TIER>` key per cell, over a fixed harness set.
- [Harness discovery](multi-harness-line-binding/tickets/3.md): An explicit `--harness` flag replaces the output-shape flags.
- [Enforcement scope](multi-harness-line-binding/tickets/4.md): Any bound tier value passes, and native-family steering stays advisory.
- [Docs and conformance](multi-harness-line-binding/tickets/5.md): The profile prints the full matrix, cross-checked against lines.env.

## Not yet specified

## Spec-writer discretion

- Exact deny-message wording and the verbose-flag name for the full-matrix report.
- Whether `bench models` discovery output should also become harness-scoped.

## Sources

## Out of scope
