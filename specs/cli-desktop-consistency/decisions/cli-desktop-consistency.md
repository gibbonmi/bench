# Consistent Bench behavior in CLI and desktop

Status: ready

## Destination

Independent CLI and desktop chats provide consistent Bench behavior through one complete outcome for setup, checks, and recovery.

## Notes

`CONTEXT.md` defines harness interface, execution environment, and required capability.
The confirmed decisions below are current reviewer choices.
The evidence reports record a failed compatibility probe and the remaining qualification gaps.
The invoking authoring session records and validates this map.

Consult craft-domain, craft-grill, and craft-synthesis during shaping.
Research uses craft-research; compatibility probes use prototype.
A map-owned asset stays in `specs/cli-desktop-consistency/decisions/cli-desktop-consistency/assets/`.

The reviewer confirms the complete shaped outcome.
All decision tickets are resolved.
The evidence gaps remain required qualification work for the specification.

## Decisions so far

- [Scope and continuity (#1)](cli-desktop-consistency/tickets/1.md): independent chats.
- [Supported environments (#2)](cli-desktop-consistency/tickets/2.md): WSL qualification.
- [Workflow parity (#3)](cli-desktop-consistency/tickets/3.md): complete Bench workflow.
- [Setup and recovery experience (#4)](cli-desktop-consistency/tickets/4.md): automatic checks.
- [Concurrent use (#5)](cli-desktop-consistency/tickets/5.md): isolated writers.
- [Configuration consistency (#6)](cli-desktop-consistency/tickets/6.md): shared Bench contract.
- [Capability failure (#7)](cli-desktop-consistency/tickets/7.md): block affected work.
- [Repair authority (#8)](cli-desktop-consistency/tickets/8.md): supported repair routes.
- [Delivery scope (#9)](cli-desktop-consistency/tickets/9.md): one compatibility outcome.

- [Installed compatibility proof (#10)](cli-desktop-consistency/tickets/10.md): negative command result.
- [Supported integration routes (#11)](cli-desktop-consistency/tickets/11.md): supported routes and explicit limits.

- [Evidence-backed guarantees (#12)](cli-desktop-consistency/tickets/12.md): capability evidence and verified recovery.

## Not yet specified

## Spec-writer discretion

## Out of scope

- Ticket #1 excludes conversation migration and synchronization.
- Ticket #2 records the environments this qualification does not add.
- Ticket #6 preserves the separation of personal settings and credentials.
- Ticket #8 fixes the authority boundary for repairs.

## Sources

- Path: `specs/cli-desktop-consistency/decisions/cli-desktop-consistency/assets/installed-compatibility-probe.md`
  Supports: Tickets #10 through #12; actual command results and unqualified workflow classes.
  Drift: Runtime, configuration, or workspace changes require affected probes again.
- Path: `specs/cli-desktop-consistency/decisions/cli-desktop-consistency/assets/supported-integration-routes.md`
  Supports: Tickets #11 and #12; supported integration routes and recovery limits.
  Drift: Bench integration, upstream documentation, or installed environment changes require a source check.
