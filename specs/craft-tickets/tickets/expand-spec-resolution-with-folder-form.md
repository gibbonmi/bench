# Expand spec resolution with folder form

Blocked by: none

## What to build

Add the folder-form spec path beside the existing flat form so folder specs can
resolve and enumerate without breaking any current caller. This is transitional
expand–contract sequencing: the later contract ticket removes live flat-form
support and establishes the final fail-closed behavior.

## Acceptance

- [x] Bare slugs resolve folder specs from the repository root and a deeper cwd.
- [x] Explicit paths remain path-first.
- [x] Folder specs appear in `Facts` with the folder name as slug.
- [x] Existing flat-form behavior remains green until the contract ticket.
