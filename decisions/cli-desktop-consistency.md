# Consistent Bench behavior in CLI and desktop

Status: shaping

## Destination

Independent CLI and desktop chats provide consistent Bench behavior through one complete outcome for setup, checks, and recovery.

## Notes

`CONTEXT.md` defines harness interface, execution environment, and required capability.
The confirmed decisions below are current reviewer choices.
The remaining tickets require observed compatibility evidence.
The invoking authoring session records and validates this map.

Consult craft-domain, craft-grill, and craft-synthesis during shaping.
Research uses craft-research; compatibility probes use prototype.
A map-owned asset stays in `decisions/cli-desktop-consistency/assets/`.

The remaining evidence tickets follow this order:

- [Installed compatibility proof (#10)](cli-desktop-consistency/tickets/10.md).
- [Supported integration routes (#11)](cli-desktop-consistency/tickets/11.md).
- [Evidence-backed guarantees (#12)](cli-desktop-consistency/tickets/12.md).

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

## Not yet specified

## Spec-writer discretion

## Out of scope

- Ticket #1 excludes conversation migration and synchronization.
- Ticket #2 records the environments this qualification does not add.
- Ticket #6 preserves the separation of personal settings and credentials.
- Ticket #8 fixes the authority boundary for repairs.

## Sources

