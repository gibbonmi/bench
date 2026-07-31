# Report the hard-cut binding migration

Blocked by: Resolve the line through the harness matrix

## What to build

`bench doctor` gains two binding rows so a maintainer migrating a binding is
told the exact rewrite without the tool touching their file. It detects the six
retired keys and reports their exact replacements, and it names any known
harness whose column is unbound along with the action that would bind it. It
never writes `lines.env`.

Covers story 8 and story 3's doctor row.

## Acceptance

- [ ] All six retired keys are named with their exact rewrites and `lines.env` is
      left byte-identical.
- [ ] A second `bench doctor` run reports the same rewrites and still changes
      nothing.
- [ ] `bench doctor` names the unbound OpenCode column and the action that would
      bind it.
