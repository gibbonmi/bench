# FT215 focused test input facts

Produced 2026-08-27 by one read-only mid-tier delegate. The coordinator
verified the cited sources and the live Go package producer at `97bb035`.

## Current producers

`bench test` accepts one optional package expression. An absent expression
selects `./...`. A bare existing directory gains a `./` prefix. Every
other expression passes to Go unchanged
(`internal/testreport/testreport.go:25-49,108-123`).

The command has no test-pattern input. `--full` changes diagnostic detail
only. The renderer requires terminal package events and reports package,
failure, and skip rows (`internal/testreport/testreport.go:83-105,256-353`).

`bench diff` owns the coherent diff subjects. Live mode combines committed,
staged, tracked-worktree, and untracked paths. An explicit base keeps that
live subject. A base and source tip form an immutable subject
(`internal/diff/diff.go:27-35,74-104`;
`internal/diff/range.go:98-115,193-212`).

Git diff uses `--no-renames`. A rename therefore becomes one deletion and
one addition. The current path projection discards change status
(`internal/diff/snapshot.go:31-60,375-394`).

The live Go 1.25 producer verified one module. `go list` exposes each
package directory, production imports, test imports, and external-test
imports. It includes test-only and canary-fixture packages. No current
`bench test` path resolves a diff through that graph.

The conformance registry accepts a named check or an ordered check set. Its
`Inputs` values are labels, not path selectors. No current producer maps a
diff path to a conformance check
(`internal/conformance/registry/registry.go:66-153`;
`internal/conformance/checks_test.go:130-215`).

## Producer-derived partitions

| input | derivable subject | current result | unresolved choice |
|---|---|---|---|
| explicit package | one Go package expression | one fresh `go test` process | flag grammar and outside-module posture |
| explicit test pattern | none | unsupported | pattern syntax and zero matches |
| live default diff | branch base through the live checkout | complete repo-relative path set | whether changed mode uses it |
| live explicit-base diff | named base through the live checkout | complete repo-relative path set | whether changed mode accepts it |
| frozen base and source tip | immutable two-commit range | complete repo-relative path set | whether changed mode accepts it |
| staged or unstaged tracked path | present in the live subject | raw path only | path-to-package rule |
| untracked path | present in the live subject | raw path, including special nodes | regular-file and refusal posture |
| renamed path | deletion plus addition | two paths, no rename identity | old, new, or both package graphs |
| deleted Go path or package | old path only | no live package is guaranteed | baseline lookup or refusal |
| non-Go or root path | raw path only | no package mapping | empty result, conformance, or refusal |
| multiple packages | union is derivable from the Go graph | current command accepts one expression | order, deduplication, and invocation form |
| reverse dependent | production and test edges are available | no resolver exists | which edge classes participate |
| no runnable package | detectable before or after Go execution | current no-terminal result exits 1 | empty success or refusal |
| named conformance check | registry name | internal selection transport exists | public focused-test grammar |
| diff-derived conformance | no selector exists | unsupported | whether to add path selectors |

## Exact edge facts

An explicit `../` or absolute package expression passes to Go. The current
command does not enforce the repository or module boundary
(`internal/testreport/testreport.go:114-122`).

A live untracked path can be a special file. The coherent diff producer
reports that path, but a future package resolver must refuse it before a read
(`internal/git/status.go:201-251`).

A deleted package cannot resolve only against the live graph. The resolver
must inspect the baseline graph or refuse the subject. Silent omission would
make a changed-input result incomplete.

The gate gives the ordinary conformance entry its root and dev tier. The
current `bench test` command does not install that environment
(`internal/gate/phases.go:107-136`;
`internal/testreport/testreport.go:32-52`).

## Reviewer decisions

- Ticket #7 chooses the live and frozen subjects and the empty-result posture.
- Ticket #8 chooses the public grammar and positional compatibility.
- Ticket #9 decides whether changed mode selects conformance checks.
- The selected policy must cover deletions, renames, non-Go paths, and
  outside-module expressions.

## Verification

The coordinator ran
`go list -f '{{.ImportPath}} {{.Dir}} {{.Imports}} {{.TestImports}} {{.XTestImports}}' ./...`
with Go 1.25.0 at `97bb035`. The command completed without a tree change.
