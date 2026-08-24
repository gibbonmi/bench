# Skip heredoc bodies in the git guard

Blocked by: none
Writes: internal/gitguard/tokenize.go, internal/gitguard/tokenize_test.go, internal/gitguard/gitguard_test.go

## What to build

The core classifier treats a heredoc body as data. A command that writes a
file whose contents mention a destructive git command is allowed. The
classifier removes each heredoc body (from the line after the `<<`/`<<-`
operator's line, through the delimiter line) before tokenization. A
destructive git command outside the body still blocks, including one on the
same line as the heredoc operator.

## Acceptance

- [ ] `cat > f <<'EOF'` with a body line `git push --force` classifies as allow.
- [ ] `cat > f <<EOF` (unquoted delimiter) with a destructive body classifies as allow.
- [ ] `<<-` with a tab-indented delimiter classifies as allow.
- [ ] `git push` after the delimiter line, or on the operator's own line, still blocks.
- [ ] `<<<` (herestring) keeps its current classification.
