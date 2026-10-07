# Authored Markdown ownership (FT358; C01)

Status: ready

## Destination

Use the approved Markdown block-reader spec for the C01 outcome.
This map records the route and its closed decisions; it creates no competing spec.

## Notes

The reviewer delegated architecture recommendations on 2026-10-06.
The approved spec remains the behavior authority.
Use craft-domain, craft-spec, and craft-seams for its source refresh.

## Decisions so far

- [Which design owns this outcome?](architecture-markdown/tickets/1.md): reuse the approved FT358 spec.
- [Which failure proofs constrain the migration?](architecture-markdown/tickets/2.md): retain grammar-level failure proofs.

## Not yet specified

## Spec-writer discretion

- Source inventory updates that preserve the approved reader and grammar contracts.

## Out of scope

- A second block reader or a second spec for FT358.
- A move of sentence and paragraph rules out of internal/prose.
- A change to the approved implementation line or build schedule.

## Sources

- Path: `roadmap/FT358.md`
  Supports: the existing roadmap owner.
  Drift: a change to the row.
- Path: `specs/markdown-block-reader/spec.md`
  Supports: the approved owner, rule deltas, and acceptance obligations.
  Drift: an approved spec amendment.
- Path: `internal/spec/spec.go`
  Supports: the current status and fence readers.
  Drift: a change to the authored-status grammar.
- Path: `internal/coverage/coverage.go`
  Supports: the current coverage reader.
  Drift: a change to coverage extraction.

