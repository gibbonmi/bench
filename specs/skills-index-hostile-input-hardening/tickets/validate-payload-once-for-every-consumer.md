# Validate payload once for every consumer

Blocked by: harden-every-skill-file-reader.md
Writes: consumer_payload.go, consumer_payload_test.go, internal/skillsindex/skillsindex.go, internal/skillsindex/skillsindex_test.go, internal/conformance/package_core_checks_test.go, internal/conformance/package_shipped_surface_test.go, internal/packagesurface/contract_documents.go, internal/packagesurface/contract_documents_test.go

## What to build

Expose one root-package byte parser that performs JSON decoding and the existing
canonical payload-row validation. Embedded `PayloadRows` sends its `go:embed` bytes
directly to that parser; every filesystem-backed skills-index, package-core,
package-shipped-surface, and packagesurface consumer first uses the published no-follow
bounded classifier and then sends the classified bytes to the same parser. Absence
remains optional only for skills-index; empty or any other present invalid state
refuses check/write and the package guards.

Keep audience, safe-source, and duplicate-source predicates only in the root package.
All filesystem consumers and HI14 land atomically because a thinner parser-first or
per-consumer migration strands the package gate: a payload FIFO still hangs an
unmigrated opener, while a copied/JSON-only parser still accepts a semantically invalid
row in another consumer.

## Acceptance

- [ ] `(covers HI8)` The canonical byte parser rejects invalid JSON, unknown audience,
  empty/unsafe source, and duplicate source; embedded bytes use it directly, absent
  filesystem allowlist stays optional, and `Check` no longer discards a present error.
- [ ] `(covers HI14)` `package-core-guard`, `package-shipped-surface`, and
  `packagesurface.ContractDocumentInputs` all complete with the same hostile/invalid
  payload dispositions, while refused skills-index writes preserve the reference.

