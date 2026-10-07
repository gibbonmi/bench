# Planning facts and record operations (C10)

Status: ready

## Destination

Extend existing preflight, record, and reader owners for repeated planning work.
Keep FT293, FT375, FT318, and FT125 as distinct deliverable outcomes.

## Notes

The reviewer delegated architecture recommendations on 2026-10-06.
A projection reports authoritative facts without granting new write authority.
Avoid: a second workflow state owner.
Use craft-cli, craft-spec, and craft-seams for each outcome.

## Decisions so far

- [How should preflight handle fence and premise knowledge?](architecture-planning/tickets/1.md): derive proposals and verify explicit authorization.
- [When may the staleness pass be omitted?](architecture-planning/tickets/2.md): skip manual audit only for complete clean mechanical evidence.
- [Which record operations should the CLI absorb?](architecture-planning/tickets/3.md): add validated assignment writes and actionable refusals.
- [How should exact readers work?](architecture-planning/tickets/4.md): compose typed readers with bounded exact pages.
- [What is the recommended order within this candidate?](architecture-planning/tickets/5.md): deliver separate capabilities through their existing owners.

## Not yet specified

## Spec-writer discretion

- Internal fact types and gatherer placement within the existing owners.
- CLI flag spelling after each spec fixes its exact response and refusal grammar.

## Out of scope

- One new orchestration engine or a duplicate plan database.
- Automatic write-fence approval.
- A semantic-digest exemption for prose changes without a separate decision.
- The capability-blocked value that FT317 must decide.
- Summarized content in an exact artifact reader.

## Sources

- Path: `roadmap/FT293.md`
  Supports: the ownership closure and premise obligations.
  Drift: a change to its required closure set.
- Path: `roadmap/FT375.md`
  Supports: the mechanical staleness rows.
  Drift: a change to its drift contract.
- Path: `roadmap/FT318.md`
  Supports: the assignment and record operation gaps.
  Drift: a change to its remaining forms.
- Path: `roadmap/FT125.md`
  Supports: the exact-reader contracts.
  Drift: a change to its reader scope.
- Path: `internal/preflight/decision.go`
  Supports: immutable facts and the existing verdict owner.
  Drift: a change to Facts or Decide.
- Path: `internal/reviewrecord/plan.go`
  Supports: plan validation and source-bound digests.
  Drift: a change to ReadPlan.
- Path: `internal/reviewrecord/recordcmd/command.go`
  Supports: the record form registry.
  Drift: a change to forms.
- Path: `internal/responsebound/owner.go`
  Supports: the shared response budget and spill owner.
  Drift: a change to complete-output preservation.

