# Enforce TEST bindings with structural and name refusals

Blocked by: Migrate the 33 behavior-owned fixtures to TEST bindings

## What to build

Stories 3 and 4 of `specs/ft91-gate-fastpath/spec.md`: the optional `TEST`
read becomes mandatory for behavior-owned fixtures. Structural defects —
missing or empty `TEST` on a behavior-owned fixture, stray `TEST` outside the
family, `TEST` at a family or package level inside it — red pre-compile beside
the existing `assertContractScopes` refusals, all reported together. Declared
owners are graded post-compile against the compiled binary's `-test.list '.*'`
output (one list call per bound package, after its compile, before any
baseline or fixture run, sweep worker budget, membership accepts only lines
beginning `Test`, nonzero exit is a sweep error naming the package); unknown
names are reported together, one diagnostic per fixture.

## Acceptance

- [ ] Missing `TEST`, empty `TEST`, stray `TEST` outside the behavior-owned
      family, and `TEST` at a family or package level inside it are four
      distinct loud reds naming the offending path, reported together.
- [ ] An owner absent from the compiled binary's `-test.list` output is a
      named refusal before any baseline or fixture run; multiple unknown
      owners report together, one diagnostic per fixture.
- [ ] No silent fallback to a full-package run exists for a behavior-owned
      fixture with a binding defect.
