# Learnings — usage journal

## 2026-08-21 — a spec's tickets left the write-spec phase failing bench preflight build [open]

`/bench-implement-spec specs/learnings-dated-line-visibility/spec.md --full`
stopped at its first entry check. `bench preflight build` returned
`rows-owned,red,"declared row(s) cited by no ticket file: DL1 … DL21"`: the one
ticket's acceptance lines are prose, and the exemplar shape
(`specs/land-executable-freshness/tickets/01-*.md`) prefixes each line with its
row ID. The mapping was already one-to-one, so the defect was pure citation.
What happened: the write-spec phase ran a verification round, accepted the
ticket, and exited, but never ran the build phase's own entry check against its
output. The build phase then refused work it does not own the fix for. Right
behavior: `/bench-write-spec` runs `bench preflight build <slug>` before it
exits, so a red lands in the phase that owns the ticket file. Proposed rule
change: `/bench-write-spec`'s exit adds that command as a required check, the
same way the build phase requires it at entry.

## 2026-08-21 — a spec anchored a rule on a literal without tracing the scaffold that emits it [open]

The FT243 spec's second rule opens on `capture/learnings.md`'s
`<!-- entries below -->` marker. The spec pinned the anchor as "the marker
appears above the first line starting `## `". Review round 1 caught that
`internal/adopt`'s `scaffoldLearnings` prints the marker *below* its worked
example `## <date> - <short title>  [open]`, so under the written predicate the
rule would never open on a freshly scaffolded journal — and that journal was the
rule's only red. Two coverage rows, DL27 and DL30, were mutually unsatisfiable.
What happened: the author read the scaffold to find the marker's literal and
read the parser to find `isTemplatePlaceholder`, but never traced the anchor
against the scaffold's actual byte order, so a contradiction between two
already-read files survived into the map. Right behavior: when a rule anchors on
a literal that a scaffold or template also emits, walk the rule over that
scaffold's real output before locking the rows. Proposed rule change:
`craft-spec` gains one line under the coverage map — a row whose fixture is
produced by a generator in the tree is written against that generator's actual
output, not against a hand-written idea of it.

## 2026-08-21 — a tilde-prefixed worktree path made cp build a stray directory [open]

`bench worktree path "<label>"` prints a portable path that begins with a
literal `~`. The session used it as `WT=$(bench worktree path …)` and then
`cp -r specs/<slug> "$WT/specs/"`, so the shell never expanded the tilde and
`cp` created `./~/.bench/worktrees/…/specs/` inside the repository. The mistake
was invisible until the next command failed for a missing spec, and it left an
untracked `~/` directory one `rm -rf` away from looking like a home-directory
wipe. What happened: a command whose output is documented as portable was
consumed as if it were a shell-ready path. Right behavior: expand the tilde
before use, or address the worktree only through
`bench worktree exec "<label>" -- <command>`, which resolves the path itself.
Proposed rule change: none for the agent yet — the reviewer's call is whether
`bench worktree path` should emit an absolute path, or name the expansion in
its own help, since a portable path that no shell accepts has one safe consumer
and many unsafe ones.
