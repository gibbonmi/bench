# Shared guard grammar (FT366; C07)

Status: ready

## Destination

Share command projection and envelope decoding while each guard retains its policy.
Fix malformed-command refusal in the degraded hook before the broader migration.

## Notes

The reviewer delegated architecture recommendations on 2026-10-06.
The guard is an honest-mistake check, not an evasion-resistant shell interpreter.
A projection identifies executable words without executing the command.
Use craft-cli, craft-spec, craft-seams, and craft-gate for the spec.

## Decisions so far

- [Which grammar belongs to the shared owner?](architecture-guards/tickets/1.md): share projections and keep guard decisions separate.
- [How much of the degraded shell hook should remain?](architecture-guards/tickets/2.md): retain the coarse recovery rim and refuse malformed input.
- [Which outcomes should be separate checkpoints?](architecture-guards/tickets/3.md): repair first, then grammar migration and scanner simplification.

## Not yet specified

## Spec-writer discretion

- Private projection types and file placement within the existing grammar owner.
- Envelope reader placement that creates no import cycle.

## Out of scope

- A complete shell interpreter or deeper automatic wrapper recursion.
- A shared authority policy for all guards.
- An automatic relaxation when the core cannot parse a command.
- A broad removal of the degraded hook's recovery behavior.

## Sources

- Path: `roadmap/FT366.md`
  Supports: the current scope and degraded-rim occurrence.
  Drift: a change to its guard requirements.
- Path: `internal/gitguard/scan.go`
  Supports: wrapper and worktree-child projection.
  Drift: a change to supported command recognition.
- Path: `internal/benchguard/benchguard.go`
  Supports: the second wrapper and envelope reader.
  Drift: a change to guard input handling.
- Path: `.bench/hooks/block-dangerous-git.sh`
  Supports: the independent degraded hook and malformed fallback.
  Drift: a change to core recovery behavior.
- Path: `internal/guards/guards.go`
  Supports: the scanner's cancellation waits.
  Drift: a change to Scan.

