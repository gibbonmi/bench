# State the tickets-only rule once

Blocked by: none
Writes: internal/landing/close.go, internal/landing/close_test.go, internal/worktree/land_identity.go, internal/worktree/land_resume.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: SR30, SR31, SR32, SR33, SR34, SR35

## What to build

The line is opus / medium. State the tickets-only predicate once in
`internal/landing` over a tree reader. The reader answers two questions: is
the spec folder a directory, and is the spec file absent. The name check
applies to both readers, so a slug with a separator or `..` is never
tickets-only.

The working-tree reader keeps the current stat calls. The git-object reader
keeps the current `cat-file -e` and `show` calls over the source commit. The
first run calls the one predicate through the working-tree reader. The resume
calls the same predicate through the git-object reader, and
`internal/worktree/land_resume.go` keeps no second `cat-file -e` and `show`
pair over the spec folder.

`Owner.Land` and `Owner.LandReviewed` keep their accepted kinds. A git reader
whose `show` fails for any reason answers spec-absent, as today.

The exit proof for this ticket is the pre-existing suite, green with its test
logic unchanged. A mechanical rename is permitted. A needed assertion change
stops the ticket and reports.

## Acceptance

- [ ] The predicate answers true for a folder with no spec file through the working-tree reader.
- [ ] The predicate answers false for a malformed name, an absent folder, and a present spec file through the working-tree reader.
- [ ] The predicate answers the same four cases through the git-object reader over a fixture commit.
- [ ] A `--spec` that names a tickets-only folder closes it, an already-removed folder still lands, and an absent folder keeps the unreadable refusal.
- [ ] A resume of an interrupted close authenticates the folder's absence and releases.
- [ ] A resume completes a spec-backed landing and accepts the slug and the path spelling.
- [ ] `bench consumers landing.TicketsOnlyFolder` lists the first-run site and the folders sweep.
- [ ] Self-probe: omit the name check from the git-object reader, and report the observed red.
