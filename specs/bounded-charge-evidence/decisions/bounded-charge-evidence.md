# Bounded charge evidence

Status: ready

## Destination

Decide how Bench delivers complete, pinned charge evidence within bounded responses.
Preserve required context for build authors and review axes.
Keep write-spec behavior unchanged.

## Notes

The reviewer supplied the defect and closed constraints on 2026-09-15.
The source baseline is `9deb0a7af31712427ff47d6fd0515e8458dbde0d`.
The [research report](bounded-charge-evidence/assets/research.md) owns the evidence and candidate comparison.
The phase uses craft-domain, craft-research, craft-grill, craft-line, and craft-synthesis.
Decision tickets own the selected contract; the research report retains evidence and rejected alternatives.

Ticket 15 records final scope approval and authorization for spec authoring and landing.
The production implementation and its acceptance proof remain the work of later phases.

## Decisions so far

- [What contract and consumer failure does the current charge expose?](bounded-charge-evidence/tickets/1.md): The research isolates the charge defect and its protected invariant.
- [What must charge completeness prove before action?](bounded-charge-evidence/tickets/2.md): Verified evidence can span bounded responses, with approval and the supplement separate.
- [Which delivery and lifecycle contract preserves that completeness?](bounded-charge-evidence/tickets/3.md): Explicit immutable preparation supports stateless bounded reads across worktrees and release.
- [Does write-spec alignment belong in this change?](bounded-charge-evidence/tickets/4.md): Build and review charges define the scope; write-spec behavior stays unchanged.
- [Does bounded retrieval preserve the exact prepared evidence bytes?](bounded-charge-evidence/tickets/5.md): The reviewer accepts the probe with explicit production validation requirements.
- [Should charge full retrieval remain supported?](bounded-charge-evidence/tickets/6.md): Remove charge --full without a compatibility route.
- [Does the artifact lifecycle preserve one immutable evidence identity?](bounded-charge-evidence/tickets/7.md): The reviewer accepts the lifecycle probe with explicit production validation requirements.
- [What encoded response budget must the charge enforce?](bounded-charge-evidence/tickets/8.md): Every response includes its complete envelope within 48,000 encoded stdout bytes.
- [What does consumer verification prove before action?](bounded-charge-evidence/tickets/9.md): Verified delivery and independent phase prerequisites govern action without universal edit blocking.
- [What artifact capacity and cleanup limits should Bench own?](bounded-charge-evidence/tickets/10.md): A 1 GiB default quota permits explicit overrides and cleanup without automatic eviction.
- [Must review captures reproduce across different environments?](bounded-charge-evidence/tickets/11.md): Freeze exact captures and record provenance without adding hermetic collectors.
- [How does the evidence identity bind bounded pages?](bounded-charge-evidence/tickets/12.md): A canonical TOON manifest binds exact source bytes through ordered page hashes.
- [Which commands own delivery, reuse, and current validation?](bounded-charge-evidence/tickets/13.md): Charge preparation and stateless evidence commands separate integrity, delivery, and current action checks.
- [Which physical container should store one evidence set?](bounded-charge-evidence/tickets/14.md): One documented packfile contains a binary header, the TOON manifest, and raw source bytes.

- [Does the complete shaped scope match the reviewer intent?](bounded-charge-evidence/tickets/15.md): The reviewer approves the scope and authorizes spec authoring, independent review, and landing.

## Not yet specified

## Spec-writer discretion

- The spec author owns engineering seams, acceptance coverage, fixtures, and gate attachment.
- Reversible implementation details can vary within the approved behavior and format contract.
- The author must present exact protocol fields and header details in the spec for review.
- No discretion permits a weaker source check, incomplete evidence, a larger response, or a different public command flow.

## Out of scope

- Production edits remain outside this shaping phase.
- Write-spec workflow changes remain outside this charge change.
- Reproducible collector execution across machines remains outside this charge change.
- Universal edit blocking across harnesses remains outside this charge change.
- Active debug-loop-guidance assignments remain untouched.

## Sources

- Path: `specs/bounded-charge-evidence/decisions/bounded-charge-evidence/assets/research.md`
  Supports: the core problem, consumer trace, alternatives, storage comparison, and accepted prototype limits.
  Drift: preparation, collector inputs, phase consumers, or storage owners change.
- Path: `internal/preflight/charge.go`
  Supports: the shared build/review projection and source identity behavior.
  Drift: the charge renderer or source policy changes.
- Path: `.agents/commands/bench-implement-spec.md`
  Supports: required evidence, approval, supplement, and action prerequisites.
  Drift: build entry, authorship, or evidence requirements change.
- Path: `.agents/commands/bench-review-implementation.md`
  Supports: review consumers and preserved shared evidence.
  Drift: review preparation, axis dispatch, or handoff changes.
- Path: `.agents/commands/bench-write-spec.md`
  Supports: the comparison with approved decision-source authoring.
  Drift: decision-source authorization, author-fork behavior, or resume behavior changes.
