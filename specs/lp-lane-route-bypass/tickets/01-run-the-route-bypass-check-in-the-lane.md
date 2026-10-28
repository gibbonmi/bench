# Run the route bypass check in the commit lane

Blocked by: none
Writes: internal/conformance/registry/route_bypass.go, internal/conformance/refusal_route_bypass_test.go, internal/gate/lane.go, internal/gate/lane_select.go, internal/gate/lane_run_test.go, projects/benchkit.md, CHANGELOG.md
Covers: none

## What to build

The fast lane of `bench commit` runs the route bypass check (RR49) when a commit changes a file that the check scans. Thus a write-verb author sees a route literal outside the registry at the commit. Before, the author saw it first at the chunk checkpoint or at the landing gate.

The conformance registry holds the list of scanned packages and the rule that names a scanned file. The RR49 scan and the new lane path class both read that one source. The new lane check `route-bypass` runs the RR49 test alone through `bench test --package`, as the `fixture-bite` check runs its test. The new path class `write-verb-source` selects the check. The lane table of the kit profile shows the new row, and the `profile-lane-table` check grades that row.

## Acceptance

- [ ] A changed or deleted production Go file in each scanned package selects the `route-bypass` check under the `write-verb-source` class.
- [ ] A test file in a scanned package, a Go file in a subdirectory that the scan does not read, and a Go file in a different package do not select the `route-bypass` check.
- [ ] A route literal in a scanned package makes `bench commit` refuse on the `route-bypass` check with the RR49 fault.
- [ ] The lane table in `projects/benchkit.md` shows the `route-bypass` row, and `bench test --check profile-lane-table` passes.
