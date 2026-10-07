# Document the reusable implement-spec run options

Blocked by: none
Writes: .agents/commands/bench-implement-spec.md, internal/anchors/
Covers: none

## What to build

The reviewer starts a full build with options that the command does not state: `--reviewer <model> <effort>`, `--consultant <model> <effort>`, and `--auto-approve`. Only `--full` and `--delegate` have a section. A later reviewer cannot find the options, and a later orchestrator cannot read what each one permits.

State each reusable option of a `--full` run once in the implementation command:

- `--reviewer <model> <effort>` sets the line of every review axis in the run.
- `--consultant <model> <effort>` sets the line of every read-only diagnostic consultation in the run.
- `--auto-approve` is the reviewer's approval, given in advance, for ticket fence expansions, spec and ticket expansions, and repair rounds past the bounded repair allowance. Under it, the orchestrator fixes a found defect at once instead of deferring it.

Keep the option text short. The option itself carries the approval, so do not add prose that explains why. Keep the file inside its current prose line budget; shorten or merge existing sentences without loss of a rule or an anchor needle. Do not raise the budget. If an anchor needle must move, move it with the sentence that it pins.

## Acceptance

- [ ] The implementation command states `--full`, `--delegate`, `--reviewer`, `--consultant`, and `--auto-approve`, and each option has one statement of its effect.
- [ ] The `--auto-approve` statement names fence expansions, spec and ticket expansions, and repair rounds, and it states the fix-at-once preference.
- [ ] The file stays inside its prose line budget, and every anchor needle for the file still passes.
