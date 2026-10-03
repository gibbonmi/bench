# Grade tracked paths only in the prose live-tree check

Blocked by: none
Writes: internal/prose/walk.go, internal/prose/walk_test.go
Covers: none

## What to build

Roadmap row FT361 owns this fix. The prose live-tree check walks every `.md`
file on disk under the graded root. It therefore grades a git-ignored file,
such as the `capture/session-handoff.md` that `bench handoff` writes. A plain
`go test ./...` in the primary checkout then goes red with no code defect.

The walk in `prose.Grade` selects only the paths that git tracks under the
graded root. An untracked or ignored `.md` file is outside the grade. The
existing skip rules for `testdata`, `node_modules`, `dist`, and a directory
link stay in force. The selection composes an existing tracked-file reader in
the tree. It does not add a second `ls-files` parser.

A graded root that is not a git work tree keeps one decided result. The
author reads each caller of `prose.Grade` and states that result in the
ticket evidence.

## Acceptance

- [ ] A git-ignored `.md` file with a prose fault under the graded root gives no diagnostic.
- [ ] A tracked `.md` file with the same fault still gives its diagnostic.
- [ ] An untracked `.md` file that is not ignored gives no diagnostic.
- [ ] The existing prose suite and the root conformance pass stay green.

## Evidence

The walk selects by the tracked set only when the graded root is the top of
its own git work tree. The walk grades any other root whole. Such a root is a
root that git cannot list, or a root below the top of an outer work tree. A
temporary directory under an in-repo `TMPDIR` is an example of the second
kind. A work tree top that tracks no file gives no diagnostic.

Each caller of `prose.Grade` gets this result:

- `runProseCheck` in `internal/testreport/command.go` grades the repository
  root in a real run. Its own test grades a plain temporary directory, which
  the walk grades whole.
- `checkProseMechanics` in `internal/conformance/prose_mechanics_test.go`
  grades the kit root, which is a work tree top in a checkout. The canary and
  learning-journal tests grade plain temporary directories, which the walk
  grades whole.

Two `TestGrade` rows hold this result. The row for a root below the top of
its work tree expects the fault in its untracked file. The row for a work
tree that tracks no file expects no diagnostic.
