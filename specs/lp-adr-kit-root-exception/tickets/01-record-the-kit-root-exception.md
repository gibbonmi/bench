# Record the kit-root exception to the one-source rule

Blocked by: none
Writes: docs/adr/0029-the-kit-rule-keeps-a-wrapper-form-and-a-go-form.md
Covers: none

## What to build

The kit rule has two production forms: the `kit_dir` rule in `bin/bench.sh` and `gate.KitForRoot`.
The wrapper differential test `TestShellWrapperKitMatchesKitForRoot` holds the two forms to one answer.
The FT335 merge-target-lane spec recorded reviewer decision D1, which accepts this pair as an exception to the one-source rule.
The retirement of that spec removed D1, and no ADR holds it.

Add ADR 0029 in `docs/adr/`.
The ADR states the kit rule, the two forms by role, the reason for the exception, and the guard of their agreement.
The ADR holds no file path and no code snippet.
The `docs/adr/` folder has no index, so no other file changes.

## Acceptance

- [ ] ADR 0029 states that the reviewer accepts the wrapper form and the Go form of the kit rule as an exception to the one-source rule.
- [ ] ADR 0029 names the wrapper differential test as the guard of the agreement between the two forms.
- [ ] ADR 0029 names the Go owner that the wrapper consumes as the rejected option.
- [ ] ADR 0029 holds no file path and no code snippet.
- [ ] The `prose`, `prose-mechanics`, and `ticket-grammar` checks pass.
