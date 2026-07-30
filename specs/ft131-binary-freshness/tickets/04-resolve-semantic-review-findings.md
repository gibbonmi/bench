# Resolve semantic review findings

Blocked by: Publish and verify a content-sealed Bench binary; Refuse stale contract subjects; Refuse stale gate phase resolution

## What to build

Close the accepted semantic-review defects at the contract launcher, freshness owner,
gate entry, build-input inventory, and shipped documentation seams without reopening
FT131's defaulted decisions.

## Acceptance

- [x] Every runtime contract consumer of `benchPath(t)` crosses one freshness-enforcing launcher, with exhaustive consumer enumeration and paired false-green/false-red proof.
- [x] Empty executable stages and sealed empty executable artifacts are refused.
- [x] Executable and seal paths reached through any symlinked ancestor are refused before following or blocking.
- [x] Gate freshness uses current-source verification so an altered or always-green selected executable cannot attest for itself, while selected-binary freshness-check-before-gate-phases ordering and rebuild recovery remain observable.
- [x] The one owned refusal diagnostic renders a copy-paste-safe rebuild action for roots and outputs containing spaces or glob characters.
- [x] Auxiliary non-Go build inputs have one production inventory consumed by both the build script and freshness owner, and declared versus unrelated input changes affect freshness correctly.
- [x] The canonical plumbing-command inventory names `freshness-check`.
- [x] The project profile describes the freshness-enforcing gate entry.
- [x] The changelog records the user-visible stale-binary refusal and recovery behavior with a typed entry.
- [x] The freshness fixture package comment describes both fresh and stale subjects.
