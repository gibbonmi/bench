# Process lifetime ownership (FT362; C02)

Status: ready

## Destination

Give every current production process-group caller one lifetime owner.
Preserve caller policy while fixing the identified shutdown defects.
This map prepares a spec and authorizes no implementation.

## Notes

The reviewer approved the owner, migration scope, and bounded failure posture on
2026-10-06. The reviewer delegated the remaining architecture choices.
The architecture review identifies this candidate as C02.

A process group is the OS signal target. Avoid: complete descendant tree.
A lifetime outcome records process completion and cleanup completion separately.
Avoid: exit code alone.

Use craft-seams, craft-domain, craft-spec, and prototype when this map resumes.
The source census is ticket 1. Refresh it when any named execution site changes.

## Decisions so far

- [Which production sites need the owner?](ft362-process-lifetime/tickets/1.md): the current census includes the bounds runner.
- [Who owns lifetime mechanics?](ft362-process-lifetime/tickets/2.md): subprocess owns mechanics; callers own policies.
- [Which behavior does the migration preserve?](ft362-process-lifetime/tickets/3.md): preserve normal completion and each caller's visible contract.
- [What happens when cleanup cannot finish?](ft362-process-lifetime/tickets/4.md): return bounded failure and retain evidence.
- [Which tests must survive?](ft362-process-lifetime/tickets/5.md): central mechanism tests complement real caller tests.
- [Can the owner support the stream shapes?](ft362-process-lifetime/tickets/6.md): the probe supports the stream strategy; durable resource protection remains an implementation obligation.

## Not yet specified

## Spec-writer discretion

- Internal names and file placement within the chosen owner.
- Ticket order within the approved scope and the resolved resource rules.

## Out of scope

- The separate shell exit-status repair.
- A new guarantee for descendants that escape their process group.
- A new native Windows execution contract.
- A change to gate phase order or evidence reuse authority.
- A migration of plain Capture callers that do not own process groups.

## Sources

- Path: `roadmap/FT362.md`
  Supports: the existing owner decision and profile contradiction.
  Drift: a change to the row's scope.
- Path: `internal/subprocess/subprocess.go`
  Supports: the existing capture owner and its explicit exclusion of group lifetime.
  Drift: a change to the package responsibility.
- Path: `specs/process-lifetime/decisions/ft362-process-lifetime/tickets/1.md`
  Supports: the current caller and policy census.
  Drift: a change to a cited execution site.
- Path: `internal/bounds/bounds.go`
  Supports: the additional production runner and the duration policy owner.
  Drift: a change to run or the resource policy.
- URL: https://raw.githubusercontent.com/golang/go/go1.25.14/src/os/exec/exec.go
  Supports: ticket 6's distinction between process completion and stream completion; read on 2026-10-06.
  Drift: a change to the declared Go toolchain or stream integration.


- Path: `specs/process-lifetime/decisions/ft362-process-lifetime/assets/prototype-evidence.md`
  Supports: the executed stream probe and its explicit production-evidence limits.
  Drift: a change to the Go toolchain or the selected stream strategy.
