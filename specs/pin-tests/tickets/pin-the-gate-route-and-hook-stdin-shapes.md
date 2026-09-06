# Pin the gate route and the hook stdin shapes

Blocked by: none
Writes: cmd/bench/gate_route_test.go (new), internal/adopt/link_hook_stdin_test.go (new), cmd/bench/main_test.go, internal/adopt/link_hook_test.go
Covers: none

## What to build

The review found two uncovered surfaces. The wrapper rejects every gate shape
that is not the bare run, `--fresh`, or a help spelling. No test holds that
grammar. The pre-push hook reads its stdin line by line. No test drives the
hook with an empty body, an unterminated last line, a zero local oid, or a
protected ref on the second line.

This ticket adds one table test for each surface. Each test gets its own file,
because the two host files sit at the structure budget. Each test reads the
shipped file. A reworded usage line, a lost usage arm, or a changed read loop turns
the lane red.

## Acceptance

- [ ] The wrapper answers `gate pin`, `gate --fresh unexpected`, and `gate --brief` with exit 2, empty stdout, and the usage line on stderr.
- [ ] The wrapper starts no gate run for a rejected gate shape.
- [ ] The hook exits 0 with empty stderr for empty stdin and for one topic line with no trailing newline.
- [ ] The hook prints the blocked line and exits 1 for a deletion push of the protected branch.
- [ ] The hook prints the blocked line and exits 1 when the second line names the protected branch.
