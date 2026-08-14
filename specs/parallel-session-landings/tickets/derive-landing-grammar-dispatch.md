# Derive landing grammar dispatch

Blocked by: none
Writes: internal/usage, internal/worktree

## What to build

Derive landing mode selection and required-value validation from the public
grammar metadata so flag-value and count knowledge has one owner.

## Acceptance

- [ ] Land and reauthorize reject missing required values through grammar-owned metadata.
- [ ] Every applicable flag-value-only PL27 case remains red against positional misclassification.
