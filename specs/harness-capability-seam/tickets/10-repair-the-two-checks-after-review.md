# Repair the two checks after review

Blocked by: 07-add-the-harness-record-conformance-check.md, 08-add-the-entry-point-parity-conformance-check.md
Writes: internal/conformance/harness_record_test.go, internal/conformance/entry_point_parity_test.go, internal/guards/guards.go, internal/guards/guards_test.go

## What to build

The review found four gaps in the two oracle checks. Each repair below is
one accepted predicate from `reviews/harness-capability-seam.md`.

Target A. `harness-record` refuses and names a live symlink at an adapter
path. The spec's edge inventory says a dangling symlink is absent and a live
symlink is refused. The check classifies the adapter path with the no-follow
classifier before it stats. A live link then yields a diagnostic that names
the path. A dangling link still counts as absent.

Target F. Two refusal branches gain tests. A config of `{"hooks":{` yields
the `is not valid JSON` diagnostic. A live symlink at `.codex/hooks.json`
yields the refusal diagnostic and no other diagnostic for that row.

Target L. The check composes the hook-wiring rule from `internal/guards`
instead of a second substring test. `guards` exports one predicate for
"this config wires this script token", and both readers call it.

Target Q. `entry-point-parity` reds an absent or unreadable static entry.
A static row whose file the no-follow read refuses yields a diagnostic that
names the path and the registry command.

## Acceptance

- [ ] A live symlink at `.bench/adapters/codex` yields a diagnostic naming that path, and a dangling one yields the absent-adapter diagnostic.
- [ ] A `.codex/hooks.json` of `{"hooks":{` yields the invalid-JSON diagnostic.
- [ ] A live symlink at `.codex/hooks.json` yields the refusal diagnostic and no other diagnostic for that row.
- [ ] `internal/conformance/harness_record_test.go` holds no `strings.Contains` wiring rule; it calls the `guards` predicate.
- [ ] A root without `scripts/release-preflight.sh` yields a parity diagnostic naming that path.
- [ ] Both checks report no diagnostic over the live root.
