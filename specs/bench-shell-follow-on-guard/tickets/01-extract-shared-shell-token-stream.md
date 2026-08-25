# Extract the shared shell token stream

Blocked by: none
Writes: internal/shellcommand (new), internal/gitguard

## What to build

Publish one shell token-stream interface for words, control operators, redirections, and simple-command spans.
Keep quote folding and heredoc-body removal behind that interface.

Migrate the destructive-Git analyzer to the interface without changing its classifications.
The `02-refuse-bench-shell-follow-ons.md` ticket consumes the interface without changing it.

## Acceptance

- [ ] The token stream distinguishes words, control operators, redirections, and simple-command spans.
- [ ] The token stream folds quotes and removes heredoc bodies while it preserves the outer heredoc operator.
- [ ] The existing destructive-Git classification matrix returns the same deny labels after migration.
