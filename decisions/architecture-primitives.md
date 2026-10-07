# Strict reads, durable files, and shell arguments (FT354; C08)

Status: ready

## Destination

Use one owner for each primitive while keeping their contracts distinct.
Separate the strict-JSON repair from file durability and printed-command changes.

## Notes

The reviewer delegated architecture recommendations on 2026-10-06.
Atomic visibility means readers observe old or new bytes.
Durability also requires the declared synchronization steps.
Avoid: atomic as a synonym for durable.
Use craft-spec, craft-seams, and craft-cli for the separate specs.

## Decisions so far

- [Which strict reader should each caller use?](architecture-primitives/tickets/1.md): use jsonfile with explicit framing.
- [What should shared file replacement guarantee?](architecture-primitives/tickets/2.md): use one durable file-replacement owner.
- [Which shell argument spelling should survive?](architecture-primitives/tickets/3.md): use sanitize's always-quoted argument contract.
- [Should these primitives share one implementation project?](architecture-primitives/tickets/4.md): split by observable contract and urgency.

## Not yet specified

## Spec-writer discretion

- Private helper names within the chosen primitive owners.
- Migration batches with complete caller inventories and unchanged sentinel categories.

## Out of scope

- One large primitive refactor that blocks the strict-JSON defect repair.
- Removal of caller-owned permission, schema, or publication policy.
- Replacement of directory exchange with file replacement.
- A new promise about unsupported filesystems or remote storage.

## Sources

- Path: `roadmap/FT354.md`
  Supports: the existing primitive ownership questions.
  Drift: a change to the row.
- Path: `internal/jsonfile/decode.go`
  Supports: complete-document and persisted-record readers.
  Drift: a change to framing or strictness.
- Path: `internal/releasepreflight/vulnerability.go`
  Supports: the exception parser's incomplete suffix check.
  Drift: a change to ValidateVulnerabilityPolicy.
- Path: `internal/gate/verdict.go`
  Supports: a file-and-directory synchronization example.
  Drift: a change to durableReplace.
- Path: `internal/reviewrecord/write.go`
  Supports: a replace operation without directory synchronization.
  Drift: a change to replace.
- Path: `internal/sanitize/sanitize.go`
  Supports: the existing always-quoted shell argument function.
  Drift: a change to ShellQuote.
- Path: `internal/axi/action.go`
  Supports: the differing bare-token rule.
  Drift: a change to ShellQuote.

