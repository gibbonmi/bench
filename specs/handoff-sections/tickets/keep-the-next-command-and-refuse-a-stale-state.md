# Keep the Next command and refuse a stale State

Blocked by: resolve-the-callers-section-in-bench-handoff.md
Writes: internal/handoff/handoff.go, internal/handoff/render.go, internal/handoff/render_test.go, internal/handoff/state_scan.go (new), internal/handoff/state_scan_test.go (new)
Covers: HS11, HS12, HS13, HS14, HS15

## What to build

Verify the premise first: `Command` calls `applyRoute` unconditionally when
`--next` is absent, and nothing scans State. Then keep a non-empty
`Next command:` value byte for byte, and call `status.RouteFor` only when the
value is empty, whitespace, or empty backticks. Add a State scan. Find
backticked runs of seven to forty hex characters, and keep those that
`git cat-file -e <token>^{commit}` accepts. When `merge-base --is-ancestor
<token> <tip>` is false, refuse with exit 1 and print the line. Use the
`git.OK` form that internal/worktree/clean.go already uses.

## Acceptance

- [ ] Without `--next`, a `Next command:` line with flags survives byte for byte on a board whose leading signal differs.
- [ ] An empty, a whitespace-only, and an empty-backtick Next command each receive the board's leading signal.
- [ ] A State that names a real commit off the tip's ancestry exits 1 with that line printed.
- [ ] A State that holds `facade` and `decade` in backticks exits 0.
- [ ] A State that names `HEAD^{tree}` exits 1.
- [ ] Self-probe: drop the `^{commit}` test, and report the tree-hash case green as the observed wrong result.
