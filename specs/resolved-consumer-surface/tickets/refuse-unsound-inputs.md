# Refuse unsound inputs fail-closed

Blocked by: add-consumers-command-surface.md
Writes: internal/consumers/

## What to build

The fail-closed rim. An ill-typed tree, a missing `go` binary, and a query
resolved to a non-Go file each emit a structured stdout refusal at exit 1.
Each refusal names its cause: the first error position, the missing binary,
or the unsupported language. Refusals append no disclosure.

## Acceptance

- [ ] CS7: a query resolved to a `.ts` file emits a structured stdout
      refusal at exit 1 naming the language.
- [ ] CS8: a fixture with one type error emits a structured stdout refusal
      at exit 1 naming that first error position.
- [ ] CS9: with `go` absent from `PATH` the command emits a structured
      stdout refusal at exit 1 naming the missing binary.
