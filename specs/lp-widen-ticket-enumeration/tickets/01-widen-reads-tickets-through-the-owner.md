# Read the signed-off tickets of a widening through the tracked-ticket owner

Blocked by: none
Writes: internal/spec/widen.go, internal/spec/widen_signedoff_test.go
Covers: none

## What to build

`bench spec widen` reads the signed-off scope from the tickets of the slug at the
default-branch tip. It selects those tickets with its own rule: a regular file below the
tickets folder with the ticket extension. The tracked-ticket enumeration
`tickets.EnumerateTracked` owns this rule for every other reader of the tickets of one
commit. The private rule has no duplicate-basename rule and no refusal of a path that
spec-TOON cannot represent, so the two readers disagree about the ticket set.

Make widen list the signed-off tickets through `tickets.EnumerateTracked`, and delete the
private rule. Widen still gives only the tracked regular files of the tickets folder to the
owner, so a symbolic link stays outside the ticket set.

## Acceptance

- [ ] Widen reads the signed-off ticket set through `tickets.EnumerateTracked`, and no
      second extension or prefix rule for that set stays in `internal/spec/widen.go`.
- [ ] A later copy of a duplicate ticket basename at the tip adds no path to the
      signed-off scope, so an entry that only that copy covers refuses as outside the scope.
- [ ] A tracked path below the tickets folder that spec-TOON cannot represent refuses the
      widening, and the ticket bytes stay unchanged.
- [ ] Every current widen test passes with no change.
