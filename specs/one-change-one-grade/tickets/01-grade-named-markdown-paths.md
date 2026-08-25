# Grade named Markdown paths through a plumbing verb

Line: sonnet / low.

Blocked by: none
Writes: internal/prose/walk.go, internal/prose/subject.go (new), internal/prose/subject_test.go (new), internal/gate/gate_prose.go (new), internal/gate/gate_prose_test.go (new), cmd/bench/main.go, cmd/bench/command_registry_test.go

## What to build

An operator grades a short list of Markdown files in about one second. The new
`bench gate-prose <root> [--] [path...]` verb grades the named paths through the
same walker, classifier, and exclusion list as the whole-tree prose check. Exit 0
is a clean list, exit 1 is a list with findings, and exit 2 is a usage error. An
unknown flag is a usage error, and it never reads as a clean list.

Add a per-subject grader to `internal/prose`. That grader owns the exclusion
list, the byte classifier, and the finding render for one named path. `Grade`
composes the same per-subject grader for the whole tree, so the prose rule keeps
one source. Register `gate-prose` in `cmd/bench/main.go` beside the `gate-go`
row, with the internal inventory and the plumbing exemption.

The lane ticket calls this verb, so the argv keeps the shape
`<bench-run-binary> gate-prose <root> -- <path>...`.

## Acceptance

- [ ] OG08 reports no finding for a named path that `.bench/prose-exclusions` lists.
- [ ] OG33 exits 1 and names the file and the line for a named file with a 27-word sentence.
